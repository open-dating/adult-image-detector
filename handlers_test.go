package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestShowForm(t *testing.T)  {
	w := httptest.NewRecorder()

	ShowForm(w)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	textBody := string(body)

	if strings.Contains(textBody, "<form") == false {
		t.Errorf("Expected html with <form tag. Got %s", textBody)
	}
}

func TestProceedImage(t *testing.T)  {
	r, _ := os.Open("./fixtures/forest.jpg")
	values := map[string]io.Reader{
		"image":  r,
	}
	form, contentType, err := createForm(values)
	if err != nil {
		t.Errorf("Error in create mulipart form %s", err.Error())
		return
	}

	req := httptest.NewRequest("POST", "/api/v1/detect", &form)
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", contentType)

	ProceedImage(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	textBody := string(body)

	if strings.Contains(textBody, "open_nsfw_score") == false {
		t.Errorf("Expected body with 'open_nsfw_score'. Got %s", textBody)
	}

	if strings.Contains(textBody, "image_name") == false {
		t.Errorf("Expected body with 'image_name'. Got %s", textBody)
	}
}

// support method for tests
func createForm(values map[string]io.Reader) (b bytes.Buffer, contentType string, err error)  {
	form := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an file
		if x, ok := r.(*os.File); ok {
			if fw, err = form.CreateFormFile(key, x.Name()); err != nil {
				return b, "", err
			}
		} else {
			// Add other fields
			if fw, err = form.CreateFormField(key); err != nil {
				return b, "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return b, "", err
		}

	}
	form.Close()
	return b, form.FormDataContentType(), nil
}
