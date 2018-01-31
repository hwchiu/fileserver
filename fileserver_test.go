package main

import (
	"bitbucket.org/linkernetworks/aurora/src/utils/fileutils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	_ "log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createTempDir(t *testing.T, prefix string) string {
	dir, err := ioutil.TempDir(".", prefix)
	assert.NoError(t, err)

	return dir
}

func createTempFile(t *testing.T, dir, name string, contents []byte) {
	f, err := os.Create(dir + "/" + name)
	assert.NoError(t, err)

	f.Write(contents)
	f.Close()
}

func TestLoadDir(t *testing.T) {
	dirPrefix := "testdir"
	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	createTempFile(t, tmpDir, "test", []byte{})

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	//Get the abosolute path for testing dir
	dir := pwd + "/" + tmpDir

	req, err := http.NewRequest("GET", "/scan"+dir, nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 200)

	//Test Files
	var fi []fileutils.FileInfo
	err = json.Unmarshal(res.Body.Bytes(), &fi)
	assert.NoError(t, err)

	assert.Equal(t, fi[0].Name, "test")
	assert.Equal(t, fi[0].Size, int64(0))
	assert.Equal(t, fi[0].Type, "")
	assert.Equal(t, fi[0].IsDir, false)

	os.RemoveAll(tmpDir)
}

func TestInvalidLoadDir(t *testing.T) {
	req, err := http.NewRequest("GET", "/scan/987654321/12345667890/1234", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 404)
}
