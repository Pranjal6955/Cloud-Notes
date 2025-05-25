package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"cloud-notes/models"
	"cloud-notes/ws"
)

type NotesHandler struct {
	hub *ws.Hub
}

func NewNotesHandler(hub *ws.Hub) *NotesHandler {
	return &NotesHandler{hub: hub}
}

func (h *NotesHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	note := models.Store.CreateNote(req.Title, req.Content)

	// Broadcast to all connected clients
	update := models.NoteUpdate{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
		Type:    "create",
	}
	h.hub.Broadcast(update)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NotesHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes := models.Store.GetAllNotes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NotesHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	note, exists := models.Store.GetNote(id)
	if !exists {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NotesHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	note, exists := models.Store.UpdateNote(id, req.Title, req.Content)
	if !exists {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Broadcast update to all connected clients
	update := models.NoteUpdate{
		ID:      note.ID,
		Title:   req.Title,
		Content: req.Content,
		Type:    "update",
	}
	h.hub.Broadcast(update)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NotesHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists := models.Store.DeleteNote(id)
	if !exists {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Broadcast deletion to all connected clients
	update := models.NoteUpdate{
		ID:   id,
		Type: "delete",
	}
	h.hub.Broadcast(update)

	w.WriteHeader(http.StatusNoContent)
}