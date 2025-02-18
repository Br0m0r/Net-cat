Overview of the Project Structure:

main.go
                                        Responsibilities:

Acts as the entry point for the application.
Parses command-line arguments to determine which port the server should listen on (defaulting to 8989 if none is specified).
Initializes the ChatServer by calling functions defined in the other files and starts the server.

server.go

                                         Responsibilities:

Contains the core server logic.
Defines the ChatServer struct that manages client connections, maintains chat history, and handles message broadcasting through a dedicated channel.
Listens for new client connections, enforcing a maximum connection limit, and uses goroutines and mutexes to manage concurrent operations.

client.go

                                        Responsibilities:

Focuses on individual client interactions.
Defines the Client struct and implements functions for handling a client’s lifecycle—sending a welcome message, prompting for a name, delivering chat history, and continuously reading and broadcasting client messages.

                                       Key Concepts to Learn:

TCP/UDP Networking: Understand how to create server-client connections using Go's net package.
Goroutines: Learn how to handle concurrency by managing multiple client connections simultaneously.
Channels: Use channels for inter-goroutine communication, particularly for message broadcasting.
Mutexes: Ensure thread-safe operations on shared resources like client maps and chat history.
Struct Organization: Modularize your code by separating responsibilities into different files and structs for clarity and maintainability.
This structured approach not only makes the code more readable and maintainable but also reinforces essential concepts in building robust, concurrent network applications in Go.



             +---------------------+
             |      main.go        |
             |---------------------|
             | - Parse arguments   |
             | - Set default port  |
             | - Create ChatServer |
             | - Start server      |
             +----------+----------+
                        │
                        ▼
             +---------------------+
             |     server.go       |
             |---------------------|
             | - Listen for TCP    |
             |   connections       |
             | - Accept connections│
             | - Enforce max limit │
             | - Spawn goroutines  |
             |   for each client   |
             | - Handle broadcasting|
             +----------+----------+
                        │
                        │ (New connection)
                        ▼
             +---------------------+
             |     client.go       |
             |---------------------|
             | - Send welcome msg  |
             |   (with Linux logo) |
             | - Prompt for name   | <-------->(forgot to add formatter.go)
             | - Send chat history |
             | - Listen for client │
             |   messages          |
             +----------+----------+
                        │
                        │ (Client sends message)
                        ▼
             +---------------------+
             |   server.go         |
             | (Broadcast channel) |
             |---------------------|
             | - Append message to |
             |   chat history      |
             | - Send message to   |
             |   all clients       |
             +---------------------+





Running the Server

Open a Terminal or VS Code Integrated Terminal

Navigate to the Project Directory

cd path/to/TCPChat

Run the Server

Default Port (8989):

go run .
Custom Port (e.g., 2525):

go run . 2525

The terminal should display:
Listening on the port :8989
or your specified port.

Connecting as a Client
Using Ncat (Recommended on Windows)
Open a New CMD Window
Connect to the Server

ncat localhost 8989

Replace 8989 with your chosen port if different.
Once connected, you should see the welcome message (including the Linux logo) and the prompt:
[ENTER YOUR NAME]:

Testing the Chat
Enter Your Name and start sending messages.
Open Additional Client Windows (using ncat or telnet) to simulate multiple users.
Observe the Chat Functionality:
Each client will receive the chat history upon connecting.
Messages are broadcasted to all connected clients in real-time.