package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if isPreflightRequest(r) {
		preflightHandler(w, r)
		return
	}

	remoteURL := r.URL.RawQuery[len("url="):]
	if remoteURL == "" {
		http.NotFound(w, r)
		return
	}
	request, err := http.NewRequest(r.Method, remoteURL, r.Body)
	copyHeader(r.Header, request.Header)
	client := &http.Client{}

	fmt.Println(request.UserAgent())
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	copyHeader(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
}

func copyHeader(source http.Header, destination http.Header) {
	for key, value := range source {
		destination.Add(key, strings.Join(value, " "))
	}
}

func main() {
	os.Setenv("HTTP_PROXY", "http://localhost:8888")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
