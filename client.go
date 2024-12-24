package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	chatLeaveChannel = make(chan bool)
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

func checkForCommands(message []byte, client net.Conn) bool {
	switch string(message) {
	case "/exit":
		client.Write([]byte("/exit"))
		chatLeaveChannel <- true

		return true
	case "/whois":
		fmt.Println("[COMMANDS]: You are: ", client.LocalAddr())

		return true
	}

	return false
}

// Handles sending messages to the server
func handleSendingMessages(client net.Conn) {
	for {
		message := requestMessageInput()

		// Prevent sending a message to others if it's a command
		if checkForCommands(message, client) {
			continue
		}

		client.Write(message)
	}
}

// Outputs the welcome message that the client receives from the server
func outputServerWelcomeMessage(client net.Conn) {
	buffer := make([]byte, 1024)
	client.Read(buffer)
	fmt.Println(string(buffer))
	fmt.Println(`List of commands: /exit, /whois`)
}

func main() {
	client, err := net.Dial("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	outputServerWelcomeMessage(client)
	go handleIncomingMessages(client)
	go handleSendingMessages(client)

	<-chatLeaveChannel
	client.Close()
	fmt.Println("Chat has ended.")
}