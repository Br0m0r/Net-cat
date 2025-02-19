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

	s.mu.Lock()
	s.clients[conn] = &Client{conn: conn, name: ""}
	numUsers := len(s.clients)
	s.mu.Unlock()

	log.Printf("New client connected. [%d users online]", numUsers)

	// ✅ Send the formatted welcome message
	welcomeMessage := FormatWelcomeMessage()
	conn.Write([]byte(welcomeMessage))

	// Channels to track user actions
	nameEntered := make(chan string, 1)
	disconnected := make(chan struct{})

	// Start a goroutine for the timeout mechanism
	go func() {
		select {
		case <-time.After(30 * time.Second): // ⏳ Timeout after 30 seconds
			s.mu.Lock()
			if client, exists := s.clients[conn]; exists && client.name == "" {
				conn.Write([]byte(FormatSystemMessage("\nTimeout: You took too long to enter a name. Disconnecting...\n")))
				delete(s.clients, conn)
			}
			s.mu.Unlock()
			conn.Close()
		case name := <-nameEntered:
			s.mu.Lock()
			if client, exists := s.clients[conn]; exists {
				client.name = name
				joinMsg := FormatSystemMessage(fmt.Sprintf("%s has joined our chat! [%d users online]", name, len(s.clients)))
				s.broadcast(joinMsg)
			}
			s.mu.Unlock()
		case <-disconnected:
			return
		}
	}()

	// ✅ Prompt for a name (loop until a valid and unique name is entered)
	scanner := bufio.NewScanner(conn)
	var name string
	conn.Write([]byte(FormatSystemMessage("[ENTER YOUR NAME]: ")))
	for scanner.Scan() {
		name = strings.TrimSpace(scanner.Text())

		// ✅ Check if the name is already taken
		s.mu.Lock()
		isTaken := false
		for _, client := range s.clients {
			if client.name == name {
				isTaken = true
				break
			}
		}
		s.mu.Unlock()

		if name == "" {
			conn.Write([]byte(FormatSystemMessage("Name cannot be empty. Please enter your name: ")))
		} else if isTaken {
			conn.Write([]byte(FormatSystemMessage("This username is already taken. Please choose another: ")))
		} else {
			break // ✅ Valid and unique name entered
		}
	}

	if name == "" { // ✅ Handle case where scanner fails (EOF)
		log.Println("Client disconnected before entering a name")
		close(disconnected)
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		return
	}

	nameEntered <- name // ✅ Send name to cancel the timeout

	// ✅ Listen for messages from this client
	for scanner.Scan() {
		msg := strings.TrimSpace(scanner.Text())
		if msg == "" {
			continue
		}

		formattedMsg := FormatChatMessage(
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			msg,
		)

		s.broadcast(formattedMsg)
	}

	// ✅ Handle client disconnection
	s.mu.Lock()
	delete(s.clients, conn)
	numUsers = len(s.clients)
	s.mu.Unlock()

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
