package photo_base

import "time"

type OperationStore interface {

	CreateOperation(operation *Operation, person *Person) (*Operation, error)

	DeleteOperation(id int64) error

	UpdateOperation(id int64, operation *Operation) (*Operation, error)

	ListOperations() ([]*Operation, error)

	ListPersonOperations(person *Person) ([]*Operation, error)

	GetOperation(id int64) (*Operation, error)
}


type Operation struct{
	ID int64
	Date time.Time
	Type string
	Person *Person
}