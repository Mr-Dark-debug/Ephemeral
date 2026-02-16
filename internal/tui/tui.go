package tui

import (
	"ephemeral/internal/config"
	"ephemeral/internal/discovery"
	"ephemeral/internal/protocol"
	"ephemeral/internal/room"
	"ephemeral/internal/transport"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define colors and styles matching the "Cyber-Minimalism" vibe
var (
	accentPurple = lipgloss.Color("#7d57c1")
	accentGreen  = lipgloss.Color("#00ff41")
	subText      = lipgloss.Color("#8b949e")
	bgDark       = lipgloss.Color("#0d1117")
	cyan         = lipgloss.Color("#00ffff")

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(accentPurple).
			Padding(0, 1).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(accentGreen).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentPurple).
			Padding(0, 1)

	myMsgStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	peerMsgStyle = lipgloss.NewStyle().
			Foreground(accentGreen)

	systemStyle = lipgloss.NewStyle().
			Foreground(subText).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(subText).
			Padding(0, 1)
)

type model struct {
	cfg       *config.Config
	roomMgr   *room.Manager
	transport *transport.Transport
	discovery *discovery.Service

	viewport  viewport.Model
	textInput textinput.Model

	width  int
	height int
	ready  bool
}

func InitialModel(cfg *config.Config, rm *room.Manager, tr *transport.Transport, disc *discovery.Service) model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.Prompt = " > "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(accentGreen)
	ti.CharLimit = 1000

	return model{
		cfg:       cfg,
		roomMgr:   rm,
		transport: tr,
		discovery: disc,
		textInput: ti,
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
			if val := m.textInput.Value(); val != "" {
				m.sendMessage(val)
				m.textInput.SetValue("")
			}
		case tea.KeyCtrlL:
			m.roomMgr.Current().Messages = nil
			m.viewport.SetContent(m.renderMessages())
		}

	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 3
		inputHeight := 3
		verticalMargin := headerHeight + footerHeight + inputHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMargin)
			m.viewport.HighPerformanceRendering = false
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargin
		}

		m.textInput.Width = msg.Width - 6
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.SetContent(m.renderMessages())

	case protocol.Envelope:
		m.roomMgr.AddMessage(msg)
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForMessage(m.transport.Incoming()))

	case discovery.Peer:
		go m.transport.Connect(msg.ID, msg.IP.String(), msg.Port)
		cmds = append(cmds, waitForPeer(m.discovery.Peers()))
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing TUI..."
	}

	header := headerStyle.Render(" EPHEMERAL ") + statusStyle.Render("🟢 Connected")
	footer := helpStyle.Render("ESC: Quit • CTRL+L: Clear • /join: Room • /nick: Name")

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n%s",
		header,
		m.viewport.View(),
		inputStyle.Width(m.width-2).Render(m.textInput.View()),
		footer,
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
			}
		case "/nick":
			if len(parts) > 1 {
				m.roomMgr.Nick = parts[1]
			}
		}
		return
	}

	env := protocol.NewEnvelope(
		fmt.Sprintf("%s-%d", m.roomMgr.PeerID, time.Now().UnixNano()),
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
	var b strings.Builder
	for _, msg := range msgs {
		t := time.Unix(msg.TS, 0).Format("15:04")
		ts := systemStyle.Render(t)

		if msg.From == m.roomMgr.PeerID {
			nick := myMsgStyle.Render("You")
			line := fmt.Sprintf("%s %s: %s", ts, nick, msg.Payload)
			// Right align local messages
			padding := m.width - lipgloss.Width(line) - 4
			if padding > 0 {
				b.WriteString(strings.Repeat(" ", padding))
			}
			b.WriteString(line + "\n")
		} else {
			nick := peerMsgStyle.Render(msg.Nick)
			b.WriteString(fmt.Sprintf("%s %s: %s\n", ts, nick, msg.Payload))
		}
	}
	return b.String()
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
