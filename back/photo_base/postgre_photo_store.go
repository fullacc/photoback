package photo_base

import (
	"github.com/go-pg/pg"
)

func NewPostgrePhotoStore(configfile *ConfigFile) (PhotoStore, error) {

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

func (ps *postgreStore) CreatePhoto(photo *Photo) (*Photo, error) {
	return photo, ps.db.Insert(photo)
}

func (ps *postgreStore) GetPhoto(id int)  (*Photo,error) {
	photo := &Photo{Id:id}
	err := ps.db.Select(photo)
	if err != nil {
		return nil,err
	}
	return photo,nil
}

func (ps *postgreStore) ListPhotos () ([]*Photo,error) {
	var photos []*Photo
	err := ps.db.Model(&photos).Select()
	if err != nil {
		return nil,err
	}
	return photos,nil
}

func (ps *postgreStore) ListPersonPhotos (id int) ([]*Photo,error) {
	var photos []*Photo
	err := ps.db.Model(&photos).Where("Person_Id = ?", id).Select()
	if err != nil {
		return nil,err
	}
	return photos,nil
}

func (ps *postgreStore) ListOperationPhotos(id int) ([]*Photo, error){
	var photos []*Photo
	err := ps.db.Model(&photos).Where("Operation_Id = ?", id).Select()
	if err != nil {
		return nil,err
	}
	return photos,nil
}

func (ps *postgreStore) UpdatePhoto(id int, photo *Photo) (*Photo, error) {
	photo1 := &Photo{Id:id}
	err := ps.db.Select(photo1)
	if err != nil {
		return nil,err
	}
	photo1 = photo
	err = ps.db.Update(photo1)
	if err != nil {
		return nil,err
	}
	return photo1, nil
}

func (ps *postgreStore) DeletePhoto(id int) error {
	photo := &Photo{Id: id}
	err := ps.db.Select(photo)
	if err != nil {
		return err
	}
	err = ps.db.Delete(photo)
	if err != nil {
		return err
	}
	return nil
}
