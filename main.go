package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Get command-line arguments (defined in main.go)
	arguments := os.Args[1:]
	var serverPort string
	if len(arguments) == 0 {
		serverPort = "8989"
	} else if len(arguments) == 1 {
		serverPort = arguments[0]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Create a ChatServer instance (NewChatServer is defined in server.go)
	chatServer, err := NewChatServer(serverPort)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	fmt.Printf("Listening on port %s\n", serverPort)
	// Start the server (Start method is defined in server.go)
	chatServer.Start()
}
