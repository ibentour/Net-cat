# Net-cat Server

This is a TCP chat server written in Go, allowing multiple clients to connect and communicate with each other in real-time. The server handles client connections, broadcasts messages, and maintains a chat log.

## Authors

* Ismail Bentour

## Table of Contents

* Features
* Requirements
* Installation
* Usage
* Project Structure
* How It Works

## Features

* Multi-client support with a maximum of 10 concurrent connections
* Real-time message broadcasting
* Username selection and validation
* Chat history logging
* Graceful server shutdown
* Command support (e.g., /exit, /change)
* ASCII art welcome message

## Requirements

Go 1.15 or higher

Installation

Clone the repository:
``` bash

```

Navigate to the project directory:
``` bash

```

Build the project:
``` bash

```


## Usage

### Starting the Server

To start the server on the default port (8989):
``` bash
```

To start the server on a specific port:
``` bash

```

Connecting to the Server
Use a TCP client (like netcat) to connect to the server:
``` bash
nc localhost 8989
```

## Project Structure

### The project consists of two main Go files:

* server.go: Contains the main function and server initialization logic.
* user.go: Handles client connections and message processing.

## How It Works

* The server listens for incoming TCP connections on the specified port.
* When a client connects, they are greeted with an ASCII art welcome message.
* The client is prompted to enter a username, which is validated for uniqueness and content.
* Once a valid username is provided, the client joins the chat room.
* Messages sent by clients are broadcasted to all connected clients.
* The server maintains a log of all chat messages in a file named chat.log.
* Clients can use special commands like /exit to leave the chat or /change to change their username.
* The server gracefully handles interrupts (Ctrl+C) and notifies all clients before shutting down.
