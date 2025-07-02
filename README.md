Net-Cat 

Net-Cat is a lightweight TCP-based group chat application written in Go. It allows multiple clients to connect, communicate in real-time, and see the conversation history.
Features

    Multiple Clients: Supports up to 10 simultaneous connections.
    Usernames: Clients must choose a unique username.
    Message Broadcasting: All messages are broadcast to every connected client with a timestamp.
    Notifications: Clients are notified when users join or leave the chat.
    Chat History: New clients receive the chat history upon connection.
    Default Port: Listens on port 8989 by default.

Installation

    Clone the Repository:
    git clone https://platform.zone01.gr/git/mfoteino/net-cat


Build the Project (Optional):

    go build -o TCPchat .

Usage
Running the Server

    Default Port (8989):
    go run .
    Custom Port (e.g., 2525):
    go run . 2525

Invalid Usage: If you provide more than one argument, the program will display:

    [USAGE]: ./TCPChat $port

Connecting as a Client

    You can use any TCP client (such as nc):
    nc localhost 8989

When connected, you'll see a welcome message and be prompted for your username.
    
                                     File Structure

    main.go: Entry point; handles command-line arguments and starts the server.
    server.go: Manages client connections, message broadcasting, and chat history.
    client.go: Handles user interactions and message processing.
    formatter.go: Provides formatted output with ANSI color codes.
