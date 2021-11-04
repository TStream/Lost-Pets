package postgres

import (
	"database/sql"
	domain "lostpets"
)

const fileSelect = `SELECT
id,
guid,
COALESCE (content_type, '') as content_type
FROM
pictures `

const addFileSQL = `INSERT INTO pictures
(id,guid,content_type) VALUES
(:id,:guid,:content_type) RETURNING id;
`

func (db *DB) GetFileMeta(id int) (*domain.FileMeta, error) {
	query := fileSelect + "WHERE id = $1 "

	file := &domain.FileMeta{}
	err := db.Get(file, query, id)
	if err == sql.ErrNoRows { //no rows in result isn't an error, just no files found
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return file, nil

}

func (db *DB) SaveFileMeta(meta *domain.FileMeta) error {
	rows, err := db.NamedQuery(addFileSQL, meta)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&meta.ID)
	} else {
		return errID
	}

	return nil
}

func (db *DB) RemoveFileMeta(id int) error {
	query := "DELETE from pictures where id=$1"

	_, err := db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}
