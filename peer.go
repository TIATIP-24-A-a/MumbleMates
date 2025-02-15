package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/TIATIP-24-A-a/MumbleMates/internal/chat"
	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
	// "github.com/TIATIP-24-A-a/MumbleMates/internal/chat"
)

type ChatNode struct {
	Node   host.Host
	Stream network.Stream
}

var name string

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
	c.Stream, err = c.Node.NewStream(context.Background(), peer.ID, protocol.ID("/chat/1.0.0"))
	return err
}

func (c *ChatNode) HandleStream(stream network.Stream) {
	defer stream.Close()
	buf := bufio.NewReader(stream)

	for {
		responseBytes, err := buf.ReadString('\n')
		event := chat.BaseEvent{}
		fmt.Println(responseBytes)
		// json.Unmarshal(responseBytes, &event)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("stream closed by remote peer")
				break
			} else {
				fmt.Println("error reading from stream:", err)
				break
			}
		}

		// println(event)

		// Print only messages received from remote peer
		if event.Type == "message" {
			messageEvent := chat.MessageEvent{}
			// json.Unmarshal(responseBytes, &messageEvent)
			fmt.Println(messageEvent.Event.Name, ": ", messageEvent.Message)
		}
	}
}

func (c *ChatNode) HandleUserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(name, " (me): ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Send the typed message to the remote peer over the stream
		if c.Stream != nil {
			payload := chat.NewMessage(name, message)
			err := json.NewEncoder(c.Stream).Encode(payload)

			// _, err := c.Stream.Write(bytes.NewBuffer(payload))
			if err != nil {
				fmt.Println("error writing to stream:", err)
				// Handle stream reset or closing
				if err.Error() == "stream reset" {
					fmt.Println("stream reset detected, closing stream.")
					c.Stream.Close()
					return
				}
			} else {
				// After sending, print the message sent
				// fmt.Println(name, " (me): ", message) // New line for "Me" message
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

func askForName() string {
	MAX_NAME_LENGTH := 20

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name (max 20 length): ")
	name, _ := reader.ReadString('\n')

	trimmedName := strings.TrimSpace(name)

	if trimmedName == "" {
		fmt.Println("Name cannot be empty. Please try again.")
		return askForName()
	}

	if len(trimmedName) > MAX_NAME_LENGTH {
		fmt.Println("Name cannot be longer than 20 characters. Please try again.")
		return askForName()
	}

	return trimmedName
}

func main() {

	name = askForName()

	fmt.Println("Hello ", name, "👋")

	chatNode, err := NewChatNode()
	if err != nil {
		panic(err)
	}

	address, err := chatNode.GetAddress()
	if err != nil {
		panic(err)
	}

	fmt.Println("libp2p node address:")
	fmt.Println(address)
	fmt.Println()

	// Start the chat node
	if err := chatNode.Start(); err != nil {
		panic(err)
	}
}
