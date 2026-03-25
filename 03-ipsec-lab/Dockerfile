# Start with a lightweight Alpine Linux base
FROM alpine:latest

# Install our required networking and python tools
RUN apk update && apk add strongswan tcpdump scapy iputils

# Disable the bypass-lan plugin so local traffic gets forced into the tunnel
RUN printf "charon {\n  plugins {\n    bypass-lan {\n      load = no\n    }\n  }\n}\n" >> /etc/strongswan.conf

# Start the IPSec daemon in the foreground
CMD ["ipsec", "start", "--nofork"]