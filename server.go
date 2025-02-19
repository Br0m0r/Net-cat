package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// ✅ Maximum number of clients allowed to connect at the same time.
const maxClients = 3

// ✅ ChatServer struct manages all server operations.
type ChatServer struct {
	// ✅ Port number where the server listens for connections.
	port string

	// ✅ TCP listener that waits for incoming connections.
	listener net.Listener

	// ✅ Map to store all active clients.
	// Key: net.Conn (client connection), Value: *Client (client struct with details).
	clients map[net.Conn]*Client

	// ✅ Stores chat history.
	// Each message is appended so new clients can receive past messages.
	history []string

	// ✅ Mutex to ensure safe access to shared resources (`clients` map and `history` slice).
	mu sync.Mutex

	// ✅ Channel used to send messages to all clients.
	// When a message is sent here, it gets broadcasted to every connected client.
	broadcastCh chan string
}

// ✅ NewChatServer initializes and returns a new ChatServer instance.
func NewChatServer(port string) (*ChatServer, error) {
	// 🔹 Attempt to start listening on the given port.
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// 🔹 If the server fails to start, return an error.
		return nil, err
	}

	// 🔹 Create a new ChatServer instance with initialized values.
	server := &ChatServer{
		port:        port,                       // ✅ Store the listening port.
		listener:    ln,                         // ✅ Assign the TCP listener.
		clients:     make(map[net.Conn]*Client), // ✅ Initialize the clients map (empty at start).
		history:     make([]string, 0),          // ✅ Initialize an empty chat history slice.
		broadcastCh: make(chan string),          // ✅ Create a channel for broadcasting messages.
	}

	return server, nil // 🔹 Return the initialized ChatServer instance.
}

// ✅ Start() - Begins accepting client connections and manages the chat system.
func (s *ChatServer) Start() {
	// 🔹 Start a separate goroutine for handling message broadcasting.
	go s.handleBroadcast()

	// 🔹 Infinite loop to continuously accept new clients.
	for {
		// 🔹 Wait for a new connection from a client.
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue // 🔹 Skip this iteration and wait for another connection.
		}

		// 🔹 Ensure thread-safe access to the `clients` map.
		s.mu.Lock()
		if len(s.clients) >= maxClients { // ✅ Check if max client limit is reached.
			conn.Write([]byte("Server is full. Try again later.\n")) // 🔹 Inform client.
			conn.Close()                                             // 🔹 Close the connection.
			s.mu.Unlock()                                            // 🔹 Release the mutex.
			return                                                   // 🔹 Exit function to stop further processing.
		}

		// ✅ Register the new client in the `clients` map.
		s.clients[conn] = &Client{conn: conn, name: ""} // 🔹 Temporarily store client with an empty name.
		s.mu.Unlock()                                   // 🔹 Release the mutex after modifying shared data.

		// ✅ Handle the client in a separate goroutine.
		go s.handleClient(conn)
	}
}

// ✅ handleBroadcast() - Listens for messages on the broadcast channel and distributes them.
func (s *ChatServer) handleBroadcast() {
	// 🔹 Loop indefinitely, processing messages sent to `s.broadcastCh`.
	for msg := range s.broadcastCh {
		// 🔹 Lock mutex to prevent race conditions while modifying shared resources.
		s.mu.Lock()
		s.history = append(s.history, msg) // ✅ Save message in history.

		// ✅ Send the message to each connected client.
		for conn, client := range s.clients {
			_, err := fmt.Fprintf(conn, "%s\n", msg) // 🔹 Send message.
			if err != nil {                          // 🔹 Log any errors while sending.
				log.Printf("Error sending message to %s: %v", client.name, err)
			}
		}
		s.mu.Unlock() // 🔹 Unlock mutex after processing message.
	}
}

// ✅ broadcast() - Sends a message to all connected clients.
func (s *ChatServer) broadcast(msg string) {
	s.broadcastCh <- msg // 🔹 Add message to the broadcast channel.
}
