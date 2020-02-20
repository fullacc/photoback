package photo_base

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)


type PhotoEndpoints interface {

	GetPhoto(idParam string) func(w http.ResponseWriter, r *http.Request)

	CreatePhoto() func(w http.ResponseWriter, r *http.Request)

	ListPhoto() func(w http.ResponseWriter, r *http.Request)

	UpdatePhoto(idParam string) func (w http.ResponseWriter,r *http.Request)

	DeletePhoto(idParam string) func(w http.ResponseWriter,r *http.Request)
}

func NewEndpointsFactory(photoStore PhotoStore) PhotoEndpoints {
	return &endpointsFactory{photoStore: photoStore}
}

type endpointsFactory struct {
	photoStore PhotoStore
}

func internal (w http.ResponseWriter,err error){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Error: " + err.Error()))
}

func (ef *endpointsFactory) GetPhoto(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Id not found"))
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			internal(w,err)
			return
		}
		photo, err := ef.photoStore.GetPhoto(int64(intid))
		if err != nil {
			internal(w,err)
			return
		}
		data, err := json.Marshal(photo)
		if err != nil {
			internal(w,err)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointsFactory) CreatePhoto() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			internal(w,err)
			return
		}

		photo := &Photo{}
		if err := json.Unmarshal(data, photo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		result, err := ef.photoStore.CreatePhoto(photo)
		if err != nil {
			internal(w,err)
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			internal(w,err)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}


func (ef *endpointsFactory) ListPhoto() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		photos, err := ef.photoStore.ListPhotos()
		if err != nil {
			internal(w,err)
			return
		}
		data, err := json.Marshal(photos)
		if err != nil {
			internal(w,err)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsFactory) UpdatePhoto(idParam string) func (w http.ResponseWriter,r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Book ID not found "))
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			internal(w,err)
			return
		}
		photo := &Photo{}
		if err := json.Unmarshal(data, photo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			internal(w,err)
			return
		}
		result, err := ef.photoStore.UpdatePhoto(int64(intid), photo)
		if err != nil {
			internal(w,err)
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			internal(w,err)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointsFactory) DeletePhoto(idParam string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		vars := mux.Vars(r)
		id,ok := vars[idParam]
		if !ok {
			w.Write([]byte("Error: Not Found"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			internal(w,err)
			return
		}
		err = ef.photoStore.DeletePhoto(int64(intid))
		if err != nil {
			internal(w,err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}