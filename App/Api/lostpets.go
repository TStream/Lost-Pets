package lostpets

import (
	"mime/multipart"
	"time"
)

type (
	Filter struct {
		Comparator string
		Value      interface{}
	}

	FilterMap map[string][]Filter

	Posting struct {
		ID       int
		GUID     string
		Pet      Pet
		Date     time.Time
		Location string
		Name     string
		Email    string
	}

	Sighting struct {
		InCustody bool
		Posting
	}

	Pet struct {
		ID        int
		PictureID int
		Name      string
		Color     string
		Marks     string
		TypeID    int
		Type      string
		Breeds    []string
		Tag       Tag
	}

	Tag struct {
		ID    int
		Shape string
		Color string
		Text  string
	}

	PetType struct {
		ID   int
		Name string
	}
)

type LostPetsRepo interface {
	GetPostingByGUID(guid string) (*Posting, error)
	GetSightingByGUID(guid string) (*Sighting, error)

	GetPostingByID(id int) (*Posting, error)
	GetSightingByID(id int) (*Sighting, error)

	GetPostingByEmail(email string) (*Posting, error)
	GetSightingByEmail(email string) (*Sighting, error)

	GetAllPostings(filters ...FilterMap) ([]Posting, error)
	GetMatchingPostings(sId int) ([]Posting, error)
	GetAllSightings(filters ...FilterMap) ([]Sighting, error)
	GetMatchingSightings(pId int) ([]Sighting, error)

	AddPosting(newPosting *Posting) error
	AddSighting(newSighting *Sighting) error

	AddMatch(pID int, sID int) error
	UpdateMatch(pID int, sID int, contactedOn time.Time) error
	RemoveMatch(pID int, sID int) error

	GetPetTypes() ([]PetType, error)
}

type FileMeta struct {
	ID          int
	GUID        string
	ContentType string
}

type FileRepo interface {
	GetFileMeta(id int) (*FileMeta, error)
	SaveFileMeta(meta *FileMeta) error
	RemoveFileMeta(id int) error
}

type FileStore interface {
	GetFile(guid string) (string, error)
	SaveFile(src multipart.File) (string, error)
	DeleteFile(guid string) error
}
