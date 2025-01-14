package addrx

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// ExtractAddr extract IP and port from net.Addr
func ExtractAddr(addr net.Addr) (net.IP, int, error) {
	switch v := addr.(type) {
	case *net.IPAddr:
		return v.IP, 0, nil
	case *net.IPNet:
		return v.IP, 0, nil
	case *net.TCPAddr:
		return v.IP, v.Port, nil
	case *net.UDPAddr:
		return v.IP, v.Port, nil
	default:
		host, port, err := net.SplitHostPort(addr.String())
		if err != nil {
			return net.IP{}, 0, err
		}
		portNum, err := strconv.Atoi(port)
		return net.ParseIP(host), portNum, err
	}
}

// GlobalUnicastIP get a global unicast IP address
func GlobalUnicastIP() (net.IP, error) {
	for _, ip := range IPs() {
		if IsGlobalUnicastIP(ip) {
			return ip, nil
		}
	}
	return nil, fmt.Errorf("no found global unicast IP")
}

// IsGlobalUnicastIP check whether the IP is a global unicast IP
func IsGlobalUnicastIP(ip net.IP) bool {
	if ip.IsUnspecified() {
		// 这个方法用于检查给定的 IP 地址是否是未指定的地址。
		// 未指定的地址: 在 IP 地址中，未指定的地址表示没有特定的网络接口或地址。具体来说：
		// IPv4: "0.0.0.0"。
		// IPv6: "::"。
		return false
	}
	if ip.IsLoopback() {
		// 这个方法用于检查给定的 IP 地址是否是回环地址。
		// 回环地址（Loopback Address）是网络中的一种特殊IP地址，主要用于测试和本地通信。
		// IPv4: "127.0.0.0/8"，最常用的回环地址是 127.0.0.1。
		// IPv6: "::1"。
		return false
	}
	if ip.IsMulticast() {
		// 这个方法用于检查给定的 IP 地址是多播地址。
		// 多播地址用于向一组主机发送数据包，而不是单个主机。
		// IPv4: 224.0.0.0 到 239.255.255.255。
		// IPv6: ff00::/8
		return false
	}
	if ip.IsLinkLocalMulticast() {
		// 这个方法用于检查给定的 IP 地址是链路本地多播地址。
		// 链路本地多播地址用于在同一网络段内的设备之间进行通信。
		// IPv4: 224.0.0.0 到 224.0.0.255。
		// IPv6: ff02::/16
		return false
	}
	if ip.IsInterfaceLocalMulticast() {
		// 这个方法用于检查给定的 IP 地址是接口本地多播地址。
		// 接口本地多播地址用于在同一网络接口内的设备之间进行通信。
		// IPv4: 224.0.0.0 到 224.0.0.255（具体到 224.0.0.252 和 224.0.0.253）。
		// IPv6: ff01::/16
		return false
	}
	if ip.IsLinkLocalUnicast() {
		// 这个方法用于检查给定的 IP 地址是链路本地单播地址。
		// 链路本地单播地址用于在同一网络段内的设备之间进行单点通信。
		// IPv4: 169.254.0.0/16。
		// IPv6: fe80::/10
		return false
	}
	// 其他情况，认为是全局单播地址
	// 全局单播地址用于在互联网上进行通信，而不是在本地网络或特定的网络段内。
	return ip.IsGlobalUnicast()
}

// GlobalUnicastIPString get a global unicast IP address string
func GlobalUnicastIPString() (string, error) {
	ip, err := GlobalUnicastIP()
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

// IPs get all IP addresses
func IPs() []net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ips []net.IP
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, _, err := ExtractAddr(addr)
			if err != nil {
				continue
			}
			ips = append(ips, ip)
		}
	}
	return ips
}

// InterfaceIPs get public IP addresses by interface name
func InterfaceIPs(name string) ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		if iface.Name != name {
			continue
		}
		for _, addr := range addrs {
			ip, _, err := ExtractAddr(addr)
			if err != nil {
				continue
			}
			ips = append(ips, ip)
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("not found the ip of interface %s", name)
	}
	return ips, nil
}

// InterfaceIPv4 get a public IPv4 address
func InterfaceIPv4(name string) ([]net.IP, error) {
	ips, err := InterfaceIPs(name)
	if err != nil {
		return nil, err
	}
	var r []net.IP
	for _, ip := range ips {
		ip = ip.To4()
		if len(ip) == 0 {
			continue
		}
		r = append(r, ip)
	}
	return r, nil
}

// IsLocalIPAddr 检测 IP 地址字符串是否是内网地址
func IsLocalIPAddr(ip string) bool {
	return IsLocalIP(net.ParseIP(ip))
}

// IsLocalIP 检测 IP 地址是否是内网地址
// 通过直接对比ip段范围效率更高
func IsLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	if ip.IsPrivate() {
		return true
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return (ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip = strings.TrimSpace(ip); ip != "" && !IsLocalIPAddr(ip) {
			return ip
		}
	}

	if ip = strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" && !IsLocalIPAddr(ip) {
		return ip
	}

	if ip = RemoteIP(r); !IsLocalIPAddr(ip) {
		return ip
	}

	return ""
}

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// IPString2Long 把ip字符串转为数值
func IPString2Long(ip string) (uint, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IPString 把数值转为ip字符串
func Long2IPString(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}

// IP2Long 把net.IP转为数值
func IP2Long(ip net.IP) (uint, error) {
	b := ip.To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IP 把数值转为net.IP
func Long2IP(i uint) (net.IP, error) {
	if i > math.MaxUint32 {
		return nil, errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip, nil
}
