Build a system where two programs send messages to each other, even when the network is unreliable.

---

### What you actually do (step-by-step)

1. **Create two programs**

* Sender
* Receiver
  Both communicate using UDP

---

2. **Define your “packet”**
   Each message you send should include:

* sequence number
* data
* simple checksum

---

3. **Send messages**

* Sender sends packets with sequence numbers (1, 2, 3…)

---

4. **Receiver logic**

* If correct packet → accept and send ACK
* If missing or wrong → ignore

---

5. **Add reliability**

* If sender doesn’t get ACK → resend packet
* Handle duplicates (ignore repeated packets)

---

6. **Add sliding window (important)**

* Send multiple packets at once (not just one)
* Keep track of which are not yet acknowledged

---

7. **Break the network (simulation)**
   Intentionally:

* drop some packets
* delay some
* reorder some

---

8. **Fix it using your protocol**
   Ensure:

* all data arrives
* in correct order
* no duplicates

---

### Final goal

You type:

```
send "hello world"
```

Even with packet loss:

* receiver still gets: `hello world` correctly

---