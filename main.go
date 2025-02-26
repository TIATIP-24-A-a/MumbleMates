package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TIATIP-24-A-a/MumbleMates/cmd"
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
	if _, err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
