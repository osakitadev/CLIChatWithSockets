package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Osaka/chat-with-sockets/client"
)

func main() {
	chatLeaveChannel := make(chan os.Signal, 1)
	signal.Notify(chatLeaveChannel, os.Interrupt)

	client := client.NewClient()

	client.Connect()
	<-chatLeaveChannel
	client.Disconnect()

	fmt.Println("Chat has ended.")
}
