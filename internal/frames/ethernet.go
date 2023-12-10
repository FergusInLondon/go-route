package frames

import "net"

var (
	ETHERTYPE_IPv4 = []byte{0x08, 0x00}
	ETHERTYPE_ARP  = []byte{0x08, 0x06}
	ETHERTYPE_IPV6 = []byte{0x86, 0xDD}
)

type EthernetFrame struct {
	Destination, Source net.HardwareAddr
	Tags                []byte
	Type                []byte
	Payload             []byte
	CRC                 []byte
}

func ParseEthernetFrame(buf []byte) *EthernetFrame {
	// This is a bit -interesting-.
	// @see https://github.com/songgao/packets/blob/master/ethernet/frame.go#L39-L46
	tagLength := 0
	if buf[12] == 0x81 && buf[13] == 0x00 {
		tagLength = 4
	} else if buf[12] == 0x88 && buf[13] == 0xA8 {
		tagLength = 8
	}

	return &EthernetFrame{
		Source:      buf[6:12],
		Destination: buf[0:6],
		Tags:        buf[12:(12 + tagLength)],
		Type:        buf[(12 + tagLength):(12 + tagLength + 2)],
		Payload:     buf[(12 + tagLength + 2):(len(buf) - 4)],
		CRC:         buf[:(len(buf) - 4)],
	}
}

type IPv6Packet struct {
	Version int8 // Actual Width: 4 Bits.
	// Differentiated Service & Explicit Congestion
	TrafficClass int8 // Actual Width: 8 Bits.
	// Value that can be used for identifying/grouping related packets.
	FlowLabel int32 // Actual Width: 20 Bits.
	// The length of the payload in bytes/octets.
	PayloadLength int16 // Actual Width: 16 Bits.
	// The type of the "next header" - i.e. the transport layer protocol
	NextHeader int8 // Actual Width: 8 Bits.
	// IPv6 equivalent of TTL
	HopLimit int8 // Actual Width: 8 Bits.
	// IPv6 Addresses
	Source, Destination []byte // Actual Width: 128 Bits / 16 bytes.
}

func ParseIPv6Packet(buf []byte) *IPv6Packet {
	// Yes, there's shouldn't be a reason for an L2 switch to have
	// awareness above L2... on the other hand, it's quite surprising
	// seeing the number of discarded frames.
	return nil
}
