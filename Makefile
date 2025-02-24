# Makefile for net-cat

BINARY = net-cat

.PHONY: build run run-custom test clean

# Build the net-cat binary.
build:
	go build -o $(BINARY) .

# Run the server on the default port (8989).
run: build
	./$(BINARY)

# Run the server on a custom port (2525).
run-custom: build
	./$(BINARY) 2525

# Run a basic usage test to verify the usage message.
test: build
	@echo "Running usage test: expecting '[USAGE]: ./TCPChat $$port'"
	@./$(BINARY) 2525 localhost

# Clean the built binary.
clean:
	rm -f $(BINARY)
