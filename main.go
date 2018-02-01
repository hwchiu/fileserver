package main

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

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := root + values["path"]

	if err := os.RemoveAll(p); err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func WriteFileHandler(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := "/" + values["path"]

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
}

func ReadFileHandler(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := root + values["path"]

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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(response)
}

func LoadDirHandler(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := root + values["path"]

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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func newRouterServer() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/scan/{path:.*}", LoadDirHandler).Methods("GET")
	router.HandleFunc("/read/{path:.*}", ReadFileHandler).Methods("GET")
	router.HandleFunc("/upload/{path:.*}", WriteFileHandler).Methods("POST")
	router.HandleFunc("/delete/{path:.*}", DeleteFileHandler).Methods("DELETE")
	return router
}

func main() {
	router := newRouterServer()
	http.ListenAndServe(":33333", router)
}
