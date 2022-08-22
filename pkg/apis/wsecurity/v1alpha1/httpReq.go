package v1alpha1

import (
	"fmt"
	"net"
	"net/http"
)

// FUTURE Add support for Trailers, review Req to find more security data
//////////////////// ReqProfile ////////////////

// Exposes ValueProfile interface
type ReqProfile struct {
	ClientIp      IpSetProfile     `json:"cip"`           // 192.168.32.1
	HopIp         IpSetProfile     `json:"hopip"`         // 1.2.3.4
	Method        SetProfile       `json:"method"`        // GET
	Proto         SetProfile       `json:"proto"`         // "HTTP/1.1"
	MediaType     MediaTypeProfile `json:"mediatype"`     // "text/html"
	ContentLength CountProfile     `json:"contentlength"` // 0
	Url           UrlProfile       `json:"url"`
	Qs            QueryProfile     `json:"qs"`
	Headers       HeadersProfile   `json:"headers"`
}

func (profile *ReqProfile) Profile(args ...interface{}) {
	req := args[0].(*http.Request)
	cip := args[1].(net.IP)

	var hopIpStr string

	profile.ClientIp.Profile(cip)

	if forwarded, ok := req.Header["X-Forwarded-For"]; ok {

		if numHops := len(forwarded); numHops > 0 {
			hopIpStr = forwarded[numHops-1]
		}
	}
	// FUTURE: Add support for rfc7239 "forwarded" header here
	profile.HopIp.Profile(hopIpStr)

	profile.MediaType.Profile(req.Header.Get("Content-Type"))
	profile.Qs.Profile(req.URL.Query())

	profile.Method.Profile(req.Method)
	profile.Proto.Profile(req.Proto)
	profile.Url.Profile(req.URL.Path)
	profile.Headers.Profile(req.Header)

	length := req.ContentLength
	if length > 0 {
		var log2length uint8

		for length > 0 {
			log2length++
			length >>= 1
		}
		profile.ContentLength = CountProfile(log2length)
	}
}

//////////////////// ReqPile ////////////////

// Exposes ValuePile interface
type ReqPile struct {
	ClientIp      IpSetPile     `json:"cip"`           // 192.168.32.1
	HopIp         IpSetPile     `json:"hopip"`         // 1.2.3.4
	Method        SetPile       `json:"method"`        // GET
	Proto         SetPile       `json:"proto"`         // "HTTP/1.1"
	MediaType     MediaTypePile `json:"mediatype"`     // "text/html"
	ContentLength CountPile     `json:"contentlength"` // 0
	Url           UrlPile       `json:"url"`
	Qs            QueryPile     `json:"qs"`
	Headers       HeadersPile   `json:"headers"`
}

func (pile *ReqPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*ReqProfile)
	pile.ClientIp.Add(&profile.ClientIp)
	pile.HopIp.Add(&profile.HopIp)
	pile.Method.Add(&profile.Method)
	pile.Proto.Add(&profile.Proto)
	pile.MediaType.Add(&profile.MediaType)
	pile.ContentLength.Add(&profile.ContentLength)
	pile.Url.Add(&profile.Url)
	pile.Qs.Add(&profile.Qs)
	pile.Headers.Add(&profile.Headers)
}

func (pile *ReqPile) Clear() {
	pile.ClientIp.Clear()
	pile.Method.Clear()
	pile.Proto.Clear()
	pile.MediaType.Clear()
	pile.ContentLength.Clear()
	pile.Url.Clear()
	pile.Qs.Clear()
	pile.Headers.Clear()
}

func (pile *ReqPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*ReqPile)
	pile.ClientIp.Merge(&otherPile.ClientIp)
	pile.Method.Merge(&otherPile.Method)
	pile.Proto.Merge(&otherPile.Proto)
	pile.MediaType.Merge(&otherPile.MediaType)
	pile.ContentLength.Merge(&otherPile.ContentLength)
	pile.Url.Merge(&otherPile.Url)
	pile.Qs.Merge(&otherPile.Qs)
	pile.Headers.Merge(&otherPile.Headers)
}

//////////////////// ReqConfig ////////////////

// Exposes ValueConfig interface
type ReqConfig struct {
	ClientIp      IpSetConfig     `json:"cip"`           // subnets for external IPs (normally empty)
	HopIp         IpSetConfig     `json:"hopip"`         // subnets for external IPs
	Method        SetConfig       `json:"method"`        // GET
	Proto         SetConfig       `json:"proto"`         // "HTTP/1.1"
	MediaType     MediaTypeConfig `json:"mediatype"`     // "text/html"
	ContentLength CountConfig     `json:"contentlength"` // 0
	Url           UrlConfig       `json:"url"`
	Qs            QueryConfig     `json:"qs"`
	Headers       HeadersConfig   `json:"headers"`
}

func (config *ReqConfig) Learn(valPile ValuePile) {
	pile := valPile.(*ReqPile)
	config.ClientIp.Learn(&pile.ClientIp)
	config.HopIp.Learn(&pile.HopIp)
	config.Method.Learn(&pile.Method)
	config.Proto.Learn(&pile.Proto)
	config.MediaType.Learn(&pile.MediaType)
	config.ContentLength.Learn(&pile.ContentLength)
	config.Headers.Learn(&pile.Headers)
	config.Qs.Learn(&pile.Qs)
	config.Url.Learn(&pile.Url)
}

func (config *ReqConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*ReqConfig)
	config.ClientIp.Fuse(&otherConfig.ClientIp)
	config.HopIp.Fuse(&otherConfig.HopIp)
	config.Method.Fuse(&otherConfig.Method)
	config.Proto.Fuse(&otherConfig.Proto)
	config.MediaType.Fuse(&otherConfig.MediaType)
	config.ContentLength.Fuse(&otherConfig.ContentLength)
	config.Headers.Fuse(&otherConfig.Headers)
	config.Qs.Fuse(&otherConfig.Qs)
	config.Url.Fuse(&otherConfig.Url)
}

func (config *ReqConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*ReqProfile)

	if ret := config.Url.Decide(&profile.Url); ret != "" {
		return fmt.Sprintf("Url: %s", ret)
	}

	if ret := config.Qs.Decide(&profile.Qs); ret != "" {
		return fmt.Sprintf("QueryString: %s", ret)
	}

	if ret := config.Headers.Decide(&profile.Headers); ret != "" {
		return fmt.Sprintf("Headers: %s", ret)
	}

	if ret := config.ClientIp.Decide(&profile.ClientIp); ret != "" {
		return fmt.Sprintf("ClientIp: %s", ret)
	}

	if ret := config.HopIp.Decide(&profile.HopIp); ret != "" {
		return fmt.Sprintf("HopIp: %s", ret)
	}

	//methodSet := make(Set)
	//methodSet[rp.Method] = true

	if ret := config.Method.Decide(&profile.Method); ret != "" {
		return fmt.Sprintf("Method: %s", ret)
	}
	//protoSet := make(Set)
	//protoSet[rp.Proto] = true

	if ret := config.Proto.Decide(&profile.Proto); ret != "" {
		return fmt.Sprintf("Proto: %s", ret)
	}

	if ret := config.MediaType.Decide(&profile.MediaType); ret != "" {
		return fmt.Sprintf("MediaType: %s", ret)
	}

	if ret := config.ContentLength.Decide(&profile.ContentLength); ret != "" {
		return fmt.Sprintf("ContentLength: %s", ret)
	}
	return ""
}
