[![Banner](banner.svg)](https://github.com/Mr-Dark-debug/Ephemeral)

![Socialify](https://socialify.git.ci/Mr-Dark-debug/Ephemeral/image?description=1&font=Jost&forks=1&issues=1&language=1&name=1&owner=1&pattern=Charlie+Brown&pulls=1&stargazers=1&theme=Light)

# Ephemeral (Meshroom)

**Ephemeral** is a production-ready, terminal-first, zero-account, local-first LAN chat system written entirely in Go. It is designed for high-privacy, instant communication on same-network environments without servers, history, or tracking.

---

## 🌟 Project Vision
Ephemeral (formerly Meshroom) is built for those who value transient communication. Whether you are at a hackathon, a conference, or on a shared home Wi-Fi, Ephemeral allows you to spin up a chat room in seconds with zero configuration and total privacy.

### Core Philosophy
- **Local-first**: No internet required. No external servers.
- **Zero-Account**: No signups, no emails, no phone numbers.
- **No-History**: Your data lives in RAM and dies in RAM.
- **Single Binary**: No complex dependencies. Just one file.

---

## ✨ Features
- **Discovery Layer**: Primary discovery via mDNS (zeroconf) with a reliable UDP broadcast fallback for restricted networks.
- **Transport**: Persistent, backpressure-safe TCP connections with JSON-Lines framing.
- **Encryption**: Optional end-to-end room-level encryption using AES-256-GCM and HKDF-SHA256.
- **Modern TUI**: A beautiful, responsive terminal interface built with Charm's `Bubble Tea` and `Lip Gloss`.
- **Responsive Design**: UI scales gracefully from small Termux screens to ultra-wide monitors.
- **Cross-Platform**: Full support for Linux, macOS, Windows, and Android (Termux).

---

<img width="1919" height="1007" alt="image" src="https://github.com/user-attachments/assets/f42cd401-5af5-4996-b3de-aedfbd75d8d9" />

---

## 📥 Installation

### 🚀 One-Line Quick Install (Recommended)
Our automated installers detect your OS/Architecture and pull the latest production binary.

**Linux / macOS / Termux:**
```bash
curl -sSL https://raw.githubusercontent.com/Mr-Dark-debug/Ephemeral/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
powershell -ExecutionPolicy ByPass -Command "iwr -useb https://raw.githubusercontent.com/Mr-Dark-debug/Ephemeral/main/scripts/install.ps1 | iex"
```

### 🛠 Manual Installation (From Source)
Requires Go 1.23+.
```bash
go install github.com/Mr-Dark-debug/Ephemeral/cmd/ephemeral@latest
```

---

## 📱 Termux (Android) Instructions
Ephemeral is fully optimized for Termux.

1.  **Install Termux** from F-Droid (not Play Store).
2.  **Move project to Home**: Due to Android filesystem limitations (FUSE), Go commands will fail in `/storage/emulated/0`. Always run from `~/`:
    ```bash
    git clone https://github.com/Mr-Dark-debug/Ephemeral ~/Ephemeral
    cd ~/Ephemeral
    go run ./cmd/ephemeral --nick myname
    ```

---

## 🛠 Usage & Commands
Launch with a simple command:
```bash
ephemeral --nick Alice
```

### Interactive Commands
Inside the TUI, type these commands in the input field:
- `/join <room> [password]`: Join a logical room. Providing a password enables AES-256-GCM encryption.
- `/leave`: Return to the `global` room.
- `/nick <newname>`: Change your display name instantly.
- `/peers`: List all discovered peers on the network.
- `/quit`: Exit the application.

### Keyboard Shortcuts
- `Ctrl+C`: Quit application.
- `Ctrl+L`: Clear the message viewport.
- `Enter`: Send message or execute command.

---

## 🧪 Demo Script (3-Step Guide)
1.  **Start Peer A**: Run `ephemeral --nick Alice`.
2.  **Start Peer B**: On another machine/terminal, run `ephemeral --nick Bob`.
3.  **Encrypted Chat**: Alice types `/join secret hunter2`. Bob types `/join secret hunter2`. They are now chatting securely.

---

## 🏗 Architecture & Documentation
Ephemeral is built with a clean, modular architecture using Dependency Injection for high testability.

- **[Design Docs](docs/design.md)**: Deep dive into the system architecture and sequence diagrams.
- **[Protocol Spec](docs/protocol.md)**: Details on the JSON-Lines wire format and message envelopes.
- **[Security Model](docs/security.md)**: Threat model and cryptographic choices.

---

## 🛡 Security & Privacy Notice
- **No Telemetry**: We do not collect analytics, crash reports, or usage data.
- **No External Connections**: Ephemeral only talks to peers on your local network.
- **Memory Safety**: No chat data is written to disk. Once you exit, the data is gone forever.
- **Privacy Warning**: Users on your LAN can see that you are running Ephemeral unless you use a VPN. Encryption only hides the *content* of your messages.

---

## ⚠️ Troubleshooting
- **No Peers Found**: Ensure you are on the same Wi-Fi subnet. Check if your firewall blocks port `9999` (TCP) and `9998` (UDP).
- **Termux RLock Error**: Move the project to `~/` (home directory) to avoid Android's restricted filesystem.
- **mDNS Issues**: On some corporate networks, mDNS is blocked. Ephemeral will automatically fallback to UDP broadcast.

---

## 📄 License
Licensed under the [MIT License](LICENSE). 2026 Ephemeral Contributors.
