package photo_base

import (
	"github.com/go-pg/pg"
)

func NewPostgreOperationStore(configfile *ConfigFile) (OperationStore, error) {
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

func (ps *postgreStore) CreateOperation(operation *Operation) (*Operation, error) {
	return operation, ps.db.Insert(operation)
}

func (ps *postgreStore) GetOperation(id int)  (*Operation,error) {
	operation := &Operation{Id:id}
	err := ps.db.Select(operation)
	if err != nil {
		return nil,err
	}
	return operation,nil
}

func (ps *postgreStore) ListOperations () ([]*Operation,error) {
	var operations []*Operation
	err := ps.db.Model(&operations).Select()
	if err != nil {
		return nil,err
	}
	return operations,nil
}

func (ps *postgreStore) ListPersonOperations (id int) ([]*Operation,error) {
	var operations []*Operation
	err := ps.db.Model(&operations).Where("Person_Id = ?", id).Select()
	if err != nil {
		return nil,err
	}
	return operations,nil
}

func (ps *postgreStore) UpdateOperation(id int, operation *Operation) (*Operation, error) {
	operation1 := &Operation{Id:id}
	err := ps.db.Select(operation1)
	if err != nil {
		return nil,err
	}
	operation1 = operation
	err = ps.db.Update(operation1)
	if err != nil {
		return nil,err
	}
	return operation1, nil
}

func (ps *postgreStore) DeleteOperation(id int) error {
	var photos []*Photo
	err := ps.db.Model(&photos).Where("Operation_Id = ?", id).Select()
	if err != nil {
		return err
	}
	_, err = ps.db.Model(&photos).Delete()
	if err != nil {
		return err
	}
	operation := &Operation{Id: id}
	err = ps.db.Select(operation)
	if err != nil {
		return err
	}
	err = ps.db.Delete(operation)
	if err != nil {
		return err
	}

	return nil
}
