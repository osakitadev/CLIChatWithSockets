package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	Conn net.Conn
}

func (c Client) getUserMessageInput() []byte {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return []byte(strings.TrimSpace(scanner.Text())) // Remove any trailing whitespaces
}

func (c *Client) handleMessagesReceiving() {
	for {
		buff := make([]byte, 1024)

		_, err := c.Conn.Read(buff)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buff))
	}
}

func (c Client) handleMessagesSending() {
	for {
		message := c.getUserMessageInput()
		c.Conn.Write(message)
	}
}

func (c *Client) Connect() {
	fmt.Println(`Type /help to see the available commands.`)
	go c.handleMessagesReceiving()
	go c.handleMessagesSending()
}

func (c Client) Disconnect() {
	c.Conn.Close()
}

func NewClient() *Client {
	client, err := net.Dial("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Conn: client,
	}
}
