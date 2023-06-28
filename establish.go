package netil

import (
	"net"
)

func DialOrDirect(addr net.Addr, dial func(addr net.Addr) (net.Conn, error)) (net.Conn, error) {
	if dial != nil {
		return dial(addr)
	}
	tcpAddr, err := ToTCPAddr(addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, tcpAddr)
}

func ListenPacketOrDirect(listenPacket func() (net.PacketConn, error)) (net.PacketConn, error) {
	if listenPacket != nil {
		return listenPacket()
	}
	return net.ListenUDP("udp", nil)
}
