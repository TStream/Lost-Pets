package http

import (
	domain "lostpets"
	"net/http"
	"path"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	postingsHandler struct {
		logger  domain.StructuredLogger
		router  *echo.Echo
		repo    domain.LostPetsRepo
		emailer emailer
	}

	apiPostingResponse struct {
		Posting  *apiPosting   `json:"posting,omitempty"`
		Postings *[]apiPosting `json:"postings,omitempty"`
	}

	apiPostPosting struct {
		apiPosting
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}

	apiPosting struct {
		ID       int      `json:"id,omitempty"`
		Pet      apiPet   `json:"pet,omitempty"`
		Date     Datetime `json:"date,omitempty"`
		Location string   `json:"location,omitempty"`
	}

	apiPet struct {
		ID        int      `json:"id,omitempty"`
		PictureID int      `json:"pictureId,omitempty"`
		Name      string   `json:"name,omitempty"`
		Color     string   `json:"color,omitempty"`
		Marks     string   `json:"marks,omitempty"`
		Type      string   `json:"type,omitempty"`
		TypeID    int      `json:"typeId,omitempty"`
		Breeds    []string `json:"breeds,omitempty"`
		Tag       apiTag   `json:"tag,omitempty"`
	}

	apiTag struct {
		ID    int    `json:"id,omitempty"`
		Shape string `json:"shape,omitempty"`
		Color string `json:"color,omitempty"`
		Text  string `json:"text,omitempty"`
	}
)

func (h *postingsHandler) initRoute(path string) {
	h.router.GET(path+"/:id", h.handleGetByID())
	h.router.GET(path+"/private/:guid", h.handleGetByGUID())
	h.router.GET(path+"/private/:guid/matches", h.handleGetAllMatches())
	h.router.GET(path, h.handleGetAll())
	h.router.POST(path, h.handleCreatePosting(path+"/private/"))
}

func (h *postingsHandler) handleGetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		postingID := c.Param("id")
		id, err := strconv.Atoi(postingID)
		if err != nil {
			return err
		}

		posting, err := h.repo.GetPostingByID(id)
		if err != nil {
			return err
		}
		if posting == nil {
			return c.NoContent(http.StatusNotFound)
		}

		resp := apiPostingResponse{
			Posting: toAPIPosting(*posting),
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *postingsHandler) handleGetByGUID() echo.HandlerFunc {
	return func(c echo.Context) error {
		postingGUID := c.Param("guid")
		posting, err := h.repo.GetPostingByGUID(postingGUID)
		if err != nil {
			return err
		}
		if posting == nil {
			return c.NoContent(http.StatusNotFound)
		}

		resp := apiPostingResponse{
			Posting: toAPIPosting(*posting),
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *postingsHandler) handleGetAll() echo.HandlerFunc {
	//TODO add query support
	return func(c echo.Context) error {
		postings, err := h.repo.GetAllPostings()
		if err != nil {
			return err
		}

		apiPostings := []apiPosting{}
		for _, p := range postings {
			apiPostings = append(apiPostings, *toAPIPosting(p))
		}

		resp := apiPostingResponse{
			Postings: &apiPostings,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
func (h *postingsHandler) handleGetAllMatches() echo.HandlerFunc {
	return func(c echo.Context) error {
		postingGUID := c.Param("guid")
		posting, err := h.repo.GetPostingByGUID(postingGUID)
		if err != nil {
			return err
		}
		matches, err := h.repo.GetMatchingSightings(posting.ID)
		if err != nil {
			return err
		}
		sightings := []apiSighting{}
		for _, m := range matches {
			sightings = append(sightings, *toAPISighting(m))
		}

		resp := apiSightingResponse{
			Sightings: &sightings,
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func (h *postingsHandler) handleCreatePosting(location string) echo.HandlerFunc {
	//TODO add validation
	return func(c echo.Context) error {
		newPosting := new(apiPostPosting)
		if err := c.Bind(newPosting); err != nil {
			return err
		}
		dPosting := toDomainPosting(*newPosting)
		err := h.repo.AddPosting(dPosting)
		if err != nil {
			return err
		}

		//run search in 'background'
		go h.searchForMatches(*dPosting)

		c.Response().Header().Set(echo.HeaderLocation, path.Join(location, dPosting.GUID))
		return c.NoContent(http.StatusCreated)
	}
}

func toDomainPosting(api apiPostPosting) *domain.Posting {
	return &domain.Posting{
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
	}
}

func toAPIPosting(d domain.Posting) *apiPosting {
	return &apiPosting{
		ID:       d.ID,
		Date:     Datetime{Time: d.Date},
		Location: d.Location,
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
func (h *postingsHandler) searchForMatches(posting domain.Posting) {
	//create filters
	filters := []domain.FilterMap{}

	locationFilter := domain.FilterMap{}
	locationFilter["location"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      posting.Location,
		},
	}
	filters = append(filters, locationFilter)

	typeFilter := domain.FilterMap{}
	typeFilter["pettype"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      posting.Pet.Type,
		},
	}
	filters = append(filters, typeFilter)

	petNameFilter := domain.FilterMap{}
	petNameFilter["petname"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      posting.Pet.Name,
		},
	}
	filters = append(filters, petNameFilter)

	petColorFilter := domain.FilterMap{}
	petColorFilter["petcolor"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      posting.Pet.Color,
		},
	}
	filters = append(filters, petColorFilter)

	petMarksFilter := domain.FilterMap{}
	petMarksFilter["petmarks"] = []domain.Filter{
		{
			Comparator: "=",
			Value:      posting.Pet.Marks,
		},
	}
	filters = append(filters, petMarksFilter)

	//search, filters should be separated by ORs
	matches, err := h.repo.GetAllSightings(filters...)
	if err != nil {
		h.logger.Error(err.Error())
	}

	for _, m := range matches {
		h.logger.Info("Adding Match p:%d, s:%d", posting.ID, m.ID)
		err := h.repo.AddMatch(posting.ID, m.ID)
		if err != nil {
			h.logger.Error(err.Error())
		}
	}
	//email will have link to 'private' page in UI which will query for found matches
	h.emailer.emailMatches("Postings", posting.Email, posting.GUID)
}
