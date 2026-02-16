package protocol

import (
	"testing"
)

func TestEnvelopeJSON(t *testing.T) {
	env := NewEnvelope("id1", "peer1", "nick1", "global", TypeChat, "hello")
	
	data, err := env.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	
	env2, err := FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}
	
	if env.ID != env2.ID {
		t.Errorf("ID mismatch: %s != %s", env.ID, env2.ID)
	}
	if env.Payload != env2.Payload {
		t.Errorf("Payload mismatch")
	}
}
