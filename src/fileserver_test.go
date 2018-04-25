package fileserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const invalidDir = "/invalidpath/ignore/me"

func newRouterServer() http.Handler {
	router := mux.NewRouter()

	root := "/"
	router.HandleFunc("/scan/{path:.*}", GetScanDirHandler(root)).Methods("GET")
	router.HandleFunc("/scan", GetScanDirHandler(root)).Methods("GET")
	router.HandleFunc("/read/{path:.*}", GetReadFileHandler(root)).Methods("GET")
	router.HandleFunc("/write/{path:.*}", GetWriteFileHandler(root)).Methods("POST")
	router.HandleFunc("/delete/{path:.*}", GetRemoveFileHandler(root)).Methods("DELETE")

	return router
}

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

func TestScanDir(t *testing.T) {
	dirPrefix := "loadDir"
	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	defer os.RemoveAll(tmpDir)
	targetfile := []string{".hidden_one", ".hidden_two", "i_am_not_hidden"}
	for _, v := range targetfile {
		createTempFile(t, tmpDir, v, []byte{})
	}

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
	var fi []FileInfo
	err = json.Unmarshal(res.Body.Bytes(), &fi)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(fi))
	assert.Equal(t, fi[0].Name, "i_am_not_hidden")
	assert.Equal(t, fi[0].Size, int64(0))
	assert.Equal(t, fi[0].Type, "")
	assert.Equal(t, fi[0].IsDir, false)
}

func TestScanDirWithHidden(t *testing.T) {
	dirPrefix := "loadDir"
	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	defer os.RemoveAll(tmpDir)
	targetfile := []string{".hidden_one", ".hidden_two", "i_am_not_hidden"}
	for _, v := range targetfile {
		createTempFile(t, tmpDir, v, []byte{})
	}

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	//Get the abosolute path for testing dir
	dir := pwd + "/" + tmpDir

	req, err := http.NewRequest("GET", "/scan"+dir+"/?hidden=1", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 200)

	//Test Files
	var fi []FileInfo
	err = json.Unmarshal(res.Body.Bytes(), &fi)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(fi))

	for i, _ := range targetfile {
		assert.Equal(t, fi[i].Name, targetfile[i])
		assert.Equal(t, fi[i].Size, int64(0))
		assert.Equal(t, fi[i].Type, "")
		assert.Equal(t, fi[i].IsDir, false)
	}
}

func TestReadFile(t *testing.T) {
	dirPrefix := "readDir"
	testFileExt := ".txt"
	testFileName := "readMe"
	testFileContents := []byte{12, 3, 4, 1, 213, 213, 13}
	testFile := testFileName + testFileExt

	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	defer os.RemoveAll(tmpDir)
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

	assert.Equal(t, testFile, fc.Name)
	assert.Equal(t, testFileContents, fc.Content)
}

func TestUploadFile(t *testing.T) {

	//Generate a simple file
	dirPrefix := "uploadDir"
	inputFile := "uploadMe.txt"
	testFileContents := []byte{1, 123, 123, 213, 13, 213, 3}
	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	defer os.RemoveAll(tmpDir)
	createTempFile(t, tmpDir, inputFile, testFileContents)

	path := tmpDir + "/" + inputFile
	file, err := os.Open(path)
	assert.NoError(t, err)
	fileContents, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	file.Close()

	//Load the file we created above and use the write API to write new file
	testFile := "readme.txt"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", testFile)
	assert.NoError(t, err)
	part.Write(fileContents)

	err = writer.Close()
	assert.NoError(t, err)

	pwd, err := os.Getwd()
	filePath := pwd + "/" + testFile
	req, err := http.NewRequest("POST", "/write"+pwd, body)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)
	defer os.RemoveAll(filePath)

	//Read the file and check the filename and type
	req, err = http.NewRequest("GET", "/read"+filePath, nil)
	assert.NoError(t, err)

	res = httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, 200, res.Code)

	//Test Files
	var fc FileContent
	err = json.Unmarshal(res.Body.Bytes(), &fc)
	assert.NoError(t, err)

	assert.Equal(t, fc.Name, testFile)
	assert.Equal(t, fc.Content, testFileContents)
}

func TestDeleteFile(t *testing.T) {
	dirPrefix := "deleteDir"
	testFile := "ignoreme"
	//Create a file under testdir
	tmpDir := createTempDir(t, dirPrefix)
	defer os.RemoveAll(tmpDir)
	createTempFile(t, tmpDir, testFile, []byte{})

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	//Get the abosolute path for testing dir
	filePath := pwd + "/" + tmpDir + "/" + testFile

	req, err := http.NewRequest("DELETE", "/delete"+filePath, nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 200)

	//Readfile again, check the file content
	req, err = http.NewRequest("GET", "/read"+filePath, nil)
	assert.NoError(t, err)

	res = httptest.NewRecorder()
	newRouterServer().ServeHTTP(res, req)

	//Test Status Code
	assert.Equal(t, res.Code, 404)
}

func TestInvalidPath(t *testing.T) {
	type testCases struct {
		Cases  string
		URL    string
		Method string
	}

	data := []testCases{
		{"Load", "/scan", "GET"},
		{"Read", "/read", "GET"},
	}

	for _, v := range data {
		t.Run(v.Cases, func(t *testing.T) {
			req, err := http.NewRequest(v.Method, v.URL+invalidDir, nil)
			assert.NoError(t, err)

			res := httptest.NewRecorder()
			newRouterServer().ServeHTTP(res, req)

			//Test Status Code
			assert.Equal(t, res.Code, 404)
		})
	}
}
