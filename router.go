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
		// allow ajax reponses from browser
		w.Header().Set("Access-Control-Allow-Origin", router.cfg.CorsOrigin)

		ProceedImage(w, r)
	case "/api/v1/pdf_detect":
		w.Header().Set("Access-Control-Allow-Origin", router.cfg.CorsOrigin)
		proceedPDF(w, r)
	default:
		ShowForm(w)
	}
}
