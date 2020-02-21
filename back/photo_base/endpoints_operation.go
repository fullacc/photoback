package photo_base

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)


type OperationEndpoints interface {

	GetOperation(idParam string) func(w http.ResponseWriter, r *http.Request)

	CreateOperation() func(w http.ResponseWriter, r *http.Request)

	ListOperations() func(w http.ResponseWriter, r *http.Request)

	ListPersonOperations(idParam string) func(w http.ResponseWriter, r *http.Request)

	UpdateOperation(idParam string) func (w http.ResponseWriter,r *http.Request)

	DeleteOperation(personid string, operationid string) func(w http.ResponseWriter,r *http.Request)
}

func NewEndpointsOperationFactory(operationStore OperationStore) OperationEndpoints {
	return &endpointsFactory{operationStore: operationStore}
}

type endpointsFactory struct {
	operationStore OperationStore
}

func (ef *endpointsFactory) GetOperation(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: operationd id not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		operation, err := ef.operationStore.GetOperation(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(operation)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointsFactory) CreateOperation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}

		operation := &Operation{}
		if err := json.Unmarshal(data, operation); err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusBadRequest)
			return
		}
		result, err := ef.operationStore.CreateOperation(operation)
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

func (ef *endpointsFactory) ListOperations() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		operations, err := ef.operationStore.ListOperations()
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(operations)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsFactory) ListPersonOperations(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: personid not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		operations, err := ef.operationStore.ListPersonOperations(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(operations)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsFactory) UpdateOperation(idParam string) func (w http.ResponseWriter,r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: operation id not found",http.StatusBadRequest)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		operation := &Operation{}
		if err := json.Unmarshal(data, operation); err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		result, err := ef.operationStore.UpdateOperation(intid, operation)
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

func (ef *endpointsFactory) DeleteOperation(personid string, operationid string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		vars := mux.Vars(r)
		id,ok := vars[operationid]
		if !ok {
			renderError(w,"Error: operation not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		err = ef.operationStore.DeleteOperation(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		err = os.RemoveAll("./photos/"+personid+"/"+operationid+"/")
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}