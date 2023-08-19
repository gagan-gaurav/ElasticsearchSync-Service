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
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")                     // Create user
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")                    // Get user
	r.HandleFunc("/users/update/{id}", handlers.UpdateUser).Methods("POST")         // Update user
	r.HandleFunc("/users/delete/{id}", handlers.DeleteUser).Methods("DELETE")       // Delete user
	r.HandleFunc("/hashtags", handlers.CreateHashtag).Methods("POST")               // Create hashtag
	r.HandleFunc("/hashtags/{id}", handlers.GetHashtag).Methods("GET")              // Get hashtag
	r.HandleFunc("/hashtags/update/{id}", handlers.UpdateHashtag).Methods("POST")   // Update hashtag
	r.HandleFunc("/hashtags/delete/{id}", handlers.DeleteHashtag).Methods("DELETE") // Delete hastag
	r.HandleFunc("/projects", handlers.CreateProject).Methods("POST")               // Create project
	r.HandleFunc("/projects/{id}", handlers.GetProject).Methods("GET")              // Get project
	r.HandleFunc("/projects/update/{id}", handlers.UpdateProject).Methods("POST")   // Update project
	r.HandleFunc("/projects/delete/{id}", handlers.DeleteProject).Methods("DELETE") // Delete project

	http.Handle("/", r)
}
