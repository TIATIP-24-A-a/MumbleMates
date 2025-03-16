package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNode(t *testing.T) {
	node, err := NewChatNode("test")

	if err != nil {
		t.Errorf("Error creating node: %s", err)
	}

	assert.NotNil(t, node)
	assert.Equal(t, "test", node.name)
}
