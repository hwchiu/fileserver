package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c9s/gomon/logger"
	fs "github.com/hwchiu/fileserver/src"

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
	router.HandleFunc("/download/{path:.*}", fs.GetDownloadFileHandler(root)).Methods("GET")
	return router
}

func main() {
	var host string
	var port string
	var documentRoot string
	var basePath string

	flag.StringVar(&documentRoot, "documentRoot", "/workspace", "the document root of the file server")
	flag.StringVar(&basePath, "basePath", "", "the url base path of the APIs")
	flag.StringVar(&host, "host", "0.0.0.0", "hostname")
	flag.StringVar(&port, "port", "33333", "port")
	flag.Parse()

	logger.Infof("Serving document root: %s at %s", documentRoot, basePath)
	router := newRouterServer(documentRoot, basePath)

	bind := net.JoinHostPort(host, port)
	logger.Infof("Listening at %s", bind)
	server := &http.Server{Addr: bind, Handler: logRequest(router)}

	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGTERM)
	go func() {
		<-sigC
		logger.Infof("caught signal SIGTERM, terminating fileserver...")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		server.Shutdown(ctx)
		os.Exit(0)
	}()

	server.ListenAndServe()
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("[%s] %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
