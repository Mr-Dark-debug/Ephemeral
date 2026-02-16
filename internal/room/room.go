package room

import (
	"ephemeral/internal/protocol"
	"sync"
)

type Room struct {
	Name      string
	Encrypted bool
	Key       []byte
	Messages  []protocol.Envelope
	Peers     map[string]bool
	mu        sync.RWMutex
}

type Manager struct {
	Rooms       map[string]*Room
	CurrentRoom string
	Nick        string
	PeerID      string
	mu          sync.RWMutex
}

func NewManager(nick, peerID string) *Manager {
	m := &Manager{
		Rooms:       make(map[string]*Room),
		CurrentRoom: "global",
		Nick:        nick,
		PeerID:      peerID,
	}
	m.Rooms["global"] = &Room{
		Name:     "global",
		Messages: make([]protocol.Envelope, 0),
		Peers:    make(map[string]bool),
	}
	return m
}

func (m *Manager) AddMessage(env protocol.Envelope) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	r, exists := m.Rooms[env.Room]
	if !exists {
		r = &Room{
			Name:     env.Room,
			Messages: make([]protocol.Envelope, 0),
			Peers:    make(map[string]bool),
		}
		m.Rooms[env.Room] = r
	}
	
	r.mu.Lock()
	r.Messages = append(r.Messages, env)
	if len(r.Messages) > 1000 {
		r.Messages = r.Messages[len(r.Messages)-1000:]
	}
	r.Peers[env.From] = true
	r.mu.Unlock()
}

func (m *Manager) GetMessages(roomName string) []protocol.Envelope {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	r, exists := m.Rooms[roomName]
	if !exists {
		return nil
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	msgs := make([]protocol.Envelope, len(r.Messages))
	copy(msgs, r.Messages)
	return msgs
}

func (m *Manager) Join(roomName string, encrypted bool, key []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.Rooms[roomName]; !exists {
		m.Rooms[roomName] = &Room{
			Name:      roomName,
			Encrypted: encrypted,
			Key:       key,
			Messages:  make([]protocol.Envelope, 0),
			Peers:     make(map[string]bool),
		}
	}
	m.CurrentRoom = roomName
}

func (m *Manager) Current() *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Rooms[m.CurrentRoom]
}
