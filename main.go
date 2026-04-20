package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
    args := os.Args[1:]
    var port string
    switch len(args) {
    case 0:
        port = "8989"
    case 1:
        port = args[0]
    default:
        fmt.Println("[USAGE]: ./TCPChat $port")
        return
    }

    server, err := NewChatServer(port)
    if err != nil {
        log.Fatal("Error starting server:", err)
    }

    fmt.Printf("Listening on the port :%s\n", port)
    server.Start()
}
