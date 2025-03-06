Net-Cat

Net-Cat is a simple TCP-based group chat application written in Go. It mimics the functionality of the classic nc (NetCat) utility and allows multiple clients to connect and communicate in real-time through a central server.

Features

TCP Server and Clients - Supports multiple simultaneous connections.

User Identification - Clients must enter a unique name upon connection.

Message Broadcasting - Messages are sent to all connected clients with timestamps.

Join/Leave Notifications - Clients are informed when someone joins or leaves.

Chat History - New clients receive past messages upon connection.

Concurrency Handling - Uses Goroutines, Channels, and Mutexes for efficient operations.

Default Port Handling - Defaults to port 8989 if none is specified.

Max Connections Limit - Restricts the chat to 10 concurrent users.

Installation

Clone the repository:

git clone https://platform.zone01.gr/git/mfoteino/net-cat
cd net-cat

Build the project (optional):

go build -o TCPchat .

Usage

Running the Server

Default Port (8989):

go run .

Custom Port (e.g., 2525):

go run . 2525

Invalid Usage:

If additional parameters are provided, the program returns:

[USAGE]: ./TCPChat $port

Connecting as a Client

Use nc (NetCat) or any TCP client:

nc localhost 8989

Upon connection, you will receive a welcome message and be prompted to enter your username.

Testing

Usage Validation:

go run . 2525 localhost

Expected output:

[USAGE]: ./TCPChat $port

Multi-Client Testing:

Start the server:

go run . 2525

Open multiple terminals and connect:

nc localhost 2525

Verify:

Clients receive a welcome message and username prompt.

Messages appear with timestamps and usernames.

New clients receive the chat history.

Disconnect notifications work.

File Structure

main.go - Entry point; parses arguments and starts the server.

server.go - Manages client connections, message broadcasting, and history.

client.go - Handles user interactions and communication with the server.

formatter.go - Provides formatted messages with ANSI color codes.

License

This project is licensed under the MIT License.

Developed by mfoteino & cm.