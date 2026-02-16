package main

import (
	"ephemeral/internal/config"
	"ephemeral/internal/discovery"
	"ephemeral/internal/room"
	"ephemeral/internal/transport"
	"ephemeral/internal/tui"
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

func main() {
	nick := flag.String("nick", "guest", "Your nickname")
	port := flag.Int("port", 9999, "Port to listen on (0 for random)")
	flag.Parse()

	cfg := config.Default()
	cfg.Nick = *nick
	cfg.Port = *port

	peerID := uuid.New().String()

	tr := transport.New(cfg.Port, peerID, cfg.Nick)
	if err := tr.Start(); err != nil {
		log.Fatalf("Failed to start transport: %v", err)
	}
	defer tr.Stop()
	
	cfg.Port = tr.Port
	
	disc := discovery.NewService(cfg.Nick, peerID, cfg.Port, true, true)
	if err := disc.Start(); err != nil {
		log.Fatalf("Failed to start discovery: %v", err)
	}
	defer disc.Stop()

	rm := room.NewManager(cfg.Nick, peerID)

	model := tui.InitialModel(cfg, rm, tr, disc)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
