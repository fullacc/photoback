package photo_base

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type PostgreConfig struct {
	User string
	Password string
	Port string
	Host string
}

type postgreStore struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Person)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	for _, model := range []interface{}{(*Operation)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
			FKConstraints: true,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	for _, model := range []interface{}{(*Photo)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
			FKConstraints: true,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

