package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

// type Client struct { // Client represents a connected client
// 	id      net.Conn
// 	name    string
// 	channel chan string
// 	done    chan struct{}
// }

// var ( // Globale Variables
// 	clients     = sync.Map{}
// 	broadcast   = make(chan string)
// 	logFile     *os.File
// 	clientLock  sync.Mutex
// 	maxClients  int32 = 10
// 	clientCount int32
// )

func main() {
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	fmt.Printf("Listening on port :%s\n", port)

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Open log file
	logFile, err = os.OpenFile("chat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		logFile = nil
	}

	defer func() {
		ln.Close()
		if logFile != nil {
			logFile.Close()
			os.Remove("chat.log") // Remove the log file
		}
	}()

	go broadcastMessages()

	// Create a signal channel to catch SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Start a goroutine to handle SIGINT
	go func() {
		<-sigChan
		fmt.Println("Server Shutting Down...")
		broadcast <- fmt.Sprintln("[Server Shutting Down!]")

		time.Sleep(50 * time.Millisecond)
		clients.Range(func(key, value interface{}) bool {
			client := value.(*Client)
			client.id.Close()
			return true
		})

		if logFile != nil {
			logFile.Close()
			os.Remove("chat.log") // Remove the log file
		}

		fmt.Println("Server OFF!")
		os.Exit(0)
	}()

	for {
		newUser, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		if atomic.LoadInt32(&clientCount) >= maxClients {
			fmt.Fprint(newUser, "Server is full! Please try again later.\n")
			newUser.Close()
			continue
		}

		atomic.AddInt32(&clientCount, 1)

		go handleUser(newUser)
	}
}

// Broadcasts messages to all clients
func broadcastMessages() {
	for message := range broadcast {
		clients.Range(func(key, value interface{}) bool {
			client := value.(*Client)
			select {
			case client.channel <- message:
			default:
				// Channel is full, block until space is available
				client.channel <- message
			}
			return true
		})
	}
}

// Handles a new client connection
// func handleUser(user net.Conn) {
// 	defer func() {
// 		atomic.AddInt32(&clientCount, -1)
// 		user.Close()
// 		fmt.Println("Current number of clients :", clientCount)
// 	}()
// 	fmt.Println("Current number of clients :", clientCount)

// 	fmt.Println("New connection from", user.RemoteAddr())

// 	fmt.Fprint(user, "Welcome to TCP-Chat!\n")
// 	fmt.Fprint(user, "         _nnnn_\n")
// 	fmt.Fprint(user, "        dGGGGMMb\n")
// 	fmt.Fprint(user, "       @p~qp~~qMb\n")
// 	fmt.Fprint(user, "       M|@||@) M|\n")
// 	fmt.Fprint(user, "       @,----.JM|\n")
// 	fmt.Fprint(user, "      JS^\\__/  qKL\n")
// 	fmt.Fprint(user, "     dZP        qKRb\n")
// 	fmt.Fprint(user, "    dZP          qKKb\n")
// 	fmt.Fprint(user, "   fZP            SMMb\n")
// 	fmt.Fprint(user, "   HZM            MMMM\n")
// 	fmt.Fprint(user, "   FqM            MMMM\n")
// 	fmt.Fprint(user, " __| \".        |\\dS\"qML\n")
// 	fmt.Fprint(user, " |    `.       | `' \\Zq\n")
// 	fmt.Fprint(user, "_)      \\.___.,|     .'\n")
// 	fmt.Fprint(user, "\\____   )MMMMMP|   .'\n")
// 	fmt.Fprint(user, "     `-'       `--'\n")

// 	// Get the username from the client
// 	username := getUsername(user)

// 	newClient := &Client{
// 		id:      user,
// 		name:    username,
// 		channel: make(chan string, 10),
// 		done:    make(chan struct{}),
// 	}

// 	// Add the new client to the list of clients
// 	clientLock.Lock()
// 	clients.Store(user, newClient)
// 	clientLock.Unlock()

// 	defer func() {
// 		if username != "Unknown!" {
// 			// Broadcast a message to all clients when a client leaves
// 			disconnectMsg := fmt.Sprintf("[%s] has left our chat...", newClient.name)
// 			broadcast <- disconnectMsg
// 		}
// 		fmt.Printf("Client %s Disconnected.\n", newClient.name)
// 		// logFile.WriteString(disconnectMsg)

// 		clientLock.Lock()
// 		clients.Delete(user)
// 		clientLock.Unlock()
// 		close(newClient.done)
// 	}()

// 	// Send chat history to the new client
// 	go func() {
// 		file, err := os.Open("chat.log")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			fmt.Fprintf(user, "%s\n", scanner.Text())
// 		}
// 	}()

// 	go printingBrodcasts(newClient)

// 	time.Sleep(20 * time.Millisecond)
// 	if username != "Unknown!" {
// 		// Send a welcome message to the new client
// 		fmt.Fprint(user, "\nConnected to the server...\nType '/exit' to quit OR '/chnage <Your New name>' to change name.\n\n")
// 	}

