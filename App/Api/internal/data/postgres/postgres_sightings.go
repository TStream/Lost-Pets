package postgres

import (
	"database/sql"
	domain "lostpets"
	"lostpets/internal"

	"github.com/lib/pq"
	_ "github.com/lib/pq" //postgres driver
)

type (
	sighting struct {
		domain.Sighting
		PetID int
	}

	internalSightingAggregate struct {
		domain.Sighting
		PictureID int
		PetID     int
		PetName   string
		PetColor  string
		Marks     string
		TypeID    int
		Type      string
		Shape     string
		TagID     int
		TagColor  string
		Text      string
		PetBreeds pq.StringArray
	}
)

const sightingSelect = `SELECT
sightings.id,
sightings.name,
in_custody,
email
guid,
date,
location,
pets.id as pet_id,
picture_id,
types.id as type_id,
types.name as type,
pets.name as pet_name,
pets.color as pet_color,
tags.id as tag_id,
marks,
shape,
tags.color as tag_color,
text,
array_remove(ARRAY_AGG(distinct(pet_breeds.name)),NULL) as pet_breeds
FROM sightings 
LEFT JOIN pets ON pets.id = sightings.pet_id
LEFT JOIN types ON types.id = pets.type_id
LEFT JOIN pet_breeds ON pets.id = pet_breeds.pet_id
LEFT JOIN tags ON tags.pet_id = pets.id `

const sightingsGroupBy = `
GROUP BY 
sightings.id,
sightings.name,
in_custody,
email,
guid,
date,
location,
pets.id,
picture_id,
types.id,
types.name,
pets.name,
pets.color,
marks,
shape,
tags.id,
tags.color,
text `

var sightingsFieldMap = map[string]string{
	"id":           "sightings.id",
	"name":         "sightings.name",
	"incustody":    "in_custody",
	"email":        "email",
	"guid":         "guid",
	"location":     "location",
	"petid":        "pets.id",
	"petpictureid": "picture_id",
	"pettype":      "types.name",
	"petname":      "pets.name",
	"petcolor":     "pets.color",
	"petmarks":     "marks",
	"petbreeds":    "pet_breeds.name",
	"pettagid":     "tags.id",
	"pettagshape":  "shape",
	"pettagcolor":  "tags.color",
	"pettagtext":   "text",
}

func (db *DB) getSighting(filters domain.FilterMap) (*domain.Sighting, error) {
	queryStr, args, err := db.buildQuery(sightingsFieldMap, filters)
	if err != nil {
		return nil, err
	}

	queryStr = queryStr + sightingsGroupBy

	aggregate := &internalSightingAggregate{}

	err = db.Get(aggregate, sightingSelect+queryStr, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	aggregate.Pet = domain.Pet{
		ID:        aggregate.PetID,
		PictureID: aggregate.PictureID,
		Name:      aggregate.PetName,
		Color:     aggregate.PetColor,
		Marks:     aggregate.Marks,
		Type:      aggregate.Type,
		TypeID:    aggregate.TypeID,
		Breeds:    aggregate.PetBreeds,
		Tag: domain.Tag{
			ID:    aggregate.TagID,
			Shape: aggregate.Shape,
			Color: aggregate.TagColor,
			Text:  aggregate.Text,
		},
	}

	return &aggregate.Sighting, nil
}

func (db *DB) GetSightingByGUID(guid string) (*domain.Sighting, error) {
	filter := domain.Filter{
		Comparator: "=",
		Value:      guid,
	}
	filters := domain.FilterMap{}
	filters["guid"] = []domain.Filter{filter}

	return db.getSighting(filters)
}

func (db *DB) GetMatchingPostings(id int) ([]domain.Posting, error) {
	matchesQuery := matchesSelect + " WHERE sightings_id = $1 "
	matches := []matches{}

	err := db.Select((&matches), matchesQuery, id)
	if err != nil {
		return nil, err
	}

	postingIDs := []int{}
	for _, m := range matches {
		postingIDs = append(postingIDs, m.PostingsID)
	}

	pFilter := domain.Filter{
		Comparator: "in",
		Value:      postingIDs,
	}
	pFilters := domain.FilterMap{}
	pFilters["id"] = []domain.Filter{pFilter}

	postings, err := db.GetAllPostings(pFilters)
	if err != nil {
		return nil, err
	}

	return postings, nil
}

func (db *DB) GetSightingByID(id int) (*domain.Sighting, error) {
	filter := domain.Filter{
		Comparator: "=",
		Value:      id,
	}
	filters := domain.FilterMap{}
	filters["id"] = []domain.Filter{filter}

	return db.getSighting(filters)
}

func (db *DB) GetSightingByEmail(email string) (*domain.Sighting, error) {
	filter := domain.Filter{
		Comparator: "=",
		Value:      email,
	}
	filters := domain.FilterMap{}
	filters["email"] = []domain.Filter{filter}

	return db.getSighting(filters)
}

func (db *DB) GetAllSightings(filters ...domain.FilterMap) ([]domain.Sighting, error) {
	queryStr, args, err := db.buildQuery(sightingsFieldMap, filters...)
	if err != nil {
		return nil, err
	}
	queryStr = queryStr + sightingsGroupBy

	aggregates := []internalSightingAggregate{}

	err = db.Select((&aggregates), sightingSelect+queryStr, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sightings := []domain.Sighting{}

	for _, a := range aggregates {
		a.Pet = domain.Pet{
			ID:        a.PetID,
			PictureID: a.PictureID,
			Name:      a.PetName,
			Color:     a.PetColor,
			Marks:     a.Marks,
			Type:      a.Type,
			Breeds:    a.PetBreeds,
			Tag: domain.Tag{
				ID:    a.TagID,
				Shape: a.Shape,
				Color: a.TagColor,
				Text:  a.Text,
			},
		}
		sightings = append(sightings, a.Sighting)
	}

	return sightings, nil
}

func (db *DB) AddSighting(newSighting *domain.Sighting) error {
	guid, err := internal.NewUUID()
	if err != nil {
		return err
	}
	newSighting.GUID = guid

	err = db.addPet(&newSighting.Pet)
	if err != nil {
		return err
	}

	dbSighting := sighting{
		Sighting: *newSighting,
		PetID:    newSighting.Pet.ID,
	}

	query := `INSERT INTO sightings(
		guid, pet_id, date, location, name, email, in_custody)
		VALUES (:guid, :pet_id, :date, :location, :name, :email, :in_custody) RETURNING id;`

	rows, err := db.NamedQuery(query, dbSighting)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&newSighting.ID)
	} else {
		return errID
	}

	return nil
}
