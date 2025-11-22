package socket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Constants for WebSocket configuration
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB for notification data
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// MessageTypeNotification represents a notification message
	MessageTypeNotification MessageType = "notification"
	// MessageTypeSystem represents a system message
	MessageTypeSystem MessageType = "system"
	// MessageTypePing represents a ping message
	MessageTypePing MessageType = "ping"
	// MessageTypePong represents a pong message
	MessageTypePong MessageType = "pong"
	// MessageTypeError represents an error message
	MessageTypeError MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client represents a WebSocket connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	mu     sync.Mutex
}

// Hub manages all WebSocket connections
type Hub struct {
	// Map of all connected clients
	clients map[*Client]bool

	// Map of clients by userID for targeted messages
	userClients map[string][]*Client

	// Channel for incoming messages to broadcast
	broadcast chan Message

	// Channel to register new clients
	register chan *Client

	// Channel to unregister clients
	unregister chan *Client

	// Mutex to protect map access
	mu sync.RWMutex

	// Context for graceful shutdown
	ctx context.Context

	// Cancel function
	cancel context.CancelFunc

	// WaitGroup for goroutines
	wg sync.WaitGroup
}

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking for production
		return true
	},
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[string][]*Client),
		broadcast:   make(chan Message, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Run starts the hub and processes client registration, unregistration, and broadcasts
func (h *Hub) Run() {
	h.wg.Add(1)
	defer h.wg.Done()

	for {
		select {
		case <-h.ctx.Done():
			log.Println("Hub shutting down...")
			h.closeAllConnections()
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	log.Println("Initiating hub shutdown...")
	h.cancel()
	h.wg.Wait()
	log.Println("Hub shutdown complete")
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	if client.userID != "" {
		h.userClients[client.userID] = append(h.userClients[client.userID], client)
		log.Printf("Client registered for user %s (total connections: %d)", client.userID, len(h.userClients[client.userID]))
	}
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		if client.userID != "" {
			// Remove client from userClients
			clients := h.userClients[client.userID]
			for i, c := range clients {
				if c == client {
					h.userClients[client.userID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}

			// If no more clients for this user, remove the user entry
			if len(h.userClients[client.userID]) == 0 {
				delete(h.userClients, client.userID)
				log.Printf("All connections closed for user %s", client.userID)
			} else {
				log.Printf("Client unregistered for user %s (remaining connections: %d)", client.userID, len(h.userClients[client.userID]))
			}
		}
	}
}

// closeAllConnections closes all client connections
func (h *Hub) closeAllConnections() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.send)
		client.conn.Close()
	}
	h.clients = make(map[*Client]bool)
	h.userClients = make(map[string][]*Client)
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message Message) {
	message.Timestamp = time.Now()
	select {
	case h.broadcast <- message:
	case <-h.ctx.Done():
		log.Println("Hub is shutting down, cannot broadcast message")
	}
}

// SendToUser sends a message to a specific user's connections
func (h *Hub) SendToUser(userID string, message Message) {
	message.Timestamp = time.Now()

	h.mu.RLock()
	clients := h.userClients[userID]
	h.mu.RUnlock()

	if len(clients) == 0 {
		log.Printf("No active connections for user %s", userID)
		return
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for user %s: %v", userID, err)
		return
	}

	for _, client := range clients {
		client.mu.Lock()
		select {
		case client.send <- jsonMessage:
			log.Printf("Message sent to user %s", userID)
		default:
			// Channel is full, unregister the client
			client.mu.Unlock()
			h.unregister <- client
			continue
		}
		client.mu.Unlock()
	}
}

// broadcastMessage sends a message to all connected clients
func (h *Hub) broadcastMessage(message Message) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.mu.RLock()
	clientCount := len(h.clients)
	h.mu.RUnlock()

	log.Printf("Broadcasting message to %d clients", clientCount)

	h.mu.RLock()
	for client := range h.clients {
		client.mu.Lock()
		select {
		case client.send <- jsonMessage:
		default:
			// Channel is full, unregister the client
			client.mu.Unlock()
			h.unregister <- client
			continue
		}
		client.mu.Unlock()
	}
	h.mu.RUnlock()
}

// GetStats returns statistics about the hub
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"total_clients": len(h.clients),
		"total_users":   len(h.userClients),
		"connections_by_user": func() map[string]int {
			stats := make(map[string]int)
			for userID, clients := range h.userClients {
				stats[userID] = len(clients)
			}
			return stats
		}(),
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", c.userID, err)
			}
			break
		}

		// Handle incoming messages (e.g., ping, acknowledgments)
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message from user %s: %v", c.userID, err)
			continue
		}

		// Handle different message types
		switch msg.Type {
		case MessageTypePing:
			// Respond with pong
			pongMsg := Message{
				Type:      MessageTypePong,
				Timestamp: time.Now(),
			}
			if jsonMsg, err := json.Marshal(pongMsg); err == nil {
				c.mu.Lock()
				select {
				case c.send <- jsonMsg:
				default:
				}
				c.mu.Unlock()
			}
		default:
			// Log other message types for debugging
			log.Printf("Received message type %s from user %s", msg.Type, c.userID)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()

	log.Printf("WebSocket connection established for user %s", userID)
}
