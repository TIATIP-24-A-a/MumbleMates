package internal

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the name input view
type nameModel struct {
	nameInput textinput.Model
}

// InitialNameModel initializes the name input view
func InitialNameModel() nameModel {
	ti := textinput.New()
	ti.CharLimit = 20
	ti.Placeholder = "e.g. Bob the Builder"
	ti.Focus()
	return nameModel{nameInput: ti}
}

// tea.Model.Init interface implementation
func (m nameModel) Init() tea.Cmd {
	return textinput.Blink
}

// tea.Model.Update interface implementation
func (m nameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			return m.switchToChat()
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

// tea.Model.View interface implementation
func (m nameModel) View() string {
	return fmt.Sprintf("What's your name?\n\n%s\n\n", m.nameInput.View())
}

// Transitions to the chat view
// Returns commands to initialize the chat view
func (m *nameModel) switchToChat() (model, tea.Cmd) {
	model := initialModel(m.nameInput.Value())

	// When transitioning to the chat view, we want model
	// to update its view based of the window size
	cmd := tea.Batch(tea.EnterAltScreen, tea.ClearScreen, model.Init())

	return model, cmd
}
