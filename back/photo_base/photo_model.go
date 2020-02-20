package photo_base

import "time"

type PhotoStore interface {

	CreatePhoto(photo *Photo, operationId int64) (*Photo, error)

	GetPhoto(id int64) (*Photo, error)

	ListPhotos() ([]*Photo, error)

	ListOperationPhotos(operation *Operation) ([]*Photo, error)

	ListPersonPhotos(person *Person) ([]*Photo, error)

	UpdatePhoto(id int64, photo *Photo) (*Photo, error)

	DeletePhoto(id int64)  error
}

type Photo struct {
	ID int64
	Date time.Time
	Operation *Operation
	Status int64
	Person *Person
	FilePath string
}