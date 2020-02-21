package photo_base

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)


type PersonEndpoints interface {

	GetPerson(idParam string) func(w http.ResponseWriter, r *http.Request)

	CreatePerson() func(w http.ResponseWriter, r *http.Request)

	ListPersons() func(w http.ResponseWriter, r *http.Request)

	UpdatePerson(idParam string) func (w http.ResponseWriter,r *http.Request)

	DeletePerson(idParam string) func(w http.ResponseWriter,r *http.Request)
}

func NewEndpointsPersonFactory(personStore PersonStore) PersonEndpoints {
	return &endpointspersonFactory{personStore: personStore}
}

type endpointspersonFactory struct {
	personStore PersonStore
}

func (ef *endpointspersonFactory) GetPerson(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: personId not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		person, err := ef.personStore.GetPerson(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(person)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointspersonFactory) CreatePerson() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}

		person := &Person{}
		if err := json.Unmarshal(data, person); err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusBadRequest)
			return
		}
		result, err := ef.personStore.CreatePerson(person)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}


func (ef *endpointspersonFactory) ListPersons() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		persons, err := ef.personStore.ListPersons()
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(persons)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointspersonFactory) UpdatePerson(idParam string) func (w http.ResponseWriter,r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: BookId not found",http.StatusBadRequest)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		person := &Person{}
		if err := json.Unmarshal(data, person); err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		result, err := ef.personStore.UpdatePerson(intid, person)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(result)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(response)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointspersonFactory) DeletePerson(idParam string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		vars := mux.Vars(r)
		id,ok := vars[idParam]
		if !ok {
			renderError(w,"Error: person not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		err = os.RemoveAll("./photos/"+idParam+"/")
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		err = ef.personStore.DeletePerson(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}