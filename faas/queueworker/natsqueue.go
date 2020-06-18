package queueworker

import (
	"encoding/json"
	"fmt"
	stan "github.com/nats-io/stan.go"
	"github.com/openfaas/nats-queue-worker/handler"
	queue "github.com/snapcoreinc/dih/faas/gateway/queue"
	"log"
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

// Queue request for processing
func (q *NATSQueue) Queue(req *queue.Request) error {
	fmt.Printf("NatsQueue - submitting request: %s.\n", req.Function)

	out, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	q.ncMutex.RLock()
	nc := q.nc
	q.ncMutex.RUnlock()

	queueName := q.Topic
	if len(req.QueueName) > 0 {
		queueName = req.QueueName
	}

	return nc.Publish(queueName, out)
}

func (q *NATSQueue) connect() error {
	log.Printf("Connect: %s\n", q.NATSURL)

	nc, err := stan.Connect(
		q.ClusterID,
		q.ClientID,
		stan.NatsURL(q.NATSURL),
		stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
			log.Printf("Disconnected from %s\n", q.NATSURL)

			q.reconnect()
		}),
	)

	if err != nil {
		return err
	}

	q.ncMutex.Lock()
	q.nc = nc
	q.ncMutex.Unlock()

	return nil
}

func (q *NATSQueue) reconnect() {
	log.Printf("Reconnect\n")

	for i := 0; i < q.maxReconnect; i++ {
		time.Sleep(time.Duration(i) * q.reconnectDelay)

		if err := q.connect(); err == nil {
			log.Printf("Reconnecting (%d/%d) to %s. OK\n", i+1, q.maxReconnect, q.NATSURL)

			return
		}

		log.Printf("Reconnecting (%d/%d) to %s failed\n", i+1, q.maxReconnect, q.NATSURL)
	}

	log.Printf("Reached reconnection limit (%d) for %s\n", q.maxReconnect, q.NATSURL)

}
