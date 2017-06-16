package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net"

	"net/http"
	"regexp"
	"time"
	"whatsmyip/assets"

	"github.com/gorilla/mux"
)

const (
	insertIPResultQuery string = `INSERT INTO "ipCheck" ("ipAddressV4", source) VALUES ($1, $2)`
	insertGeoLocQuery   string = `INSERT INTO "geoLoc"
		("ipAddressV4", provider, city, country, "countryCode", region, timezone, "zipCode", latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	countSimilarRecentGeoLoc string = `SELECT COUNT(1) FROM "geoLoc"
		WHERE "ipAddressV4" = $1 and instant > $2`
)

var (
	ipAddressPattern *regexp.Regexp = regexp.MustCompile(`(([0-9]{1,3}\.)){3}[0-9]{1,3}`)
	templateFiles                   = []string{"tpl/ip.html"}
	ipTemplate       *template.Template
)

func init() {
	tplContent, err := assets.Asset("ip.html")
	if err != nil {
		log.Fatalln(err)
	}
	ipTemplate, err = template.New("ip").Parse(string(tplContent))
	if err != nil {
		log.Fatalln(err)
	}
}

// ipAddress is the represention of the IP address of the remote caller, providing also the source from where it was evaluated.
type ipAddress struct {
	IPAddressV4 string `json:"ipAddressV4" xml:"ipAddressV4"`
	Source      string `json:"source" xml:"source"`
	Error       string `json:"error" xml:"error"`
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

// HandleIPRequest handles all the incoming requests to determine the source IP address.
func HandleIPRequest(router *mux.Router, rootPath string) {
	router.
		Methods(http.MethodGet).
		Path(rootPath).
		HeadersRegexp("Accept", ".*((application/((xhtml+)?xml|json|javascript))|(text/x?html)).*").
		HandlerFunc(provideIP).
		Name("ip")
}

// provideIP processes the request to known the IP address of the remote caller of a GET call.
func provideIP(w http.ResponseWriter, req *http.Request) {
	result, error := evaluateIPAddress(req)
	if error != nil {
		result.Error = error.Error()
	}
	go result.fetchGeoAndPersist()

	writeReponse(w, req, ipTemplate, result)
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

// fetchGeoAndPersist aims at getting the geographical location of the IP address and save it in the database.
func (ipAddress *ipAddress) fetchGeoAndPersist() {
	_, err := Db.Exec(insertIPResultQuery, ipAddress.IPAddressV4, ipAddress.Source)
	if err != nil {
		log.Println("Error", err)
		return
	}

	// Fetch the address only if the IP address is valid.
	if !ipAddressPattern.Match([]byte(ipAddress.IPAddressV4)) {
		log.Debugln("The address cannot be used for geo localization")
		return
	}

	// Lookup of a recent similar geo localisation for the same IP Address.
	weekBefore := time.Now().Add(time.Hour * 24 * 7)
	row := Db.QueryRow(countSimilarRecentGeoLoc, ipAddress.IPAddressV4, weekBefore)
	var count int
	if err = row.Scan(&count); err != nil {
		log.Println("Error", err)
		return
	}

	if count == 0 {
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
		log.Debugf("Result for Geo localisation: %v\n", *geoLoc)

		// Persist the geo localization.
		_, err = Db.Exec(insertGeoLocQuery, ipAddress.IPAddressV4, geoLoc.Provider, geoLoc.City, geoLoc.Country, geoLoc.CountryCode, geoLoc.Region, geoLoc.Timezone, geoLoc.ZipCode, geoLoc.Latitude, geoLoc.Longitude)
		if err != nil {
			log.Println("Error", err)
			return
		}
	} else {
		log.Println("No need to refresh the geo localisation", err)
	}

}
