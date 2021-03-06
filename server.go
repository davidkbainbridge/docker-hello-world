package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"
)

type StatusUpdate struct {
	Status int `json:"status"`
}

func main() {
	// I know this is not thread safe
	var healthStatus int = 200

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			log.Printf("[INFO] Received health request, returning %d", healthStatus)
			w.WriteHeader(healthStatus)
		case http.MethodPost:
			status := StatusUpdate{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&status)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Printf("[INFO] Health status updated to %d", status.Status)
				healthStatus = status.Status
				w.WriteHeader(http.StatusOK)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

		}
	})
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[INFO] Container Start Hook Called")
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[INFO] Container Stop Hook Called")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "HOST: %s\n", hostname)
		fmt.Fprintf(w, "ADDRESSES:\n")
		addrs, _ := net.InterfaceAddrs()
		for _, addr := range addrs {
			fmt.Fprintf(w, "    %s\n", addr.String())
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
