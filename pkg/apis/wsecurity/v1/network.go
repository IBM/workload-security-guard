package v1

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type IpPile struct {
	List []string
	m    map[string]bool
}

type CidrSet []net.IPNet

type IpSet struct {
	list []net.IP
	m    map[string]bool
}

func (ipp *IpPile) Add(ips *IpSet) {
	if ipp.m == nil {
		ipp.m = make(map[string]bool, len(ips.list))
	}
	for ip := range ips.m {
		if !ipp.m[ip] {
			ipp.m[ip] = true
			ipp.List = append(ipp.List, ip)
		}
	}
}
func (ipp *IpPile) Clear() {
	ipp.m = make(map[string]bool)
	ipp.List = make([]string, 0, 1)
}

func (ipp *IpPile) Append(a *IpPile) {
	for _, ip := range a.List {
		if !ipp.m[ip] {
			ipp.m[ip] = true
			ipp.List = append(ipp.List, ip)
		}
	}
}

func (cidr CidrSet) Decide(ips *IpSet) string {
	var ok bool
	if ips == nil {
		return ""
	}
	fmt.Printf("func (%v) Decide(%v)\n", cidr, ips)
	for _, ip := range ips.list {
		if ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() {
			continue
		}
		ok = false
		for _, subnet := range cidr {
			if subnet.Contains(ip) {
				ok = true
				break
			}
		}
		if ok {
			continue
		}
		return fmt.Sprintf("IP %s not allowed", ip.String())
	}
	return ""
}

func nextForignIp(data_in []byte) (net.IP, []byte) {
	var i int
	var b byte
	var data []byte = data_in
	var ipstr string
NextLine:
	for {
		ipstr = ""
		// 1. Move forward in data, set ipstr to the next candidate
		for i, b = range data {
			if b == 0xA { // get to end of line
				data = data[i+1:]
				// 1a. moved to a new line
				for i, b = range data {
					if b == 0x3A { // get to colon after sl
						data = data[i+1:]
						// 1b. moved after first colon
						for i, b = range data {
							if b == 0x3A { // get to colon after localIp
								data = data[i+6:]
								// 1c. moved till after second colon + 6 bytes - it is where the ip starts
								for i, b = range data {
									if b == 0x3A { // get to colon after remoteIp
										ipstr = string(data[:i])
										data = data[i+1:]
										// 1d. moved till after third colon and placed the ip in ipstr
										break
									}
								}
								break
							}
						}
						break
					}
				}
				break
			}
		}

		// 2. Try to process ipstr
		//    We return nil if no more IPs or if ip has bad format
		//fmt.Printf("Proccessing IP %s\n", ipstr)

		var ip net.IP
		if len(ipstr) == 8 { //ipv4
			ip = make(net.IP, net.IPv4len)
			v, err := strconv.ParseUint(ipstr, 16, 32)
			if err != nil {
				return nil, nil
			}
			binary.LittleEndian.PutUint32(ip, uint32(v))
			//fmt.Printf("Proccessed IPv4 %s\n", ip.String())
		} else if len(ipstr) == 32 { //ipv6
			ip = make(net.IP, net.IPv6len)
			for i := 0; i < 16; i += 4 {
				u, err := strconv.ParseUint(ipstr[0:8], 16, 32)
				if err != nil {
					return nil, nil
				}
				binary.LittleEndian.PutUint32(ip[i:i+4], uint32(u))
				ipstr = ipstr[8:] //skip 8 bytes
			}
			//fmt.Printf("Proccessed IPv6 %s\n", ip.String())
		} else {
			fmt.Printf("Proccessed skipped IP structrue\n")
			return nil, nil
		}

		// 3. Success!! If ip of interest  - back to caller, else move to next line

		if ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() {
			continue NextLine
		}

		fmt.Printf("Adding IP %s (not unspecified, private or loopback!)\n", ip.String())
		return ip, data
	}

	//const grpLen = 4
	//i, j := 0, 4
	//i := 0

	//for i := 0; i < 16; i += 4 {
	//for len(ipstr) != 0 {
	//grp := ipstr[0:8] // next 8 bytes of IP
	//	u, err := strconv.ParseUint(ipstr[0:8], 16, 32)
	//	if err != nil {
	//fmt.Printf("err %v\n", err)
	//		return
	//	}
	//binary.LittleEndian.PutUint32(ip[i:j], uint32(u))
	//	binary.LittleEndian.PutUint32(ip[i:i+4], uint32(u))
	//i, j = i+grpLen, j+grpLen
	//i += 4
	//	ipstr = ipstr[8:] //skip 8 bytes
	//}

}

