# net-cat

net-cat is a Go-based group chat application that mimics key features of the classic NetCat (`nc`) utility. It uses a server-client architecture to allow multiple clients to connect to a central server and exchange messages in real time.

## Features

- **TCP Server and Clients:** Supports multiple TCP connections.
- **Default Port Handling:** Listens on port `8989` if no port is specified.
- **Usage Validation:** Displays a usage message if extra parameters are provided.
- **Client Name Requirement:** Prompts each client to enter a unique, non-empty username.
- **Chat History:** New clients receive the entire chat history upon joining.
- **Message Broadcasting:** Sends messages with a timestamp and sender's name to all connected clients.
- **Join/Leave Notifications:** Notifies all clients when someone joins or leaves.
- **Concurrency:** Utilizes goroutines, channels, and mutexes for concurrent operations.


## Installation

1. **Clone the repository:**

   ```bash
   git clone https://platform.zone01.gr/git/mfoteino/net-cat
   cd net-cat

    Build the project (optional):

    You can build the binary with:

    go build -o net-cat .

Usage
Running the Server

    Default Port (8989):

    Run the server with:

go run .

Custom Port (e.g., 2525):

Run the server with:

go run . 2525

Invalid Usage:
If extra parameters are provided (e.g., go run . 2525 localhost), the server displays a usage message:

    [USAGE]: ./TCPChat $port

Connecting as a Client

Use NetCat (nc) or any TCP client to connect to the server. For example, if the server is running on port 8989:

nc localhost 8989

Upon connection, you will see a welcome message with an ASCII art logo and a prompt to enter your username.
Testing

To verify key functionalities, follow these steps:

    Usage Test:
    Run the server with extra parameters to confirm that the usage message is triggered:

go run . 2525 localhost

Expected Output:

[USAGE]: ./TCPChat $port

Multi-Client Test:

    Open one terminal and run the server on a custom port:

go run . 2525

Open two or more separate terminals and connect using:

        nc localhost 2525

    Verify that:
        Each client receives the welcome message and is prompted for a username.
        When a client enters a valid name, all connected clients are notified of the new join.
        Messages are broadcast to all clients with proper timestamps and sender names.
        A new client receives the complete chat history upon connection.
        When a client disconnects, the remaining clients receive a notification.

    Multi-Computer Test:
    Ensure that clients on different computers can connect and exchange messages.

These tests cover functional audit checkpoints such as verifying server port handling, client connectivity, message formatting, broadcast functionality, and proper notifications on join/leave events.
Project Structure

    main.go:
    Entry point of the application; parses command-line arguments and starts the server.

    server.go:
    Contains the ChatServer struct and methods for accepting connections, broadcasting messages, and managing chat history.

    client.go:
    Manages client-specific operations such as prompting for usernames, handling message input/output, and managing disconnections.

    formatter.go:
    Provides helper functions to format chat messages, system messages, and the welcome banner using ANSI color codes.

License

This project is licensed under the MIT License.

mfoteino , cm

