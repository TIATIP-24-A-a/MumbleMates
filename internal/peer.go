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

func (c *ChatNode) GetAddress() (string, error) {
	peerInfo := peerstore.AddrInfo{
		ID:    c.Node.ID(),
		Addrs: c.Node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return "", err
	}
	return addrs[0].String(), nil
}

func (c *ChatNode) ConnectToPeer(pi peerstore.AddrInfo) error {
	if err := c.Node.Connect(context.Background(), pi); err != nil {
		return err
	}

	stream, err := c.Node.NewStream(context.Background(), pi.ID, PROTOCOL_ID)

	if err != nil {
		return err
	}

	c.peers[pi.ID] = stream

	return nil
}

func (c *ChatNode) HandleStream(stream network.Stream) {
	defer stream.Close()
	buf := bufio.NewReader(stream)

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

func (c *ChatNode) Start() error {
	c.Node.SetStreamHandler(PROTOCOL_ID, c.HandleStream)

	if err := setupMDNSDiscovery(c); err != nil {
		return err
	}

	return nil
}

func (c *ChatNode) Stop() error {
	if err := c.Node.Close(); err != nil {
		return err
	}

	close(c.Events)

	return nil
}

func (c *ChatNode) HandlePeerFound(pi peerstore.AddrInfo) {
	isSelf := pi.ID == c.Node.ID()
	if isSelf {
		return
	}

	isConnected := c.peers[pi.ID] != nil
	if isConnected {
		return
	}

	if err := c.ConnectToPeer(pi); err != nil {
		fmt.Println("error connecting to peer:", err)
	}
}
