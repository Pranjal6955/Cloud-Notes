package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Collaborators []string `json:"collaborators"`
}

type NoteUpdate struct {
	ID      string `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Type    string `json:"type"` // "title", "content", "create", "delete"
}

// In-memory storage for notes (replace with database in production)
type NoteStore struct {
	notes map[string]*Note
	mutex sync.RWMutex
}

var Store = &NoteStore{
	notes: make(map[string]*Note),
}

func (ns *NoteStore) CreateNote(title, content string) *Note {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	note := &Note{
		ID:          uuid.New().String(),
		Title:       title,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Collaborators: []string{},
	}

	ns.notes[note.ID] = note
	return note
}

func (ns *NoteStore) GetNote(id string) (*Note, bool) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	note, exists := ns.notes[id]
	return note, exists
}

func (ns *NoteStore) GetAllNotes() []*Note {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	notes := make([]*Note, 0, len(ns.notes))
	for _, note := range ns.notes {
		notes = append(notes, note)
	}
	return notes
}

func (ns *NoteStore) UpdateNote(id, title, content string) (*Note, bool) {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	note, exists := ns.notes[id]
	if !exists {
		return nil, false
	}

	if title != "" {
		note.Title = title
	}
	if content != "" {
		note.Content = content
	}
	note.UpdatedAt = time.Now()

	return note, true
}

func (ns *NoteStore) DeleteNote(id string) bool {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	_, exists := ns.notes[id]
	if exists {
		delete(ns.notes, id)
	}
	return exists
}