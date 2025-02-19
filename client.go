package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// Client holds individual client details.
type Client struct {
	// conn represents the client's TCP connection.
	conn net.Conn
	// name holds the client's chosen username.
	name string
}

// handleClient manages a new client by handling the welcome sequence,
// registering the client’s name, sending chat history, and continuously
// reading and broadcasting the client's messages.
func (s *ChatServer) handleClient(conn net.Conn) {
	defer conn.Close()

	// Track the new client immediately in the map with an empty name
	s.mu.Lock()
	s.clients[conn] = &Client{conn: conn, name: ""}
	numUsers := len(s.clients) // Get current user count
	s.mu.Unlock()

	// Notify server log
	log.Printf("New client connected. [%d users online]", numUsers)

	// Send the welcome message
	s.sendWelcome(conn)

	// Channel to track when a name is entered
	nameEntered := make(chan string)

	// Start a goroutine for the timeout mechanism
	go func() {
		select {
		case <-time.After(30 * time.Second): // ⏳ Timeout after 30 seconds
			conn.Write([]byte("\nTimeout: You took too long to enter a name. Disconnecting...\n"))
			conn.Close()
		case name := <-nameEntered: // ✅ User entered a name before timeout
			client := s.clients[conn]
			client.name = name
			joinMsg := FormatSystemMessage(fmt.Sprintf("%s has joined our chat! [%d users online]", name, numUsers))
			s.broadcast(joinMsg)
		}
	}()

	// Prompt for a name (runs in the main goroutine)
	name, err := s.promptName(conn)
	if err != nil {
		log.Println("Error reading name:", err)
		return
	}

	// Send name to cancel timeout and continue
	nameEntered <- name

	// Listen for messages from this client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.TrimSpace(msg) == "" {
			continue // Ignore empty messages
		}

		formattedMsg := FormatChatMessage(
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			msg,
		)
		s.broadcast(formattedMsg)
	}

	// Handle client disconnection
	s.mu.Lock()
	delete(s.clients, conn)   // Remove from client list
	numUsers = len(s.clients) // Update user count after removal
	s.mu.Unlock()

	// Log the disconnect
	log.Printf("Client disconnected: %s [%d users online]", name, numUsers)

	// ✅ Only broadcast if the client had a name
	if name != "" {
		leaveMsg := FormatSystemMessage(fmt.Sprintf("%s has left our chat... [%d users online]", name, numUsers))
		s.broadcast(leaveMsg)
	}
}

// sendWelcome writes the welcome message and ASCII logo to the new connection.
func (s *ChatServer) sendWelcome(conn net.Conn) {
	// Use the formatter to create a colored welcome message.
	welcome := FormatWelcomeMessage()
	conn.Write([]byte(welcome))
}

// promptName asks the client for their name until a non-empty value is received.
func (s *ChatServer) promptName(conn net.Conn) (string, error) {
	// Prompt the client to enter their name.
	_, err := conn.Write([]byte("[ENTER YOUR NAME]: "))
	if err != nil {
		return "", err
	}

	// Create a scanner to read input from the client's connection.
	scanner := bufio.NewScanner(conn)
	// Continuously scan for input.
	for scanner.Scan() {
		// Trim any whitespace from the input to get the name.
		name := strings.TrimSpace(scanner.Text())
		// If a non-empty name is provided, return it.
		if name != "" {
			return name, nil
		}
		// If the name is empty, prompt the client again.
		conn.Write([]byte("Name cannot be empty. Please enter your name: "))
	}
	// Return an error if the scanner encountered an issue.
	return "", scanner.Err()
}

// sendHistory writes the chat history to the newly connected client.
func (s *ChatServer) sendHistory(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, msg := range s.history {
		fmt.Fprintf(conn, "%s\n", msg)
	}
}
