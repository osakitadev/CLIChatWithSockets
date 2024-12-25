package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Osaka/chat-with-sockets/config"
)

type Server struct {
	connectedClients map[string]net.Conn
	listener         net.Listener
	port             string
	protocol         string
	rateLimiter      RateLimiter
}

/*
Handles the command sent by the client if it is a command.

When the function returns true, means that the message should
be replicated to all the connected clients. If it returns false
means that it was a command and the message should not be replicated.
*/
func (s *Server) handleCommand(message []byte, client net.Conn) bool {
	for commandName, commandFn := range Commands {
		if string(message) != commandName {
			continue
		}

		commandFn(message, client, s)
		return false
	}

	return true
}

// Replicates the message to all the connected clients but not to the sender
func (s Server) replicateMessage(message []byte, sender net.Conn) {
	for clientId, client := range s.connectedClients {
		if client == sender {
			continue
		}

		fmt.Fprintf(client, "[Client %v]: %v", clientId, string(message))
	}
}

/*
Handles the client connection, reads the messages and replicates them to
all the connected clients. It also helps handling the commands.
*/
func (s Server) handleClient(client net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := client.Read(buffer)

		// When the client disconnects, it will return an error
		// so instead of doing log.Fatal, i return, to prevent
		// crashing the server and the connected clients.
		if err != nil {
			return
		}

		if s.rateLimiter.IsRateLimited(client) {
			continue
		}

		message := buffer[:n] // Cut off the trailing bytes for correct checking

		if !s.handleCommand(message, client) {
			continue
		}

		s.replicateMessage(message, client)
		s.rateLimiter.IncrementMessageCount(client)
	}
}

// Handles the incoming clients
func (s *Server) handleAcceptIncomingClients() {
	for {
		client, err := s.listener.Accept()

		if err != nil {
			log.Fatal("Couldn't accept client connection", err)
		}

		log.Printf("[%v] has joined the chat\n", client.RemoteAddr())

		s.connectedClients[client.RemoteAddr().String()] = client
		s.rateLimiter.AddToQueue(client)
		go s.handleClient(client)
	}
}

////////////////////////////////////////////////////////////////////////////////////

func (s *Server) Start() {
	log.Printf("Listening on port %v with protocol %v\n", s.port, s.protocol)

	listener, err := net.Listen(s.protocol, "localhost"+s.port)

	if err != nil {
		log.Fatal("Couldn't start the server", err)
	}

	s.listener = listener

	go s.handleAcceptIncomingClients()
}

func (s Server) Stop() {
	log.Println("Stopping the server")

	for _, client := range s.connectedClients {
		client.Close()
	}

	s.listener.Close()
}

func NewServer() *Server {
	return &Server{
		port:             config.ServerPort,
		protocol:         config.Protocol,
		connectedClients: map[string]net.Conn{},
		rateLimiter:      NewRateLimiter(),
	}
}
