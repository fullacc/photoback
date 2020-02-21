package photo_base

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)


type PhotoEndpoints interface {

	GetPhoto(idParam string) func(w http.ResponseWriter, r *http.Request)

	CreatePhoto(personid string,operationid string) func(w http.ResponseWriter, r *http.Request)

	ListPhotos() func(w http.ResponseWriter, r *http.Request)

	ListPersonPhotos(idParam string) func(w http.ResponseWriter, r *http.Request)

	ListOperationPhotos(idParam string) func(w http.ResponseWriter, r *http.Request)

	UpdatePhoto(idParam string) func (w http.ResponseWriter,r *http.Request)

	DeletePhoto(idParam string) func(w http.ResponseWriter,r *http.Request)
}

func NewEndpointsPhotoFactory(photoStore PhotoStore) PhotoEndpoints {
	return &endpointsphotoFactory{photoStore: photoStore}
}

type endpointsphotoFactory struct {
	photoStore PhotoStore
}

func (ef *endpointsphotoFactory) GetPhoto(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: Id not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		photo, err := ef.photoStore.GetPhoto(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(photo)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusCreated)
	}
}

func (ef *endpointsphotoFactory) CreatePhoto(personid string, operationid string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}

		const maxUploadSize = 10 * 1024

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err = r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w,"Error: FILE_TOO_BIG",http.StatusBadRequest)
			return
		}

		fileType := r.PostFormValue("type")
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "Error: INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "Error: INVALID_FILE", http.StatusBadRequest)
			return
		}
		filetype := http.DetectContentType(fileBytes)
		if filetype != "image/jpeg" && filetype != "image/jpg" &&
			filetype != "image/gif" && filetype != "image/png" &&
			filetype != "application/pdf" {
			renderError(w, "Error: INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileEndings, err := mime.ExtensionsByType(fileType)
		if err != nil {
			renderError(w, "Error: CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}

		fileName := xid.New().String()
		uploadPath := "./photos/"+personid+"/"+operationid+"/"

		err = os.MkdirAll(uploadPath,0777)
		if err != nil {
			renderError(w,"Error: CANT_CREATE_DIRECTORY",http.StatusInternalServerError)
			return
		}

		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "Error: CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			renderError(w, "Error: CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		photo := &Photo{}

		if err := json.Unmarshal(data, photo); err != nil {
			renderError(w,"Error: " + err.Error(),http.StatusBadRequest)
			return
		}
		photo.Uid = fileName
		photo.FilePath = newPath
		result, err := ef.photoStore.CreatePhoto(photo)
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


func (ef *endpointsphotoFactory) ListPhotos() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		photos, err := ef.photoStore.ListPhotos()
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(photos)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsphotoFactory) ListPersonPhotos(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: Id not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		photos, err := ef.photoStore.ListPersonPhotos(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(photos)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsphotoFactory) ListOperationPhotos(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter,r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: Id not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		photos, err := ef.photoStore.ListOperationPhotos(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(photos)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

func (ef *endpointsphotoFactory) UpdatePhoto(idParam string) func (w http.ResponseWriter,r *http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars[idParam]
		if !ok {
			renderError(w,"Error: PhotoId not found",http.StatusBadRequest)
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		photo := &Photo{}
		if err := json.Unmarshal(data, photo); err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		result, err := ef.photoStore.UpdatePhoto(intid, photo)
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

func (ef *endpointsphotoFactory) DeletePhoto(idParam string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		vars := mux.Vars(r)
		id,ok := vars[idParam]
		if !ok {
			renderError(w,"Error: Photo not found",http.StatusBadRequest)
			return
		}
		intid, err := strconv.Atoi(id)
		if err!=nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		photo,err := ef.photoStore.GetPhoto(intid)
		os.Remove(photo.FilePath)
		err = ef.photoStore.DeletePhoto(intid)
		if err != nil {
			renderError(w,"Error: "+err.Error(),http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}