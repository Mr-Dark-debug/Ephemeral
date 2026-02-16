package transport

import (
	"context"
	"encoding/json"
	"ephemeral/internal/protocol"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Transport struct {
	Port       int
	ID         string
	Nick       string
	
	listener   net.Listener
	peers      map[string]*PeerConn
	peersLock  sync.RWMutex
	
	incomingCh chan protocol.Envelope
	ctx        context.Context
	cancel     context.CancelFunc
}

type PeerConn struct {
	ID   string
	Conn net.Conn
	Enc  *json.Encoder
	Dec  *json.Decoder
}

func New(port int, id, nick string) *Transport {
	ctx, cancel := context.WithCancel(context.Background())
	return &Transport{
		Port:       port,
		ID:         id,
		Nick:       nick,
		peers:      make(map[string]*PeerConn),
		incomingCh: make(chan protocol.Envelope, 100),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (t *Transport) Start() error {
	addr := fmt.Sprintf(":%d", t.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	t.listener = ln
	
	if t.Port == 0 {
		if tcpAddr, ok := ln.Addr().(*net.TCPAddr); ok {
			t.Port = tcpAddr.Port
		}
	}
	
	go t.acceptLoop()
	return nil
}

func (t *Transport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.ctx.Done():
				return
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}
		go t.handleConn(conn, nil, "")
	}
}

func (t *Transport) handleConn(conn net.Conn, dec *json.Decoder, knownPeerID string) {
	if dec == nil {
		dec = json.NewDecoder(conn)
	}

	peerID := knownPeerID

	for {
		var env protocol.Envelope
		if err := dec.Decode(&env); err != nil {
			conn.Close()
			if peerID != "" {
				t.removePeer(peerID)
			}
			return
		}
		
		if peerID == "" {
			peerID = env.From
			enc := json.NewEncoder(conn)
			t.addPeer(peerID, conn, enc, dec)
		}
		
		t.incomingCh <- env
	}
}

func (t *Transport) Connect(peerID, ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	
	hello := protocol.NewEnvelope(
		fmt.Sprintf("%s-%d", t.ID, time.Now().UnixNano()),
		t.ID,
		t.Nick,
		"global",
		protocol.TypePresence,
		"",
	)
	
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	
	t.addPeer(peerID, conn, enc, dec)
	
	if err := enc.Encode(hello); err != nil {
		conn.Close()
		return err
	}
	
	go t.handleConn(conn, dec, peerID)
	
	return nil
}

func (t *Transport) addPeer(id string, conn net.Conn, enc *json.Encoder, dec *json.Decoder) {
	t.peersLock.Lock()
	defer t.peersLock.Unlock()
	t.peers[id] = &PeerConn{
		ID:   id,
		Conn: conn,
		Enc:  enc,
		Dec:  dec,
	}
}

func (t *Transport) removePeer(id string) {
	t.peersLock.Lock()
	defer t.peersLock.Unlock()
	delete(t.peers, id)
}

func (t *Transport) Broadcast(env protocol.Envelope) {
	t.peersLock.RLock()
	defer t.peersLock.RUnlock()
	
	for _, p := range t.peers {
		go func(p *PeerConn) {
			p.Enc.Encode(env)
		}(p)
	}
}

func (t *Transport) Incoming() <-chan protocol.Envelope {
	return t.incomingCh
}

func (t *Transport) Stop() {
	t.cancel()
	if t.listener != nil {
		t.listener.Close()
	}
	t.peersLock.Lock()
	defer t.peersLock.Unlock()
	for _, p := range t.peers {
		p.Conn.Close()
	}
}
