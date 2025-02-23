package chat

import (
	"testing"
)

func TestCreateNode(t *testing.T) {
	node, err := NewChatNode("test")

	if err != nil {
		t.Errorf("Error creating node: %s", err)
	}

	if node == nil {
		t.Errorf("Node shoud not be nil")
	} else if node.name != "test" {
		t.Errorf("Node name should be 'test', got %s", node.name)
	}
}
