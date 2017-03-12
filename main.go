package main

import (
	// Standard library packages
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)



// https://blog.golang.org/context/userip/userip.go
func getIP(w http.ResponseWriter, req *http.Request) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Errorf("The user IP: %q is not IP:port", req.RemoteAddr)
		return
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forwardedFor := req.Header.Get("X-Forwarded-For")
	if forwardedFor == "" {
		userIP := net.ParseIP(ip)
		if userIP == nil {
			//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
			fmt.Fprintf(w, "userip: %q is not IP:port", req.RemoteAddr)
			return
		}
		fmt.Fprintf(w, "{\"ipAddressV4\":\"%s\", \"source\":\"Remote address\"}", userIP)
	} else {
		fmt.Fprintf(w, "{\"ipAddressV4\":\"%s\", \"source\":\"Header X-Forwarded-For\"}", forwardedFor)
	}
}

func main() {
	port := os.Getenv("PORT");
	http.HandleFunc("/ip", getIP)
	fmt.Fprintf(os.Stdout, "Service starting on port %s\n", port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
