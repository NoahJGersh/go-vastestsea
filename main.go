package main

import (
	"log"
	"net/http"
)

func main() {
	// Construct mux
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /vastestsea/lang", languagesHandler)
	serveMux.HandleFunc("GET /vastestsea/lang/{language}", languageHandler)

	// Run server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	log.Fatal(server.ListenAndServe())
}
