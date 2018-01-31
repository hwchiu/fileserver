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

const invalidDir = "/123456789/9876554321/123456789"

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
	dirPrefix := "loadDir"
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
	req, err := http.NewRequest("GET", "/scan"+invalidDir, nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 404)
}

func TestReadFile(t *testing.T) {
	dirPrefix := "readDir"
	testFileExt := ".txt"
	testFileName := "readMe"
	testFileContents := []byte{12, 3, 4, 1, 213, 213, 13}
	testFile := testFileName + testFileExt

	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	createTempFile(t, tmpDir, testFile, testFileContents)

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	//Get the abosolute path for testing dir
	filePath := pwd + "/" + tmpDir + "/" + testFile

	req, err := http.NewRequest("GET", "/read"+filePath, nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 200)

	//Test Files
	var fc FileContent
	err = json.Unmarshal(res.Body.Bytes(), &fc)
	assert.NoError(t, err)

	assert.Equal(t, fc.Name, testFile)
	assert.Equal(t, fc.Ext, testFileExt)
	assert.Equal(t, fc.Type, "text/plain; charset=utf-8")
	assert.Equal(t, fc.Content, testFileContents)

	os.RemoveAll(tmpDir)
}

func TestInvalidReadfile(t *testing.T) {
	req, err := http.NewRequest("GET", "/read"+invalidDir+"/a", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 404)
}
