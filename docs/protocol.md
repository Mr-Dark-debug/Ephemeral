# Protocol Specification - Ephemeral

## Wire Format
Ephemeral uses **JSON-Lines** over TCP. Each message is a single JSON object terminated by a newline character (`
`).

## Message Envelope
```json
{
  "v": 1,
  "id": "<uuid-v4>",
  "from": "<peer-id>",
  "nick": "<display-name>",
  "room": "global",
  "ts": 1670000000,
  "type": "chat",
  "payload": "<string or base64 encrypted data>",
  "sig": "<optional-signature>"
}
```

### Fields:
- `v`: Protocol version (integer).
- `id`: Unique message identifier for deduplication.
- `from`: Unique peer identifier of the sender.
- `nick`: Current nickname of the sender.
- `room`: The logical room name.
- `ts`: Unix timestamp.
- `type`: Message category (`chat`, `presence`, `control`, `ack`).
- `payload`: The actual message content.
- `sig`: HMAC signature for authenticity (optional).

## Framing Rules
- Max message size: 4096 bytes.
- Connections: Long-lived TCP.
- Reconnect: Clients should attempt to reconnect on discovery refresh if a connection is lost.
