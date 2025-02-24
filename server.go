package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const maxClients = 10

// ChatServer is defined in server.go and manages all chat server operations.
type ChatServer struct {
	port        string
	listener    net.Listener
	clients     map[net.Conn]*Client
	history     []string
	mutex       sync.Mutex
	broadcastCh chan string
}

// NewChatServer creates and returns a ChatServer instance. (Defined in server.go)
func NewChatServer(port string) (*ChatServer, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	return &ChatServer{
		port:        port,
		listener:    listener,
		clients:     make(map[net.Conn]*Client),
		history:     make([]string, 0),
		broadcastCh: make(chan string),
	}, nil
}

// Start accepts client connections and launches the broadcast handler. (Defined in server.go)
func (server *ChatServer) Start() {
	go server.handleBroadcast() // handleBroadcast is defined below in server.go
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		server.mutex.Lock()
		if len(server.clients) >= maxClients {
			connection.Write([]byte("Server is full. Try again later.\n"))
			connection.Close()
			server.mutex.Unlock()
			continue
		}

		server.clients[connection] = &Client{conn: connection, name: ""}
		server.mutex.Unlock()

		// handleClient is defined in client.go but belongs to ChatServer.
		go server.handleClient(connection)
	}
}

// handleBroadcast sends messages from the broadcast channel to all connected clients. (Defined in server.go)
func (server *ChatServer) handleBroadcast() {
	for message := range server.broadcastCh {
		server.mutex.Lock()
		server.history = append(server.history, message)
		for connection, connectedClient := range server.clients {
			if connectedClient.name == "" {
				continue
			}
			if _, err := fmt.Fprintf(connection, "%s\n", message); err != nil {
				log.Printf("Error sending message to %s: %v", connectedClient.name, err)
			}
		}
		server.mutex.Unlock()
	}
}

// broadcast sends a message to all clients via the broadcast channel. (Defined in server.go)
func (server *ChatServer) broadcast(message string) {
	server.broadcastCh <- message
}
