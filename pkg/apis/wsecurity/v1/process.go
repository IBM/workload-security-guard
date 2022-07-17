package v1

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type ProcessPile struct {
	ResponseTime   []uint8 `json:"responsetime"`
	CompletionTime []uint8 `json:"completiontime"`
	Tcp4Peers      IpPile  `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers      IpPile  `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers  IpPile  `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers      IpPile  `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers      IpPile  `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers  IpPile  `json:"udplite6peers"` // from /proc/net/udpline6
}

type ProcessConfig struct {
	tcp4Peers     CidrSet
	udp4Peers     CidrSet
	udplite4Peers CidrSet
	tcp6Peers     CidrSet
	udp6Peers     CidrSet
	udplite6Peers CidrSet

	ResponseTime   U8MinmaxSlice `json:"responsetime"`
	CompletionTime U8MinmaxSlice `json:"completiontime"`
	Tcp4Peers      []string      `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers      []string      `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers  []string      `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers      []string      `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers      []string      `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers  []string      `json:"udplite6peers"` // from /proc/net/udpline6
}

type ProcessProfile struct {
	ResponseTime   uint8 `json:"responsetime"`
	CompletionTime uint8 `json:"completiontime"`

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

func (p *ProcessPile) Add(pp *ProcessProfile) {
	fmt.Println("Add")
	p.CompletionTime = append(p.CompletionTime, pp.CompletionTime)
	p.ResponseTime = append(p.ResponseTime, pp.ResponseTime)
	p.Tcp4Peers.Add(pp.Tcp4Peers)
	p.Udp4Peers.Add(pp.Udp4Peers)
	p.Udplite4Peers.Add(pp.Udplite4Peers)
	p.Tcp6Peers.Add(pp.Tcp6Peers)
	p.Udp6Peers.Add(pp.Udp6Peers)
	p.Udplite6Peers.Add(pp.Udplite6Peers)
}

func (p *ProcessPile) Clear() {
	fmt.Println("Clear")
	p.CompletionTime = make([]uint8, 0, 1)
	p.ResponseTime = make([]uint8, 0, 1)
	p.Tcp4Peers.Clear()
	p.Udp4Peers.Clear()
	p.Udplite4Peers.Clear()
	p.Tcp6Peers.Clear()
	p.Udp6Peers.Clear()
	p.Udplite6Peers.Clear()
}
func (p *ProcessPile) Append(a *ProcessPile) {
	fmt.Println("Append")
	p.CompletionTime = append(p.CompletionTime, a.CompletionTime...)
	p.ResponseTime = append(p.ResponseTime, a.ResponseTime...)
	p.Tcp4Peers.Append(&a.Tcp4Peers)
	p.Udp4Peers.Append(&a.Udp4Peers)
	p.Udplite4Peers.Append(&a.Udplite4Peers)
	p.Tcp6Peers.Append(&a.Tcp6Peers)
	p.Udp6Peers.Append(&a.Udp6Peers)
	p.Udplite6Peers.Append(&a.Udplite6Peers)
}

func (config *ProcessConfig) Reconcile() {
	fmt.Printf("Reconcile Process %v\n", config)
	config.tcp4Peers = GetCidrsFromList(config.Tcp4Peers)
	config.udp4Peers = GetCidrsFromList(config.Udp4Peers)
	config.udplite4Peers = GetCidrsFromList(config.Udplite4Peers)
	config.tcp6Peers = GetCidrsFromList(config.Tcp6Peers)
	config.udp6Peers = GetCidrsFromList(config.Udp6Peers)
	config.udplite6Peers = GetCidrsFromList(config.Udplite6Peers)
}

func (config *ProcessConfig) Learn(p *ProcessPile) {
	fmt.Println("Learn")
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
	config.CompletionTime.Learn(p.CompletionTime)
	config.ResponseTime.Learn(p.ResponseTime)
	//fmt.Printf("Config %v\n", config)
}

func (config *ProcessConfig) Merge(m *ProcessConfig) {
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
	config.CompletionTime.Merge(m.CompletionTime)
	config.ResponseTime.Merge(m.ResponseTime)
	fmt.Printf("Config %v\n", config)
}

func (config *ProcessConfig) Normalize() {
	fmt.Println("Normalize")
	config.ResponseTime = append(config.ResponseTime, U8Minmax{0, 0})
	config.CompletionTime = append(config.CompletionTime, U8Minmax{0, 0})
	//config.Tcp4Peers = new(CidrSet)
	//config.Udp4Peers = new(CidrSet)
	//config.Udplite4Peers = new(CidrSet)
	//config.Tcp6Peers = new(CidrSet)
	//config.Udp6Peers = new(CidrSet)
	//config.Udplite6Peers = new(CidrSet)

}

func (config *ProcessConfig) Decide(pp *ProcessProfile) string {
	fmt.Println("Decide Process")
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
	ret = config.ResponseTime.Decide(pp.ResponseTime)
	if ret != "" {
		return fmt.Sprintf("ResponseTime: %s", ret)
	}
	ret = config.CompletionTime.Decide(pp.CompletionTime)
	if ret != "" {
		return fmt.Sprintf("CompletionTime: %s", ret)
	}
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

func (config *ProcessConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ResponseTime: %s,\n", config.ResponseTime.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  CompletionTime: %s,\n", config.CompletionTime.Marshal()))
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

// Allow typical values - use for development but not in production
func (config *ProcessConfig) AddTypicalVal() {
	config.ResponseTime = make([]U8Minmax, 1)
	config.ResponseTime[0].Max = 60
	config.CompletionTime = make([]U8Minmax, 1)
	config.CompletionTime[0].Max = 120
}

// Profile timestamps and /proc
func (pp *ProcessProfile) Profile(reqTime time.Time, respTime time.Time, endTime time.Time) {
	fmt.Println("Process Profile")
	pp.Tcp4Peers = IpSetFromProc("tcp")
	pp.Udp4Peers = IpSetFromProc("udp")
	pp.Udplite4Peers = IpSetFromProc("udplite")
	pp.Tcp6Peers = IpSetFromProc("tcp6")
	pp.Udp6Peers = IpSetFromProc("udp6")
	pp.Udplite6Peers = IpSetFromProc("udplite6")

	completionTime := endTime.Sub(reqTime).Seconds()
	if completionTime > 255 {
		pp.CompletionTime = 255
	} else {
		pp.CompletionTime = uint8(completionTime)
	}

	responseTime := respTime.Sub(reqTime).Seconds()
	if responseTime > 255 {
		pp.ResponseTime = 255
	} else {
		pp.ResponseTime = uint8(responseTime)
	}
}

func (pp *ProcessProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ResponseTime: %d\n", pp.ResponseTime))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  CompletionTime: %d\n", pp.CompletionTime))
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

func (pp *ProcessPile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ResponseTime: %v\n", pp.ResponseTime))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  CompletionTime: %v\n", pp.CompletionTime))
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
