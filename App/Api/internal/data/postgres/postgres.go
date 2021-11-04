package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	domain "lostpets"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" //postgres driver
)

type (
	Config struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"dbName"`
		SSLMode  string `json:"sslmode"`
	}

	DB struct {
		*sqlx.DB
	}

	Repo struct {
		db *DB
	}

	breed struct {
		ID    int
		PetID int
		Name  string
	}

	tag struct {
		domain.Tag
		PetID int
	}

	matches struct {
		PostingsID    int
		SightingsID   int
		LastContacted *time.Time
	}

	query struct {
		QueryStrs []string
		Args      []interface{}
	}
)

var errorSentinel = errors.New("invalid in filter")

var errID = errors.New("postgresDb: ID was not returned after insert")

const connTemplate = "user=%s password=%s host=%s dbname=%s port=%d sslmode=%s TimeZone=UTC"

func NewDBConnection(config Config) (*DB, error) {
	connStr := fmt.Sprintf(connTemplate, config.Username, config.Password, config.Host, config.Name, config.Port, config.SSLMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.MapperFunc(underscore)
	return &DB{DB: db}, nil
}

func (db *DB) NewLostPetsRepo() *Repo {
	return &Repo{db: db}
}

const matchesSelect = `SELECT
postings_id,
sightings_id,
last_contacted
FROM matches `

func (db *DB) addPet(pet *domain.Pet) error {

	query := `INSERT INTO pets(
		picture_id, type_id, name, color, marks)
		VALUES (:picture_id, :type_id, :name, :color, :marks) RETURNING id;`

	rows, err := db.NamedQuery(query, pet)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&pet.ID)
	} else {
		return errID
	}

	// add attributes
	err = db.addBreeds(pet)
	if err != nil {
		return err
	}

	// add certifications
	err = db.addTag(pet)
	if err != nil {
		return err
	}

	return nil

}

func (db *DB) addBreeds(pet *domain.Pet) error {

	// add breeds
	breeds := []breed{}
	for _, a := range pet.Breeds {
		breeds = append(breeds, breed{Name: a, PetID: int(pet.ID)})
	}
	if len(breeds) != 0 {
		query := `INSERT INTO pet_breeds(pet_id, name) VALUES (:pet_id, :name)`
		_, err := db.NamedExec(query, breeds)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			default:
				return err
			}
		}
	}

	return nil
}

func (db *DB) addTag(pet *domain.Pet) error {

	tag := tag{
		Tag:   pet.Tag,
		PetID: pet.ID,
	}

	query := `INSERT INTO tags(
		pet_id, shape, text, color)
		VALUES (:pet_id, :shape, :text, :color) RETURNING id;`

	rows, err := db.NamedQuery(query, tag)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&pet.Tag.ID)
	} else {
		return errID
	}

	return nil
}

func (db *DB) AddMatch(pID int, sID int) error {
	query := `INSERT INTO matches(
		postings_id, sightings_id)
		VALUES ($1,$2)`

	rows, err := db.Query(query, pID, sID)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}
func (db *DB) UpdateMatch(pID int, sID int, contactedOn time.Time) error {
	query := `UPDATE matches SET
	last_contacted=$1
	WHERE postings_id = $2 AND sightings_id = $3`

	rows, err := db.Query(query, contactedOn, pID, sID)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}
func (db *DB) RemoveMatch(pID int, sID int) error {
	query := `DELETE FROM matches WHERE postings_id = $1 AND sightings_id = $2`

	rows, err := db.Query(query, pID, sID)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (db *DB) GetPetTypes() ([]domain.PetType, error) {
	query := "SELECT id, name FROM types ORDER BY id ASC"
	types := []domain.PetType{}

	err := db.Select((&types), query)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return types, nil
}

func (db *DB) buildQuery(fieldMap map[string]string, filters ...domain.FilterMap) (string, []interface{}, error) {
	querys := query{
		QueryStrs: []string{}, // array of AND querys
		Args:      []interface{}{},
	}

	for _, f := range filters {
		queryStr, params, err := getFilters(fieldMap, f)

		if err != nil {
			return "", nil, err
		}

		queryStr, args, err := sqlx.Named(queryStr, params)
		if err != nil {
			return "", nil, err
		}
		queryStr, args, err = sqlx.In(queryStr, args...)

		if err != nil {
			return "", nil, err
		}

		queryStrF := db.Rebind(queryStr)

		querys.QueryStrs = append(querys.QueryStrs, "("+queryStrF+")")
		querys.Args = append(querys.Args, args...)

	}
	joinQueryStr := ""
	if len(querys.QueryStrs) > 0 {
		//example of joined = "(breed = $1 and type = $2)OR(breed = $1 and color = $2)OR(type = $1)"
		joinStr := strings.Join(querys.QueryStrs, " OR ")
		//now need to correct argumnent placings
		m := regexp.MustCompile(`\$[0-9]+`) //looking for all '$<numbner>'
		match := m.FindStringIndex(joinStr) //match is int array [0] = start, [1] = end
		argNum := 1
		for len(match) > 0 {
			//temp use '##' as to not get matched again
			joinStr = joinStr[:match[0]] + "##" + strconv.Itoa(argNum) + joinStr[match[1]:]
			argNum = argNum + 1
			match = m.FindStringIndex(joinStr)
		}
		//replace '##' with '$'
		joinQueryStr = strings.Replace(joinStr, "##", "$", -1)
		joinQueryStr = " WHERE " + joinQueryStr
	}
	return joinQueryStr, querys.Args, nil
}
