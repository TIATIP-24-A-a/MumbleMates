package event

import (
	"time"

	"github.com/TIATIP-24-A-a/MumbleMates/internal/peer"
	"github.com/google/uuid"
)

const (
	MessageEventType = "message"
	ConnectEventType = "connect"
)

const (
	ConnectedStatus    = "connected"
	DisconnectedStatus = "disconnected"
	Idle               = "idle"
	TypingStatus       = "typing"
)

type Event struct {
	ID        uuid.UUID     `json:"id"`
	Type      string        `json:"type"`
	Timestamp time.Time     `json:"timestamp"`
	PeerInfo  peer.PeerInfo `json:"peerInfo"`
	Payload   interface{}   `json:"payload"`
}

func newEvent(eventType string, peerInfo peer.PeerInfo, payload any) *Event {
	return &Event{
		ID:        uuid.New(),
		Type:      eventType,
		Timestamp: time.Now(),
		PeerInfo:  peerInfo,
		Payload:   payload,
	}
}
