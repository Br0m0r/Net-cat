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

// handleClient manages the client's connection, name registration, history delivery, and message handling.
func (server *ChatServer) handleClient(connection net.Conn) {
	defer connection.Close()

	// 1. Locking server to get the number of active clients.
	//------------------------------------------------------
	server.mutex.Lock()
	activeClients := len(server.clients)
	server.mutex.Unlock()
	//------------------------------------------------------

	log.Printf("New client connected. [%d users online]", activeClients)

	// 2. Send welcome message.
	//------------------------------------------------------
	welcomeMessage := FormatWelcomeMessage()
	connection.Write([]byte(welcomeMessage))
	//------------------------------------------------------

	// 3. Create channels for name input and disconnection handling.
	//------------------------------------------------------
	nameChan := make(chan string, 1)      // Stores single string value(buffered)
	disconnectChan := make(chan struct{}) // Signaling events (unbuffered)
	//------------------------------------------------------

	// 4. Timeout mechanism for name entry (runs in a separate goroutine).
	//------------------------------------------------------
	go func() {
		select {
		case <-time.After(20 * time.Second): // Timeout after 30 seconds
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
		case <-disconnectChan: //B.when triggered the timeout go func exits
			return
		}
	}()
	//------------------------------------------------------

	// 5. Prompt client for their name.
	//------------------------------------------------------
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
	//------------------------------------------------------

	// 6. Handle client disconnection if no name was entered.
	//------------------------------------------------------
	if clientName == "" {
		log.Println("Client disconnected before entering a name")
		close(disconnectChan) // A. this triggers  case <-disconnectChan
		server.mutex.Lock()
		delete(server.clients, connection)
		server.mutex.Unlock()
		return //C. this exits the handleClient
	}
	//------------------------------------------------------

	// 7. Store the entered name into nameChan.
	//------------------------------------------------------
	nameChan <- clientName
	//------------------------------------------------------

	// 8. Listen for messages from the client.
	//------------------------------------------------------
	for scanner.Scan() {
		messageText := strings.TrimSpace(scanner.Text())
		if messageText == "" {
			continue
		}

		// Format message with timestamp and broadcast it.
		formattedMessage := FormatChatMessage(
			time.Now().Format("2006-01-02 15:04:05"),
			clientName,
			messageText,
		)
		server.broadcast(formattedMessage)
	}
	//------------------------------------------------------

	// 9. Handle client disconnection.
	//------------------------------------------------------
	server.mutex.Lock()
	delete(server.clients, connection)
	activeClients = len(server.clients)
	server.mutex.Unlock()

	log.Printf("Client disconnected: %s [%d users online]", clientName, activeClients)
	if clientName != "" {
		leaveMessage := FormatSystemMessage(fmt.Sprintf("%s has left our chat... [%d users online]", clientName, activeClients))
		server.broadcast(leaveMessage)
	}
	//------------------------------------------------------
}
