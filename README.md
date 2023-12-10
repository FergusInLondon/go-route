# vnet - a demonstrative virtual network
Because it reaches a point when it's no longer appropriate to answer "*magic*" when someone asks you a networking-related question.

## Description

This is a simple demonstration of tunneling layer 2 network traffic (i.e. raw Ethernet II / DIX frames) over UDP - in a *vaguely similar* vein to L2TP. It contains no consideration to confidentiality of transmitted data, nor authentication of connected nodes.

The client contains a TAP interface which forwards all Ethernet Frames to the server via UDP.

The server handles inbound Ethernet Frames from the UDP socket, and attempts to route it to the correct client based upon the MAC address.

### Technical Description

Currently implemented as a Go equivalent to the excellent [build-your-own-zerotier](https://github.com/peiyuanix/build-your-own-zerotier/tree/master). Although I'm currently adding in some parsing for higher layer frames - i.e. IPv6 - so as to handle multicast and address resolution.