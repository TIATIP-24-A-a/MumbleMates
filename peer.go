package main

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"

    "github.com/libp2p/go-libp2p"
    peerstore "github.com/libp2p/go-libp2p/core/peer"
    "github.com/libp2p/go-libp2p/core/protocol"
    network "github.com/libp2p/go-libp2p/core/network"
    multiaddr "github.com/multiformats/go-multiaddr"
)

func handleStream(stream network.Stream) {
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
        // Print only messages received from remote peer
        fmt.Println("You:", strings.TrimSpace(message)) // display as "You: <message>"
    }
}

func handleUserInput(stream network.Stream) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Me: ")
        message, _ := reader.ReadString('\n')
        message = strings.TrimSpace(message)
        if message == "" {
            continue
        }

        // Send the typed message to the remote peer over the stream
        if stream != nil {
            _, err := stream.Write([]byte(message + "\n"))
            if err != nil {
                fmt.Println("error writing to stream:", err)
                // Handle stream reset or closing
                if err.Error() == "stream reset" {
                    fmt.Println("stream reset detected, closing stream.")
                    stream.Close()
                    return
                }
            } else {
                // After sending, print the message sent
                fmt.Println("Me:", message) // display as "Me: <message>"
            }
        }
    }
}

func main() {
    // start a libp2p node that listens on a random local TCP port
    node, err := libp2p.New(
        libp2p.ListenAddrStrings(),
    )
    if err != nil {
        panic(err)
    }

    // print the node's PeerInfo in multiaddr format
    peerInfo := peerstore.AddrInfo{
        ID:    node.ID(),
        Addrs: node.Addrs(),
    }
    addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
    if err != nil {
        panic(err)
    }
    fmt.Println("libp2p node address:", addrs[0])

    var stream network.Stream

    // if a remote peer has been passed on the command line, connect to it
    if len(os.Args) > 1 {
        addr, err := multiaddr.NewMultiaddr(os.Args[1])
        if err != nil {
            panic(err)
        }
        peer, err := peerstore.AddrInfoFromP2pAddr(addr)
        if err != nil {
            panic(err)
        }
        if err := node.Connect(context.Background(), *peer); err != nil {
            panic(err)
        }
        fmt.Println("connected to", addr)

        // open a stream to the remote peer
        stream, err = node.NewStream(context.Background(), peer.ID, protocol.ID("/chat/1.0.0"))
        if err != nil {
            panic(err)
        }
    } else {
        // handle incoming streams
        node.SetStreamHandler(protocol.ID("/chat/1.0.0"), handleStream)
    }

    // Start the user input handler in a separate goroutine
    go handleUserInput(stream)

    // wait for a SIGINT or SIGTERM signal
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    <-ch
    fmt.Println("Received signal, shutting down...")

    // shut the node down
    if err := node.Close(); err != nil {
        panic(err)
    }
}
