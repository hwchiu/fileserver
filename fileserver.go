package fileserver

import (
	"bitbucket.org/linkernetworks/aurora/src/logger"
	"bitbucket.org/linkernetworks/aurora/src/query"
	"bitbucket.org/linkernetworks/aurora/src/utils/fileutils"

	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
)

const root = "/workspace"

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

func GetRemoveFileHandler(root string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		RemoveFileHandler(root, w, r)
	}
}

func GetWriteFileHandler(root string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteFileHandler(root, w, r)
	}
}

func GetReadFileHandler(root string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ReadFileHandler(root, w, r)
	}
}

func GetScanDirHandler(root string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ScanDirHandler(root, w, r)
	}
}
func RemoveFileHandler(root string, w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := path.Join(root, values["path"])

	if err := os.RemoveAll(p); err != nil {
		logger.Errorf("remove error: %v", err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func WriteFileHandler(root string, w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)

	p := path.Join(root, values["path"])

	var fc FileContent
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&fc); err != nil {
		logger.Errorf("json error: %v", err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	filePath := p + "/" + fc.Name
	if err := ioutil.WriteFile(filePath, fc.Content, 0644); err != nil {
		logger.Errorf("write error: %v", err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ReadFileHandler(root string, w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := path.Join(root, values["path"])

	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		logger.Errorf("read error: %v", err)
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
}

func ScanDirHandler(root string, w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	p := root
	if subPath, ok := values["path"]; ok {
		p = path.Join(root, subPath)
	}

	//default behavior is ignore the hidden files
	excludePrefix := "."
	query := query.New(r.URL.Query())
	if value, ok := query.Str("hidden"); ok {
		if value == "1" {
			excludePrefix = ""
		}
	}

	infos, err := fileutils.ScanDir(p, excludePrefix)
	if err != nil {
		logger.Errorf("scan dir error: %v", err)
		writeError(w, err, http.StatusNotFound)
		return
	}

	response, err := json.Marshal(infos)
	if err != nil {
		logger.Errorf("json error: %v", err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
