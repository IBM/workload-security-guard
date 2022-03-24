package v1

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

type ProcessPile struct {
	ResponseTime   []uint8 `json:"responsetime"`
	CompletionTime []uint8 `json:"completiontime"`
	Tcp4Peers      *IpSet  `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers      *IpSet  `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers  *IpSet  `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers      *IpSet  `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers      *IpSet  `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers  *IpSet  `json:"udplite6peers"` // from /proc/net/udpline6
}

type ProcessConfig struct {
	ResponseTime   U8MinmaxSlice `json:"responsetime"`
	CompletionTime U8MinmaxSlice `json:"completiontime"`
	Tcp4Peers      IpnetSet      `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers      IpnetSet      `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers  IpnetSet      `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers      IpnetSet      `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers      IpnetSet      `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers  IpnetSet      `json:"udplite6peers"` // from /proc/net/udpline6
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

func (config *ProcessConfig) Normalize() {
	config.ResponseTime = append(config.ResponseTime, U8Minmax{0, 0})
	config.CompletionTime = append(config.CompletionTime, U8Minmax{0, 0})
	config.Tcp4Peers = make([]net.IPNet, 0)
	config.Udp4Peers = make([]net.IPNet, 0)
	config.Udplite4Peers = make([]net.IPNet, 0)
	config.Tcp6Peers = make([]net.IPNet, 0)
	config.Udp6Peers = make([]net.IPNet, 0)
	config.Udplite6Peers = make([]net.IPNet, 0)

}

func (config *ProcessConfig) Decide(pp *ProcessProfile) string {
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
	ret = config.Tcp4Peers.Decide(pp.Tcp4Peers)
	if ret != "" {
		return fmt.Sprintf("Tcp4Peers: %s", ret)
	}
	ret = config.Udp4Peers.Decide(pp.Udp4Peers)
	if ret != "" {
		return fmt.Sprintf("Udp4Peers: %s", ret)
	}
	ret = config.Udplite4Peers.Decide(pp.Udplite4Peers)
	if ret != "" {
		return fmt.Sprintf("Udplite4Peers: %s", ret)
	}
	ret = config.Tcp6Peers.Decide(pp.Tcp6Peers)
	if ret != "" {
		return fmt.Sprintf("Tcp6Peers: %s", ret)
	}
	ret = config.Udp6Peers.Decide(pp.Udp6Peers)
	if ret != "" {
		return fmt.Sprintf("Udp6Peers: %s", ret)
	}
	ret = config.Udplite6Peers.Decide(pp.Udplite6Peers)
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
