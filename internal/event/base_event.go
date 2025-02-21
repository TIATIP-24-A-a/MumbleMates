package event

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Event interface {
	GetID() uuid.UUID
	GetType() string
	GetName() string
	ToJSON() ([]byte, error)
	FromJSON(data []byte) error
}

type BaseEvent struct {
	ID      uuid.UUID       `json:"id"`
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
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

func (e *BaseEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
