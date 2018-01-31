package main

import (
	"bitbucket.org/linkernetworks/aurora/src/utils/fileutils"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
)

type FileContent struct {
	Name    string `json:"name"`
	Ext     string `json:"ext"`
	Type    string `json:"type"`
	Content []byte `json:"content"`
}

func ReadFileHandler(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)

	p := "/" + values["path"]

	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
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

	infos, err := fileutils.ScanDir("/" + values["path"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(infos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/scan/{path:.*}", LoadDirHandler).Methods("GET")
	router.HandleFunc("/read/{path:.*}", ReadFileHandler).Methods("GET")

	http.ListenAndServe(":33333", router)
}
