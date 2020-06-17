package queueworker

import (
	"time"
)

type NATSConfig interface {
	GetClientID() string
	GetMaxReconnect() int
	GetReconnectDelay() time.Duration
}
