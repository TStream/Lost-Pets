package postgres

import (
	"database/sql"
	"fmt"
	domain "lostpets"
	"lostpets/internal"

	"github.com/lib/pq"
	_ "github.com/lib/pq" //postgres driver
)

type (
	posting struct {
		domain.Posting
		PetID int
	}

	internalPostingAggregate struct {
		domain.Posting
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

const postingSelect = `SELECT
postings.id,
postings.name,
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
FROM postings 
LEFT JOIN pets ON pets.id = postings.pet_id
LEFT JOIN types ON types.id = pets.type_id
LEFT JOIN pet_breeds ON pets.id = pet_breeds.pet_id
LEFT JOIN tags ON tags.pet_id = pets.id `

const postingsGroupBy = `
GROUP BY 
postings.id,
postings.name,
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
text  `

var postingFieldMap = map[string]string{
	"id":           "postings.id",
	"name":         "postings.name",
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

func (db *DB) getPosting(filters domain.FilterMap) (*domain.Posting, error) {
	queryStr, args, err := db.buildQuery(postingFieldMap, filters)
	if err != nil {
		return nil, err
	}
	queryStr = queryStr + postingsGroupBy

	aggregate := &internalPostingAggregate{}

	err = db.Get(aggregate, postingSelect+queryStr, args...)
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
		TypeID:    aggregate.TypeID,
		Type:      aggregate.Type,
		Breeds:    aggregate.PetBreeds,
		Tag: domain.Tag{
			ID:    aggregate.TagID,
			Shape: aggregate.Shape,
			Color: aggregate.TagColor,
			Text:  aggregate.Text,
		},
	}

	return &aggregate.Posting, nil
}

func (db *DB) GetPostingByGUID(guid string) (*domain.Posting, error) {
	filter := domain.Filter{
		Comparator: "=",
		Value:      guid,
	}
	filters := domain.FilterMap{}
	filters["guid"] = []domain.Filter{filter}

	return db.getPosting(filters)

}

func (db *DB) GetMatchingSightings(id int) ([]domain.Sighting, error) {
	matchesQuery := matchesSelect + "WHERE postings_id = $1"
	matches := []matches{}

	err := db.Select((&matches), matchesQuery, id)
	if err != nil {
		return nil, err
	}
	sightingIDs := []int{}
	for _, m := range matches {
		sightingIDs = append(sightingIDs, m.SightingsID)
	}

	sFilter := domain.Filter{
		Comparator: "in",
		Value:      sightingIDs,
	}
	sFilters := domain.FilterMap{}
	sFilters["id"] = []domain.Filter{sFilter}

	sightings, err := db.GetAllSightings(sFilters)
	if err != nil {
		return nil, err
	}

	return sightings, nil
}

func (db *DB) GetPostingByID(id int) (*domain.Posting, error) {
	filter := domain.Filter{
		Comparator: "=",
		Value:      id,
	}
	filters := domain.FilterMap{}
	filters["id"] = []domain.Filter{filter}

	return db.getPosting(filters)
}

func (db *DB) GetPostingByEmail(email string) (*domain.Posting, error) {
	filter := domain.Filter{
		Comparator: "=",
		Value:      email,
	}
	filters := domain.FilterMap{}
	filters["email"] = []domain.Filter{filter}

	return db.getPosting(filters)
}

func (db *DB) GetAllPostings(filters ...domain.FilterMap) ([]domain.Posting, error) {
	queryStr, args, err := db.buildQuery(postingFieldMap, filters...)
	if err != nil {
		return nil, err
	}
	fmt.Println("querying postings")
	fmt.Println(queryStr)
	queryStr = queryStr + postingsGroupBy

	aggregates := []internalPostingAggregate{}

	err = db.Select((&aggregates), postingSelect+queryStr, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	postings := []domain.Posting{}

	for _, a := range aggregates {
		a.Pet = domain.Pet{
			ID:        a.PetID,
			PictureID: a.PictureID,
			Name:      a.PetName,
			Color:     a.PetColor,
			Marks:     a.Marks,
			Type:      a.Type,
			TypeID:    a.TypeID,
			Breeds:    a.PetBreeds,
			Tag: domain.Tag{
				ID:    a.TagID,
				Shape: a.Shape,
				Color: a.TagColor,
				Text:  a.Text,
			},
		}
		postings = append(postings, a.Posting)
	}

	return postings, nil
}

func (db *DB) AddPosting(newPosting *domain.Posting) error {

	guid, err := internal.NewUUID()
	if err != nil {
		return err
	}
	newPosting.GUID = guid

	err = db.addPet(&newPosting.Pet)
	if err != nil {
		return err
	}

	dbPosting := posting{
		Posting: *newPosting,
		PetID:   newPosting.Pet.ID,
	}

	query := `INSERT INTO postings(
		guid, pet_id, date, location, name, email)
		VALUES (:guid, :pet_id, :date, :location, :name, :email) RETURNING id;`

	rows, err := db.NamedQuery(query, dbPosting)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&newPosting.ID)
	} else {
		return errID
	}

	return nil
}
