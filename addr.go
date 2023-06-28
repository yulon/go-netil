package netil

import (
	"net"
	"strconv"
)

type IPPortAddr struct {
	IP   net.IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

func (a *IPPortAddr) Network() string { return "ip" }

func ipEmptyString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}

func (a *IPPortAddr) String() string {
	if a == nil {
		return ""
	}
	ip := ipEmptyString(a.IP)
	if a.Zone != "" {
		return net.JoinHostPort(ip+"%"+a.Zone, strconv.Itoa(a.Port))
	}
	return net.JoinHostPort(ip, strconv.Itoa(a.Port))
}

type DomainAddr string

func (a DomainAddr) Network() string {
	return "domain"
}

func (a DomainAddr) String() string {
	return string(a)
}

func ToTCPAddr(addr net.Addr) (*net.TCPAddr, error) {
	switch addr.(type) {
	case *net.TCPAddr:
		return addr.(*net.TCPAddr), nil

	case *IPPortAddr:
		ipAddr := addr.(*IPPortAddr)
		return &net.TCPAddr{IP: ipAddr.IP, Port: ipAddr.Port, Zone: ipAddr.Zone}, nil

	case *net.UDPAddr:
		udpAddr := addr.(*net.UDPAddr)
		return &net.TCPAddr{IP: udpAddr.IP, Port: udpAddr.Port, Zone: udpAddr.Zone}, nil
	}
	return net.ResolveTCPAddr("tcp", addr.String())
}

func CopyToTCPAddr(addr net.Addr) (*net.TCPAddr, error) {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if ok {
		return &net.TCPAddr{IP: tcpAddr.IP, Port: tcpAddr.Port, Zone: tcpAddr.Zone}, nil
	}
	return ToTCPAddr(addr)
}

func ToUDPAddr(addr net.Addr) (*net.UDPAddr, error) {
	switch addr.(type) {
	case *net.UDPAddr:
		return addr.(*net.UDPAddr), nil

	case *IPPortAddr:
		ipAddr := addr.(*IPPortAddr)
		return &net.UDPAddr{IP: ipAddr.IP, Port: ipAddr.Port, Zone: ipAddr.Zone}, nil

	case *net.TCPAddr:
		tcpAddr := addr.(*net.TCPAddr)
		return &net.UDPAddr{IP: tcpAddr.IP, Port: tcpAddr.Port, Zone: tcpAddr.Zone}, nil
	}
	return net.ResolveUDPAddr("udp", addr.String())
}

func CopyToUDPAddr(addr net.Addr) (*net.UDPAddr, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if ok {
		return &net.UDPAddr{IP: udpAddr.IP, Port: udpAddr.Port, Zone: udpAddr.Zone}, nil
	}
	return ToUDPAddr(addr)
}

func ToIPPortAddr(addr net.Addr) (*IPPortAddr, error) {
	switch addr.(type) {
	case *IPPortAddr:
		return addr.(*IPPortAddr), nil

	case *net.UDPAddr:
		udpAddr := addr.(*net.UDPAddr)
		return &IPPortAddr{IP: udpAddr.IP, Port: udpAddr.Port, Zone: udpAddr.Zone}, nil

	case *net.TCPAddr:
		tcpAddr := addr.(*net.TCPAddr)
		return &IPPortAddr{IP: tcpAddr.IP, Port: tcpAddr.Port, Zone: tcpAddr.Zone}, nil
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	return &IPPortAddr{IP: tcpAddr.IP, Port: tcpAddr.Port, Zone: tcpAddr.Zone}, nil
}

func CopyToIPPortAddr(addr net.Addr) (*IPPortAddr, error) {
	ipPortAddr, ok := addr.(*IPPortAddr)
	if ok {
		return &IPPortAddr{IP: ipPortAddr.IP, Port: ipPortAddr.Port, Zone: ipPortAddr.Zone}, nil
	}
	return ToIPPortAddr(addr)
}

func PublicAddr(privateAddr, refPublicAddr net.Addr) (net.Addr, error) {
	privateTCPAddr, err := ToIPPortAddr(privateAddr)
	if err != nil {
		return nil, err
	}
	if privateTCPAddr.IP.IsGlobalUnicast() {
		return privateTCPAddr, nil
	}
	refIPPortAddr, err := ToIPPortAddr(refPublicAddr)
	if err != nil {
		return nil, err
	}
	return refIPPortAddr, nil
}
