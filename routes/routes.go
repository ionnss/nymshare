package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router) {
	// Serve static files and assets
	fileServer := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Pages
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/index.html")
	}).Methods("GET")

	r.HandleFunc("/channel/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/channel.html")
	}).Methods("GET")

	r.HandleFunc("/share/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/share.html")
	}).Methods("GET")

}
