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
	multiaddr "github.com/multiformats/go-multiaddr"
)

type ChatNode struct {
	Node   host.Host
	Stream network.Stream
	name   string
}

const END_BYTE = byte('\n')

func NewChatNode(name string) (*ChatNode, error) {
	node, err := libp2p.New(libp2p.ListenAddrStrings())
	if err != nil {
		return nil, err
	}
	return &ChatNode{
		Node: node,
		name: name,
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
	c.Stream, err = c.Node.NewStream(context.Background(), peer.ID, protocol.ID("/chat/1.0.0"))
	return err
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
	encoder := json.NewEncoder(c.Stream)
	for {
		fmt.Print(c.name, " (me): ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Send the typed message to the remote peer over the stream
		if c.Stream != nil {
			message := event.NewMessage(c.name, message)
			err := encoder.Encode(message)
			if err != nil {
				fmt.Println("error writing to stream:", err)
				// Handle stream reset or closing
				if err.Error() == "stream reset" {
					fmt.Println("stream reset detected, closing stream.")
					c.Stream.Close()
					return
				}
			}
		}
	}
}

func (c *ChatNode) Start() error {
	// if a remote peer has been passed on the command line, connect to it
	if len(os.Args) > 1 {
		addr := os.Args[1]
		if err := c.ConnectToPeer(addr); err != nil {
			return err
		}
	} else {
		// handle incoming streams
		c.Node.SetStreamHandler(protocol.ID("/chat/1.0.0"), c.HandleStream)
	}

	// Start the user input handler in a separate goroutine
	go c.HandleUserInput()

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
