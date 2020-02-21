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
	Id int `json:"id,omitempty"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	DateOfBirth time.Time `json:"dateofbirth"`
	Phone string `json:"phone"`
	City string `json:"city"`
	Comment string `json:"comment,omitempty"`
}
