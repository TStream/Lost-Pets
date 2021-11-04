package http

import (
	domain "lostpets"
	"net/http"
	"path"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	sightingsHandler struct {
		logger  domain.StructuredLogger
		router  *echo.Echo
		repo    domain.LostPetsRepo
		emailer emailer
	}

	apiSightingResponse struct {
		Sighting  *apiSighting   `json:"sighting,omitempty"`
		Sightings *[]apiSighting `json:"sightings,omitempty"`
	}

	apiPostSighting struct {
		apiSighting
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}

	apiSighting struct {
		ID        int      `json:"id,omitempty"`
		Pet       apiPet   `json:"pet,omitempty"`
		Date      Datetime `json:"date,omitempty"`
		Location  string   `json:"location,omitempty"`
		InCustody bool     `json:"inCustody"`
	}
)

func (h *sightingsHandler) initRoute(path string) {
	h.router.GET(path+"/:id", h.handleGetByID())
	h.router.GET(path+"/private/:guid", h.handleGetByGUID())
	h.router.GET(path+"/private/:guid/matches", h.handleGetAllMatches())
	h.router.GET(path, h.handleGetAll())
	h.router.POST(path, h.handleCreateSighting(path+"/private/"))
}

func (h *sightingsHandler) handleGetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		sightingID := c.Param("id")
		id, err := strconv.Atoi(sightingID)
		if err != nil {
			return err
		}

		sighting, err := h.repo.GetSightingByID(id)
		if err != nil {
			return err
		}
		if sighting == nil {
			return c.NoContent(http.StatusNotFound)
		}

		resp := apiSightingResponse{
			Sighting: toAPISighting(*sighting),
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *sightingsHandler) handleGetByGUID() echo.HandlerFunc {
	return func(c echo.Context) error {
		sightingGUID := c.Param("guid")
		sighting, err := h.repo.GetSightingByGUID(sightingGUID)
		if err != nil {
			return err
		}
		if sighting == nil {
			return c.NoContent(http.StatusNotFound)
		}

		resp := apiSightingResponse{
			Sighting: toAPISighting(*sighting),
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *sightingsHandler) handleGetAll() echo.HandlerFunc {
	//TODO add query support
	return func(c echo.Context) error {
		sightings, err := h.repo.GetAllSightings()
		if err != nil {
			return err
		}

		apiSightings := []apiSighting{}
		for _, s := range sightings {
			apiSightings = append(apiSightings, *toAPISighting(s))
		}
		resp := apiSightingResponse{
			Sightings: &apiSightings,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *sightingsHandler) handleGetAllMatches() echo.HandlerFunc {
	return func(c echo.Context) error {
		sGUID := c.Param("guid")
		sighting, err := h.repo.GetSightingByGUID(sGUID)
		if err != nil {
			return err
		}

		matches, err := h.repo.GetMatchingPostings(sighting.ID)
		if err != nil {
			return err
		}
		apiPostings := []apiPosting{}
		for _, m := range matches {
			apiPostings = append(apiPostings, *toAPIPosting(m))
		}
		resp := apiPostingResponse{
			Postings: &apiPostings,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *sightingsHandler) handleCreateSighting(location string) echo.HandlerFunc {
	//TODO add validation
	return func(c echo.Context) error {
		newSighting := new(apiPostSighting)
		if err := c.Bind(newSighting); err != nil {
			return err
		}
		dSighting := toDomainSighting(*newSighting)
		err := h.repo.AddSighting(dSighting)
		if err != nil {
			return err
		}

		//run search in 'background'
		go h.searchForMatches(*dSighting)

		c.Response().Header().Set(echo.HeaderLocation, path.Join(location, dSighting.GUID))
		return c.NoContent(http.StatusCreated)
	}
}

func toDomainSighting(api apiPostSighting) *domain.Sighting {
	return &domain.Sighting{
		InCustody: api.InCustody,
		Posting: domain.Posting{
			ID:       api.ID,
			Date:     api.Date.Time,
			Location: api.Location,
			Name:     api.Name,
			Email:    api.Email,
			Pet: domain.Pet{
				ID:        api.Pet.ID,
				PictureID: api.Pet.PictureID,
				Name:      api.Pet.Name,
				Color:     api.Pet.Color,
				Breeds:    api.Pet.Breeds,
				Marks:     api.Pet.Marks,
				Type:      api.Pet.Type,
				TypeID:    api.Pet.TypeID,
				Tag: domain.Tag{
					ID:    api.Pet.Tag.ID,
					Shape: api.Pet.Tag.Shape,
					Color: api.Pet.Tag.Color,
					Text:  api.Pet.Tag.Text,
				},
			},
		},
	}
}

func toAPISighting(d domain.Sighting) *apiSighting {
	return &apiSighting{
		InCustody: d.InCustody,
		ID:        d.ID,
		Date:      Datetime{Time: d.Date},
		Location:  d.Location,
		Pet: apiPet{
			ID:        d.Pet.ID,
			PictureID: d.Pet.PictureID,
			Name:      d.Pet.Name,
			Color:     d.Pet.Color,
			Breeds:    d.Pet.Breeds,
			Marks:     d.Pet.Marks,
			Type:      d.Pet.Type,
			TypeID:    d.Pet.TypeID,
			Tag: apiTag{
				ID:    d.Pet.Tag.ID,
				Shape: d.Pet.Tag.Shape,
				Color: d.Pet.Tag.Color,
				Text:  d.Pet.Tag.Text,
			},
		},
	}
}

//TODO Clean up and come up with a better way to search
func (h *sightingsHandler) searchForMatches(sighting domain.Sighting) {
	//create filters
	filters := []domain.FilterMap{}

	locationFilter := domain.FilterMap{}
	locationFilter["location"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      sighting.Location,
		},
	}
	filters = append(filters, locationFilter)

	typeFilter := domain.FilterMap{}
	typeFilter["pettype"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      sighting.Pet.Type,
		},
	}
	filters = append(filters, typeFilter)

	petNameFilter := domain.FilterMap{}
	petNameFilter["petname"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      sighting.Pet.Name,
		},
	}
	filters = append(filters, petNameFilter)

	petColorFilter := domain.FilterMap{}
	petColorFilter["petcolor"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      sighting.Pet.Color,
		},
	}
	filters = append(filters, petColorFilter)

	petMarksFilter := domain.FilterMap{}
	petMarksFilter["petmarks"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      sighting.Pet.Marks,
		},
	}
	filters = append(filters, petMarksFilter)

	//search, filters should be separated by ORs
	matches, err := h.repo.GetAllPostings(filters...)
	if err != nil {
		h.logger.Error(err.Error())
	}

	for _, m := range matches {
		h.logger.Info("Adding Match p:%d, s:%d", m.ID, sighting.ID)
		err := h.repo.AddMatch(m.ID, sighting.ID)
		if err != nil {
			h.logger.Error(err.Error())
		}
	}
	//email will have link to 'private' page in UI which will query for found matches
	h.emailer.emailMatches("Sighting", sighting.Email, sighting.GUID)
}
