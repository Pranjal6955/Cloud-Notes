package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"cloud-notes/handlers"
	"cloud-notes/ws"
)

func SetupRoutes(router *mux.Router, hub *ws.Hub) {
	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Notes handler
	notesHandler := handlers.NewNotesHandler(hub)

	// Notes routes
	api.HandleFunc("/notes", notesHandler.GetNotes).Methods("GET")
	api.HandleFunc("/notes", notesHandler.CreateNote).Methods("POST")
	api.HandleFunc("/notes/{id}", notesHandler.GetNote).Methods("GET")
	api.HandleFunc("/notes/{id}", notesHandler.UpdateNote).Methods("PUT")
	api.HandleFunc("/notes/{id}", notesHandler.DeleteNote).Methods("DELETE")

	// WebSocket route
	router.HandleFunc("/ws", hub.HandleWebSocket)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// CORS preflight
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})
}