package tui

import (
	"ephemeral/internal/config"
	"ephemeral/internal/discovery"
	"ephemeral/internal/protocol"
	"ephemeral/internal/room"
	"ephemeral/internal/transport"
	"fmt"
	"net"
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

	metadataStyle = lipgloss.NewStyle().
			Foreground(subText).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(accentGreen).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentPurple).
			Padding(0, 1)

	suggestionStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#21262d")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1)

	selectedSuggestionStyle = lipgloss.NewStyle().
					Background(accentPurple).
					Foreground(lipgloss.Color("#ffffff")).
					Padding(0, 1)

	suggestionBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentPurple).
				Background(lipgloss.Color("#21262d"))

	keycapStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#30363d")).
			Padding(0, 1).
			MarginRight(1)

	keycapDescStyle = lipgloss.NewStyle().
			Foreground(subText).
			MarginRight(2)

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

var availableCommands = []string{"/join", "/nick", "/clear", "/help", "/ip"}

type model struct {
	cfg       *config.Config
	roomMgr   *room.Manager
	transport *transport.Transport
	discovery *discovery.Service

	viewport  viewport.Model
	textInput textinput.Model

	suggestionMenuOpen   bool
	selectedCommandIndex int
	filteredCommands     []string

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
		if m.suggestionMenuOpen {
			switch msg.Type {
			case tea.KeyUp:
				m.selectedCommandIndex--
				if m.selectedCommandIndex < 0 {
					m.selectedCommandIndex = len(m.filteredCommands) - 1
				}
				return m, nil
			case tea.KeyDown:
				m.selectedCommandIndex++
				if m.selectedCommandIndex >= len(m.filteredCommands) {
					m.selectedCommandIndex = 0
				}
				return m, nil
			case tea.KeyTab, tea.KeyEnter:
				if len(m.filteredCommands) > 0 {
					m.textInput.SetValue(m.filteredCommands[m.selectedCommandIndex] + " ")
					m.textInput.CursorEnd()
					m.suggestionMenuOpen = false
					return m, nil
				}
			case tea.KeyEsc:
				m.suggestionMenuOpen = false
				return m, nil
			}
		}

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
		verticalMargin := 11

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

	// Update suggestions
	val := m.textInput.Value()
	if strings.HasPrefix(val, "/") && !strings.Contains(val, " ") {
		m.filteredCommands = nil
		for _, cmd := range availableCommands {
			if strings.HasPrefix(cmd, val) {
				m.filteredCommands = append(m.filteredCommands, cmd)
			}
		}
		m.suggestionMenuOpen = len(m.filteredCommands) > 0
		if m.selectedCommandIndex >= len(m.filteredCommands) {
			m.selectedCommandIndex = 0
		}
	} else {
		m.suggestionMenuOpen = false
	}

	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing TUI..."
	}

	header := m.renderHeader()
	suggestions := m.renderSuggestions()
	footer := m.renderFooter()

	inputView := inputStyle.Width(m.width - 2).Render(m.textInput.View())

	// Position suggestions above input
	var mainContent string
	if m.suggestionMenuOpen {
		// We might need to adjust viewport height if suggestions overlap too much,
		// but for now let's just stack them.
		mainContent = fmt.Sprintf("%s\n%s\n%s", m.viewport.View(), suggestions, inputView)
	} else {
		mainContent = fmt.Sprintf("%s\n\n%s", m.viewport.View(), inputView)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		header,
		mainContent,
		footer,
	)
}

func (m model) renderHeader() string {
	left := headerStyle.Render(" EPHEMERAL ")

	if m.width < 60 {
		right := statusStyle.Render("🟢 Connected")
		w := m.width - lipgloss.Width(left)
		if w < 0 {
			return left
		}
		return lipgloss.JoinHorizontal(lipgloss.Center, left, lipgloss.PlaceHorizontal(w, lipgloss.Right, right))
	}

	right := metadataStyle.Render(fmt.Sprintf("v1.0.0 | 👥 %d Online | 🔌 Connected", m.discovery.OnlineCount()))
	w := m.width - lipgloss.Width(left)
	return lipgloss.JoinHorizontal(lipgloss.Center, left, lipgloss.PlaceHorizontal(w, lipgloss.Right, right))
}

func (m model) renderSuggestions() string {
	if !m.suggestionMenuOpen || len(m.filteredCommands) == 0 {
		return ""
	}

	var items []string
	for i, cmd := range m.filteredCommands {
		if i == m.selectedCommandIndex {
			items = append(items, selectedSuggestionStyle.Render(" > "+cmd+" "))
		} else {
			items = append(items, suggestionStyle.Render("   "+cmd+" "))
		}
	}

	return suggestionBorderStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

func (m model) renderFooter() string {
	esc := keycapStyle.Render("ESC") + keycapDescStyle.Render("Quit")
	ctrlL := keycapStyle.Render("CTRL+L") + keycapDescStyle.Render("Clear")
	tab := keycapStyle.Render("TAB") + keycapDescStyle.Render("Autocomplete")

	return lipgloss.JoinHorizontal(lipgloss.Top, esc, ctrlL, tab)
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
		case "/clear":
			m.roomMgr.Current().Messages = nil
			m.viewport.SetContent(m.renderMessages())
		case "/help":
			m.roomMgr.AddMessage(protocol.NewEnvelope("sys", "system", "System", m.roomMgr.CurrentRoom, protocol.TypeChat, "Available commands: /join <room>, /nick <name>, /clear, /help, /ip"))
			m.viewport.SetContent(m.renderMessages())
			m.viewport.GotoBottom()
		case "/ip":
			addrs, _ := net.InterfaceAddrs()
			var ip string
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
						break
					}
				}
			}
			m.roomMgr.AddMessage(protocol.NewEnvelope("sys", "system", "System", m.roomMgr.CurrentRoom, protocol.TypeChat, fmt.Sprintf("Your Local IP: %s", ip)))
			m.viewport.SetContent(m.renderMessages())
			m.viewport.GotoBottom()
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
