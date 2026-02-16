# Security Model - Ephemeral

## Threat Model
Ephemeral is designed for privacy on local networks. It protects against:
- **Eavesdropping**: Passive attackers on the same Wi-Fi cannot read encrypted room traffic.
- **Tampering**: AES-GCM provides authentication; modified packets will fail decryption.
- **Replay Attacks**: In-memory ID cache prevents re-processing of the same message.

It does **not** protect against:
- **Compromised Endpoints**: If a peer's terminal or OS is compromised, keys can be extracted from memory.
- **Traffic Analysis**: An observer can see that IP A is talking to IP B.

## Cryptographic Choices
- **Key Derivation**: HKDF-SHA256. We use a salt (application-specific or user-defined) and the room passphrase to derive a 256-bit key.
- **Encryption**: AES-256-GCM. Provides high-performance authenticated encryption.
- **Nonces**: 12-byte random nonces generated via `crypto/rand` for every message. Nonces are never reused with the same key.

## Data Persistence
- **Zero-History**: No chat logs are ever written to disk.
- **In-Memory Only**: Keys and messages exist only in volatile memory and are wiped when the process exits.
