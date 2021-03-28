package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type UploadedFileInfo struct {
	FilePath           string
	Filename           string
	FileExt            string
	SaveAsFilename     string
	disableOpenNsfw    bool
	disableAnAlgorithm bool
	debug              bool
	password           string
}

func getImagesFromPDF(fp, dir, password string) ([]string, string, error) {
	var extractedImages []string

	if err := api.ExtractImagesFile(fp, dir, nil, &pdfcpu.Configuration{
		Reader15:         true,
		DecodeAllStreams: true,
		UserPW:           password,
	}); err != nil {
		return nil, "", err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}

	for _, v := range files {
		if isImage(v.Name()) {
			extractedImages = append(extractedImages, filepath.Join(dir, v.Name()))
		}
	}

	return extractedImages, dir, nil
}

func isImage(name string) bool {
	parts := strings.Split(name, ".")
	fileExt := strings.ToLower(parts[len(parts)-1])

	return fileExt == "jpg" || fileExt == "jpeg" || fileExt == "png" || fileExt == "gif"
}

// procced multipart form and save file
func HandleUploadFileForm(r *http.Request, dir string) (parsedForm UploadedFileInfo, err error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		return parsedForm, err
	}
	defer file.Close()

	parts := strings.Split(handler.Filename, ".")
	fileExt := parts[len(parts)-1]

	parsedForm.SaveAsFilename = time.Now().Format(time.RFC3339) + "_" + uuid.New().String() + "." + fileExt

	filePath, err := filepath.Abs(filepath.Join(dir, parsedForm.SaveAsFilename))
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
	parsedForm.FileExt = fileExt

	parsedForm.disableAnAlgorithm = r.FormValue("disableAnAlgorithm") != ""
	parsedForm.disableOpenNsfw = r.FormValue("disableOpenNsfw") != ""
	parsedForm.debug = r.FormValue("debug") != ""
	parsedForm.password = r.FormValue("password")
	return parsedForm, nil
}

// remove file
func RemoveFile(filePath string) {
	err := os.RemoveAll(filePath)
	if err != nil {
		panic(err)
	}
}
