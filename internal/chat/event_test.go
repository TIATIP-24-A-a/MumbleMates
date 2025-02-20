package chat

import (
	"encoding/json"
	"testing"
)

func TestNewMessage(t *testing.T) {
	name := "Albert"
	message := "Greetings!"

	messageEvent := NewMessage(name, message)

	if messageEvent.GetID().String() == "" {
		t.Errorf("Expected ID to be set, got empty string")
	}

	if messageEvent.GetType() != "message" {
		t.Errorf("Expected type to be 'message', got %s", messageEvent.GetType())
	}

	if messageEvent.GetName() != name {
		t.Errorf("Expected name to be %s, got %s", name, messageEvent.GetName())
	}

	var payload MessagePayload
	if err := json.Unmarshal(messageEvent.Payload, &payload); err != nil {
		t.Errorf("Failed to unmarshal payload: %s", err)
	}

	if payload.Message != message {
		t.Errorf("Expected message to be %s, got %s", message, payload.Message)
	}
}
