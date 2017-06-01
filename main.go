package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"flag"
	"log"
	"os"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if isPreflightRequest(r) {
		preflightHandler(w, r)
		return
	}

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
		log.Fatalln("Error on fetching", err)
		return
	}

	destinationHeader := w.Header()
	copyHeader(&response.Header, &destinationHeader)
	w.WriteHeader(response.StatusCode)

	defer response.Body.Close()
	_, err = io.Copy(w, response.Body)
	if err != nil {
		log.Fatalln("Error on copying response", err)
		return
	}
}

func copyHeader(source *http.Header, destination *http.Header) {
	for key, value := range *source {
		destination.Add(key, strings.Join(value, " "))
	}
}

func main() {
	portNumber := flag.Int("port", 8080, "Port to listen to")

	flag.Parse()

	port := strconv.Itoa(*portNumber)

	log.SetOutput(os.Stdout)
	fmt.Println("Running cors proxy on http://localhost:" + port)
	fmt.Println("Use http://localhost:" + port + "/?url= to proxy url calls")
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:"+port, nil)
}
