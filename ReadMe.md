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
- **Allowed Packages:** Uses only permitted packages such as `io`, `log`, `os`, `fmt`, `net`, `sync`, `time`, `bufio`, `errors`, `strings`, and `reflect`.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://platform.zone01.gr/git/mfoteino/net-cat
   cd net-cat

    Build the project (optional):

    make build

Usage
Running the Server

You can run the server directly using Go or via the Makefile targets.

    Default Port (8989):

make run

Custom Port (2525):

make run-custom

Invalid Usage:
Running with extra parameters (e.g., ./net-cat 2525 localhost) displays:

    [USAGE]: ./TCPChat $port

Connecting as a Client

Use NetCat (nc) to connect to the server. For example, if the server is running on port 8989:

nc localhost 8989

Upon connection, the client will see a welcome message with an ASCII art logo and a prompt to enter a username.
Testing

To verify key functionalities, follow these steps:

    Usage Test:
    Run the usage test to confirm that extra parameters trigger the usage message.

make test

Audit Reference:
"Try running ./TCPChat 2525 localhost. Did the server respond with usage, as above?"

Multi-Client Test:

    Open one terminal and run:

make run-custom

Open two or more separate terminals and connect using:

        nc localhost 2525

    Verify that:
        Clients receive the welcome message and prompt for their name.
        When a client enters a valid name, all clients are notified.
        Messages are broadcast with proper timestamps and sender names.
        A new client sees all previous messages.
        Remaining clients receive notifications if a client disconnects.

    Audit Reference:
    These steps cover several functional audit checkpoints regarding client connection, message broadcasting, chat history, and disconnection notifications.

    Multi-Computer Test:
    Ensure that clients on different computers can connect and exchange messages.

Makefile Targets and Audit Mapping

    make build
    What it does: Compiles the net-cat binary.
    Audit Check: Validates project compilation (addresses "Invalid compilation" concerns).

    make run
    What it does: Runs the server on the default port (8989).
    Audit Check:
    "Try running ./TCPChat. Is the server listening for connections on the default port?"

    make run-custom
    What it does: Runs the server on port 2525.
    Audit Check:
    "Try running ./TCPChat 2525. Is the server listening for connections on the port 2525?"

    make test
    What it does: Executes a usage test with extra arguments.
    Audit Check:
    "Try running ./TCPChat 2525 localhost. Did the server respond with usage, as above?"

    make clean
    What it does: Removes the compiled net-cat binary.
    Audit Check:
    Ensures a clean build environment, following good project practices.

Run a target by executing:

make <target>

For example:

make build
make run
make run-custom
make test
make clean

Project Structure

    main.go:
    Entry point; parses command-line arguments and starts the server.

    server.go:
    Contains the ChatServer struct and methods for handling connections, broadcasting messages, and managing chat history.

    client.go:
    Manages client-specific operations such as prompting for usernames and handling message exchanges.

    formatter.go:
    Provides functions to format chat messages, system messages, and the welcome banner using ANSI color codes.

License

This project is licensed under the MIT License.

