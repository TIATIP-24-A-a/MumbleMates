package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/TIATIP-24-A-a/MumbleMates/internal/event"
	"github.com/TIATIP-24-A-a/MumbleMates/internal/peer"
	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

type ChatNode struct {
	Node     host.Host
	name     string
	peers    map[peerstore.ID]network.Stream
	Events   chan event.Event
	PeerInfo peer.PeerInfo
}

const (
	END_BYTE    = byte('\n')
	SERVICE_TAG = "mumblemates-chat"
	PROTOCOL_ID = protocol.ID("/chat/1.0.0")
)

// Create a new chat node
func NewChatNode(name string) (*ChatNode, error) {
	node, err := libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		return nil, err
	}
	return &ChatNode{
		Node:   node,
		Events: make(chan event.Event),
		peers:  make(map[peerstore.ID]network.Stream),
		name:   name,
	}, nil
}

// Connect to a peer
func (c *ChatNode) connectToPeer(pi peerstore.AddrInfo) error {
	if err := c.Node.Connect(context.Background(), pi); err != nil {
		return err
	}

	stream, err := c.Node.NewStream(context.Background(), pi.ID, PROTOCOL_ID)

	if err != nil {
		return err
	}

	c.peers[pi.ID] = stream

	// Notify the peer that we have connected
	connectEvent := event.NewConnection(c.name)
	err = c.sendEventToPeer(pi.ID, *connectEvent)
	if err != nil {
		return err
	}

	return nil
}

// Handle incoming streams
func (c *ChatNode) handleStream(stream network.Stream) {
	defer stream.Close()
	buf := bufio.NewReader(stream)

	remoteId := stream.Conn().RemotePeer()
	if _, ok := c.peers[remoteId]; !ok {
		pi := peerstore.AddrInfo{
			ID:    remoteId,
			Addrs: []multiaddr.Multiaddr{stream.Conn().RemoteMultiaddr()},
		}
		if err := c.connectToPeer(pi); err != nil {
			fmt.Println("error connecting to peer:", err)
			return
		}
	}

	for {
		responseBytes, err := buf.ReadString(END_BYTE)
		if err != nil {
			fmt.Println("error reading from stream:", err)
			break
		}

		if !json.Valid([]byte(responseBytes)) {
			fmt.Println("Received invalid JSON: ", responseBytes)
			continue
		}

		var baseEvent event.Event
		err = json.Unmarshal([]byte(responseBytes), &baseEvent)
		if err != nil {
			fmt.Println("error unmarshalling event:", err)
			continue
		}

		c.Events <- baseEvent
	}
}

// Send an event to all connected peers
func (c *ChatNode) SendEvent(e event.Event) error {
	for _, stream := range c.peers {
		encoder := json.NewEncoder(stream)
		err := encoder.Encode(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// Send an event to a specific peer
func (c *ChatNode) sendEventToPeer(id peerstore.ID, e event.Event) error {
	stream := c.peers[id]
	if stream == nil {
		return fmt.Errorf("peer not found")
	}

	encoder := json.NewEncoder(stream)
	err := encoder.Encode(e)
	if err != nil {
		return err
	}

	return nil
}

// Start the chat node
func (c *ChatNode) Start() error {
	c.Node.SetStreamHandler(PROTOCOL_ID, c.handleStream)

	if err := setupMDNSDiscovery(c); err != nil {
		return err
	}

	return nil
}

// Stop the chat node
func (c *ChatNode) Stop() error {
	if err := c.Node.Close(); err != nil {
		return err
	}

	close(c.Events)

	return nil
}

// tea.Model.Init interface implementation
func (c *ChatNode) HandlePeerFound(pi peerstore.AddrInfo) {
	isSelf := pi.ID == c.Node.ID()
	if isSelf {
		return
	}

	isConnected := c.peers[pi.ID] != nil
	if isConnected {
		return
	}

	if err := c.connectToPeer(pi); err != nil {
		fmt.Println("error connecting to peer:", err)
	}
}
