package fileserver

import (
	"bitbucket.org/linkernetworks/aurora/src/utils/fileutils"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

const root = "/"

type FileContent struct {
	Name    string `json:"name"`
	Ext     string `json:"ext"`
	Type    string `json:"type"`
	Content []byte `json:"content"`
}

func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func RemoveFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Ready to delete the file")
	values := mux.Vars(r)
	p := root + values["path"]

	log.Println("target path is ", p)
	if err := os.RemoveAll(p); err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Delete file success")
}
func WriteFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Ready to upload the file")
	values := mux.Vars(r)
	p := "/" + values["path"]

	log.Println("target path is ", p)
	var fc FileContent
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&fc); err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	filePath := p + "/" + fc.Name
	if err := ioutil.WriteFile(filePath, fc.Content, 0644); err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Upload file success")
}

func ReadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Ready to read the file")
	values := mux.Vars(r)
	p := root + values["path"]

	log.Println("target path is ", p)
	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusNotFound)
		return
	}

	response, err := json.Marshal(FileContent{
		Name:    path.Base(p),
		Ext:     path.Ext(p),
		Type:    mime.TypeByExtension(path.Ext(p)),
		Content: bytes,
	})
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Read file success")
}

func ScanDirHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Ready to load the dir")
	values := mux.Vars(r)
	p := root + values["path"]

	log.Println("target path is ", p)
	infos, err := fileutils.ScanDir(p)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusNotFound)
		return
	}

	response, err := json.Marshal(infos)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Load dir success")
}
