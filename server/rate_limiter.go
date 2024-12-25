package server

import (
	"log"
	"net"
	"time"

	"github.com/Osaka/chat-with-sockets/config"
)

type RateLimitedClient struct {
	lastMessageTime time.Time
	messageCount    int
	rateLimited     bool
}

type RateLimiter struct {
	rateLimitedClients map[string]*RateLimitedClient
}

func (r *RateLimiter) handleRateLimitReset() {
	for {
		time.Sleep(time.Second)

		for _, client := range r.rateLimitedClients {
			client.messageCount = 0
		}
	}
}

func (r *RateLimiter) IncrementMessageCount(client net.Conn) {
	clientMap := r.rateLimitedClients[client.RemoteAddr().String()]
	clientMap.messageCount++

	if clientMap.messageCount >= config.RateLimitMessagesThreshold && time.Since(clientMap.lastMessageTime) < 1*time.Second {
		log.Println("Rate limiting client", client.RemoteAddr().String())

		clientMap.rateLimited = true
		r.SendRateLimitMessage(client)
		time.Sleep(config.RateLimitTimeoutSeconds)
		clientMap.messageCount = 0
		clientMap.rateLimited = false
	}

	clientMap.lastMessageTime = time.Now()
}

func (r *RateLimiter) IsRateLimited(client net.Conn) bool {
	return r.rateLimitedClients[client.RemoteAddr().String()].rateLimited
}

func (r RateLimiter) SendRateLimitMessage(client net.Conn) {
	client.Write([]byte(config.RateLimitTimeoutMessage))
}

func (r *RateLimiter) AddToQueue(client net.Conn) {
	r.rateLimitedClients[client.RemoteAddr().String()] = &RateLimitedClient{
		lastMessageTime: time.Now(),
		messageCount:    0,
		rateLimited:     false,
	}

	go r.handleRateLimitReset()
}

func NewRateLimiter() RateLimiter {
	return RateLimiter{
		rateLimitedClients: make(map[string]*RateLimitedClient),
	}
}
