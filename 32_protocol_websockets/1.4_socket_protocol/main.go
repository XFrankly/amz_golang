package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println(net.ParseIP("192.0.2.1"))
	fmt.Println(net.ParseIP("2001:db8::68"))
	ips := net.ParseIP("10.2.1.26")
	fmt.Println(net.ParseIP("192.0.2.22.11"))
	fmt.Println(ips, ips.DefaultMask())
	fmt.Println(ips.Equal(ips)) // ipv4 是否相同
	//// 是否全球单播地址 在RFC4632 RFC4291中定义,即公网地址， 和是否ipv6本地多播地址
	// // IPv4 定向广播地址除外。即使 ip 在 IPv4 私有地址空间或本地 IPv6 单播地址空间中，它也会返回 true。
	ipv6Global := net.ParseIP("2000::")
	ipv6UniqLocal := net.ParseIP("2000::")
	ipv6Multi := net.ParseIP("FF00::")
	ipv6InterfaceLocalMulti := net.ParseIP("ff01::1")

	ipv4Private := net.ParseIP("10.255.0.0")
	ipv4Public := net.ParseIP("8.8.8.8")
	ipv4Broadcase := net.ParseIP("255.255.255.255")

	fmt.Println(ipv6Global, ipv6Global.IsGlobalUnicast(), ipv6Global.IsInterfaceLocalMulticast())
	fmt.Println(ipv6UniqLocal, ipv6UniqLocal.IsGlobalUnicast(), ipv6UniqLocal.IsInterfaceLocalMulticast())
	fmt.Println(ipv6Multi, ipv6Multi.IsGlobalUnicast(), ipv6Multi.IsInterfaceLocalMulticast())
	fmt.Println(ipv6InterfaceLocalMulti, ipv6InterfaceLocalMulti.IsGlobalUnicast(), ipv6InterfaceLocalMulti.IsInterfaceLocalMulticast())

	fmt.Println("\n")
	fmt.Println(ipv4Private, ipv4Private.IsGlobalUnicast(), ipv4Private.IsInterfaceLocalMulticast())
	fmt.Println(ipv4Public, ipv4Public.IsGlobalUnicast(), ipv4Public.IsInterfaceLocalMulticast())
	fmt.Println(ipv4Broadcase, ipv4Broadcase.IsGlobalUnicast(), ipv4Broadcase.IsInterfaceLocalMulticast())
}
