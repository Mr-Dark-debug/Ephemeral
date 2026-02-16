package discovery

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/grandcat/zeroconf"
)

const (
	MDNSServiceType = "_meshroom._tcp"
	MDNSDomain      = "local."
	UDPBroadcastPort = 9998
)

type Peer struct {
	ID   string
	Nick string
	IP   net.IP
	Port int
}

type Service struct {
	Nick      string
	Port      int
	PeerID    string
	MDNSEnabled bool
	UDPEnabled  bool
	
	peers     map[string]Peer
	newPeerCh chan Peer
	
	mdnsServer *zeroconf.Server
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewService(nick, peerID string, port int, mdns, udp bool) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		Nick:        nick,
		Port:        port,
		PeerID:      peerID,
		MDNSEnabled: mdns,
		UDPEnabled:  udp,
		peers:       make(map[string]Peer),
		newPeerCh:   make(chan Peer, 10),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (s *Service) Start() error {
	if s.MDNSEnabled {
		if err := s.startMDNS(); err != nil {
			return err
		}
	}
	if s.UDPEnabled {
		go s.startUDPListener()
		go s.startUDPBroadcaster()
	}
	return nil
}

func (s *Service) Stop() {
	s.cancel()
	if s.mdnsServer != nil {
		s.mdnsServer.Shutdown()
	}
}

func (s *Service) Peers() <-chan Peer {
	return s.newPeerCh
}

func (s *Service) startMDNS() error {
	meta := []string{
		"nick=" + s.Nick,
		"id=" + s.PeerID,
	}
	server, err := zeroconf.Register(
		s.Nick,
		MDNSServiceType,
		MDNSDomain,
		s.Port,
		meta,
		nil,
	)
	if err != nil {
		return err
	}
	s.mdnsServer = server

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			isSelf := false
			for _, field := range entry.Text {
				if field == "id="+s.PeerID {
					isSelf = true
					break
				}
			}
			if isSelf {
				continue
			}

			var nick string
			var id string
			for _, field := range entry.Text {
				if len(field) > 5 && field[:5] == "nick=" {
					nick = field[5:]
				}
				if len(field) > 3 && field[:3] == "id=" {
					id = field[3:]
				}
			}

			if id != "" && len(entry.AddrIPv4) > 0 {
				peer := Peer{
					ID:   id,
					Nick: nick,
					IP:   entry.AddrIPv4[0],
					Port: entry.Port,
				}
				s.handleFoundPeer(peer)
			}
		}
	}(entries)

	if err := resolver.Browse(s.ctx, MDNSServiceType, MDNSDomain, entries); err != nil {
		return err
	}

	return nil
}

type UDPDiscoveryPacket struct {
	Cmd  string `json:"cmd"`
	Nick string `json:"nick"`
	ID   string `json:"id"`
	Port int    `json:"port"`
}

func (s *Service) startUDPListener() {
	addr := &net.UDPAddr{
		Port: UDPBroadcastPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("UDP listen error: %v", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			var pkt UDPDiscoveryPacket
			if err := json.Unmarshal(buf[:n], &pkt); err != nil {
				continue
			}

			if pkt.ID == s.PeerID {
				continue
			}

			if pkt.Cmd == "DISCOVER" {
				peer := Peer{
					ID:   pkt.ID,
					Nick: pkt.Nick,
					IP:   remoteAddr.IP,
					Port: pkt.Port,
				}
				s.handleFoundPeer(peer)
			}
		}
	}
}

func (s *Service) startUDPBroadcaster() {
	addr := &net.UDPAddr{
		Port: UDPBroadcastPort,
		IP:   net.ParseIP("255.255.255.255"),
	}
	
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("UDP dial error: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	pkt := UDPDiscoveryPacket{
		Cmd:  "DISCOVER",
		Nick: s.Nick,
		ID:   s.PeerID,
		Port: s.Port,
	}
	data, _ := json.Marshal(pkt)

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			conn.Write(data)
		}
	}
}

func (s *Service) handleFoundPeer(p Peer) {
	if _, exists := s.peers[p.ID]; !exists {
		s.peers[p.ID] = p
		select {
		case s.newPeerCh <- p:
		default:
		}
	}
}
