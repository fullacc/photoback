package photo_base

import "time"

type PersonStore interface{

	CreatePerson(person *Person) (*Person, error)

	DeletePerson(id int) error

	UpdatePerson(id int, person *Person) (*Person, error)

	ListPersons() ([]*Person, error)

	GetPerson(id int) (*Person, error)
}

type Person struct{
	Id int
	Name string
	Surname string
	DateOfBirth time.Time
	Phone string
	City string
	Comment string
}
