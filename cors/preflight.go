package cors

import (
	"log"
	"net/http"
)

func PreflightHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isPreflightRequest(r) {
			handler(w, r)
			return
		}
		
		log.Println("Handle preflight request for", r.URL)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", r.Method)
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.WriteHeader(http.StatusOK)
	}
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
