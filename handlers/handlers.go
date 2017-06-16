// Package handlers is in charge of handling the different incoming requests.
package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"whatsmyip/logger"
)

var (
	// Db is the open connection to the database. It has to be set by the main code.
	Db  *sql.DB
	log = logger.GetLogger()
)

func writeReponse(w http.ResponseWriter, req *http.Request, template *template.Template, body interface{}) {
	var contentType string
	accept := req.Header.Get("Accept")
	log.Println("Accept", accept)
	switch {
	case strings.Contains(accept, "text/html"):
		contentType = "text/html"
		fallthrough
	case strings.Contains(accept, "application/xhtml+xml"):
		if contentType == "" {
			contentType = "application/xhtml+xml"
		}
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		template.Execute(w, body)
	default:
		writeApiReponse(w, req, accept, body)
	}
}

func writeApiReponse(w http.ResponseWriter, req *http.Request, accept string, body interface{}) {
	var responseBody, contentType string

	// Check what the client expects as format.
	switch {
	case strings.Contains(accept, "application/xml"):
		contentType = "application/xml"
		xmlBody, _ := xml.Marshal(body)
		responseBody = string(xmlBody)
	case strings.Contains(accept, "application/json"):
		contentType = "application/json"
		fallthrough
	case strings.Contains(accept, "application/javascript"):
		fallthrough
	default:
		if contentType == "" { // Default response type is JSON.
			contentType = "application/javascript"
		}
		jsonBody, _ := json.Marshal(body)
		responseBody = string(jsonBody)
	}
	// Write the adequate headers.
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write the body to the response.
	fmt.Fprintf(w, responseBody)
}
