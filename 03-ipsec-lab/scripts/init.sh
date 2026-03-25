#!/bin/sh

echo "Updating package repositories..."
apk update

# Install strongSwan, tcpdump, and Scapy (Python packet manipulation library)
echo "Installing dependencies..."
apk add strongswan tcpdump py3-scapy

# Start the strongSwan background service
echo "Starting IPSec daemon..."
ipsec start

echo "Initialization complete. IPSec is running in the background."