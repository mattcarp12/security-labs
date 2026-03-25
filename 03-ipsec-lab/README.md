# Lab 03: IPsec, IKEv2, and Cryptographic Tunnels

This lab demonstrates a working Site-to-Site IPsec VPN using Pre-Shared Keys (PSK). It utilizes Docker containers to simulate two distinct networks ("Alice" and "Bob"), strongSwan for IKEv2 negotiation, and Python (Scapy) to mathematically verify Encapsulating Security Payload (ESP) encryption on the wire.

## Learning Objectives (Security+)
* **IPsec Tunnel Mode:** Encapsulating entire IP packets for Site-to-Site routing.
* **IKEv2 & PSK:** Negotiating cryptographic keys and authenticating endpoints.
* **ESP (Protocol 50):** Verifying data confidentiality by sniffing the network interface.

## Prerequisites
* Docker & Docker Compose
* Linux host (or WSL2 on Windows)

## Directory Structure
```text
lab-01-ipsec/
├── docker-compose.yml
├── Dockerfile
├── configs/
│   ├── alice/
│   │   ├── ipsec.conf
│   │   └── ipsec.secrets
│   └── bob/
│       ├── ipsec.conf
│       └── ipsec.secrets
└── scripts/
    └── sniff.py
```

## Usage Instructions

### 1. Build and Start the Infrastructure
This will build the custom Alpine images, install strongSwan and Python, and automatically bring up the IPsec tunnel.
```bash
docker-compose up -d --build
```

### 2. Verify the Tunnel is Established
Check Alice's IPsec status to ensure the Security Associations (SAs) are successfully installed.
```bash
docker exec ipsec-alice ipsec status
```
*(Look for `ESTABLISHED` and `INSTALLED, TUNNEL` in the output).*

### 3. Start the Packet Sniffer
Open a new terminal window and start the Python Scapy script on Bob to listen to incoming traffic on `eth0`.
```bash
docker exec -it ipsec-bob python3 /scripts/sniff.py
```

### 4. Generate Traffic
In your original terminal, have Alice send 4 ICMP ping packets to Bob's IP address.
```bash
docker exec ipsec-alice ping -c 4 10.0.0.20
```

## Expected Output & Verification
In the sniffer terminal, you should see output similar to:
```text
[SECURE] ESP Packet detected! 10.0.0.10 -> 10.0.0.20
```

Because the IPsec tunnel is active, the Linux kernel intercepts the standard cleartext ICMP ping and wraps it in an encrypted ESP packet (Protocol 50) before it hits the wire. The payload is successfully secured in transit.

## Cleanup
To tear down the lab and remove the virtual networks:
```bash
docker-compose down
```
