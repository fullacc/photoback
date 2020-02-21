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
	Id int
	Date time.Time
	OperationId int
	Status int
	PersonId int
	FilePath string
	Uid string
}