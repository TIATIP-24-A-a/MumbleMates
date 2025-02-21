package chat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/TIATIP-24-A-a/MumbleMates/internal/event"
	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type ChatNode struct {
	Node  host.Host
	name  string
	peers map[peerstore.ID]network.Stream
}

const END_BYTE = byte('\n')
const SERVICE_TAG = "mumblemates-chat"
const PROTOCOL_ID = protocol.ID("/chat/1.0.0")

func NewChatNode(name string) (*ChatNode, error) {
	node, err := libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		return nil, err
	}
	return &ChatNode{
		Node:  node,
		peers: make(map[peerstore.ID]network.Stream),
		name:  name,
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

		var baseEvent event.BaseEvent
		err = json.Unmarshal([]byte(responseBytes), &baseEvent)
		if err != nil {
			fmt.Println("error unmarshalling event:", err)
			continue
		}

		// Print only messages received from remote peer
		if baseEvent.GetType() == "message" {
			var messagePayload event.MessagePayload
			err = json.Unmarshal(baseEvent.Payload, &messagePayload)
			if err != nil {
				fmt.Println("error unmarshalling message payload:", err)
				continue
			}

			fmt.Printf("%s: %s\n", baseEvent.Name, messagePayload.Message)
		} else {
			fmt.Println("Unknown event type: ", baseEvent.Type)
		}
	}
}

func (c *ChatNode) HandleUserInput() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(c.name, " (me): ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Send the typed message to the every peer
		for peerId, stream := range c.peers {
			encoder := json.NewEncoder(stream)
			message := event.NewMessage(c.name, message)
			err := encoder.Encode(message)
			if err != nil {
				fmt.Println("error writing to stream:", err)
				// Handle stream reset or closing

				if err.Error() == "write on closed stream" {
					fmt.Println("stream closed detected, closing stream.")
					stream.Close()
					delete(c.peers, peerId)
					return
				}

				if err.Error() == "stream reset" {
					fmt.Println("stream reset detected, closing stream.")
					stream.Close()
					delete(c.peers, peerId)
					return
				}
			}
		}
	}
}

func (c *ChatNode) Start() error {
	c.Node.SetStreamHandler(PROTOCOL_ID, c.HandleStream)

	// Start the user input handler in a separate goroutine
	go c.HandleUserInput()

	if err := setupMDNSDiscovery(c); err != nil {
		return err
	}

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := c.Node.Close(); err != nil {
		return err
	}

	return nil
}

func setupMDNSDiscovery(chatNode *ChatNode) error {
	service := mdns.NewMdnsService(chatNode.Node, SERVICE_TAG, &mdnsNotifee{c: *chatNode})
	return service.Start()
}

type mdnsNotifee struct {
	c ChatNode
}

func (n *mdnsNotifee) HandlePeerFound(pi peerstore.AddrInfo) {
	isSelf := pi.ID == n.c.Node.ID()
	if isSelf {
		return
	}

	isConnected := n.c.peers[pi.ID] != nil
	if isConnected {
		fmt.Println("Already connected to peer:", pi.ID)
		return
	}

	if err := n.c.ConnectToPeer(pi); err != nil {
		fmt.Println("error connecting to peer:", err)
	} else {
		fmt.Println("connected to peer:", pi.ID)
	}
}
