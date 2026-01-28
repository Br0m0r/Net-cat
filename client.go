package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    "time"
)

type Client struct {
    conn net.Conn
    name string
}

// handleClient manages one client connection from handshake to disconnect.
func (s *ChatServer) handleClient(conn net.Conn) {
    defer conn.Close()

    s.mu.Lock()
    numUsers := len(s.clients)
    s.mu.Unlock()

    log.Printf("New client connected. [%d users online]", numUsers)

    conn.Write([]byte(FormatWelcomeMessage()))

    nameEntered := make(chan string, 1)
    disconnected := make(chan struct{})

    // Enforce a name entry timeout while allowing disconnects.
    go func() {
        select {
        case <-time.After(30 * time.Second):
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
            return
        }
    }()

    scanner := bufio.NewScanner(conn)
    var name string
    conn.Write([]byte(FormatSystemMessage("[ENTER YOUR NAME]: ")))

    // Name input loop: require non-empty, unique username.
    for scanner.Scan() {
        name = strings.TrimSpace(scanner.Text())

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
            break
        }
    }

    if name == "" {
        log.Println("Client disconnected before entering a name")
        close(disconnected)
        s.mu.Lock()
        delete(s.clients, conn)
        s.mu.Unlock()
        return
    }

    nameEntered <- name

    // Message loop: ignore empty lines and broadcast formatted messages.
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

    s.mu.Lock()
    delete(s.clients, conn)
    numUsers = len(s.clients)
    s.mu.Unlock()

    log.Printf("Client disconnected: %s [%d users online]", name, numUsers)

    if name != "" {
        leaveMsg := FormatSystemMessage(fmt.Sprintf("%s has left our chat... [%d users online]", name, numUsers))
        s.broadcast(leaveMsg)
    }
}

func (s *ChatServer) sendWelcome(conn net.Conn) {
    conn.Write([]byte(FormatWelcomeMessage()))
}

func (s *ChatServer) promptName(conn net.Conn) (string, error) {
    if _, err := conn.Write([]byte("[ENTER YOUR NAME]: ")); err != nil {
        return "", err
    }

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        name := strings.TrimSpace(scanner.Text())
        if name != "" {
            return name, nil
        }
        conn.Write([]byte("Name cannot be empty. Please enter your name: "))
    }

    return "", scanner.Err()
}

func (s *ChatServer) sendHistory(conn net.Conn) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for _, msg := range s.history {
        fmt.Fprintf(conn, "%s\n", msg)
    }
}
