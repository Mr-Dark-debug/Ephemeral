# Ephemeral (Meshroom)

Ephemeral is a production-ready, terminal-first, zero-account, local-first LAN chat system written entirely in Go. Designed for privacy and instant communication without servers, history, or tracking.

## 🚀 Vision
Ephemeral (formerly Meshroom) enables users on the same Wi-Fi network to chat instantly using a single static binary. It prioritizes zero-persistence and robust peer discovery.

## ✨ Features
- **Local-first**: Operates entirely within your local network (LAN). No internet required.
- **Auto-Discovery**: Primary discovery via mDNS (zeroconf) with a reliable UDP broadcast fallback.
- **Secure by Design**: Optional room-level encryption using AES-256-GCM and HKDF-SHA256 key derivation.
- **Ephemeral Storage**: All chat state and keys are kept in-memory. Nothing is written to disk by default.
- **Modern TUI**: Reactive terminal UI built with Charm's Bubble Tea and Lip Gloss.
- **Cross-Platform**: Native support for Linux, macOS, Windows, and Android (Termux).

## 📥 Installation

### From Source
Requires Go 1.23+.
```bash
git clone https://github.com/yourusername/ephemeral
cd ephemeral
go install ./cmd/ephemeral
```

### Termux (Android) Instructions
1.  Install **Termux** (F-Droid version recommended).
2.  Update packages: `pkg update && pkg upgrade`
3.  Install Go: `pkg install golang`
4.  Clone and run:
    ```bash
    git clone https://github.com/yourusername/ephemeral
    cd ephemeral
    go run ./cmd/ephemeral --nick myname
    ```

## 🛠 Usage
Launch the application:
```bash
ephemeral --nick alice
```
Or specify a port:
```bash
ephemeral --nick bob --port 9999
```

### Keyboard Shortcuts
- `Ctrl+C`: Quit
- `Enter`: Send message / Execute command
- `Alt+Screen`: TUI runs in alternate buffer (restores terminal on exit)

### Commands
- `/join <room> [password]`: Join a room. If a password is provided, the room is encrypted.
- `/nick <newname>`: Change your display name.
- `/peers`: List currently discovered peers.
- `/quit`: Exit Ephemeral.

## 🏗 Architecture
- **Language**: Go (Idiomatic, interface-driven design).
- **Discovery**: mDNS (via `github.com/grandcat/zeroconf`) + custom UDP Broadcaster.
- **Transport**: Persistent TCP connections with JSON-Lines framing.
- **Encryption**: AES-256-GCM for confidentiality and integrity.
- **UI Framework**: Bubble Tea (TEA) for state management.

## 🧪 Development & Testing
Ephemeral maintains high code quality with extensive testing.

### Run Tests
```bash
go test -v ./...
```
This runs unit tests for crypto, protocol, and config, as well as integration tests for peer-to-peer message exchange.

### Build Cross-Platform
```bash
./scripts/build.sh
```

## 🛡 Security & Privacy
- **No Persistence**: Chat history is lost once the application closes.
- **No Telemetry**: Ephemeral does not phone home or track usage.
- **Replay Protection**: Unique message UUIDs and in-memory caches prevent re-sending of old packets.
- **Threat Model**: Protects against local network eavesdropping when encryption is enabled. Does not protect against local machine physical access or compromised OS.

## ⚠️ Limitations
- **LAN Only**: Does not work across different subnets or over the internet without a VPN/Tunnel.
- **mDNS Restrictions**: Some public or corporate Wi-Fi networks block multicast/broadcast traffic, which may hinder discovery.
- **Scaling**: Optimized for small to medium groups (up to ~50 peers).

## 📄 License
MIT License. See `LICENSE` for details.
