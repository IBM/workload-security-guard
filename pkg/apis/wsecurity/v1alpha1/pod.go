package v1alpha1

import (
	"fmt"
)

//////////////////// PodProfile ////////////////

// Exposes ValueProfile interface

type PodProfile struct {
	// from local /proc/net (same net namespace)
	Tcp4Peers     IpSetProfile `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     IpSetProfile `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers IpSetProfile `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     IpSetProfile `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     IpSetProfile `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers IpSetProfile `json:"udplite6peers"` // from /proc/net/udpline6

	// The below requires sharing Process Namespace...
	ProcessCount CountProfile  `json:"processcount"` // from /proc/[1-9] - require ssh to container
	Processes    []string      `json:"processes"`    // from /proc/*/cmdline - require ssh to container
	FdCount      CountProfile  `json:"fdcount"`      // from /proc/<PID>/fd - require ssh to container
	IoRchar      KeyValProfile `json:"iorchar"`      // from /proc/<PID>/io - cycled on 32 bits - require ssh to container
	IoWchar      KeyValProfile `json:"iowchar"`      // from /proc/<PID>/io - cycled on 32 bits - require ssh to container
}

func (profile *PodProfile) Profile(args ...interface{}) {
	profile.Tcp4Peers.Profile(IpNetFromProc("tcp"))
	profile.Udp4Peers.Profile(IpNetFromProc("udp"))
	profile.Udplite4Peers.Profile(IpNetFromProc("udplite"))
	profile.Tcp6Peers.Profile(IpNetFromProc("tcp6"))
	profile.Udp6Peers.Profile(IpNetFromProc("udp6"))
	profile.Udplite6Peers.Profile(IpNetFromProc("udplite6"))
}

//////////////////// PodPile ////////////////

// Exposes ValuePile interface
type PodPile struct {
	Tcp4Peers     IpSetPile `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     IpSetPile `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers IpSetPile `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     IpSetPile `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     IpSetPile `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers IpSetPile `json:"udplite6peers"` // from /proc/net/udpline6
}

func (pile *PodPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*PodProfile)

	pile.Tcp4Peers.Add(&profile.Tcp4Peers)
	pile.Udp4Peers.Add(&profile.Udp4Peers)
	pile.Udplite4Peers.Add(&profile.Udplite4Peers)
	pile.Tcp6Peers.Add(&profile.Tcp6Peers)
	pile.Udp6Peers.Add(&profile.Udp6Peers)
	pile.Udplite6Peers.Add(&profile.Udplite6Peers)
}

func (pile *PodPile) Clear() {
	pile.Tcp4Peers.Clear()
	pile.Udp4Peers.Clear()
	pile.Udplite4Peers.Clear()
	pile.Tcp6Peers.Clear()
	pile.Udp6Peers.Clear()
	pile.Udplite6Peers.Clear()
}

func (pile *PodPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*PodPile)

	pile.Tcp4Peers.Merge(&otherPile.Tcp4Peers)
	pile.Udp4Peers.Merge(&otherPile.Udp4Peers)
	pile.Udplite4Peers.Merge(&otherPile.Udplite4Peers)
	pile.Tcp6Peers.Merge(&otherPile.Tcp6Peers)
	pile.Udp6Peers.Merge(&otherPile.Udp6Peers)
	pile.Udplite6Peers.Merge(&otherPile.Udplite6Peers)
}

//////////////////// PodConfig ////////////////

// Exposes ValueConfig interface

type PodConfig struct {
	Tcp4Peers     IpSetConfig `json:"tcp4peers"`     // from /proc/net/tcp
	Udp4Peers     IpSetConfig `json:"udp4peers"`     // from /proc/net/udp
	Udplite4Peers IpSetConfig `json:"udplite4peers"` // from /proc/udpline
	Tcp6Peers     IpSetConfig `json:"tcp6peers"`     // from /proc/net/tcp6
	Udp6Peers     IpSetConfig `json:"udp6peers"`     // from /proc/net/udp6
	Udplite6Peers IpSetConfig `json:"udplite6peers"` // from /proc/net/udpline6
}

func (config *PodConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*PodProfile)

	var ret string
	ret = config.Tcp4Peers.Decide(&profile.Tcp4Peers)
	if ret != "" {
		return fmt.Sprintf("Tcp4Peers: %s", ret)
	}
	ret = config.Udp4Peers.Decide(&profile.Udp4Peers)
	if ret != "" {
		return fmt.Sprintf("Udp4Peers: %s", ret)
	}
	ret = config.Udplite4Peers.Decide(&profile.Udplite4Peers)
	if ret != "" {
		return fmt.Sprintf("Udplite4Peers: %s", ret)
	}
	ret = config.Tcp6Peers.Decide(&profile.Tcp6Peers)
	if ret != "" {
		return fmt.Sprintf("Tcp6Peers: %s", ret)
	}
	ret = config.Udp6Peers.Decide(&profile.Udp6Peers)
	if ret != "" {
		return fmt.Sprintf("Udp6Peers: %s", ret)
	}
	ret = config.Udplite6Peers.Decide(&profile.Udplite6Peers)
	if ret != "" {
		return fmt.Sprintf("Udplite6Peers: %s", ret)
	}
	return ""
}

func (config *PodConfig) Learn(valPile ValuePile) {
	pile := valPile.(*PodPile)

	config.Tcp4Peers.Learn(&pile.Tcp4Peers)
	config.Udp4Peers.Learn(&pile.Udp4Peers)
	config.Udplite4Peers.Learn(&pile.Udplite4Peers)
	config.Tcp6Peers.Learn(&pile.Tcp6Peers)
	config.Udp6Peers.Learn(&pile.Udp6Peers)
	config.Udplite6Peers.Learn(&pile.Udplite6Peers)
}

func (config *PodConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*PodConfig)

	config.Tcp4Peers.Fuse(&otherConfig.Tcp4Peers)
	config.Udp4Peers.Fuse(&otherConfig.Udp4Peers)
	config.Udplite4Peers.Fuse(&otherConfig.Udplite4Peers)
	config.Tcp6Peers.Fuse(&otherConfig.Tcp6Peers)
	config.Udp6Peers.Fuse(&otherConfig.Udp6Peers)
	config.Udplite6Peers.Fuse(&otherConfig.Udplite6Peers)
}
