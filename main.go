package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type FileOperation struct {
	Path     string `json:"path"`
	FileName string `json:"fileName"`
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)
	log.Println(path)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/list/{path:.*}", ListFiles).Methods("GET")

	http.ListenAndServe(":33333", router)
}
