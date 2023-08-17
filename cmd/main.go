package main

import (
	"fmt"
	"fold/internal/database"
	"fold/internal/routes"
	"net/http"
)

func main() {
	database.MakeDatabaseConnection()
	routes.SetRouter()

	// Start the server at port 8080
	port := "8080"
	fmt.Printf("Server started on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
