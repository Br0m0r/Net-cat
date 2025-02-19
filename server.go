package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// Maximum number of simultaneous client connections.
const maxClients = 3

// ChatServer holds all the server information and data required to manage the chat.
type ChatServer struct {
	// port stores the port number on which the server listens for incoming connections.
	port string
	// listener represents the TCP listener that accepts new client connections.
	listener net.Listener
	// clients is a map that keeps track of all connected clients.
	// The key is the client's connection (net.Conn) and the value is a pointer to the corresponding Client struct.
	clients map[net.Conn]*Client
	// history holds the entire chat history. Each message is appended to this slice.
	history []string
	// mu is a mutex used to ensure that access to shared resources like the clients map and history slice is thread-safe.
	mu sync.Mutex
	// broadcastCh is a channel used for message broadcasting.
	// When a message is sent to this channel, a dedicated goroutine forwards it to all connected clients.
	broadcastCh chan string
}

// NewChatServer creates and returns a new ChatServer listening on the specified port.
func NewChatServer(port string) (*ChatServer, error) {
	// Start listening on the specified TCP port.
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// Return an error if the server cannot start listening.
		return nil, err
	}

	// Initialize a new ChatServer with the required fields.
	server := &ChatServer{
		port:        port,                       // Set the listening port.
		listener:    ln,                         // Assign the TCP listener.
		clients:     make(map[net.Conn]*Client), // Initialize the clients map.
		history:     make([]string, 0),          // Initialize the chat history slice.
		broadcastCh: make(chan string),          // Create the broadcast channel.
	}

	return server, nil
}

// Start begins accepting client connections and handling broadcasts.
func (s *ChatServer) Start() {
	// Start a goroutine that listens for messages on the broadcast channel
	// and sends them to all connected clients.
	go s.handleBroadcast()

	// Continuously accept new client connections.
	for {
		// Wait for a new connection.
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		s.mu.Lock()
		if len(s.clients) >= maxClients { // ❌ This might not be accurate under heavy load
			s.mu.Unlock()
			conn.Write([]byte("Server is full. Try again later.\n"))
			conn.Close()
			continue
		}
		s.mu.Unlock()

		go s.handleClient(conn) // ❌ Client is not added to the map immediately

	}
}

// handleBroadcast listens on the broadcast channel and sends messages to all connected clients.
func (s *ChatServer) handleBroadcast() {
	// Iterate over messages received on the broadcast channel.
	for msg := range s.broadcastCh {
		// Lock the mutex to safely access shared resources.
		s.mu.Lock()
		// Append the new message to the chat history.
		s.history = append(s.history, msg)
		// Loop through each connected client.
		for conn, client := range s.clients {
			// Send the message to the client.
			_, err := fmt.Fprintf(conn, "%s\n", msg)
			// Log an error if there's an issue sending the message.
			if err != nil {
				log.Printf("Error sending message to %s: %v", client.name, err)
			}
		}
		// Unlock the mutex after processing the message.
		s.mu.Unlock()
	}
}

// broadcast sends a message to all clients by placing it onto the broadcast channel.
func (s *ChatServer) broadcast(msg string) {
	s.broadcastCh <- msg
}
