package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// procced multipart form and save file
func SaveUploadFile(r *http.Request) (filePath string, fileName string, err error)  {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	fileName = normalizeFileName(handler.Filename)

	filePath, err = filepath.Abs("./uploads/" + fileName)
	if err != nil {
		return "", "", err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return filePath, handler.Filename, nil
}

// remove backslahes from image name
func normalizeFileName(fileName string) string  {
	// TODO generate uuid?
	fileName = strings.Replace(fileName, "\\", "", -1)
	fileName = strings.Replace(fileName, "/", "", -1)
	return fileName
}

// remove file
func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}