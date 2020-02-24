package photo_base

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"net/http"
)

type postgreStore struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Person)(nil),(*Operation)(nil),(*Photo)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func renderError(w http.ResponseWriter,msg string,statuscode int) {
	w.WriteHeader(statuscode)
	w.Write([]byte(msg))
}
