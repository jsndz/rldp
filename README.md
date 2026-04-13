**Project:** Reliable Data Link Protocol over UDP (Go)

---

### Objective

Build a custom **Layer-2–style reliable protocol** on top of UDP that guarantees ordered, reliable delivery despite packet loss, duplication, and reordering.

---

### System Overview

Multiple nodes communicate over UDP. Each node implements a full data link protocol:

* framing
* sequencing
* acknowledgments
* retransmission
* flow control

---

### What needs to be built

#### 1. Frame Structure

Define a binary frame format:

* sequence number
* acknowledgment number
* flags (DATA, ACK, RETRANSMIT)
* payload
* checksum

---

#### 2. Sender Module

Responsible for:

* sending frames from buffer
* maintaining sliding window
* tracking unacknowledged frames
* starting retransmission timers

---

#### 3. Receiver Module

Responsible for:

* receiving frames
* validating checksum
* handling:

  * in-order frames
  * out-of-order frames
* sending ACKs

---

#### 4. Sliding Window Protocol

Implement:

* window size control
* sequence number wrap-around
* Selective Repeat (preferred)

---

#### 5. Reliability Logic

* retransmit on timeout
* discard duplicates
* reorder packets before delivery
* ensure in-order delivery to application

---

#### 6. Network Simulator (important)

Simulate unreliable network:

* packet loss (drop %)
* delay
* duplication
* reordering

This is critical for learning.

---

#### 7. Timer System

* per-packet timers
* dynamic timeout (optional advanced)

---

#### 8. Application Layer Interface

Expose simple API:

* `Send(data []byte)`
* `OnReceive(callback)`

---

#### 9. Metrics & Logging

Track:

* RTT
* retransmissions
* throughput
* packet loss rate

---

### Folder Structure (suggested)

* `/cmd/node` → run node
* `/protocol` → core logic
* `/transport` → UDP + simulation
* `/frame` → encoding/decoding
* `/metrics` → stats

---

### Minimum Working Version (MVP)

Must support:

* send message from A → B
* simulate packet loss
* still deliver correctly and in order

---

### Advanced Extensions

* congestion control (AIMD)
* adaptive timeout (RTT-based)
* message batching
* multiple clients
* persistent queue

---

### End Result

You will have:

* a working reliable transport protocol
* deep understanding of:

  * TCP-like behavior
  * distributed system communication
  * failure handling

---
