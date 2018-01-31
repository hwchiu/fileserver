package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"time"
)

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Type    string    `json:"type"`
	ModTime time.Time `json:"mtime"`
	IsDir   bool      `json:"isDir"`
}

type FileContent struct {
	Name    string `json:"name"`
	Ext     string `json:"ext"`
	Type    string `json:"type"`
	Content []byte `json:"content"`
}

func ScanDir(p string) ([]FileInfo, error) {
	fileInfos := []FileInfo{}
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return fileInfos, err
	}

	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Name:    file.Name(),
			Size:    file.Size(),
			ModTime: file.ModTime(),
			IsDir:   file.IsDir(),
			Type:    mime.TypeByExtension(path.Ext(file.Name())),
		})
	}

	return fileInfos, nil
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

	infos, err := ScanDir("/" + values["path"])
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
