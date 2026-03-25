#!/usr/bin/env python3
from scapy.all import sniff, IP, ESP, ICMP

def analyze_packet(packet):
    # Check if the packet has been wrapped in an IPSec ESP header
    if packet.haslayer(ESP):
        print(f"[SECURE] ESP Packet detected! {packet[IP].src} -> {packet[IP].dst}")
        print(f"         Payload is encrypted. SPI (Security Parameter Index): {hex(packet[ESP].spi)}\n")
    
    # Check if it's a standard cleartext ping
    elif packet.haslayer(ICMP):
        print(f"[INSECURE] Cleartext ICMP (Ping)! {packet[IP].src} -> {packet[IP].dst}")
        # If the ping has raw data inside it, print it out to prove it's readable
        if hasattr(packet[ICMP], 'load'):
            print(f"           Exposed Payload: {packet[ICMP].load}\n")

print("Listening for 10 packets on eth0...")
sniff(iface="eth0", filter="ip", prn=analyze_packet, count=10)