package protocol

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	TypeChat     MessageType = "chat"
	TypePresence MessageType = "presence"
	TypeControl  MessageType = "control"
	TypeAck      MessageType = "ack"
)

type Envelope struct {
	V       int         `json:"v"`
	ID      string      `json:"id"`
	From    string      `json:"from"`
	Nick    string      `json:"nick"`
	Room    string      `json:"room"`
	TS      int64       `json:"ts"`
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
	Sig     string      `json:"sig,omitempty"`
}

func NewEnvelope(id, from, nick, room string, msgType MessageType, payload string) Envelope {
	return Envelope{
		V:       1,
		ID:      id,
		From:    from,
		Nick:    nick,
		Room:    room,
		TS:      time.Now().Unix(),
		Type:    msgType,
		Payload: payload,
	}
}

func (e Envelope) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func FromJSON(data []byte) (Envelope, error) {
	var e Envelope
	err := json.Unmarshal(data, &e)
	return e, err
}
