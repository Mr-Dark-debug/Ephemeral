# Design Documentation - Ephemeral

Ephemeral is a decentralized, local-first chat application designed for immediate use on a LAN without any central coordination.

## Architecture

The system is composed of several decoupled modules:

1.  **Discovery Layer**: Uses mDNS (Multicast DNS) as the primary mechanism. Peers advertise `_meshroom._tcp` on the `.local` domain. A UDP broadcast fallback (port 9998) is used for networks that block multicast.
2.  **Transport Layer**: Reliable TCP connections. Once a peer is discovered, a persistent TCP connection is established.
3.  **Protocol Layer**: JSON-Lines based messaging. Each message is an independent JSON object followed by a newline.
4.  **Room Manager**: Logic-based rooms. Users "join" a room by filtering and broadcasting messages with specific room tags.
5.  **Crypto Module**: Handles passphrase-based key derivation (HKDF-SHA256) and authenticated encryption (AES-256-GCM).
6.  **TUI Layer**: Reactive terminal interface using Bubble Tea.

## Sequence Diagrams

### Discovery & Connection
1. Peer A starts and registers `_meshroom._tcp` via mDNS.
2. Peer B starts and browses for `_meshroom._tcp`.
3. Peer B discovers Peer A's IP and Port.
4. Peer B initiates a TCP connection to Peer A.
5. Peer A accepts and identifies Peer B via the initial `presence` message.

### Messaging
1. User types message in TUI.
2. Message is wrapped in a JSON `Envelope`.
3. If the room is encrypted, the payload is encrypted using the room's derived key.
4. The envelope is serialized and sent over all active TCP peer connections.
5. Recipients deserialize, decrypt (if necessary), and display the message.

## Scalability
The current mesh flooding model is $O(N^2)$ in terms of connections and bandwidth. It is optimized for small to medium groups (up to 50-100 peers) on a local network.
