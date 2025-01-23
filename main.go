package main

import (
	"log"
	"net/http"
	"nymshare/db"
	"nymshare/routes"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Share, Nym!")

	// Connect to database
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer conn.Close()

	// Execute migrations
	if err := db.ExecuteMigrations(conn); err != nil {
		log.Fatal("Error executing migrations:", err)
	}

	log.Println("Database successfully initialized")

	// Configure routes
	r := mux.NewRouter()
	routes.ConfigureRoutes(r)

	// Start server
	log.Println("Server initialized at port 8080 in https://localhost")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error initiating server:", err)
	}

}
