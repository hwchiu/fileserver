package fileserver

import (
	"io/ioutil"
	"mime"
	"path"
	"regexp"
	"time"
)

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Type    string    `json:"type"`
	ModTime time.Time `json:"mtime"`
	IsDir   bool      `json:"isDir"`
}

func ScanDir(p string, excludePatterns []string) ([]FileInfo, error) {
	fileInfos := []FileInfo{}
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return fileInfos, err
	}

	exPattern := []regexp.Regexp{}
	for _, v := range excludePatterns {
		exPattern = append(exPattern, *regexp.MustCompile(v))
	}

CheckFile:
	for _, file := range files {
		for _, v := range exPattern {
			if v.MatchString(file.Name()) {
				continue CheckFile
			}
		}
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
