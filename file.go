package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/google/uuid"
	"strings"
	"time"
)

type UploadedFileInfo struct {
	FilePath 			string
	Filename 			string
	FileExt  			string
	SaveAsFilename      string
	disableOpenNsfw 	bool
	disableAnAlgorithm 	bool
	debug               bool
}

// procced multipart form and save file
func HandleUploadFileForm(r *http.Request) (parsedForm UploadedFileInfo, err error)  {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		return parsedForm, err
	}
	defer file.Close()

	parts := strings.Split(handler.Filename, ".")
	fileExt := parts[len(parts) - 1]

	parsedForm.SaveAsFilename = time.Now().Format(time.RFC3339) + "_" + uuid.New().String() + "." + fileExt

	filePath, err := filepath.Abs("./uploads/" + parsedForm.SaveAsFilename)
	if err != nil {
		return parsedForm, err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return parsedForm, err
	}
	defer f.Close()
	io.Copy(f, file)

	parsedForm.Filename = handler.Filename
	parsedForm.FilePath = filePath

	parsedForm.disableAnAlgorithm = r.FormValue("disableAnAlgorithm") != ""
	parsedForm.disableOpenNsfw = r.FormValue("disableOpenNsfw") != ""
	parsedForm.debug = r.FormValue("debug") != ""

	return parsedForm, nil
}

// remove file
func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}