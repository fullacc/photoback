package photo_base

import "time"

type PersonStore interface{

	CreatePerson(person *Person) (*Person, error)

	DeletePerson(id int64) error

	UpdatePerson(id int64, person *Person) (*Person, error)

	ListPersons() ([]*Person, error)

	GetPerson(id int64) (*Person, error)
}

type Person struct{
	ID int64
	Name string
	Surname string
	DateOfBirth time.Time
	Phone string
	City string
	Comment string
}
