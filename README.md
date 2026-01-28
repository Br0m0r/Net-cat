Net-Cat (TCP Chat)

Overview
A lightweight TCP group chat in Go that mimics netcat-style usage with usernames, timestamps, join/leave notices, and history on join.

Features
- TCP server with multiple clients (max 10)
- Unique, non-empty usernames required
- Timestamped messages: [YYYY-MM-DD HH:MM:SS][username]:message
- Join/leave notifications
- History replay for new clients
- Default port 8989

Usage
- Default port:
  go run .
- Custom port:
  go run . 2525
- Invalid usage:
  [USAGE]: ./TCPChat $port

Client
- Connect with nc:
  nc localhost 8989
- You will see the ASCII logo and be prompted:
  [ENTER YOUR NAME]:

Files
- main.go: CLI parsing, server start
- server.go: accept loop, broadcast, history
- client.go: per-client session, name handshake
- formatter.go: ANSI colors + ASCII art
