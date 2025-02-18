package main

import (
	"fmt"
	"log"
	"os"
)

// main is the entry point of the application.
// It parses command-line arguments, initializes the ChatServer,
// and starts the server.
func main() {
	// Parse command-line arguments.
	// If no port is provided, use default port 8989.
	// If one argument is provided, use that as the port.
	// Otherwise, display a usage message.
	args := os.Args[1:]
	var port string
	if len(args) == 0 {
		port = "8989"
	} else if len(args) == 1 {
		port = args[0]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Create a new ChatServer on the specified port.
	server, err := NewChatServer(port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	fmt.Printf("Listening on the port :%s\n", port)
	// Start the server.
	server.Start()
}
