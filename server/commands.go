package server

import (
	"fmt"
	"net"
)

type Command map[string]func(message []byte, client net.Conn, s *Server)

var Commands Command = Command{
	"/quit": func(message []byte, client net.Conn, s *Server) {
		s.replicateMessage([]byte(fmt.Sprintf("%v has left the chat\n", client.RemoteAddr())), client)
		delete(s.connectedClients, client.RemoteAddr().String())
		client.Close()
	},

	"/online": func(message []byte, client net.Conn, s *Server) {
		fmt.Fprintf(client, "[COMMANDS]: There are %v clients connected", len(s.connectedClients))
	},

	"/help": func(message []byte, client net.Conn, s *Server) {
		// Idk if I should do it like this, but who cares
		fmt.Fprintln(client, `
/exit - Exits the chat
/online - Shows the number of clients connected
/help - Shows this message`)
	},
}
