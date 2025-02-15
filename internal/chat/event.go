package chat

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Event interface {
	GetID() uuid.UUID
	GetType() string
	GetName() string
	ToJSON() ([]byte, error)
}

type BaseEvent struct {
	ID   uuid.UUID
	Type string
	Name string
}

func (e *BaseEvent) GetID() uuid.UUID {
	return e.ID
}

func (e *BaseEvent) GetType() string {
	return e.Type
}

func (e *BaseEvent) GetName() string {
	return e.Name
}

func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

type MessageEvent struct {
	Event   BaseEvent
	Message string
}

func NewMessage(name, message string) *MessageEvent {
	return &MessageEvent{
		Event: BaseEvent{
			ID:   uuid.New(),
			Type: "message",
			Name: name,
		},
		Message: message,
	}
}
