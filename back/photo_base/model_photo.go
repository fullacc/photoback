package photo_base

import "time"

type PhotoStore interface {

	CreatePhoto(photo *Photo) (*Photo, error)

	GetPhoto(id int) (*Photo, error)

	ListPhotos() ([]*Photo, error)

	ListOperationPhotos(id int) ([]*Photo, error)

	ListPersonPhotos(id int) ([]*Photo, error)

	UpdatePhoto(id int, photo *Photo) (*Photo, error)

	DeletePhoto(id int)  error
}

type Photo struct {
	Id int `json:"id,omitempty"`
	Date time.Time `json:"date"`
	OperationId int `json:"operationid"`
	Status int `json:"status"`
	PersonId int `json:"personid"`
	FilePath string
	Uid string
}