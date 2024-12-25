package config

import "time"

const (
	// If the user sends RateLimitMessagesThreshold messages in less than 1 second, the server will ratelimit the user
	RateLimitMessagesThreshold = 5
	// Seconds to wait before the user can send messages again
	RateLimitTimeoutSeconds = 4 * time.Second
	// Message to send to the user when they are ratelimited
	RateLimitTimeoutMessage = "[SERVER]: You are sending messages too fast. You got ratelimited"
)
