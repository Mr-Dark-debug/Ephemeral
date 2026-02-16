# Ephemeral (Meshroom)

Ephemeral is a terminal-first, zero-account, no-server, no-history LAN chat system written entirely in Go.

## Features

- **Local-first**: Works on LAN without internet.
- **Zero-config**: Auto-discovery via mDNS (zeroconf).
- **Secure**: Optional AES-256-GCM encryption with HKDF key derivation.
- **Private**: No chat history stored on disk. No telemetry.
- **Cross-platform**: Linux, macOS, Windows, Android (Termux).

## 📥 Installation

### One-Line Quick Install
For the fastest setup, use our automated installer:

**Linux / macOS / Termux:**
```bash
curl -sSL https://raw.githubusercontent.com/yourusername/ephemeral/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
powershell -ExecutionPolicy ByPass -Command "iwr -useb https://raw.githubusercontent.com/yourusername/ephemeral/main/scripts/install.ps1 | iex"
```

### Manual Installation (Go)
Requires Go 1.23+. This will install the `ephemeral` binary to your `$GOPATH/bin`.
```bash
go install github.com/yourusername/ephemeral/cmd/ephemeral@latest
```

### Termux (Android) Note
If you are running Go commands in `/storage/emulated/0`, you will encounter `RLock: function not implemented` due to filesystem limitations. **Move the project to your Termux home directory (`~/`) to resolve this.**

## Usage

Start the chat:

```bash
ephemeral --nick Alice
```

Commands inside TUI:
- `/join <room> [password]`: Join a room (optionally encrypted).
- `/nick <name>`: Change nickname.
- `/quit`: Exit.

## Architecture

- **Discovery**: mDNS (grandcat/zeroconf).
- **Transport**: TCP JSON-lines with backpressure.
- **UI**: Bubble Tea + Lip Gloss.
- **Crypto**: AES-256-GCM, HKDF-SHA256.

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
./scripts/build.sh
```

## Security

- Keys are derived from passphrases using HKDF.
- Messages are encrypted with AES-256-GCM.
- Replay protection is implemented via unique message IDs (in-memory).
- No data is persisted to disk.

## License

MIT
