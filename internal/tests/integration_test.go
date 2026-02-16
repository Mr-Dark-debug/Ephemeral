package tests

import (
	"ephemeral/internal/protocol"
	"ephemeral/internal/transport"
	"testing"
	"time"
)

func TestTransportExchange(t *testing.T) {
	trA := transport.New(0, "peerA", "Alice")
	if err := trA.Start(); err != nil {
		t.Fatalf("Start A failed: %v", err)
	}
	defer trA.Stop()

	trB := transport.New(0, "peerB", "Bob")
	if err := trB.Start(); err != nil {
		t.Fatalf("Start B failed: %v", err)
	}
	defer trB.Stop()

	if trA.Port == 0 || trB.Port == 0 {
		t.Fatalf("Ports not assigned: A=%d, B=%d", trA.Port, trB.Port)
	}

	if err := trA.Connect("peerB", "127.0.0.1", trB.Port); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	env := protocol.NewEnvelope("msg1", "peerA", "Alice", "global", protocol.TypeChat, "Hello Bob")
	trA.Broadcast(env)

	timeout := time.After(2 * time.Second)
	found := false
	for {
		select {
		case msg := <-trB.Incoming():
			if msg.Type == protocol.TypeChat {
				if msg.Payload != "Hello Bob" {
					t.Errorf("Expected 'Hello Bob', got '%s'", msg.Payload)
				}
				if msg.From != "peerA" {
					t.Errorf("Expected From peerA, got %s", msg.From)
				}
				found = true
				goto CheckReverse
			}
		case <-timeout:
			t.Fatal("Timeout waiting for message on B")
		}
	}

CheckReverse:
	if !found {
		t.Fatal("Message not found")
	}

	env2 := protocol.NewEnvelope("msg2", "peerB", "Bob", "global", protocol.TypeChat, "Hi Alice")
	trB.Broadcast(env2)

	timeout = time.After(2 * time.Second)
	found = false
	for {
		select {
		case msg := <-trA.Incoming():
			if msg.Type == protocol.TypeChat {
				if msg.Payload != "Hi Alice" {
					t.Errorf("Expected 'Hi Alice', got '%s'", msg.Payload)
				}
				found = true
				goto Done
			}
		case <-timeout:
			t.Fatal("Timeout waiting for message on A")
		}
	}
Done:
	if !found {
		t.Fatal("Message not found on A")
	}
}
