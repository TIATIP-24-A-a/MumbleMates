package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	chat "github.com/TIATIP-24-A-a/MumbleMates/internal"
)

func askForName() string {
	const MAX_NAME_LENGTH = 20

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
	name := askForName()
	fmt.Println("Hello ", name, "👋")

	chatNode, err := chat.NewChatNode(name)
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
