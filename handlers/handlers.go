package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

var (
	// Db is the open connection to the database. It has to be set by the main code.
	Db *sql.DB
	// DebugMode indicates if debug messages have to be logged. It has to be set by the main code.
	DebugMode bool
)

func writeResult(w http.ResponseWriter, req *http.Request, body interface{}) {
	var responseBody, contentType string

	// Check what the client expects as format.
	switch accept := req.Header.Get("Accept"); accept {
	case "application/xml":
		contentType = "application/xml"
		xmlBody, _ := xml.Marshal(body)
		responseBody = string(xmlBody)
	case "application/json":
		contentType = "application/json"
		fallthrough
	case "application/javascript":
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
