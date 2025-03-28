package main

import (
	"fmt"
	"os"

	"github.com/TIATIP-24-A-a/MumbleMates/cmd"
)

func main() {
	if _, err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
