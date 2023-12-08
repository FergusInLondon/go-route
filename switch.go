package main

import (
	"context"
	"errors"
	"fmt"
	"net"
)

const (
	// @see https://en.wikipedia.org/wiki/Ethernet_frame#Structure
	ethernetFrameLength = 1522
	errSocketClosed     = "socket closure requested"
)

var (
	broadcastMacAddr = convertMacAddress([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
)

type LayerTwoSwitchParams struct {
	ListenAddr []byte
	ListenPort int
}

type LayerTwoSwitch struct {
	ctx  context.Context
	stop context.CancelCauseFunc
	conn interface {
		ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
		WriteToUDP(b []byte, addr *net.UDPAddr) (int, error)
		Close() error
	} // i.e. *net.UDPConn
	macTable map[uint64]*net.UDPAddr
}

func NewLayerTwoSwitch(
	ctx context.Context, params *LayerTwoSwitchParams,
) *LayerTwoSwitch {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: params.ListenAddr, Port: params.ListenPort, Zone: "",
	})
	if err != nil {
		panic(err)
	}

	vswitchCtx, vswitchStop := context.WithCancelCause(ctx)
	return &LayerTwoSwitch{vswitchCtx, vswitchStop, conn, make(map[uint64]*net.UDPAddr)}
}

// This feels like a quick way, but maybe isn't.
func convertMacAddress(raw []byte) uint64 {
	addrVal := uint64(0)

	for _, b := range raw {
		addrVal = (addrVal << 8) | uint64(b)
	}

	return addrVal
}

func addr_check(frameAddr, cacheAddr *net.UDPAddr) bool {
	return frameAddr.IP.Equal(cacheAddr.IP) && frameAddr.Port == cacheAddr.Port
}

func (vswitch *LayerTwoSwitch) handleFrame(len int, srcAddr *net.UDPAddr, data []byte) {
	// Maintain ARP Cache / MAC Table for Source
	srcMacAddr := convertMacAddress(data[0:6])
	srcArpAddr, haveSrcArpAddr := vswitch.macTable[srcMacAddr]
	if !(haveSrcArpAddr && addr_check(srcArpAddr, srcAddr)) {
		vswitch.macTable[srcMacAddr] = srcAddr
	}

	dstMacAddr := convertMacAddress(data[6:12])

	// Handle Broadcast Frames
	if dstMacAddr == broadcastMacAddr {
		for macAddr, udpAddr := range vswitch.macTable {
			if macAddr != srcMacAddr {
				vswitch.sendTo(data[0:len], udpAddr)
			}
		}

		return
	}

	// Handle Direct Frames
	dstUdpAddr, haveAddress := vswitch.macTable[dstMacAddr]
	if !haveAddress {
		// should probably log.
		return
	}

	vswitch.sendTo(data[0:len], dstUdpAddr)
}

func (vswitch *LayerTwoSwitch) sendTo(frame []byte, dest *net.UDPAddr) {
	len, err := vswitch.conn.WriteToUDP(frame, dest)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("tx", len)
}

func (vswitch *LayerTwoSwitch) Listen() error {
	for {
		buf := make([]byte, ethernetFrameLength)
		len, frameAddr, err := vswitch.conn.ReadFromUDP(buf)
		if err != nil {
			go vswitch.handleFrame(len, frameAddr, buf)
		}

		// @todo handle Read Errors.
	}
}

func (vswitch *LayerTwoSwitch) Close() {
	if err := context.Cause(vswitch.ctx); err != nil {
		fmt.Println("vswitchet already closed", err)
		return
	}

	vswitch.conn.Close()
	vswitch.stop(errors.New(errSocketClosed))
}
