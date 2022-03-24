package v1

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

type IpnetSet []net.IPNet

type IpSet struct {
	m map[string]bool
	s []net.IP
}

func (ipnets IpnetSet) Decide(ips *IpSet) string {
	var ok bool
	if ips == nil {
		return ""
	}
	for _, ip := range ips.s {
		ok = false
		for _, subnet := range ipnets {
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

func nextForignIp(data []byte) (remoteIp net.IP, moreData []byte) {
	var i int
	var b byte
	for i, b = range data {
		if b == 0xA { // get to end of line
			break
		}
	}
	data = data[i:]
	for i, b = range data {
		if b == 0x3A { // get to colon after sl
			break
		}
	}
	if len(data) < i+20 {
		return
	}
	data = data[i+1:]
	for i, b = range data {
		if b == 0x3A { // get to colon after localIp
			break
		}
	}
	if len(data) < i+20 {
		return
	}
	data = data[i+6:]

	for i, b = range data {
		if b == 0x3A { // get to colon after remoteIp
			break
		}
	}
	moreData = data[i:]

	ipstr := string(data[:i])
	fmt.Printf("Proccessing IP %s\n", ipstr)
	var ip net.IP
	if len(ipstr) == 8 { //ipv4
		ip = make(net.IP, net.IPv4len)
		v, err := strconv.ParseUint(ipstr, 16, 32)
		if err != nil {
			return
		}
		binary.LittleEndian.PutUint32(ip, uint32(v))
		fmt.Printf("Proccessed IPv4 %s\n", ip.String())
	} else if len(ipstr) == 32 { //ipv6
		ip = make(net.IP, net.IPv6len)
		for i := 0; i < 16; i += 4 {
			u, err := strconv.ParseUint(ipstr[0:8], 16, 32)
			if err != nil {
				return
			}
			binary.LittleEndian.PutUint32(ip[i:i+4], uint32(u))
			ipstr = ipstr[8:] //skip 8 bytes
		}
		fmt.Printf("Proccessed IPv6 %s\n", ip.String())
	} else {
		fmt.Printf("Proccessed skipped IP structrue\n")
		return
	}
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() {
		return
	}
	remoteIp = ip
	fmt.Printf("Adding IP %s (not unspecified, private or loopback!)\n", ip.String())
	return
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
	ips.s = make([]net.IP, 1)
	ips.m[ip.String()] = true
	ips.s[0] = ip
	return
}

func IpSetFromProc(protocol string) (ips *IpSet) {
	procfile := "/proc/net/" + protocol
	data, err := ioutil.ReadFile(procfile)
	if err != nil {
		fmt.Printf("error while reading %s: %s\n", procfile, err.Error())
		// Used for development and debugging on macos - remove TODO
		procfile = "/tmp" + procfile
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
	ips.s = make([]net.IP, len(ips.m))
	i := 0
	for ipstr := range ips.m {
		ips.s[i] = net.ParseIP(ipstr)
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

func (in IpnetSet) DeepCopyInto(out *IpnetSet) {
	cpy := make([]net.IPNet, len(in))
	for i, v := range in {
		copy(cpy[i].IP, v.IP)
		copy(cpy[i].Mask, v.Mask)
	}
	*out = cpy
	//*out = (IpnetSet)(cpy)
}
