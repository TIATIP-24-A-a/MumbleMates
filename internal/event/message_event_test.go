package event

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	name := "Albert"
	message := "Greetings!"

	messageEvent := NewMessage(name, message)

	if messageEvent.ID.String() == "" {
		t.Errorf("Expected ID to be set, got empty string")
	}

	if messageEvent.Type != "message" {
		t.Errorf("Expected type to be 'message', got %s", messageEvent.Type)
	}

	if messageEvent.Payload != message {
		t.Errorf("Expected message to be %s, got %s", message, messageEvent.Payload)
	}
}