func IpSetFromIp(ip net.IP) (ips *IpSet) {
	ips = new(IpSet)
	ips.m = make(map[string]bool)
	ips.list = make([]net.IP, 0, 1)
	ips.m[ip.String()] = true
	if ip != nil {
		ips.list = append(ips.list, ip)
	}
	return
}

func IpSetFromProc(protocol string) (ips *IpSet) {
	procfile := "/proc/net/" + protocol
	data, err := ioutil.ReadFile(procfile)
	if err != nil {
		fmt.Printf("error while reading %s: %s\n", procfile, err.Error())
		// Used for development and debugging on macos - remove TODO
		procfile = "/Users/davidhadasmac16" + procfile
		data, err = ioutil.ReadFile(procfile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// create a set for collecting the Ips
	ips = new(IpSet)
	ips.m = make(map[string]bool)
	ip, data := nextForignIp(data)
	for data != nil {
		ips.m[ip.String()] = true
		ip, data = nextForignIp(data)
	}
	ips.list = make([]net.IP, len(ips.m))
	i := 0
	for ipstr := range ips.m {
		ips.list[i] = net.ParseIP(ipstr)
	}
	return
}

/*
	y := "1.2.4.6"
	x := "1.2.3.16/22"
	ipy := net.ParseIP(y)
	ipx, subnetx, _ := net.ParseCIDR(x)
	if ipx == nil {
		ipx = net.ParseIP(x)
		fmt.Printf("ipy == ipx %v\n", ipx.Equal(ipy))
	} else {
		fmt.Printf("ipy In Subnetx %v\n", subnetx.Contains(ipy))
	}
*/

func (in CidrSet) DeepCopyInto(out *CidrSet) {
	copy((*out), in)
	cpy := make([]net.IPNet, len(in))
	for i, v := range in {
		copy(cpy[i].IP, v.IP)
		copy(cpy[i].Mask, v.Mask)
	}
	(*out) = cpy
}

func (cidrs CidrSet) Strings() (s []string) {
	s = make([]string, len(cidrs))
	for i, cidr := range cidrs {
		s[i] = cidr.String()
	}
	return
}

func GetCidrsFromIpList(list []string) CidrSet {
	cidr := make([]net.IPNet, len(list))
	var n int
	var ipNet *net.IPNet
	var err error
	for _, v := range list {
		if strings.Contains(v, ":") {
			_, ipNet, err = net.ParseCIDR(v + "/128")
		} else {
			_, ipNet, err = net.ParseCIDR(v + "/32")
		}
		if err == nil {
			fmt.Printf("CIDRS found CIDR %v\n", ipNet)

			cidr[n].IP = make(net.IP, len(ipNet.IP))
			cidr[n].Mask = make(net.IPMask, len(ipNet.Mask))
			copy(cidr[n].IP, ipNet.IP)
			copy(cidr[n].Mask, ipNet.Mask)
			n++
			continue
		}
		fmt.Printf("Ilegal cidr %s is skipped during GetCidrsFromList\n", v)
	}
	fmt.Printf("CIDRS %v\n", cidr)
	return cidr
}

func GetCidrsFromList(list []string) CidrSet {
	cidr := make([]net.IPNet, len(list))
	var n int
	for _, v := range list {
		_, ipNet, err := net.ParseCIDR(v)
		if err == nil {
			fmt.Printf("CIDRS found CIDR %v\n", ipNet)

			cidr[n].IP = make(net.IP, len(ipNet.IP))
			cidr[n].Mask = make(net.IPMask, len(ipNet.Mask))
			copy(cidr[n].IP, ipNet.IP)
			copy(cidr[n].Mask, ipNet.Mask)
			n++
			continue
		}
		fmt.Printf("Ilegal cidr %s is skipped during GetCidrsFromList\n", v)
	}
	fmt.Printf("CIDRS %v\n", cidr)
	return cidr

}
