package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// Client struct is defined in client.go and holds details for each client.
type Client struct {
	conn net.Conn
	name string
}

// handleClient manages the client's connection, name registration, history delivery, and message handling. (Defined in client.go)
func (server *ChatServer) handleClient(connection net.Conn) {
	defer connection.Close()

	server.mutex.Lock()
	activeClients := len(server.clients)
	server.mutex.Unlock()

	log.Printf("New client connected. [%d users online]", activeClients)

	// Send the welcome message (using FormatWelcomeMessage defined in formatter.go)
	welcomeMessage := FormatWelcomeMessage()
	connection.Write([]byte(welcomeMessage))

	nameChan := make(chan string, 1)
	disconnectChan := make(chan struct{})

	// Timeout mechanism for name entry.
	go func() {
		select {
		case <-time.After(30 * time.Second):
			server.mutex.Lock()
			if client, exists := server.clients[connection]; exists && client.name == "" {
				connection.Write([]byte(FormatSystemMessage("\nTimeout: You took too long to enter a name. Disconnecting...\n")))
				delete(server.clients, connection)
			}
			server.mutex.Unlock()
			connection.Close()
		case enteredName := <-nameChan:
			server.mutex.Lock()
			if client, exists := server.clients[connection]; exists {
				client.name = enteredName
				joinMessage := FormatSystemMessage(fmt.Sprintf("%s has joined our chat! [%d users online]", enteredName, len(server.clients)))
				server.broadcast(joinMessage)
				// Send chat history to the new client.
				if len(server.history) > 0 {
					connection.Write([]byte(FormatSystemMessage("\n--- Chat History ---\n")))
					for _, pastMessage := range server.history {
						connection.Write([]byte(pastMessage + "\n"))
					}
					connection.Write([]byte(FormatSystemMessage("\n--- End of History ---\n\n")))
				}
			}
			server.mutex.Unlock()
		case <-disconnectChan:
			return
		}
	}()

	scanner := bufio.NewScanner(connection)
	var clientName string
	connection.Write([]byte(FormatSystemMessage("[ENTER YOUR NAME]: ")))
	for scanner.Scan() {
		clientName = strings.TrimSpace(scanner.Text())

		server.mutex.Lock()
		nameTaken := false
		for _, connectedClient := range server.clients {
			if connectedClient.name == clientName {
				nameTaken = true
				break
			}
		}
		server.mutex.Unlock()

		if clientName == "" {
			connection.Write([]byte(FormatSystemMessage("Name cannot be empty. Please enter your name: ")))
		} else if nameTaken {
			connection.Write([]byte(FormatSystemMessage("This username is already taken. Please choose another: ")))
		} else {
			break
		}
	}

	if clientName == "" {
		log.Println("Client disconnected before entering a name")
		close(disconnectChan)
		server.mutex.Lock()
		delete(server.clients, connection)
		server.mutex.Unlock()
		return
	}

	nameChan <- clientName

	// Listen for messages from this client.
	for scanner.Scan() {
		messageText := strings.TrimSpace(scanner.Text())
		if messageText == "" {
			continue
		}

		// FormatChatMessage is defined in formatter.go.
		formattedMessage := FormatChatMessage(
			time.Now().Format("2006-01-02 15:04:05"),
			clientName,
			messageText,
		)
		server.broadcast(formattedMessage)
	}

	server.mutex.Lock()
	delete(server.clients, connection)
	activeClients = len(server.clients)
	server.mutex.Unlock()

	log.Printf("Client disconnected: %s [%d users online]", clientName, activeClients)
	if clientName != "" {
		leaveMessage := FormatSystemMessage(fmt.Sprintf("%s has left our chat... [%d users online]", clientName, activeClients))
		server.broadcast(leaveMessage)
	}
}
