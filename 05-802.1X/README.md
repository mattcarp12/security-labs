# Lab 05: 802.1X, EAP, and RADIUS Authentication

This lab demonstrates Layer 2 Port-Based Network Access Control (PNAC) using the IEEE 802.1X standard. It utilizes a Dockerized FreeRADIUS container acting as the Authentication Server, and a custom Go application acting as both the Authenticator (Switch) and Supplicant (Client).

## Learning Objectives (Security+)
* **The 802.1X Triad:** Understanding the flow between Supplicant, Authenticator, and Authentication Server.
* **RADIUS over UDP:** Utilizing UDP ports 1812 (Auth) and 1813 (Accounting).
* **EAP (Extensible Authentication Protocol):** Constructing EAP envelopes and understanding the difference between PAP and EAP.
* **Cryptographic Integrity:** Using the `Message-Authenticator` (HMAC-MD5) attribute to prevent Man-in-the-Middle downgrade attacks (RFC 3579).

## Prerequisites
* Docker & Docker Compose
* Go 1.20+
* Local network environment (127.0.0.1)

## Directory Structure
```text
lab-02-8021x/
├── docker-compose.yml
├── freeradius/
│   ├── clients.conf    # Defines the Pre-Shared Secret for the Authenticator
│   └── users           # The credential database
└── supplicant/
    ├── go.mod
    ├── go.sum
    └── main.go         # The Go program crafting raw RADIUS/EAP packets
```

## Usage Instructions

### 1. Start the FreeRADIUS Server
Boot the backend database in the foreground with extreme debugging (`-X`) enabled so you can watch the authentication process in real-time.
```bash
docker-compose up -d
```
*To view the real-time server logs: `docker logs -f 8021x-radius`*

### 2. Initialize the Go Module
If you haven't already, pull down the required RADIUS libraries for the Go code:
```bash
cd supplicant
go mod tidy
```

### 3. Run the Authentication Request
Fire the raw EAP packet at the RADIUS server. 
```bash
go run main.go
```

## Expected Output
* **Success (PAP):** `[SUCCESS] RADIUS Code 2: Access-Accept`
* **Zero Trust Rejection (PAP):** `[REJECTED] RADIUS Code 3: Access-Reject`
* **EAP Handshake:** `[CHALLENGE] RADIUS Code 11: Access-Challenge` (Server demands cryptographic proof).

## Cleanup
```bash
docker-compose down
```
