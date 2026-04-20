package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const maxClients = 10

type ChatServer struct {
    port        string
    listener    net.Listener
    clients     map[net.Conn]*Client
    history     []string
    mu          sync.Mutex
    broadcastCh chan string
}

func NewChatServer(port string) (*ChatServer, error) {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return nil, err
    }

    server := &ChatServer{
        port:        port,
        listener:    ln,
        clients:     make(map[net.Conn]*Client),
        history:     make([]string, 0),
        broadcastCh: make(chan string),
    }

    return server, nil
}

func (s *ChatServer) Start() {
    go s.handleBroadcast()

    for {
        conn, err := s.listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }

        s.mu.Lock()
        if len(s.clients) >= maxClients {
            conn.Write([]byte("Server is full. Try again later.\n"))
            conn.Close()
            s.mu.Unlock()
            continue
        }

        s.clients[conn] = &Client{conn: conn, name: ""}
        s.mu.Unlock()

        go s.handleClient(conn)
    }
}

func (s *ChatServer) handleBroadcast() {
    for msg := range s.broadcastCh {
        s.mu.Lock()
        s.history = append(s.history, msg)
        for conn, client := range s.clients {
            if _, err := fmt.Fprintf(conn, "%s\n", msg); err != nil {
                log.Printf("Error sending message to %s: %v", client.name, err)
            }
        }
        s.mu.Unlock()
    }
}

func (s *ChatServer) broadcast(msg string) {
    s.broadcastCh <- msg
}
