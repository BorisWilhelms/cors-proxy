package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if isPreflightRequest(r) {
		preflightHandler(w, r)
		return
	}
	log.Println("Incoming request", r.URL.String())
	if r.URL.RawQuery == "" {
		http.NotFound(w, r)
		return
	}

	remoteURL := r.URL.RawQuery[len("url="):]
	if remoteURL == "" {
		http.NotFound(w, r)
		return
	}

	log.Println("Fetching url", remoteURL)
	request, err := http.NewRequest(r.Method, remoteURL, r.Body)
	copyHeader(&r.Header, &request.Header)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln("Error on fetching", err.Error())
		return
	}

	destinationHeader := w.Header()
	copyHeader(&response.Header, &destinationHeader)
	w.WriteHeader(response.StatusCode)

	defer response.Body.Close()
	_, err = io.Copy(w, response.Body)
	if err != nil {
		log.Fatalln("Error on copying response", err.Error())
		return
	}
}

func copyHeader(source *http.Header, destination *http.Header) {
	for key, value := range *source {
		destination.Add(key, strings.Join(value, " "))
	}
}

func main() {
	log.SetOutput(os.Stdout)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
