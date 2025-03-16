package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TIATIP-24-A-a/MumbleMates/internal/event"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

const gap = "\n\n"

type (
	errMsg func() error
)

// Main model for the TUI
type model struct {
	viewport      viewport.Model
	messages      []Message
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	receiverStyle lipgloss.Style
	eventStyle    lipgloss.Style
	infoStyle     lipgloss.Style

	name string

	chatNode ChatNode
	err      error
}

type Message struct {
	id      uuid.UUID
	content string
	sender  string
	time    time.Time
}

// Initial model
func initialModel(name string) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(20)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(20, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	chatNode, err := NewChatNode(name)
	if err != nil {
		log.Fatal(err)
	}

	chatNode.Start()

	return model{
		textarea:      ta,
		messages:      []Message{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		receiverStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		eventStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		infoStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		name:          name,
		chatNode:      *chatNode,
		err:           nil,
	}
}

// Waits for an event to be sent as a command
func waitForEvent(events chan event.Event) tea.Cmd {
	return func() tea.Msg {
		return <-events
	}
}

// tea.Model.Init interface implementation
func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		waitForEvent(m.chatNode.Events),
	)
}

// tea.Model.Update interface implementation
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd    tea.Cmd
		vpCmd    tea.Cmd
		eventCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	eventCmd = waitForEvent(m.chatNode.Events)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onResize(msg)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() == "" {
				break
			}

			formattedTime := m.eventStyle.Render(time.Now().Format("15:04"))
			receiver := m.senderStyle.Render(fmt.Sprintf("%s (You):", m.name))
			content := fmt.Sprintf("%s [%s] %s", receiver, formattedTime, m.textarea.Value())

			message := Message{
				id:      uuid.New(),
				content: content,
				sender:  m.name,
				time:    time.Now(),
			}
			messageEvent := event.NewMessage(m.name, m.textarea.Value())

			m.messages = append(m.messages, message)
			m.chatNode.SendEvent(*messageEvent)

			m.refreshMessagesView()
			m.textarea.Reset()
		}
	case event.Event:
		switch msg.Type {
		case event.MessageEventType:
			messageEvent := msg.Payload.(string)

			formattedTime := m.eventStyle.Render(msg.Timestamp.Format("15:04"))
			receiver := m.receiverStyle.Render(msg.PeerInfo.Name + ":")
			content := fmt.Sprintf("%s [%s] %s", receiver, formattedTime, messageEvent)

			m.messages = append(m.messages, Message{
				id:      msg.ID,
				content: content,
				sender:  msg.PeerInfo.Name,
				time:    msg.Timestamp,
			})

			m.refreshMessagesView()

		case event.ConnectEventType:
			message := fmt.Sprintf("[%s] %s has connected", msg.Timestamp.Format("15:04"), msg.PeerInfo.Name)
			styled := m.eventStyle.Render(message)
			m.messages = append(m.messages, Message{
				id:      msg.ID,
				content: styled,
				sender:  msg.PeerInfo.Name,
				time:    msg.Timestamp,
			})
			m.refreshMessagesView()
		}

	case tea.QuitMsg:
		m.chatNode.Stop()
		return m, nil
	// We handle errors just like any other message
	case errMsg:
		m.err = msg()
		m.chatNode.Stop()
		return m, nil
	}

	return m, tea.Batch(vpCmd, tiCmd, eventCmd)
}

// tea.Model.View interface implementation
func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}

// Updates the viewport content
func (m *model) onResize(msg tea.WindowSizeMsg) {
	m.viewport.Width = msg.Width
	m.textarea.SetWidth(msg.Width)
	m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
	m.refreshMessagesView()
}

// Rerenders the messages in the viewport
func (m *model) refreshMessagesView() {
	if len(m.messages) == 0 {
		return
	}

	var messages []string
	for _, m := range m.messages {
		messages = append(messages, m.content)
	}
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(messages, "\n")))
	m.viewport.GotoBottom()
}
