package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/BorisWilhelms/cors-proxy/cors"

	"flag"
	"log"
	"os"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request", r.URL)
	if r.URL.RawQuery == "" || !strings.Contains(r.URL.String(), "url=") {
		http.NotFound(w, r)
		return
	}

	remoteURL := r.URL.RawQuery[len("url="):]
	if remoteURL == "" {
		http.NotFound(w, r)
		return
	}

	response, err := fetchRemote(r.Method, remoteURL, &r.Header, r.Body)
	if err != nil {
		log.Println("Error on fetching", err)
		return
	}
	sendResponse(w, response)
	if err != nil {
		log.Println("Error on copying response", err)
		return
	}
}

func sendResponse(w http.ResponseWriter, response *http.Response) error {
	destinationHeader := w.Header()
	copyHeader(&response.Header, &destinationHeader)
	w.WriteHeader(response.StatusCode)

	defer response.Body.Close()
	_, err := io.Copy(w, response.Body)
	return err
}

func fetchRemote(method, url string, header *http.Header, body io.Reader) (*http.Response, error) {
	log.Println("Fetching url", url)

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	copyHeader(header, &request.Header)

	client := &http.Client{}
	return client.Do(request)
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
	http.HandleFunc("/", cors.PreflightHandler(handler))
	http.ListenAndServe("localhost:"+port, nil)
}
