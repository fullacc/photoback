package photo_base

import "time"

type OperationStore interface {

	CreateOperation(operation *Operation) (*Operation, error)

	DeleteOperation(id int) error

	UpdateOperation(id int, operation *Operation) (*Operation, error)

	ListOperations() ([]*Operation, error)

	ListPersonOperations(id int) ([]*Operation, error)

	GetOperation(id int) (*Operation, error)
}


type Operation struct{
	Id int `json:"id,omitempty"`
	Date time.Time `json:"date"`
	Type string `json:"type"`
	PersonId int `json:"personid"`
}