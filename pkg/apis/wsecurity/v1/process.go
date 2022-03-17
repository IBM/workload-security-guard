package v1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

type ProcessPile struct {
	ResponseTime   []uint8 `json:"responsetime"`
	CompletionTime []uint8 `json:"completiontime"`
	Tcp4Peers      Set     `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers      Set     `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers  Set     `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers      Set     `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers      Set     `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers  Set     `json:"udplite6peers"` // from /proc/net/udpline6
}

type ProcessConfig struct {
	ResponseTime   U8MinmaxSlice `json:"responsetime"`
	CompletionTime U8MinmaxSlice `json:"completiontime"`
	Tcp4Peers      Set           `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers      Set           `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers  Set           `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers      Set           `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers      Set           `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers  Set           `json:"udplite6peers"` // from /proc/net/udpline6
}

type ProcessProfile struct {
	ResponseTime   uint8 `json:"responsetime"`
	CompletionTime uint8 `json:"completiontime"`

	// from local /proc/net (same net namespace)
	Tcp4Peers     Set `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     Set `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers Set `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     Set `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     Set `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers Set `json:"udplite6peers"` // from /proc/net/udpline6

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
}

func (config *ProcessConfig) Decide(pp *ProcessProfile) string {
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

func nextForignIp(data []byte) (remoteIp string, moreData []byte) {
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
	if len(ipstr) == 8 {

		v, err := strconv.ParseUint(ipstr, 16, 32)
		if err != nil {
			//fmt.Printf("err %v\n", err)
			return
		}
		ip := make(net.IP, net.IPv4len)
		binary.LittleEndian.PutUint32(ip, uint32(v))

		remoteIp = ip.String()
		return
	}
	if len(ipstr) == 32 {
		ip := make(net.IP, net.IPv6len)
		const grpLen = 4
		i, j := 0, 4
		for len(ipstr) != 0 {
			grp := ipstr[0:8]
			u, err := strconv.ParseUint(grp, 16, 32)
			binary.LittleEndian.PutUint32(ip[i:j], uint32(u))
			if err != nil {
				//fmt.Printf("err %v\n", err)
				return
			}
			i, j = i+grpLen, j+grpLen
			ipstr = ipstr[8:]
		}

		remoteIp = ip.String()
	}
	return
}

func periodicalNet(protocol string) (m map[string]bool) {
	m = make(map[string]bool)
	procfile := "/proc/net/" + protocol
	data, err := ioutil.ReadFile(procfile)
	if err != nil {
		// Used for development and debugging on macos - remove TODO
		procfile = "/tmp" + procfile
		data, err = ioutil.ReadFile(procfile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	var ip string
	ip, data = nextForignIp(data)
	for data != nil {
		m[ip] = true
		ip, data = nextForignIp(data)
	}
	return
}

// Profile timestamps and /proc
func (pp *ProcessProfile) Profile(reqTime time.Time, respTime time.Time, endTime time.Time) {
	pp.Tcp4Peers = Set(periodicalNet("tcp"))
	pp.Udp4Peers = Set(periodicalNet("udp"))
	pp.Udplite4Peers = Set(periodicalNet("udplite"))
	pp.Tcp6Peers = Set(periodicalNet("tcp6"))
	pp.Udp6Peers = Set(periodicalNet("udp6"))
	pp.Udplite6Peers = Set(periodicalNet("udplite6"))

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
