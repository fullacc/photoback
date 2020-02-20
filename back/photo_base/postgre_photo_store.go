package book_store

import (
	"bufio"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/segmentio/encoding/json"
	"io/ioutil"
	"os"
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

func NewPostgreBookStore(filename string) (BookStore, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		return nil, err
	}
	var configfile ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		return nil, err
	}
	file.Close()

	db := pg.Connect(&pg.Options{
		Database: configfile.Name,
		Addr: configfile.DbHost + ":" + configfile.DbPort,
		User: "postgres",
		Password: configfile.Password,
	})

	//err = createSchema(db)
	if err != nil {
		return nil, err
	}
	return &postgreStore{db: db}, nil
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Book)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *postgreStore) Create(book *Book) (*Book, error) {
	return book, ps.db.Insert(book)
}

func (ps *postgreStore) GetBook(id int64)  (*Book,error) {
	book := &Book{ID:id}
	err := ps.db.Select(book)
	if err != nil {
		return nil,err
	}
	return book,nil
}

func (ps *postgreStore) ListBooks () ([]*Book,error) {
	var books []*Book
	err := ps.db.Model(&books).Select()
	if err != nil {
		return nil,err
	}
	return books,nil
}

func (ps *postgreStore) UpdateBook(id int64, book *Book) (*Book, error) {
	book1 := &Book{ID:id}
	err := ps.db.Select(book1)
	if err != nil {
		return nil,err
	}
	book1 = book
	err = ps.db.Update(book1)
	if err != nil {
		return nil,err
	}
	return book1,nil
}

func (ps *postgreStore) DeleteBook(id int64) error {
	book := &Book{ID:id}
	err := ps.db.Select(book)
	if err != nil {
		return err
	}
	err = ps.db.Delete(book)
	if err != nil {
		return err
	}
	return nil
}

func (ps *postgreStore) SaveBooks(filepath string) error {
	return nil
}