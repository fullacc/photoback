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
	Id int
	Date time.Time
	Type string
	PersonId int
}