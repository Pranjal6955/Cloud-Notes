package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"collaborative-notes/config"
	"collaborative-notes/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Hub struct {
	clients    map[string]map[*Client]bool // noteId -> clients
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	noteID string
	userID uint
}

type Message struct {
	Type      string      `json:"type"`
	NoteID    string      `json:"note_id"`
	UserID    uint        `json:"user_id"`
	Operation *Operation  `json:"operation,omitempty"`
	Content   string      `json:"content,omitempty"`
	Version   int         `json:"version,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type Operation struct {
	Type     string `json:"type"`     // insert, delete, retain
	Position int    `json:"position"`
	Content  string `json:"content,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.noteID] == nil {
				h.clients[client.noteID] = make(map[*Client]bool)
			}
			h.clients[client.noteID][client] = true
			log.Printf("Client connected to note %s", client.noteID)

			// Send current note content to the new client
			go h.sendCurrentContent(client)

		case client := <-h.unregister:
			if clients, exists := h.clients[client.noteID]; exists {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.noteID)
					}
					log.Printf("Client disconnected from note %s", client.noteID)
				}
			}

		case message := <-h.broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Broadcast to all clients of the specific note
			if clients, exists := h.clients[msg.NoteID]; exists {
				for client := range clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

func (h *Hub) sendCurrentContent(client *Client) {
	noteID, err := strconv.ParseUint(client.noteID, 10, 32)
	if err != nil {
		return
	}

	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		return
	}

	message := Message{
		Type:    "sync",
		NoteID:  client.noteID,
		Content: note.Content,
		Version: note.Version,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	select {
	case client.send <- data:
	default:
		close(client.send)
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request, noteID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Extract user ID from query params or headers
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		conn.Close()
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		noteID: noteID,
		userID: uint(userID),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		msg.UserID = c.userID
		msg.NoteID = c.noteID

		// Handle different message types
		switch msg.Type {
		case "operation":
			c.handleOperation(&msg)
		case "cursor":
			// Broadcast cursor position to other clients
			data, _ := json.Marshal(msg)
			c.hub.broadcast <- data
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

func (c *Client) handleOperation(msg *Message) {
	if msg.Operation == nil {
		return
	}

	noteID, err := strconv.ParseUint(c.noteID, 10, 32)
	if err != nil {
		return
	}

	// Save operation to database
	operation := models.Operation{
		NoteID:   uint(noteID),
		UserID:   c.userID,
		Type:     msg.Operation.Type,
		Position: msg.Operation.Position,
		Content:  msg.Operation.Content,
		Length:   msg.Operation.Length,
		Version:  msg.Version,
	}

	if err := config.DB.Create(&operation).Error; err != nil {
		log.Printf("Error saving operation: %v", err)
		return
	}

	// Apply operation to note content
	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		return
	}

	newContent := applyOperation(note.Content, msg.Operation)
	note.Content = newContent
	note.Version++

	if err := config.DB.Save(&note).Error; err != nil {
		log.Printf("Error updating note: %v", err)
		return
	}

	// Broadcast operation to other clients
	msg.Version = note.Version
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	c.hub.broadcast <- data
}

func applyOperation(content string, op *Operation) string {
	switch op.Type {
	case "insert":
		if op.Position > len(content) {
			return content + op.Content
		}
		return content[:op.Position] + op.Content + content[op.Position:]
	case "delete":
		if op.Position >= len(content) {
			return content
		}
		end := op.Position + op.Length
		if end > len(content) {
			end = len(content)
		}
		return content[:op.Position] + content[end:]
	}
	return content
}
