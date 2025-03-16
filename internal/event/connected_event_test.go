package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConntectedEvent(t *testing.T) {
	name := "Albert"

	connctedEvent := NewConnection(name)

	assert.NotNil(t, connctedEvent)
	assert.NotNil(t, connctedEvent.ID)
	assert.Nil(t, connctedEvent.Payload)
	assert.Equal(t, ConnectEventType, connctedEvent.Type)
}
