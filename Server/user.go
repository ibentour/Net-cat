package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct { // Client represents a connected client
	id      net.Conn
	name    string
	channel chan string
	done    chan struct{}
}

var ( // Globale Variables
	clients     = sync.Map{}
	broadcast   = make(chan string)
	logFile     *os.File
	clientLock  sync.Mutex
	maxClients  int32 = 10
	clientCount int32
)

// Handles a new client connection
func handleUser(user net.Conn) {
	defer func() {
		atomic.AddInt32(&clientCount, -1)
		user.Close()
		fmt.Println("Current number of clients :", clientCount)
	}()
	fmt.Println("Current number of clients :", clientCount)

	fmt.Println("New connection from", user.RemoteAddr())

	fmt.Fprint(user, "Welcome to TCP-Chat!\n")
	fmt.Fprint(user, "         _nnnn_\n")
	fmt.Fprint(user, "        dGGGGMMb\n")
	fmt.Fprint(user, "       @p~qp~~qMb\n")
	fmt.Fprint(user, "       M|@||@) M|\n")
	fmt.Fprint(user, "       @,----.JM|\n")
	fmt.Fprint(user, "      JS^\\__/  qKL\n")
	fmt.Fprint(user, "     dZP        qKRb\n")
	fmt.Fprint(user, "    dZP          qKKb\n")
	fmt.Fprint(user, "   fZP            SMMb\n")
	fmt.Fprint(user, "   HZM            MMMM\n")
	fmt.Fprint(user, "   FqM            MMMM\n")
	fmt.Fprint(user, " __| \".        |\\dS\"qML\n")
	fmt.Fprint(user, " |    `.       | `' \\Zq\n")
	fmt.Fprint(user, "_)      \\.___.,|     .'\n")
	fmt.Fprint(user, "\\____   )MMMMMP|   .'\n")
	fmt.Fprint(user, "     `-'       `--'\n")

	// Get the username from the client
	username := getUsername(user)

	newClient := &Client{
		id:      user,
		name:    username,
		channel: make(chan string, 10),
		done:    make(chan struct{}),
	}

	// Add the new client to the list of clients
	clientLock.Lock()
	clients.Store(user, newClient)
	clientLock.Unlock()

	defer func() {
		if username != "Unknown!" {
			// Broadcast a message to all clients when a client leaves
			disconnectMsg := fmt.Sprintf("[%s] has left our chat...", newClient.name)
			broadcast <- disconnectMsg
		}
		fmt.Printf("Client %s Disconnected.\n", newClient.name)
		// logFile.WriteString(disconnectMsg)

		clientLock.Lock()
		clients.Delete(user)
		clientLock.Unlock()
		close(newClient.done)
	}()

	// Send chat history to the new client
	go func() {
		file, err := os.Open("chat.log")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Fprintf(user, "%s\n", scanner.Text())
		}
	}()

	go printingBrodcasts(newClient)

	time.Sleep(20 * time.Millisecond)
	if username != "Unknown!" {
		// Send a welcome message to the new client
		fmt.Fprint(user, "\nConnected to the server...\nType '/exit' to quit OR '/change' to change name.\n\n")
	}

	// Broadcast a message to all clients when a new client joins
	time.Sleep(30 * time.Millisecond)
	if username != "Unknown!" {
		// Send a welcome message to the new client
		welcome := fmt.Sprintf("[%s] has joined our chat...", username)
		broadcast <- welcome
	}

	// time.Sleep(1000 * time.Millisecond)
	readingFunc(newClient)
}

// Retrieves the username from the client
func getUsername(user net.Conn) string {
	scanner := bufio.NewScanner(user)
	fmt.Fprint(user, "[ENTER YOUR NAME]: ")
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		filtered := filterNonPrintableASCII(username)
		if filtered != "" && filtered != "Server" && !strings.HasPrefix(username, "/change") &&
			len(filtered) != 0 && !strings.Contains(username, "\x1b") {

			clientLock.Lock()
			var exists bool
			clients.Range(func(key, value interface{}) bool {
				client := value.(*Client)
				if client.name == username {
					exists = true
					return false
				}
				return true
			})
			clientLock.Unlock()

			if exists {
				clearLastLine(user)
				clearLastLine(user)
				fmt.Fprintf(user, "Username already exists! Please try again...\n")
			} else {
				return username
			}

		} else {
			clearLastLine(user)
			clearLastLine(user)
			fmt.Fprintf(user, "Invalid name! Please try again...\n")
		}
		fmt.Fprint(user, "[ENTER YOUR NAME]: ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(user, "Bye Bye! :)\n")
		return "Unknown!"
	}
	return "Unknown!"
}

// Sends messages to a client
func printingBrodcasts(client *Client) {
	for {
		select {
		case message := <-client.channel:
			if message == "[Server Shutting Down!]\n" {
				fmt.Fprint(client.id, "\n"+message)
			} else if message != "" {
				fmt.Fprintf(client.id, "\n\033[1A\033[K%s\n", message)
				timePrompt(client)
			}
		case <-client.done:
			return
		}
	}
}

// Receives messages from a client
func readingFunc(client *Client) {
	scanner := bufio.NewScanner(client.id)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "" || len(text) == 0 || strings.Contains(text, "\x1b") {
			clearLastLine(client.id)
			timePrompt(client)
			continue
		}

		if strings.HasPrefix(text, "/change") {
			oldName := client.name
			newName := getUsername(client.id)
			client.name = newName
			msg := fmt.Sprintf("[%s] has changed their name to [%s]", oldName, newName)
			broadcast <- msg
			logFile.WriteString(msg)
			continue
		}

		if text == "/exit" {
			fmt.Fprintln(client.id, "Bye Bye! :)")
			client.id.Close()
			return
		}

		filteredText := strings.TrimSpace(filterNonPrintableASCII(text))
		if filteredText == "" {
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		msg := fmt.Sprintf("[%s][%s]: %s", timestamp, client.name, filteredText)
		clearLastLine(client.id)
		broadcast <- msg
		logFile.WriteString(msg + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Closing the client %s Connection...\n", client.name)
	}
}

// Filters out non-printable ASCII characters from a string
func filterNonPrintableASCII(input string) string {
	var result strings.Builder
	for _, char := range input {
		if char >= 32 && char <= 126 {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func clearLastLine(user net.Conn) {
	fmt.Fprint(user, "\033[1A") // Move cursor up one line
	fmt.Fprint(user, "\033[K")  // Clear from cursor to end of line
}

func timePrompt(client *Client) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(client.id, "[%s][%s]: ", timestamp, client.name)
}
