# vnet - a demonstrative virtual network
Because it reaches a point when it's no longer appropriate to answer "*magic*" when someone asks you a networking-related question.

## Description

This is a simple demonstration of tunneling layer 2 network traffic (i.e. raw Ethernet II / DIX frames) over UDP - in a *vaguely similar* vein to L2TP. It contains no consideration to confidentiality of transmitted data, nor authentication of connected nodes.

### Technical Description

Currently implemented as a Go equivalent to the excellent [build-your-own-zerotier](https://github.com/peiyuanix/build-your-own-zerotier/tree/master).
