package main

import (
	"log"
	"net/http"
)

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle preflight request for", r.URL)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET HEAD POST PUT DELETE TRACE CONNECT")
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.WriteHeader(http.StatusOK)
}

func isPreflightRequest(r *http.Request) bool {
	if r.Method != "OPTIONS" {
		return false
	}
	var accessControlRequestMethod = r.Header.Get("Access-Control-Request-Method")
	if accessControlRequestMethod != "" {
		return true
	}

	return false
}
