package main

import (
	"os"
	"os/signal"

	"github.com/Osaka/chat-with-sockets/server"
)

func main() {
	closeServerChannel := make(chan os.Signal, 1)
	server := server.NewServer()

	server.Start()
	signal.Notify(closeServerChannel, os.Interrupt)
	<-closeServerChannel
	server.Stop()
}
