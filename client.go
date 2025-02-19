package main

import (
	"bufio"
	"fmt"
	"io"
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
	// Ensure the connection is closed when the function exits.
	defer conn.Close()

	// Send a welcome message (with ASCII logo and colors) to the client.
	s.sendWelcome(conn)

	// Prompt the client to enter their name.
	name, err := s.promptName(conn)
	if err != nil {
		log.Println("Error reading name:", err)
		return
	}

	// Send the stored chat history to the new client.
	s.sendHistory(conn)

	// Notify all connected clients that a new client has joined.
	joinMsg := FormatSystemMessage(fmt.Sprintf("%s has joined our chat...", name))
	s.broadcast(joinMsg)

	// Create a scanner to read messages from the client's connection.
	scanner := bufio.NewScanner(conn)
	// Continuously scan for new messages.
	for scanner.Scan() {
		msg := scanner.Text()
		// Ignore empty messages.
		if strings.TrimSpace(msg) == "" {
			continue
		}
		// Format the message with a timestamp and the client's name.
		formattedMsg := FormatChatMessage(
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			msg,
		)
		// Broadcast the formatted message to all clients.
		s.broadcast(formattedMsg)
	}
	// Log any errors encountered while scanning messages (excluding EOF).
	if err := scanner.Err(); err != nil && err != io.EOF {
		log.Println("Error reading from client:", err)
	}

	// When the client disconnects, remove them from the clients map.
	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
	// Broadcast a message notifying that this client has left.
	leaveMsg := FormatSystemMessage(fmt.Sprintf("%s has left our chat...", name))
	s.broadcast(leaveMsg)
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
