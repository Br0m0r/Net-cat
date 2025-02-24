# net-cat: Study Guide & Code Flow Overview

This document is intended for the coding team as an internal study guide. It explains the architecture, key concepts, and code flow for the net-cat project. The goal is to help you understand how each component fits into the overall design, review essential networking and concurrency concepts, and clarify why we use them in this exercise.

---

## Project Architecture & File Responsibilities

### main.go
- **Role:**  
  Acts as the entry point of the application.
- **Responsibilities:**  
  - Parse command-line arguments.
  - Determine the server port (defaulting to 8989 if none is provided).
  - Create and start the `ChatServer`.
- **Key Concept:**  
  **Command-line Argument Parsing** — This allows the program to be flexible, letting users specify a port if needed.

### server.go
- **Role:**  
  Implements the core server functionality.
- **Responsibilities:**  
  - Define the `ChatServer` struct to manage:
    - Active client connections.
    - Chat history (a log of messages).
    - Message broadcasting using a dedicated channel.
  - Listen for incoming TCP connections.
  - Accept new connections and enforce a maximum client limit.
  - Launch a new goroutine for each client connection.
- **Key Concepts:**  
  - **TCP Connections:**  
    A TCP (Transmission Control Protocol) connection is a reliable, ordered, and error-checked communication channel between two endpoints. In net-cat, TCP connections ensure that messages sent by one client are delivered accurately and in order to the server and then to other clients.
  - **Concurrency with Goroutines:**  
    Goroutines are lightweight threads managed by the Go runtime. They allow the server to handle many clients concurrently without the overhead of traditional threads. In this exercise, each client connection runs in its own goroutine, enabling simultaneous communication.
  - **Synchronization with Mutexes:**  
    Mutexes (mutual exclusion locks) prevent concurrent access to shared resources. In net-cat, a mutex protects the clients map and chat history from race conditions when accessed by multiple goroutines simultaneously.

### client.go
- **Role:**  
  Manages individual client interactions.
- **Responsibilities:**  
  - Define the `Client` struct (storing the client’s TCP connection and username).
  - Send a welcome message (formatted via `formatter.go`) that includes an ASCII art logo.
  - Prompt the client for a unique, non-empty username.
  - Deliver the complete chat history to new clients.
  - Continuously read messages from the client and pass them to the server for broadcasting.
  - Handle client disconnections and notify remaining clients.
- **Key Concepts:**  
  - **User Input & Validation:**  
    Ensuring that every client enters a valid and unique username is crucial for proper identification in the chat.
  - **Real-Time Message Handling:**  
    The continuous loop that reads client messages allows for real-time communication. This is essential for a live chat system.

### formatter.go
- **Role:**  
  Provides utility functions for styling output messages.
- **Responsibilities:**  
  - Format the welcome message with an ASCII art logo.
  - Format chat messages to include a timestamp, the sender’s username, and the message content.
  - Format system messages (e.g., notifications for join/leave events) with distinct ANSI colors.
- **Key Concepts:**  
  **ANSI Escape Codes:**  
  These codes enable colored text output in terminal applications. They are used in net-cat to visually differentiate various types of messages (welcome messages, chat messages, system notifications), thereby enhancing readability and user experience.

---

## Application Flow Visualization

The following diagram illustrates the flow of the application from startup to message broadcasting:

            +--------------------------+
            |         main.go          |
            |--------------------------|
            | - Parse command-line     |
            |   arguments              |
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
            |   (with ASCII logo)      |
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


---

## Key Concepts and Terminology

### TCP Connection
- **Definition:**  
  A TCP connection establishes a reliable, bidirectional communication channel between two endpoints (typically a client and a server). TCP ensures that data is transmitted in order, without loss or duplication.
- **Why We Use It:**  
  In net-cat, TCP guarantees that every message sent by one client reaches the server and, subsequently, all other connected clients accurately. This reliability is essential for a group chat application where message order and integrity are critical.
- **Key Properties:**  
  - **Connection-oriented:** A dedicated connection is established before any data is transferred.
  - **Reliable:** Built-in error checking and correction ensure data integrity.
  - **Ordered:** Data packets arrive in the same order they were sent.

### Goroutines
- **Definition:**  
  Goroutines are functions or methods that run concurrently with other functions or methods. They are managed by the Go runtime and are much lighter than traditional threads.
- **Why We Use Them:**  
  In net-cat, each client is handled in its own goroutine, allowing the server to manage multiple client connections simultaneously. This concurrency ensures that one slow client does not block the others, making the chat system responsive.
- **Key Benefits:**  
  - **Lightweight:** Thousands of goroutines can run concurrently with minimal resource usage.
  - **Simple Concurrency Model:** They simplify concurrent programming by abstracting many of the complexities of thread management.

### Channels
- **Definition:**  
  Channels are a Go language construct that allow goroutines to communicate with each other and synchronize their execution. They provide a safe way to pass data between concurrent processes.
- **Why We Use Them:**  
  In net-cat, a dedicated broadcast channel is used to send messages from the server to all connected clients. Channels help ensure that messages are delivered in a thread-safe manner.
- **Key Benefits:**  
  - **Safe Communication:** Channels prevent race conditions by enforcing controlled access to shared data.
  - **Synchronization:** They can be used to coordinate the execution of multiple goroutines.

### Mutexes
- **Definition:**  
  A mutex (mutual exclusion lock) is a synchronization primitive that is used to prevent concurrent access to shared resources.
- **Why We Use Them:**  
  In net-cat, mutexes protect the clients map and chat history from concurrent modifications. Without mutexes, simultaneous access by multiple goroutines could lead to data corruption or unexpected behavior.
- **Key Benefits:**  
  - **Data Integrity:** Mutexes ensure that only one goroutine can modify a shared resource at any given time.
  - **Simplicity:** They provide a straightforward way to handle synchronization issues.

### ANSI Escape Codes
- **Definition:**  
  ANSI escape codes are sequences of characters used to control text formatting, color, and other output options in terminal emulators.
- **Why We Use Them:**  
  In net-cat, these codes are used in the `formatter.go` file to create visually distinct messages:
  - **Welcome Message:** Features an ASCII art logo with color for an enhanced visual impact.
  - **Chat Messages:** Display timestamps and usernames in different colors to improve readability.
  - **System Notifications:** Use distinct colors to differentiate join/leave messages from regular chat messages.
- **Key Benefits:**  
  - **Improved Readability:** Colored output helps users quickly identify different types of messages.
  - **Aesthetic Appeal:** Enhances the overall user experience in the terminal.

---

## Summary

The net-cat project is designed to demonstrate and reinforce fundamental networking and concurrency concepts in Go. By dividing the code into modular components—**main.go** (application entry point), **server.go** (server logic), **client.go** (client management), and **formatter.go** (message styling)—we achieve a clean, scalable, and maintainable codebase.

Understanding how these components interact, from establishing reliable TCP connections, handling user input and real-time messaging, to synchronizing concurrent operations with goroutines, channels, and mutexes, is key to mastering both the theoretical and practical aspects of network programming in Go.

This study guide serves as a reference for the team to quickly grasp the project's structure and the underlying principles that make net-cat a robust, real-time chat application.

