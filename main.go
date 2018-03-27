package main

import (
	"bitbucket.org/linkernetworks/aurora/src/aurora"
	fs "bitbucket.org/linkernetworks/aurora/src/fileserver"
	"flag"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

func newRouterServer() http.Handler {
	router := mux.NewRouter()
	root := "/workspace"
	router.HandleFunc("/scan/{path:.*}", fs.GetScanDirHandler(root)).Methods("GET")
	router.HandleFunc("/scan", fs.GetScanDirHandler(root)).Methods("GET")
	router.HandleFunc("/read/{path:.*}", fs.GetReadFileHandler(root)).Methods("GET")
	router.HandleFunc("/write/{path:.*}", fs.GetWriteFileHandler(root)).Methods("POST")
	router.HandleFunc("/delete/{path:.*}", fs.GetRemoveFileHandler(root)).Methods("DELETE")
	return router
}

func main() {
	var host string
	var port string
	var version bool = false

	flag.BoolVar(&version, "version", false, "version")
	flag.StringVar(&host, "h", "", "hostname")
	flag.StringVar(&port, "p", "33333", "port")
	flag.Parse()

	if version {
		aurora.PrintVersion()
		return
	}

	bind := net.JoinHostPort(host, port)

	router := newRouterServer()
	http.ListenAndServe(bind, router)
}
