package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"whatsmyip/assets"

	"github.com/gorilla/mux"
)

var (
	declaredRootPath string
)

// HandleStaticRequest handles all the incoming requests to fetch static resources.
func HandleStaticRequest(router *mux.Router, rootPath string) {
	declaredRootPath = rootPath
	router.
		Methods(http.MethodGet).
		Path(rootPath + "/{folder}/{file:.*(?:.js|.css)}").
		HandlerFunc(deliverStaticContent).
		Name("static")
}

// mapRequestAndInvoke check if the request can be processed and run the adequate function.
func deliverStaticContent(w http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	assetFolder := pathVars["folder"]
	assetFile := pathVars["file"]
	if DebugMode {
		log.Printf("Serving the file %s/%s\n", assetFolder, assetFile)
	}

	asset, err := assets.Asset(assetFolder + "/" + assetFile)
	if err != nil {
		if DebugMode {
			log.Printf("The file %s/%s was not found.\n", assetFolder, assetFile)
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		var contentType string
		switch {
		case strings.HasSuffix(assetFile, ".js"):
			contentType = "text/javascript"
		case strings.HasSuffix(assetFile, ".css"):
			contentType = "text/css"
		}
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(asset)
	}
}
