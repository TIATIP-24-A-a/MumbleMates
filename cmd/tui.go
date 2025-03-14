package cmd

import (
	chat "github.com/TIATIP-24-A-a/MumbleMates/internal"
	tea "github.com/charmbracelet/bubbletea"
)

func Run() (tea.Model, error) {
	return tea.NewProgram(chat.InitialNameModel()).Run()
}
