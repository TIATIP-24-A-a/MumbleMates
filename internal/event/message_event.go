package event

import (
	"github.com/TIATIP-24-A-a/MumbleMates/internal/peer"
)

type MessagePayload struct {
	Message string `json:"message"`
}

// NewMessage creates a new message event
func NewMessage(name string, message string) *Event {
	return newEvent(
		MessageEventType,
		peer.PeerInfo{
			Name: name,
		},
		message,
	)
}
