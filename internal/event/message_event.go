package event

import (
	"encoding/json"

	"github.com/google/uuid"
)

type MessagePayload struct {
	Message string `json:"message"`
}

func NewMessage(name string, message string) *BaseEvent {
	payload, _ := json.Marshal(MessagePayload{Message: message})

	return &BaseEvent{
		ID:      uuid.New(),
		Type:    "message",
		Name:    name,
		Payload: payload,
	}
}
