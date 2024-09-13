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
