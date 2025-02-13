package main

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"

	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type ChatNode struct {
	Node   host.Host
	Stream network.Stream
}

func NewChatNode() (*ChatNode, error) {
	node, err := libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		return nil, err
	}
	return &ChatNode{
		Node: node,
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

func (c *ChatNode) ConnectToPeer(address string) error {
	addr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return err
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}
	if err := c.Node.Connect(context.Background(), *peer); err != nil {
		return err
	}
	fmt.Println("connected to", address)
	c.Stream, err = c.Node.NewStream(context.Background(), peer.ID, "/chat/1.0.0")
	return err
}

func (c *ChatNode) HandleStream(stream network.Stream) {
	defer stream.Close()
	buf := bufio.NewReader(stream)
	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("stream closed by remote peer")
				break
			} else {
				fmt.Println("error reading from stream:", err)
				break
			}
		}
		fmt.Println("You:", strings.TrimSpace(message))
	}
}

func TestNewChatNode(t *testing.T) {
	node, err := NewChatNode()
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.NotNil(t, node.Node)
}

func TestHandleStream(t *testing.T) {
	node1, err := NewChatNode()
	assert.NoError(t, err)

	node2, err := NewChatNode()
	assert.NoError(t, err)

	node1.Node.SetStreamHandler("/chat/1.0.0", node1.HandleStream)

	address1, err := node1.GetAddress()
	assert.NoError(t, err)

	err = node2.ConnectToPeer(address1)
	assert.NoError(t, err)

	message := "Hello, Node1!"
	_, err = node2.Stream.Write([]byte(message + "\n"))
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func TestGetAddress(t *testing.T) {
	node, err := NewChatNode()
	assert.NoError(t, err)
	address, err := node.GetAddress()
	assert.NoError(t, err)
	assert.NotEmpty(t, address)
}
