package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	chat "github.com/TIATIP-24-A-a/MumbleMates/internal"
	"github.com/TIATIP-24-A-a/MumbleMates/internal/event"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

const gap = "\n\n"

var program *tea.Program

type (
	errMsg error
)

type model struct {
	viewport      viewport.Model
	messages      []Message
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	receiverStyle lipgloss.Style
	eventStyle    lipgloss.Style

	chatNode chat.ChatNode
	err      error
}

type Message struct {
	id      uuid.UUID
	content string
	sender  string
	time    time.Time
}

func initialModel() model {
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

	// TODO: Prompt for name
	chatNode, err := chat.NewChatNode("User")
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
		chatNode:      *chatNode,
		err:           nil,
	}
}

func (m model) Init() tea.Cmd {
	go listenForStreamEvents(m)

	return textarea.Blink
}

func listenForStreamEvents(m model) {
	for msg := range m.chatNode.Events {
		program.Send(msg)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			var messages []string
			for _, m := range m.messages {
				messages = append(messages, m.content+"\n")
			}
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() == "" {
				break
			}

			message := Message{
				id:      uuid.New(),
				content: m.senderStyle.Render("You: ") + m.textarea.Value(),
				sender:  "You",
				time:    time.Now(),
			}
			messageEvent := event.NewMessage("User", m.textarea.Value())
			m.messages = append(m.messages, message)
			m.chatNode.SendEvent(*messageEvent)

			m.RefreshMessagesView()
			m.textarea.Reset()
		}
	case event.Event:
		switch msg.Type {
		case "message":
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

			m.RefreshMessagesView()
		}

	case tea.QuitMsg:
		m.chatNode.Stop()
		return m, nil
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		m.chatNode.Stop()
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}

func (m model) RefreshMessagesView() {
	var messages []string
	for _, m := range m.messages {
		messages = append(messages, m.content)
	}
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(messages, "\n")))
	m.viewport.GotoBottom()
}

func Run() (tea.Model, error) {
	program = tea.NewProgram(initialModel())
	return program.Run()
}
