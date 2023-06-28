package netil

import "net"

func NarrowIP(ip net.IP) net.IP {
	if len(ip) == 0 {
		return nil
	}
	if len(ip) == net.IPv4len {
		return ip
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return ip
	}
	return ip4
}

func IPFromAddr(addr net.Addr) net.IP {
	if addr == nil {
		return nil
	}
	ippAddr, err := ToIPPortAddr(addr)
	if err != nil {
		return nil
	}
	return ippAddr.IP
}
