package photo_base

import (
	"github.com/go-pg/pg"
)

func NewPostgrePersonStore(configfile *ConfigFile) (PersonStore, error) {
	db := pg.Connect(&pg.Options{
		Database: configfile.Name,
		Addr: configfile.DbHost + ":" + configfile.DbPort,
		User: "postgres",
		Password: configfile.Password,
	})

	err := createSchema(db)
	if err != nil {
		return nil, err
	}
	return &postgreStore{db: db}, nil
}

func (ps *postgreStore) CreatePerson(person *Person) (*Person, error) {
	return person, ps.db.Insert(person)
}

func (ps *postgreStore) GetPerson(id int)  (*Person,error) {
	person := &Person{Id:id}
	err := ps.db.Select(person)
	if err != nil {
		return nil,err
	}
	return person,nil
}

func (ps *postgreStore) ListPersons () ([]*Person,error) {
	var persons []*Person
	err := ps.db.Model(&persons).Select()
	if err != nil {
		return nil,err
	}
	return persons,nil
}

func (ps *postgreStore) UpdatePerson(id int, person *Person) (*Person, error) {
	person1 := &Person{Id:id}
	err := ps.db.Select(person1)
	if err != nil {
		return nil,err
	}
	person1 = person
	err = ps.db.Update(person1)
	if err != nil {
		return nil,err
	}
	return person1, nil
}

func (ps *postgreStore) DeletePerson(id int) error {
	var operations []*Operation
	err := ps.db.Model(&operations).Where("PersonId = ?",id).Select()
	var photos []*Photo
	for _,v := range operations {
		photos = nil
		err = ps.db.Model(&photos).Where("OperationId = ?", v.Id).Select()
		if err != nil {
			return err
		}
		_, err = ps.db.Model(&photos).Delete()
		if err != nil {
			return err
		}
		err = ps.db.Delete(v)
		if err != nil {
			return err
		}
	}
	person := &Person{Id: id}
	err = ps.db.Select(person)
	if err != nil {
		return err
	}
	err = ps.db.Delete(person)
	if err != nil {
		return err
	}
	return nil
}
