package routes

import (
	"fold/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetRouter() {
	// Create a new mux router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/hashtags", handlers.CreateHashtag).Methods("POST")
	r.HandleFunc("/hashtags/{id}", handlers.GetHashtag).Methods("GET")
	r.HandleFunc("/projects", handlers.CreateProject).Methods("POST")
	r.HandleFunc("/projects/{id}", handlers.GetProject).Methods("GET")

	http.Handle("/", r)
}
