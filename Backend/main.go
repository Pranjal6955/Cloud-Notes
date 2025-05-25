package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"cloud-notes/routes"
	"cloud-notes/ws"
)

func main() {
	// Initialize WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize router
	router := mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(router, hub)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://cloud-notes.local"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Get port from environment or default to
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}