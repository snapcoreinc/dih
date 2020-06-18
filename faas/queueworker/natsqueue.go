package queueworker

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"log"
	"sync"
	"time"
	"github.com/openfaas/nats-queue-worker/handler"
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

// CreateNATSQueue ready for asynchronous processing
func CreateNATSQueue(address string, port int, clusterName, channel string, clientConfig handler.NATSConfig) (*NATSQueue, error) {
	var err error
	natsURL := fmt.Sprintf("nats://%s:%d", address, port)
	log.Printf("Opening connection to %s\n", natsURL)

	clientID := clientConfig.GetClientID()

	// If 'channel' is empty, use the previous default.
	if channel == "" {
		channel = "faas-request"
	}

	queue1 := NATSQueue{
		ClientID:       clientID,
		ClusterID:      clusterName,
		NATSURL:        natsURL,
		Topic:          channel,
		maxReconnect:   clientConfig.GetMaxReconnect(),
		reconnectDelay: clientConfig.GetReconnectDelay(),
		ncMutex:        &sync.RWMutex{},
	}

	err = queue1.connect()

	return &queue1, err
}
