package main

import (
	"fmt"
	"net/http"
)

// router routes requests
type router struct {
	cfg params
}

// handler http responses
func (rtr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/v1/detect":
		switch r.Method {
		case http.MethodPost:
			// allow ajax reponses from browser
			w.Header().Set("Access-Control-Allow-Origin", rtr.cfg.CorsOrigin)

			proceedImage(w, r)
		default:
			ShowImageForm(w)
		}

	case "/api/v1/pdf_detect":
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Access-Control-Allow-Origin", rtr.cfg.CorsOrigin)
			proceedPDF(w, r)
		case http.MethodGet:
			ShowPDFForm(w)
		default:
			HandleError(w, fmt.Errorf("bad request. make http POST request instead"))
			return
		}

	default:
		HandleError(w, fmt.Errorf(`
		<html>
          <body>
			Bad request. Invalid endpoint. Avalaibled endpoints:
            <div><a href="/api/v1/detect">/api/v1/detect</a> - for detect in image</div>
            <div><a href="/api/v1/pdf_detect">/api/v1/pdf_detect</a> - for detect in pdf</div>
          </body>
		</html>
		`))
		return
	}
}
