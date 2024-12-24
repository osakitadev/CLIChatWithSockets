package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

var connectedClients = []net.Conn{}

func replicateMessage(message []byte, sender net.Conn) {
	for _, client := range connectedClients {
		if client == sender {
			continue
		}

		fmt.Fprintf(client, "[Client %v]: %v", client.RemoteAddr(), string(message))
	}
}

func handleClientConnection(conn net.Conn) {
	conn.Write([]byte("Welcome to the chat, say hello to everyone!"))

	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)

		if err != nil {
			return
		}

		message := buff[:n] // Cut off the trailing bytes for correct checking and some memory earnings
		if bytes.Equal(message, []byte("/exit")) {
			// FIXME: Remove the client from the list of connected clients, it wasn't working before idk why
			// deleteClient(conn)
			log.Println(conn.RemoteAddr(), "disconnected from the chat")
			conn.Close()
			continue
		}

		replicateMessage(message, conn)
	}
}

func handleAcceptIncomingClients(server net.Listener) {
	for {
		client, err := server.Accept()

		if err != nil {
			log.Fatal("Couldn't accept client connection", err)
		}

		connectedClients = append(connectedClients, client)
		log.Println(client.RemoteAddr(), "has been added to the list of clients")
		go handleClientConnection(client)
	}
}

func disconnectClients() {
	for _, client := range connectedClients {
		client.Write([]byte("Server is shutting down, goodbye!"))
		client.Close()
	}
}

func main() {
	log.Println("Listening for new clients...")

	server, err := net.Listen("tcp", ":8080")
	closeServerChannel := make(chan os.Signal, 1)

	signal.Notify(closeServerChannel, os.Interrupt)

	if err != nil {
		log.Fatal("Couldn't start server:", err)
	}

	go handleAcceptIncomingClients(server)
	<-closeServerChannel
	disconnectClients()
	server.Close()
}
