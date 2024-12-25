package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

// Requests the user for a message input
func requestMessageInput() []byte {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return []byte(strings.TrimSpace(scanner.Text())) // Remove any trailing whitespaces
}

// Handles incoming messages from the server
func handleIncomingMessages(client net.Conn) {
	for {
		buff := make([]byte, 1024)

		_, err := client.Read(buff)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buff))
	}
}

// Handles sending messages to the server
func handleSendingMessages(client net.Conn) {
	for {
		message := requestMessageInput()
		client.Write(message)
	}
}

func main() {
	chatLeaveChannel := make(chan os.Signal, 1)

	signal.Notify(chatLeaveChannel, os.Interrupt)

	client, err := net.Dial("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`Type /help to see the available commands.`)
	go handleIncomingMessages(client)
	go handleSendingMessages(client)

	<-chatLeaveChannel
	client.Close()
	fmt.Println("Chat has ended.")
}
