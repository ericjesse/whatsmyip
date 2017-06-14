package handlers

import (
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	// Db is the open connection to the database. It has to be set by the main code.
	Db *sql.DB
	// DebugMode indicates if debug messages have to be logged. It has to be set by the main code.
	DebugMode bool
)

// gzipResponseWriter is a writer to handle the GZIP compressed responses.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write is the implementation of the Writer interface for gzipResponseWriter.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// makeGzipHandler wraps the HTTP response writer by a GZIP writer if the client expects it.
func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	// The handler is created inline.
	return func(w http.ResponseWriter, req *http.Request) {
		// Check if the client expects a GZIP response.
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			fmt.Fprintln(os.Stdout, "No GZIP response is expected")
			fn(w, req)
		} else {
			// Define the response as one containing a GZIP body.
			w.Header().Set("Content-Encoding", "gzip")
			// Create a new GZIP writer surrounding the HTTP response writer.
			gz := gzip.NewWriter(w)
			// Make sure the writer is closed after the handler is done.
			gzr := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
			defer gz.Close()
			// Run the handler of the HTTP request.
			fn(gzr, req)
		}
	}
}

func writeResult(w http.ResponseWriter, req *http.Request, body interface{}) {
	var responseBody, contentType string
	accept := req.Header.Get("Accept") // Check what the client expects as format.
	if strings.Contains(accept, "xml") {
		// XML has to be explicitly requested.
		contentType = "application/xml"
		xmlBody, _ := xml.Marshal(body)
		responseBody = string(xmlBody)
	} else {
		// JSON is the default.
		contentType = "application/json"
		jsonBody, _ := json.Marshal(body)
		responseBody = string(jsonBody)
	}
	// Write the adequate headers.
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write the body to the response.
	fmt.Fprintf(w, responseBody)
}
