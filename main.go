package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Share, Nym!")

	// Connect to database

	// Execute migrations

	// Configure routes
	r := mux.NewRouter()
	routes.Routes(r)

	// Start server
	log.Println("Server initialized at port 8080 in https://localhost")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error initiating server:", err)
	}

}
