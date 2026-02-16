package tui

import (
	"ephemeral/internal/config"
	"ephemeral/internal/discovery"
	"ephemeral/internal/protocol"
	"ephemeral/internal/room"
	"ephemeral/internal/transport"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	cfg       *config.Config
	roomMgr   *room.Manager
	transport *transport.Transport
	discovery *discovery.Service
	
	viewport  viewport.Model
	textInput textinput.Model
	
	err       error
	width     int
	height    int
}

func InitialModel(cfg *config.Config, rm *room.Manager, tr *transport.Transport, disc *discovery.Service) model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to Ephemeral!\n")

	return model{
		cfg:       cfg,
		roomMgr:   rm,
		transport: tr,
		discovery: disc,
		textInput: ti,
		viewport:  vp,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		waitForMessage(m.transport.Incoming()),
		waitForPeer(m.discovery.Peers()),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textInput.Value() != "" {
				m.sendMessage(m.textInput.Value())
				m.textInput.SetValue("")
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 5
		m.textInput.Width = msg.Width

	case protocol.Envelope:
		m.roomMgr.AddMessage(msg)
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForMessage(m.transport.Incoming()))
	
	case discovery.Peer:
		go m.transport.Connect(msg.ID, msg.IP.String(), msg.Port)
		cmds = append(cmds, waitForPeer(m.discovery.Peers()))

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	
	cmds = append(cmds, tiCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		"Ephemera Chat",
		m.viewport.View(),
		m.textInput.View(),
	)
}

func (m *model) sendMessage(text string) {
	if strings.HasPrefix(text, "/") {
		parts := strings.Fields(text)
		cmd := parts[0]
		switch cmd {
		case "/join":
			if len(parts) > 1 {
				m.roomMgr.Join(parts[1], false, nil)
				m.viewport.SetContent(m.renderMessages())
			}
		case "/nick":
			if len(parts) > 1 {
				m.roomMgr.Nick = parts[1]
			}
		case "/quit":
			// handled by ctrl-c
		}
		return
	}

	env := protocol.NewEnvelope(
		fmt.Sprintf("%s-%d", m.cfg.Nick, 0),
		m.roomMgr.PeerID,
		m.roomMgr.Nick,
		m.roomMgr.CurrentRoom,
		protocol.TypeChat,
		text,
	)
	
	m.roomMgr.AddMessage(env)
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
	
	m.transport.Broadcast(env)
}

func (m *model) renderMessages() string {
	msgs := m.roomMgr.GetMessages(m.roomMgr.CurrentRoom)
	var s strings.Builder
	for _, msg := range msgs {
		s.WriteString(fmt.Sprintf("[%d] %s: %s\n", msg.TS, msg.Nick, msg.Payload))
	}
	return s.String()
}

func waitForMessage(ch <-chan protocol.Envelope) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func waitForPeer(ch <-chan discovery.Peer) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
