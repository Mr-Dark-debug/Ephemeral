[![Banner](banner.svg)](https://github.com/Mr-Dark-debug/Ephemeral)

![Socialify](https://socialify.git.ci/Mr-Dark-debug/Ephemeral/image?description=1&font=Jost&forks=1&issues=1&language=1&name=1&owner=1&pattern=Charlie+Brown&pulls=1&stargazers=1&theme=Light)

# Ephemeral (Meshroom)

Ephemeral is a production-ready, terminal-first, zero-account, local-first LAN chat system written entirely in Go. Designed for privacy and instant communication without servers, history, or tracking.

## 🚀 Vision
Ephemeral enables users on the same Wi-Fi network to chat instantly using a single static binary. It prioritizes zero-persistence and robust peer discovery.

## ✨ Features
- **Local-first**: Operates entirely within your local network (LAN). No internet required.
- **Auto-Discovery**: mDNS (zeroconf) with a reliable UDP broadcast fallback.
- **Secure by Design**: Optional room-level encryption using AES-256-GCM and HKDF-SHA256.
- **Ephemeral Storage**: All chat state and keys are kept in-memory. Nothing is written to disk.
- **Modern TUI**: Built with Charm's Bubble Tea and Lip Gloss.
- **Cross-Platform**: Linux, macOS, Windows, and Android (Termux).

## 📥 Installation

### One-Line Quick Install
For the fastest setup, use our automated installer:

**Linux / macOS / Termux:**
```bash
curl -sSL https://raw.githubusercontent.com/Mr-Dark-debug/Ephemeral/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
powershell -ExecutionPolicy ByPass -Command "iwr -useb https://raw.githubusercontent.com/Mr-Dark-debug/Ephemeral/main/scripts/install.ps1 | iex"
```

### Manual Installation (Go)
Requires Go 1.23+. This will install the `ephemeral` binary to your `$GOPATH/bin`.
```bash
go install github.com/Mr-Dark-debug/Ephemeral/cmd/ephemeral@latest
```

### Termux (Android) Note
If you are running Go commands in `/storage/emulated/0`, you will encounter `RLock: function not implemented`. **Move the project to your Termux home directory (`~/`) to build from source.**

## 🛠 Usage
Launch the application:
```bash
ephemeral --nick alice
```

### Keyboard Shortcuts
- `Ctrl+C`: Quit
- `Enter`: Send message / Execute command
- `Ctrl+L`: Clear screen

### Commands
- `/join <room> [password]`: Join a room (encrypted if password provided).
- `/nick <newname>`: Change display name.
- `/quit`: Exit.

## 🏗 Architecture
- **Discovery**: mDNS (`github.com/grandcat/zeroconf`) + UDP Fallback.
- **Transport**: Persistent TCP with JSON-Lines framing.
- **Encryption**: AES-256-GCM, HKDF-SHA256.
- **UI Framework**: Bubble Tea.

## 🧪 Development
```bash
go test -v ./...
./scripts/build.sh
```

## 🛡 Security & Privacy
- **No Persistence**: History is lost on exit.
- **No Telemetry**: Ephemeral does not phone home.
- **Replay Protection**: Unique message UUIDs and in-memory caches.

## 📄 License
MIT License.
