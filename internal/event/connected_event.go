package event

import (
	"github.com/TIATIP-24-A-a/MumbleMates/internal/peer"
)

// NewConnection creates a new connection event
func NewConnection(name string) *Event {
	return newEvent(
		ConnectEventType,
		peer.PeerInfo{
			Name: name,
		},
		nil,
	)
}
