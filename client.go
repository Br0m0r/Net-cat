package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// ✅ Client struct represents an individual user in the chat system.
type Client struct {
	// ✅ `conn` stores the client's TCP connection to the server.
	conn net.Conn
	// ✅ `name` holds the username chosen by the client.
	name string
}

// ✅ handleClient() - Manages a single client session.
func (s *ChatServer) handleClient(conn net.Conn) {
	// ✅ Ensure the connection is closed when the function exits.
	defer conn.Close()

	// ✅ Lock mutex to safely access shared client data.
	s.mu.Lock()
	numUsers := len(s.clients) // ✅ Get the current number of connected clients.
	s.mu.Unlock()

	log.Printf("New client connected. [%d users online]", numUsers)

	// ✅ Send a formatted welcome message (includes ASCII art).
	welcomeMessage := FormatWelcomeMessage() // Calls formatter function.
	conn.Write([]byte(welcomeMessage))       // ✅ Sends the welcome message to the client.

	// ✅ Create channels to track name input and disconnection status.
	nameEntered := make(chan string, 1) // Stores the name when entered.
	disconnected := make(chan struct{}) // Detects if the client disconnects.

	// ✅ Start a separate goroutine to enforce a timeout.
	go func() {
		select {
		case <-time.After(30 * time.Second): // ⏳ If no name entered in 30 seconds.
			s.mu.Lock()
			if client, exists := s.clients[conn]; exists && client.name == "" {
				conn.Write([]byte(FormatSystemMessage("\nTimeout: You took too long to enter a name. Disconnecting...\n")))
				delete(s.clients, conn) // ✅ Remove the client from the active list.
			}
			s.mu.Unlock()
			conn.Close() // ✅ Disconnect the client.
		case name := <-nameEntered:
			// ✅ If name entered before timeout, store it.
			s.mu.Lock()
			if client, exists := s.clients[conn]; exists {
				client.name = name // ✅ Assign name to the client.
				joinMsg := FormatSystemMessage(fmt.Sprintf("%s has joined our chat! [%d users online]", name, len(s.clients)))
				s.broadcast(joinMsg) // ✅ Notify all other users.

				// ✅ Send past chat history to the new user.
				if len(s.history) > 0 {
					conn.Write([]byte(FormatSystemMessage("\n--- Chat History ---\n")))
					for _, msg := range s.history {
						conn.Write([]byte(msg + "\n"))
					}
					conn.Write([]byte(FormatSystemMessage("\n--- End of History ---\n\n")))
				}
			}
			s.mu.Unlock()
		case <-disconnected:
			// ✅ If the client disconnected before entering a name, exit.
			return
		}
	}()

	// ✅ Prompt the client to enter their username.
	scanner := bufio.NewScanner(conn)
	var name string
	conn.Write([]byte(FormatSystemMessage("[ENTER YOUR NAME]: ")))

	// 🔹 Loop until the client provides a valid, unique username.
	for scanner.Scan() {
		name = strings.TrimSpace(scanner.Text()) // ✅ Read input & remove spaces.

		// ✅ Ensure the username is unique.
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
			// 🔹 If name is empty, ask again.
			conn.Write([]byte(FormatSystemMessage("Name cannot be empty. Please enter your name: ")))
		} else if isTaken {
			// 🔹 If name is taken, ask again.
			conn.Write([]byte(FormatSystemMessage("This username is already taken. Please choose another: ")))
		} else {
			break // ✅ Name is valid, exit loop.
		}
	}

	// ✅ Handle case where client disconnects before entering a name.
	if name == "" {
		log.Println("Client disconnected before entering a name")
		close(disconnected)
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		return
	}

	// ✅ Send the name to cancel the timeout.
	nameEntered <- name

	// ✅ Start listening for messages from this client.
	for scanner.Scan() {
		msg := strings.TrimSpace(scanner.Text()) // ✅ Read the message.
		if msg == "" {
			continue // 🔹 Ignore empty messages.
		}

		// ✅ Format the message with a timestamp and username.
		formattedMsg := FormatChatMessage(
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			msg,
		)

		// ✅ Send the formatted message to all clients.
		s.broadcast(formattedMsg)
	}

	// ✅ Handle client disconnection.
	s.mu.Lock()
	delete(s.clients, conn) // 🔹 Remove from active clients.
	numUsers = len(s.clients)
	s.mu.Unlock()

	log.Printf("Client disconnected: %s [%d users online]", name, numUsers)

	// ✅ Notify others when a client leaves.
	if name != "" {
		leaveMsg := FormatSystemMessage(fmt.Sprintf("%s has left our chat... [%d users online]", name, numUsers))
		s.broadcast(leaveMsg)
	}
}

// ✅ sendWelcome() - Sends the welcome message with ASCII logo.
func (s *ChatServer) sendWelcome(conn net.Conn) {
	// 🔹 Calls the formatter function to generate the welcome message.
	welcome := FormatWelcomeMessage()
	conn.Write([]byte(welcome)) // 🔹 Sends the welcome message.
}

// ✅ promptName() - Asks the client for a valid name.
func (s *ChatServer) promptName(conn net.Conn) (string, error) {
	// 🔹 Ask the user to enter their name.
	_, err := conn.Write([]byte("[ENTER YOUR NAME]: "))
	if err != nil {
		return "", err // 🔹 If connection fails, return an error.
	}

	// 🔹 Read the client's input.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text()) // ✅ Remove spaces.

		// ✅ If valid, return name.
		if name != "" {
			return name, nil
		}

		// 🔹 If name is empty, ask again.
		conn.Write([]byte("Name cannot be empty. Please enter your name: "))
	}

	// 🔹 Return an error if scanner fails.
	return "", scanner.Err()
}

// ✅ sendHistory() - Sends stored chat history to a newly connected client.
func (s *ChatServer) sendHistory(conn net.Conn) {
	s.mu.Lock()         // 🔹 Lock access to `s.history`.
	defer s.mu.Unlock() // ✅ Ensure mutex is released after function runs.
	for _, msg := range s.history {
		fmt.Fprintf(conn, "%s\n", msg) // 🔹 Send each past message.
	}
}
