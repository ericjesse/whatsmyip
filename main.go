package main

import (
	// Standard library packages
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

type ipAddress struct {
	IpAddressV4 string "" // Needs to start with a lowercase to be marshalled.
	Source      string ""
	Error       string ""
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Surround the HTTP response writer by a GZIP writer if the client expects it.
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
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			defer gz.Close()
			// Run the handler of the HTTP request.
			fn(gzr, req)
		}
	}
}

// Determine what the IP address of the remote caller is, based upon the HTTP request details.
func determineIP(req *http.Request) (ipAddress, error) {
	result := ipAddress{}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return result, errors.New(fmt.Sprintf("The user IP: %q is not IP:port", req.RemoteAddr))
	} else {
		// If the client is behind a non-anonymous proxy, the IP address is in the X-Forwarded-For header.
		// req.Header.Get is case-insensitive.
		result.IpAddressV4 = req.Header.Get("X-Forwarded-For")
		if result.IpAddressV4 == "" {
			// If no header can be read, directly extract the address for the request.
			userIP := net.ParseIP(ip)
			if userIP == nil {
				return result, errors.New(fmt.Sprintf("The user IP: %q is not IP:port", req.RemoteAddr))
			}
			result.IpAddressV4 = userIP.String()
			result.Source = "Remote address"
		} else {
			result.Source = "X-Forwarded-For"
		}
	}
	return result, nil
}

// After the IP address was determined, send it with the expected format
// to the HTTP response.
func determineAndSendIP(w http.ResponseWriter, req *http.Request) {
	result, error := determineIP(req)
	if error != nil {
		result.Error = error.Error()
	}

	var responseBody, contentType string
	accept := req.Header.Get("Accept") // Check what the client expects as format.
	if strings.Contains(accept, "xml") {
		// XML has to be explicitly requested.
		contentType = "application/xml"
		xmlBody, _ := xml.Marshal(&result)
		responseBody = string(xmlBody)
	} else {
		// JSON is the default.
		contentType = "application/json"
		jsonBody, _ := json.Marshal(result)
		responseBody = string(jsonBody)
	}
	// Write the adequate headers.
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write the body to the response.
	fmt.Fprintf(w, responseBody)
	return
}

func main() {
	var port string
	if len(os.Args) >= 2 {
		// Get the port from the command line argument.
		port = os.Args[1]
	} else {
		// Get the port to use from the environment variables.
		port = os.Getenv("PORT")
	}

	if port == "" {
		// The service cannot run without port to listen.
		log.Fatal("The port to listen could not be found")
	} else {
		// Surround the handler with a GZIP writter handler.
		http.HandleFunc("/ip", makeGzipHandler(determineAndSendIP))

		fmt.Fprintf(os.Stdout, "Service starting on port %s\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
