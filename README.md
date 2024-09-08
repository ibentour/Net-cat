# Net-cat Server

This is a TCP chat server written in Go, allowing multiple clients to connect and communicate with each other in real-time. The server handles client connections, broadcasts messages, and maintains a chat log.

## Features

*   Multiple client connections (up to 10 concurrent connections)
*   Real-time message broadcasting to all connected clients
*   Chat log maintenance (stored in a file named "chat.log")
*   Client disconnection handling
*   Server shutdown handling (with a 30ms delay to allow clients to receive the shutdown message)
*   Username management (clients can change their usernames)
*   Message filtering (non-printable ASCII characters are removed)
*   Code Overview
*   The code is organized into several functions, each responsible for a specific task:

### main()

- Initializes the server, listening on a specified port (default: 8989)
- Handles SIGINT signals to shut down the server gracefully
- Opens a log file for chat logging
- Accepts new client connections and handles them in separate goroutines

### handleUser()
- Handles a new client connection
- Sends a welcome message and chat history to the client
- Broadcasts a message to all clients when a new client joins
- Receives and processes messages from the client
- Handles client disconnections

### getUsername()

- Retrieves a username from a client
- Checks for username validity and uniqueness
- Returns the username or an error message

### sendMessagesToClient()

- Sends messages to a client from the broadcast channel
- Handles client disconnections

### receiveMessagesFromClient()

- Receives messages from a client
- Processes and filters messages (removing non-printable ASCII characters)
- Broadcasts messages to all clients
- Handles client disconnections

### isArrowKey()

- Checks if a message contains an arrow key press (Up, Down, Left, or Right)
- filterNonPrintableASCII()
- Removes non-printable ASCII characters from a string

### broadcastMessages()

- Broadcasts messages to all connected clients

## Usage

To run the server, execute the following command:

``` bash
go run server.go
```
To specify a custom port, use the following command:

``` bash
go run server.go <port>
```
Replace <port> with the desired port number.