// 	// Broadcast a message to all clients when a new client joins
// 	time.Sleep(30 * time.Millisecond)
// 	if username != "Unknown!" {
// 		// Send a welcome message to the new client
// 		welcome := fmt.Sprintf("[%s] has joined our chat...", username)
// 		broadcast <- welcome
// 	}

// 	// time.Sleep(1000 * time.Millisecond)
// 	readingFunc(newClient)
// }

// // Retrieves the username from the client
// func getUsername(user net.Conn) string {
// 	scanner := bufio.NewScanner(user)
// 	fmt.Fprint(user, "[ENTER YOUR NAME]: ")
// 	for scanner.Scan() {
// 		username := strings.TrimSpace(scanner.Text())
// 		filtered := filterNonPrintableASCII(username)
// 		if filtered != "" && filtered != "Server" && !strings.HasPrefix(username, "/change") &&
// 			len(filtered) != 0 && !strings.Contains(username, "\x1b") {

// 			clientLock.Lock()
// 			var exists bool
// 			clients.Range(func(key, value interface{}) bool {
// 				client := value.(*Client)
// 				if client.name == username {
// 					exists = true
// 					return false
// 				}
// 				return true
// 			})
// 			clientLock.Unlock()

// 			if exists {
// 				clearLastLine(user)
// 				clearLastLine(user)
// 				fmt.Fprintf(user, "Username already exists! Please try again...\n")
// 			} else {
// 				return username
// 			}

// 		} else {
// 			clearLastLine(user)
// 			clearLastLine(user)
// 			fmt.Fprintf(user, "Invalid name! Please try again...\n")
// 		}
// 		fmt.Fprint(user, "[ENTER YOUR NAME]: ")
// 	}
// 	if err := scanner.Err(); err != nil {
// 		fmt.Fprintf(user, "Bye Bye! :)\n")
// 		return "Unknown!"
// 	}
// 	return "Unknown!"
// }

// // Sends messages to a client
// func printingBrodcasts(client *Client) {
// 	for {
// 		select {
// 		case message := <-client.channel:
// 			if message == "[Server Shutting Down!]\n" {
// 				fmt.Fprint(client.id, "\n"+message)
// 			} else if message != "" {
// 				fmt.Fprintf(client.id, "\n\033[1A\033[K%s\n", message)
// 				timePrompt(client)
// 			}
// 		case <-client.done:
// 			return
// 		}
// 	}
// }

// // Receives messages from a client
// func readingFunc(client *Client) {
// 	scanner := bufio.NewScanner(client.id)
// 	for scanner.Scan() {
// 		text := scanner.Text()
// 		if strings.TrimSpace(text) == "" || len(text) == 0 || strings.Contains(text, "\x1b") {
// 			clearLastLine(client.id)
// 			timePrompt(client)
// 			continue
// 		}

// 		if strings.HasPrefix(text, "/change") {
// 			oldName := client.name
// 			newName := getUsername(client.id)
// 			client.name = newName
// 			msg := fmt.Sprintf("[%s] has changed their name to [%s]\n", oldName, newName)
// 			broadcast <- msg
// 			logFile.WriteString(msg)
// 			continue
// 		}

// 		if text == "/exit" {
// 			fmt.Fprintln(client.id, "Bye Bye! :)")
// 			client.id.Close()
// 			return
// 		}

// 		filteredText := strings.TrimSpace(filterNonPrintableASCII(text))
// 		if filteredText == "" {
// 			continue
// 		}

// 		timestamp := time.Now().Format("2006-01-02 15:04:05")
// 		msg := fmt.Sprintf("[%s][%s]: %s", timestamp, client.name, filteredText)
// 		clearLastLine(client.id)
// 		broadcast <- msg
// 		logFile.WriteString(msg + "\n")
// 	}

// 	if err := scanner.Err(); err != nil {
// 		fmt.Printf("Closing the client %s Connection...\n", client.name)
// 	}
// }

// // Filters out non-printable ASCII characters from a string
// func filterNonPrintableASCII(input string) string {
// 	var result strings.Builder
// 	for _, char := range input {
// 		if char >= 32 && char <= 126 {
// 			result.WriteRune(char)
// 		}
// 	}
// 	return result.String()
// }

// // Broadcasts messages to all clients
// func broadcastMessages() {
// 	for message := range broadcast {
// 		clients.Range(func(key, value interface{}) bool {
// 			client := value.(*Client)
// 			select {
// 			case client.channel <- message:
// 			default:
// 				// Channel is full, block until space is available
// 				client.channel <- message
// 			}
// 			return true
// 		})
// 	}
// }

// func clearLastLine(user net.Conn) {
// 	fmt.Fprint(user, "\033[1A") // Move cursor up one line
// 	fmt.Fprint(user, "\033[K")  // Clear from cursor to end of line
// }

// func timePrompt(client *Client) {
// 	timestamp := time.Now().Format("2006-01-02 15:04:05")
// 	fmt.Fprintf(client.id, "[%s][%s]: ", timestamp, client.name)
// }
