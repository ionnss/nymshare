package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router) {
	// Serve static files and assets
	fileServer := http.FileServer(http.Dir("."))
	r.PathPrefix("/static/").Handler(fileServer)

	// Pages
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/index.html")
	}).Methods("GET")
}
