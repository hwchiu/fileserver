package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/c9s/gomon/logger"
	fs "github.com/linkernetworks/fileserver/src"

	"github.com/gorilla/mux"
)

func newRouterServer(root string, basepath string) http.Handler {
	router := mux.NewRouter()
	if len(basepath) > 0 {
		router = router.PathPrefix(basepath).Subrouter()
	}
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
	var documentRoot string
	var basePath string
	var version bool = false

	flag.BoolVar(&version, "version", false, "version")
	flag.StringVar(&documentRoot, "documentRoot", "/workspace", "the document root of the file server")
	flag.StringVar(&basePath, "basePath", "", "the url base path of the APIs")
	flag.StringVar(&host, "host", "", "hostname")
	flag.StringVar(&port, "port", "33333", "port")
	flag.Parse()

	logger.Infof("Serving document root: %s at %s", documentRoot, basePath)
	router := newRouterServer(documentRoot, basePath)

	bind := net.JoinHostPort(host, port)
	logger.Infof("Listening at %s", bind)

	http.ListenAndServe(bind, logRequest(router))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("[%s] %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
