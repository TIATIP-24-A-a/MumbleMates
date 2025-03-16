package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	name := "Albert"
	message := "Greetings!"

	messageEvent := NewMessage(name, message)

	assert.NotNil(t, messageEvent)
	assert.NotNil(t, messageEvent.ID)
	assert.Equal(t, message, messageEvent.Payload)
	assert.Equal(t, "message", messageEvent.Type)
}
