package v1

import (
	"bytes"
	"fmt"
	"strings"
)

type PodPile struct {
	Tcp4Peers     IpPile `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     IpPile `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers IpPile `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     IpPile `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     IpPile `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers IpPile `json:"udplite6peers"` // from /proc/net/udpline6
}

type PodConfig struct {
	tcp4Peers     CidrSet
	udp4Peers     CidrSet
	udplite4Peers CidrSet
	tcp6Peers     CidrSet
	udp6Peers     CidrSet
	udplite6Peers CidrSet

	Tcp4Peers     []string `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     []string `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers []string `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     []string `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     []string `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers []string `json:"udplite6peers"` // from /proc/net/udpline6
}

type PodProfile struct {
	// from local /proc/net (same net namespace)
	Tcp4Peers     *IpSet `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     *IpSet `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers *IpSet `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     *IpSet `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     *IpSet `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers *IpSet `json:"udplite6peers"` // from /proc/net/udpline6

	// The below requires sharing Process Namespace...
	ProcessCount uint32        `json:"processcount"` // from /proc/[1-9] - require ssh to container
	Processes    []string      `json:"processes"`    // from /proc/*/cmdline - require ssh to container
	FdCount      uint32        `json:"fdcount"`      // from /proc/<PID>/fd - require ssh to container
	IoRchar      KeyValProfile `json:"iorchar"`      // from /proc/<PID>/io - cycled on 32 bits - require ssh to container
	IoWchar      KeyValProfile `json:"iowchar"`      // from /proc/<PID>/io - cycled on 32 bits - require ssh to container
}

func (p *PodPile) Add(pp *PodProfile) {
	p.Tcp4Peers.Add(pp.Tcp4Peers)
	p.Udp4Peers.Add(pp.Udp4Peers)
	p.Udplite4Peers.Add(pp.Udplite4Peers)
	p.Tcp6Peers.Add(pp.Tcp6Peers)
	p.Udp6Peers.Add(pp.Udp6Peers)
	p.Udplite6Peers.Add(pp.Udplite6Peers)
}

func (p *EnvelopPile) Clear() {
	p.CompletionTime = make([]uint8, 0, 1)
	p.ResponseTime = make([]uint8, 0, 1)
}

func (p *PodPile) Clear() {
	p.Tcp4Peers.Clear()
	p.Udp4Peers.Clear()
	p.Udplite4Peers.Clear()
	p.Tcp6Peers.Clear()
	p.Udp6Peers.Clear()
	p.Udplite6Peers.Clear()
}

func (p *PodPile) Append(a *PodPile) {
	p.Tcp4Peers.Append(&a.Tcp4Peers)
	p.Udp4Peers.Append(&a.Udp4Peers)
	p.Udplite4Peers.Append(&a.Udplite4Peers)
	p.Tcp6Peers.Append(&a.Tcp6Peers)
	p.Udp6Peers.Append(&a.Udp6Peers)
	p.Udplite6Peers.Append(&a.Udplite6Peers)
}

func (config *PodConfig) Reconcile() {
	config.tcp4Peers = GetCidrsFromList(config.Tcp4Peers)
	config.udp4Peers = GetCidrsFromList(config.Udp4Peers)
	config.udplite4Peers = GetCidrsFromList(config.Udplite4Peers)
	config.tcp6Peers = GetCidrsFromList(config.Tcp6Peers)
	config.udp6Peers = GetCidrsFromList(config.Udp6Peers)
	config.udplite6Peers = GetCidrsFromList(config.Udplite6Peers)
}

func (config *PodConfig) Learn(p *PodPile) {
	config.tcp4Peers = GetCidrsFromIpList(p.Tcp4Peers.List)
	config.udp4Peers = GetCidrsFromIpList(p.Udp4Peers.List)
	config.udplite4Peers = GetCidrsFromIpList(p.Udplite4Peers.List)
	config.tcp6Peers = GetCidrsFromIpList(p.Tcp6Peers.List)
	config.udp6Peers = GetCidrsFromIpList(p.Udp6Peers.List)
	config.udplite6Peers = GetCidrsFromIpList(p.Udplite6Peers.List)
	config.Tcp4Peers = config.tcp4Peers.Strings()
	config.Udp4Peers = config.udp4Peers.Strings()
	config.Udplite4Peers = config.udplite4Peers.Strings()
	config.Tcp6Peers = config.tcp6Peers.Strings()
	config.Udp6Peers = config.udp6Peers.Strings()
	config.Udplite6Peers = config.udplite6Peers.Strings()
}

func (config *PodConfig) Merge(m *PodConfig) {
	config.tcp4Peers = GetCidrsFromIpList(m.Tcp4Peers)
	config.udp4Peers = GetCidrsFromIpList(m.Udp4Peers)
	config.udplite4Peers = GetCidrsFromIpList(m.Udplite4Peers)
	config.tcp6Peers = GetCidrsFromIpList(m.Tcp6Peers)
	config.udp6Peers = GetCidrsFromIpList(m.Udp6Peers)
	config.udplite6Peers = GetCidrsFromIpList(m.Udplite6Peers)
	config.Tcp4Peers = config.tcp4Peers.Strings()
	config.Udp4Peers = config.udp4Peers.Strings()
	config.Udplite4Peers = config.udplite4Peers.Strings()
	config.Tcp6Peers = config.tcp6Peers.Strings()
	config.Udp6Peers = config.udp6Peers.Strings()
	config.Udplite6Peers = config.udplite6Peers.Strings()
}

func (config *PodConfig) Decide(pp *PodProfile) string {
	//TBD move from string to net.IP and net.subnet, depend on a new IpSet rather than Set
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

	var ret string
	ret = config.tcp4Peers.Decide(pp.Tcp4Peers)
	if ret != "" {
		return fmt.Sprintf("Tcp4Peers: %s", ret)
	}
	ret = config.udp4Peers.Decide(pp.Udp4Peers)
	if ret != "" {
		return fmt.Sprintf("Udp4Peers: %s", ret)
	}
	ret = config.udplite4Peers.Decide(pp.Udplite4Peers)
	if ret != "" {
		return fmt.Sprintf("Udplite4Peers: %s", ret)
	}
	ret = config.tcp6Peers.Decide(pp.Tcp6Peers)
	if ret != "" {
		return fmt.Sprintf("Tcp6Peers: %s", ret)
	}
	ret = config.udp6Peers.Decide(pp.Udp6Peers)
	if ret != "" {
		return fmt.Sprintf("Udp6Peers: %s", ret)
	}
	ret = config.udplite6Peers.Decide(pp.Udplite6Peers)
	if ret != "" {
		return fmt.Sprintf("Udplite6Peers: %s", ret)
	}
	return ""
}

func (config *PodConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Tcp4Peers: %v,\n", config.Tcp4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udp4Peers: %v,\n", config.Udp4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udplite4Peers: %v,\n", config.Udplite4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Tcp6Peers: %v,\n", config.Tcp6Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udp6Peers: %v,\n", config.Udp6Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udplite6Peers: %v,\n", config.Udplite6Peers))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

// Profile timestamps and /proc
func (pp *PodProfile) Profile() {
	pp.Tcp4Peers = IpSetFromProc("tcp")
	pp.Udp4Peers = IpSetFromProc("udp")
	pp.Udplite4Peers = IpSetFromProc("udplite")
	pp.Tcp6Peers = IpSetFromProc("tcp6")
	pp.Udp6Peers = IpSetFromProc("udp6")
	pp.Udplite6Peers = IpSetFromProc("udplite6")
}

func (pp *PodProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Tcp4Peers: %v\n", pp.Tcp4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udp4Peers: %v\n", pp.Udp4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udplite4Peers: %v\n", pp.Udplite4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Tcp6Peers: %v\n", pp.Tcp6Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udp6Peers: %v\n", pp.Udp6Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udplite6Peers: %v\n", pp.Udplite6Peers))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (pp *PodPile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Tcp4Peers: %v\n", pp.Tcp4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udp4Peers: %v\n", pp.Udp4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udplite4Peers: %v\n", pp.Udplite4Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Tcp6Peers: %v\n", pp.Tcp6Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udp6Peers: %v\n", pp.Udp6Peers))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Udplite6Peers: %v\n", pp.Udplite6Peers))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}
