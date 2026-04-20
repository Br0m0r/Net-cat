# Net-Cat (TCP Chat)

A TCP group chat server written in Go, inspired by the netcat (`nc`) utility. Clients connect via any TCP client, choose a username, and exchange timestamped messages in real time.

## Features

| Feature | Description |
|---|---|
| Group chat | Up to 10 concurrent clients |
| Usernames | Unique, non-empty names required on connect |
| Timestamps | Messages formatted as `[YYYY-MM-DD HH:MM:SS][user]: text` |
| Notifications | Join/leave announcements with online count |
| History | New clients receive the full message history |
| Timeout | Clients that idle at the name prompt are disconnected after 30s |

## Project Structure

| File | Purpose |
|---|---|
| `main.go` | Entry point and CLI argument parsing |
| `server.go` | Listener, connection management, broadcast loop |
| `client.go` | Per-client session lifecycle |
| `formatter.go` | ANSI color formatting and ASCII welcome banner |

## Requirements

- Go 1.22 or later

## Run

```
go run .              # listens on port 8989
go run . 2525         # listens on port 2525
```

## Connecting

Use any TCP client such as netcat:

```
nc localhost 8989
```

## Usage

```
[USAGE]: ./TCPChat $port
```

If no port is provided, the server defaults to `8989`. Only one argument (the port) is accepted.
