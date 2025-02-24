# Project Overview: net-cat

The **net-cat** project is a Go-based group chat application inspired by the classic NetCat (`nc`) utility. It uses a server-client architecture to enable multiple clients to connect to a central server and exchange messages in real time. This document explains the responsibilities of each file and illustrates the flow of the application.

## File Responsibilities

### main.go
- **Purpose:** Entry point of the application.
- **Key Functions:**
  - Parses command-line arguments.
  - Sets the default port (`8989`) if none is provided.
  - Instantiates a `ChatServer` (defined in `server.go`).
  - Starts the server.

### server.go
- **Purpose:** Implements core server functionality.
- **Key Functions:**
  - Defines the `ChatServer` struct, which manages:
    - Client connections.
    - Chat history.
    - Message broadcasting via a dedicated channel.
  - Listens for incoming TCP connections.
  - Accepts new connections and enforces a maximum client limit.
  - Spawns goroutines to handle each client connection.

### client.go
- **Purpose:** Manages individual client interactions.
- **Key Functions:**
  - Defines the `Client` struct to hold a client's connection and username.
  - Sends a welcome message (formatted using functions in `formatter.go`).
  - Prompts the client for a unique, non-empty username.
  - Delivers the complete chat history to new clients.
  - Continuously reads messages from the client and passes them to the server for broadcasting.
  - Handles client disconnections and notifies remaining clients.

### formatter.go
- **Purpose:** Provides helper functions for formatting messages.
- **Key Functions:**
  - Formats the welcome message (with an ASCII art logo).
  - Formats chat messages to include a timestamp, username, and the message content.
  - Formats system messages (such as join/leave notifications) using ANSI color codes.

## Application Flow

Below is a visualization of the application's flow, showing where each component fits:

            +--------------------------+
            |         main.go          |
            |--------------------------|
            | - Parse arguments        |
            | - Set default port (8989) |
            | - Create ChatServer      |
            | - Start server           |
            +-------------+------------+
                          │
                          ▼
            +--------------------------+
            |        server.go         |
            |--------------------------|
            | - Listen for TCP         |
            |   connections            |
            | - Accept new connections |
            | - Enforce max client     |
            |   limit                  |
            | - Spawn goroutines for   |
            |   each client            |
            +-------------+------------+
                          │
                          │ (New connection)
                          ▼
            +--------------------------+
            |        client.go         |
            |--------------------------|
            | - Send welcome message   |
            |   (with Linux logo)      |
            | - Prompt for username    |
            | - Send chat history      |
            | - Read client messages   |
            +-------------+------------+
                          │
                          │ (Client sends message)
                          ▼
            +--------------------------+
            |        server.go         |
            |  (Broadcasting Process)  |
            |--------------------------|
            | - Append message to      |
            |   chat history           |
            | - Broadcast message to   |
            |   all clients            |
            +--------------------------+


## Key Concepts Demonstrated

- **TCP Networking:**  
  Uses Go’s `net` package to establish and manage TCP connections between the server and clients.

- **Concurrency:**  
  Employs goroutines so that each client connection is handled concurrently, ensuring the server can manage multiple connections simultaneously.

- **Channels and Mutexes:**  
  - **Channels:** A dedicated broadcast channel is used to send messages from the server to all connected clients.
  - **Mutexes:** Ensure thread-safe operations when accessing shared resources like the clients map and chat history.

- **Modular Code Organization:**  
  The project is split into multiple files:
  - `main.go` handles application startup.
  - `server.go` contains the server logic.
  - `client.go` manages client interactions.
  - `formatter.go` provides message formatting utilities.  
  This separation improves clarity and maintainability.

This structure not only reinforces essential programming concepts in Go (such as networking, concurrency, and synchronization) but also ensures a clean and scalable codebase for building robust, real-time chat applications.

