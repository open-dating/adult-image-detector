package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strconv"
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
    // forest.jpg
	body, err := uploadFixtureAndGetResult("./fixtures/forest.jpg")
	if err != nil {
		t.Error("Image upload err:", err.Error())
	}
	result := ImageScoringResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Errorf("Expected body with json. Got %s", string(body))
	}
	if result.AnAlgorithmForNudityDetection != false {
		t.Errorf("Expected AnAlgorithmForNudityDetection false for forest got AnAlgorithmForNudityDetection %s", strconv.FormatBool(result.AnAlgorithmForNudityDetection))
	}
	if result.OpenNsfwScore < 0.15 == false {
		t.Errorf("Expected OpenNsfwScore < 0.15 for forest got OpenNsfwScore %f", result.OpenNsfwScore)
	}

    // big_boobs.cropped.png
	body, err = uploadFixtureAndGetResult("./fixtures/big_boobs.cropped.png")
	if err != nil {
		t.Error("Image upload err:", err.Error())
	}
	result = ImageScoringResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Errorf("Expected body with json. Got %s", string(body))
	}
	if result.AnAlgorithmForNudityDetection != true {
		t.Errorf("Expected AnAlgorithmForNudityDetection true for big_boobs.cropped.png got AnAlgorithmForNudityDetection %s", strconv.FormatBool(result.AnAlgorithmForNudityDetection))
	}
	if result.OpenNsfwScore > 0.5 == false {
		t.Errorf("Expected OpenNsfwScore > 0.5 for big_boobs.cropped.png got OpenNsfwScore %f", result.OpenNsfwScore)
	}
}

// support method for tests image recognition result
func uploadFixtureAndGetResult(filePath string) (body []byte, err error) {
	r, _ := os.Open(filePath)
	values := map[string]io.Reader{
		"image":  r,
	}
	form, contentType, err := createForm(values)
	if err != nil {
		return body, err
	}

	req := httptest.NewRequest("POST", "/api/v1/detect", &form)
	req.Header.Set("Content-Type", contentType)

	w := httptest.NewRecorder()
	ProceedImage(w, req)

	resp := w.Result()
	return ioutil.ReadAll(resp.Body)
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
