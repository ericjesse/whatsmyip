// Package handlers is in charge of handling the different incoming requests.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	insertIPResultQuery string = `INSERT INTO "ipCheck" ("ipAddressV4", source) VALUES ($1, $2)`
	insertGeoLocQuery   string = `INSERT INTO "geoLoc"
		("ipAddressV4", provider, city, country, "countryCode", region, timezone, "zipCode", latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
)

// ipAddress is the represention of the IP address of the remote caller, providing also the source from where it was evaluated.
type ipAddress struct {
	IPAddressV4 string `json:"ipAddressV4"`
	Source      string `json:"source"`
	Error       string `json:"error"`
}

type geoLoc struct {
	Provider    string  `json:"isp"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"regionName"`
	Timezone    string  `json:"timezone"`
	ZipCode     string  `json:"zip"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
}

// fetchGeoAndPersist aims at getting the geographical location of the IP address and save it in the database.
func (ipAddress *ipAddress) fetchGeoAndPersist() {
	_, err := Db.Exec(insertIPResultQuery, ipAddress.IPAddressV4, ipAddress.Source)
	if err != nil {
		log.Println("Error", err)
		return
	}

	res, err := http.Get("http://ip-api.com/json/" + ipAddress.IPAddressV4)
	if err != nil {
		log.Println("Error", err)
		return
	}

	// Unmarshal the response to the structure.
	defer res.Body.Close()
	var geoLoc = new(geoLoc)
	err = json.NewDecoder(res.Body).Decode(geoLoc)
	if err != nil {
		log.Println("Error", err)
		return
	}
	if DebugMode {
		log.Printf("Result for Geo localisation: %v\n", *geoLoc)
	}

	// Persist the geo localization.
	_, err = Db.Exec(insertGeoLocQuery, ipAddress.IPAddressV4, geoLoc.Provider, geoLoc.City, geoLoc.Country, geoLoc.CountryCode, geoLoc.Region, geoLoc.Timezone, geoLoc.ZipCode, geoLoc.Latitude, geoLoc.Longitude)
	if err != nil {
		log.Println("Error", err)
		return
	}

}

// HandleIPRequest handles all the incoming requests to determine the source IP address.
func HandleIPRequest() http.HandlerFunc {
	return makeGzipHandler(mapRequestAndInvoke)
}

// mapRequestAndInvoke check if the request can be processed and run the adequate function.
func mapRequestAndInvoke(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getIP(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// getIP processes the request to known the IP address of the remote caller of a GET call.
func getIP(w http.ResponseWriter, req *http.Request) {
	result, error := evaluateIPAddress(req)
	if error != nil {
		result.Error = error.Error()
	}
	go result.fetchGeoAndPersist()
	writeResult(w, req, result)
}

// evaluateIPAddress determines what the IP address of the remote caller is, based upon the HTTP request details.
func evaluateIPAddress(req *http.Request) (*ipAddress, error) {
	result := ipAddress{}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return &result, fmt.Errorf("The user IP: %q is not IP:port", req.RemoteAddr)
	}

	// If the client is behind a non-anonymous proxy, the IP address is in the X-Forwarded-For header.
	// req.Header.Get is case-insensitive.
	result.IPAddressV4 = req.Header.Get("X-Forwarded-For")
	if result.IPAddressV4 == "" {
		// If no header can be read, directly extract the address for the request.
		userIP := net.ParseIP(ip)
		if userIP == nil {
			return &result, fmt.Errorf("The user IP: %q is not IP:port", req.RemoteAddr)
		}
		result.IPAddressV4 = userIP.String()
		result.Source = "Remote address"
	} else {
		result.Source = "X-Forwarded-For"
	}
	return &result, nil
}
