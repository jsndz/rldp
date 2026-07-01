# Reliable Data Link Protocol (RLDP) in Go

A custom Layer-2–style reliable transport protocol implemented on top of UDP in Go. This project simulates network frames, sequence numbering, acknowledgments, sliding-window flow control, retransmission timers, and out-of-order buffering to ensure ordered, reliable, and duplicate-free message delivery over an unreliable network.

It also contains an experimental raw sockets sub-module designed to run directly over Ethernet frames at the Data Link layer.

---

## Features

- **Binary Framing**: Custom frame structure with Big-Endian serialization.
- **Error Detection**: IEEE CRC32 checksum verification for frame integrity.
- **Stop-and-Wait Reliability**: Single-frame sender with a 10-second timeout/retransmission mechanism.
- **Sliding Window Protocol**: Concurrently transmits up to 5 unacknowledged frames, sliding the sender window on ACK reception.
- **Selective Repeat / Out-of-Order Buffering**: The receiver buffers out-of-order frames, ignores duplicates, and reorders frames for sequential delivery.
- **Raw Socket Transport (Experimental)**: Bypasses the IP/UDP stack to communicate directly using Ethernet frames via raw sockets with custom EtherType `0x88B5`.

---

## Project Structure

- [types/types.go](file:///home/jaison/code/projects/rldp/types/types.go): Defines the core [Frame](file:///home/jaison/code/projects/rldp/types/types.go#L8) struct.
- [protocol/framing.go](file:///home/jaison/code/projects/rldp/protocol/framing.go): Handles [Encoding](file:///home/jaison/code/projects/rldp/protocol/framing.go#L11), [Decoding](file:///home/jaison/code/projects/rldp/protocol/framing.go#L23), and IEEE CRC32 [VerifyChecksum](file:///home/jaison/code/projects/rldp/protocol/framing.go#L43) verification.
- [transport/](file:///home/jaison/code/projects/rldp/transport): Implements UDP-based reliable communication.
  - [sender.go](file:///home/jaison/code/projects/rldp/transport/sender.go): Implements Stop-and-Wait [Send](file:///home/jaison/code/projects/rldp/transport/sender.go#L13) logic.
  - [receiver.go](file:///home/jaison/code/projects/rldp/transport/receiver.go): Implements [Receive](file:///home/jaison/code/projects/rldp/transport/receiver.go#L12) logic with duplicate filtering and out-of-order packet buffering.
- [rawtransport/](file:///home/jaison/code/projects/rldp/rawtransport): Implements low-level layer-2 Ethernet communication.
  - [sender.go](file:///home/jaison/code/projects/rldp/rawtransport/sender.go): Sends raw packet frames.
  - [receiver.go](file:///home/jaison/code/projects/rldp/rawtransport/receiver.go): Receives raw packet frames by binding to a socket file descriptor.
- [cmd/node/main.go](file:///home/jaison/code/projects/rldp/cmd/node/main.go): Main CLI entry point supporting `receive`, `send` (Stop-and-Wait), and `batch-send` (Sliding Window).

---

## Frame Structure

Each packet transmitted over the network is serialized in binary format with the following fields:

| Field | Size / Type | Description |
| :--- | :--- | :--- |
| **Sequence Number (Seq)** | `uint32` (4 bytes) | Identifies the sequence number of the transmitted frame. |
| **Acknowledgment Number (Ack)** | `uint32` (4 bytes) | Confirms the received sequence number. |
| **Frame Type** | `uint8` (1 byte) | Denotes the packet type (`DATA` = 1, `ACK` = 2). |
| **Payload** | `[]byte` (Variable) | The actual message payload. |
| **Checksum** | `uint32` (4 bytes) | IEEE CRC32 checksum computed over the frame header and payload. |

---

## Usage Instructions

Compile and run the protocol node using standard Go commands.

### 1. Start the Receiver
The receiver listens for incoming UDP packets on port `:8000`, verifies their checksums, handles duplicates, buffers out-of-order packets, and processes them sequentially:
```bash
go run cmd/node/main.go receive
```

### 2. Send Messages Interactively (Stop-and-Wait)
Run the node in interactive send mode. Note that the command-line usage prints `<dst-mac>` but actually expects a UDP host:port address (e.g. `127.0.0.1:8000`):
```bash
go run cmd/node/main.go send 127.0.0.1:8000 "initial"
```
Once started, you can input messages line-by-line. Type `exit` to terminate the session.

### 3. Send Messages in Batches (Sliding Window)
Send multiple packets simultaneously using a sliding window protocol of size 5:
```bash
go run cmd/node/main.go batch-send 127.0.0.1:8000 "hello" "world" "reliable" "data" "link" "protocol"
```
The sender will dispatch the packets concurrently and handle sliding the window upon receiving ACKs. If any packet times out (3 seconds), it retransmits the unacknowledged window.

### 4. Raw Socket (Ethernet Frame Mode)
If you wish to test the legacy layer-2 raw transport mode (requires root privileges and virtual/physical interfaces like `wlo1` or `eth0` configured in code):
- Run raw sender: code path calls [Send](file:///home/jaison/code/projects/rldp/rawtransport/sender.go#L11).
- Run raw receiver: code path calls [Receive](file:///home/jaison/code/projects/rldp/rawtransport/receiver.go#L10).
```bash
sudo go run cmd/node/main.go ...
```

---

## Legacy Evolution
Initially, raw sockets binding directly to the Network Interface Card (NIC) was attempted using custom EtherType `0x88B5` (see [rawtransport/](file:///home/jaison/code/projects/rldp/rawtransport) and [logs.md](file:///home/jaison/code/projects/rldp/logs.md)). Due to loopback constraints and local OS interface filtering, the project switched to UDP for standard transport simulation while keeping the framing, sliding-window flow control, and checksum validation logic fully custom.
