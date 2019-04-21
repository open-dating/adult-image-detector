package main

import (
	"net/http"
)

type Router struct {
	cfg Params
}

// handler http responses
func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/v1/detect":
		ProceedImage(w, r)
	default:
		ShowForm(w)
	}
}