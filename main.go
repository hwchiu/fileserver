package main

import (
	fs "bitbucket.org/linkernetworks/aurora/src/fileserver"
	"flag"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

func newRouterServer() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/scan/{path:.*}", fs.ScanDirHandler).Methods("GET")
	router.HandleFunc("/read/{path:.*}", fs.ReadFileHandler).Methods("GET")
	router.HandleFunc("/write/{path:.*}", fs.WriteFileHandler).Methods("POST")
	router.HandleFunc("/delete/{path:.*}", fs.RemoveFileHandler).Methods("DELETE")
	return router
}

func main() {
	var host string
	var port string

	flag.StringVar(&host, "h", "", "hostname")
	flag.StringVar(&port, "p", "33333", "port")
	flag.Parse()

	bind := net.JoinHostPort(host, port)

	router := newRouterServer()
	http.ListenAndServe(bind, router)
}
