package queueworker

import (
	stan "github.com/nats-io/stan.go"
	"sync"
	"time"
)

// NATSQueue queue for work
type NATSQueue struct {
	nc             stan.Conn
	ncMutex        *sync.RWMutex
	maxReconnect   int
	reconnectDelay time.Duration

	// ClientID for NATS Streaming
	ClientID string

	// ClusterID in NATS Streaming
	ClusterID string

	// NATSURL URL to connect to NATS
	NATSURL string

	// Topic to respond to
	Topic string
}
