package v1

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
)

type RespPile struct {
	Headers HeadersPile `json:"headers"`
}

type RespConfig struct {
	Headers HeadersConfig `json:"headers"`
}

type ReqPile struct {
	ClientIp      []net.IP `json:"cip"`           // 192.168.32.1
	HopIp         []net.IP `json:"hopip"`         // 1.2.3.4
	Method        []string `json:"method"`        // GET
	Proto         []string `json:"proto"`         // "HTTP/1.1"
	ContentLength []uint8  `json:"contentlength"` // 0

	Url     UrlPile     `json:"url"`
	Qs      QueryPile   `json:"qs"`
	Headers HeadersPile `json:"headers"`
}

type ReqConfig struct {
	ClientIp      IpnetSet      `json:"cip"`           // subnets for external IPs (normally empty)
	HopIp         IpnetSet      `json:"hopip"`         // subnets for external IPs
	Method        Set           `json:"method"`        // GET
	Proto         Set           `json:"proto"`         // "HTTP/1.1"
	ContentLength U8MinmaxSlice `json:"contentlength"` // 0

	Url     UrlConfig     `json:"url"`
	Qs      QueryConfig   `json:"qs"`
	Headers HeadersConfig `json:"headers"`
}

type ReqProfile struct {
	ClientIp      net.IP          `json:"cip"`           // 192.168.32.1
	HopIp         net.IP          `json:"hopip"`         // 1.2.3.4
	Method        string          `json:"method"`        // GET
	Proto         string          `json:"proto"`         // "HTTP/1.1"
	ContentLength uint8           `json:"contentlength"` // 0
	Url           *UrlProfile     `json:"url"`
	Qs            *QueryProfile   `json:"qs"`
	Headers       *HeadersProfile `json:"headers"`
}

type RespProfile struct {
	Headers *HeadersProfile `json:"headers"`
	// Trailers...
}

type UrlPile struct {
	Val      *SimpleValPile `json:"val"`
	Segments []uint8        `json:"segments"`
}

type UrlProfile struct {
	Scheme   string            `json:"scheme"` // http
	Val      *SimpleValProfile `json:"val"`
	Segments uint8             `json:"segments"`
}

type UrlConfig struct {
	Val      SimpleValConfig `json:"val"`
	Segments U8MinmaxSlice   `json:"segments"`
}

type QueryPile struct {
	Kv *KeyValPile `json:"kv"`
}

type QueryProfile struct {
	Kv *KeyValProfile `json:"kv"`
}

type QueryConfig struct {
	Kv KeyValConfig `json:"kv"`
}

type HeadersPile struct {
	Kv *KeyValPile `json:"kv"`
}

type HeadersProfile struct {
	Kv *KeyValProfile `json:"kv"`
}

type HeadersConfig struct {
	Kv KeyValConfig `json:"kv"`
}

func (p *UrlPile) Add(u *UrlProfile) {
	p.Segments = append(p.Segments, u.Segments)
	p.Val.Add(u.Val)
}

func (p *UrlPile) Clear() {
	p.Segments = make([]uint8, 0, 1)
	p.Val = new(SimpleValPile)
}

func (u *UrlProfile) Profile(path string) {
	segments := strings.Split(path, "/")
	numSegments := len(segments)
	if (numSegments > 0) && segments[0] == "" {
		segments = segments[1:]
		numSegments--
	}
	if (numSegments > 0) && segments[numSegments-1] == "" {
		numSegments--
		segments = segments[:numSegments]

	}
	cleanPath := strings.Join(segments, "")
	u.Val = new(SimpleValProfile)
	u.Val.Profile(cleanPath)

	if numSegments > 0xFF {
		numSegments = 0xFF
	}
	u.Segments = uint8(numSegments)
	fmt.Printf("Path %s, segments %v, len %d\n", path, segments, numSegments)
}

func (u *UrlProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Val: %s", u.Val.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Segments: %d", u.Segments))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *UrlConfig) Normalize() {
	config.Val.Normalize()
	config.Segments = append(config.Segments, U8Minmax{0, 0})
}

func (config *UrlConfig) Decide(u *UrlProfile) string {
	if str := config.Segments.Decide(u.Segments); str != "" {
		return fmt.Sprintf("Segmengs: %s", str)
	}

	if str := config.Val.Decide(u.Val); str != "" {
		return fmt.Sprintf("KeyVal: %s", str)
	}
	return ""
}

func (config *UrlConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  KeyVal: %s", config.Val.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Segments: %s", config.Segments.Marshal()))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

// Allow typical URL values - use for development but not in production
func (config *UrlConfig) AddTypicalVal() {
	config.Val.Spaces = make([]U8Minmax, 0, 1)
	config.Val.Unicodes = make([]U8Minmax, 0, 1)
	config.Val.NonReadables = make([]U8Minmax, 0, 1)
	config.Val.Letters = make([]U8Minmax, 0, 1)
	config.Val.Digits = make([]U8Minmax, 0, 1)
	config.Val.Sequences = make([]U8Minmax, 0, 1)
	//config.Val.Words = make([]U8Minmax, 0, 1)
	//config.Val.Numbers = make([]U8Minmax, 0, 1)

	config.Val.Spaces[0].Max = 0
	config.Val.Unicodes[0].Max = 0
	config.Val.NonReadables[0].Max = 0
	config.Val.Letters[0].Max = 64
	config.Val.Digits[0].Max = 64
	config.Val.Sequences[0].Max = 16
	//config.Val.Words[0].Max = 16
	//config.Val.Numbers[0].Max = 16
	//config.Val.FlagsL = 1 << SlashSlot
	config.Val.Flags = 1 << SlashSlot
	config.Segments = make([]U8Minmax, 0, 1)
	config.Segments[0].Max = 8
}

func (p *QueryPile) Add(q *QueryProfile) {
	p.Kv.Add(q.Kv)
}

func (p *QueryPile) Clear() {
	p.Kv.Clear()
}

func (q *QueryProfile) Profile(m map[string][]string) {
	q.Kv = new(KeyValProfile)
	q.Kv.Profile(m)
}

func (config *QueryConfig) Normalize() {
	config.Kv.Normalize()
}
func (q *QueryProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Kv: %s", q.Kv.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *QueryConfig) Decide(q *QueryProfile) string {
	str := config.Kv.Decide(q.Kv)
	if str == "" {
		return str
	}
	return fmt.Sprintf("KeyVal: %s", str)
}

// Allow typical query string values - use for development but not in production
func (config *QueryConfig) AddTypicalVal() {
	config.Kv.OtherKeynames = NewSimpleValConfig(0, 0, 0, 16, 16, 0, 4)
	config.Kv.OtherVals = NewSimpleValConfig(0, 0, 0, 32, 32, 0, 16)
}

func (config *QueryConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Kv: %s", config.Kv.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (p *HeadersPile) Add(h *HeadersProfile) {
	p.Kv.Add(h.Kv)
}
func (p *HeadersPile) Clear() {
	p.Kv.Clear()
}

func (h *HeadersProfile) Profile(m map[string][]string) {
	h.Kv = new(KeyValProfile)
	h.Kv.Profile(m)
}

func (h *HeadersProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Kv: %s", h.Kv.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *HeadersConfig) Normalize() {
	config.Kv.Normalize()
}

func (config *HeadersConfig) Decide(h *HeadersProfile) string {
	str := config.Kv.Decide(h.Kv)
	if str == "" {
		return str
	}
	return fmt.Sprintf("KeyVal: %s", str)
}

func (config *HeadersConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  KeyVal: %s", config.Kv.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

// Allow typical values - use for development but not in production
func (config *HeadersConfig) AddTypicalVal() {
	config.Kv.OtherKeynames = NewSimpleValConfig(0, 0, 0, 16, 16, 2, 4)
	config.Kv.OtherVals = NewSimpleValConfig(0, 0, 0, 32, 32, 8, 16)
	//config.Kv.OtherVals.FlagsL = 1<<MinusSlot | 1<<AsteriskSlot | 1<<SlashSlot | 1<<CommentsSlot | 1<<PeriodSlot
	config.Kv.OtherVals.Flags = 1<<MinusSlot | 1<<AsteriskSlot | 1<<SlashSlot | 1<<CommentsSlot | 1<<PeriodSlot
}

func (p *RespPile) Add(rp *RespProfile) {
	p.Headers.Add(rp.Headers)

}

func (p *RespPile) Clear() {
	p.Headers.Clear()
}

func (rp *RespProfile) Profile(resp *http.Response) {
	rp.Headers = new(HeadersProfile)
	rp.Headers.Profile(resp.Header)
}

func (rp *RespProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Headers: %s", rp.Headers.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (p *ReqPile) Add(rp *ReqProfile) {
	p.ClientIp = append(p.ClientIp, rp.ClientIp)
	p.Method = append(p.Method, rp.Method)
	p.Proto = append(p.Proto, rp.Proto)
	p.ContentLength = append(p.ContentLength, rp.ContentLength)
	p.Url.Add(rp.Url)
	p.Qs.Add(rp.Qs)
	p.Headers.Add(rp.Headers)
}

func (p *ReqPile) Clear() {
	p.ClientIp = make([]net.IP, 0, 1)
	p.Method = make([]string, 0, 1)
	p.Proto = make([]string, 0, 1)
	p.ContentLength = make([]uint8, 0, 1)
	p.Url.Clear()
	p.Qs.Clear()
	p.Headers.Clear()
}

func (rp *ReqProfile) Profile(req *http.Request, cip net.IP) {
	var forwarded []string
	var ok bool
	var hopipstr string
	var hopip net.IP

	// TBD - Add support for rfc7239 "forwarded" header
	forwarded, ok = req.Header["X-Forwarded-For"]
	if ok {
		numhops := len(forwarded)
		if numhops > 0 {
			hopipstr = forwarded[numhops-1]
		}
	}
	if len(hopipstr) > 0 {
		hopip = net.ParseIP(hopipstr)
		rp.HopIp = hopip
	}
	fmt.Printf("HOP-IP %v %s %s %s \n", hopip, hopipstr, req.Header["X-Forwarded-For"], req.Header["Forwarded"])
	rp.ClientIp = cip
	rp.Method = req.Method
	rp.Proto = req.Proto
	log2length := uint8(0)
	length := req.ContentLength
	for length > 0 {
		log2length++
		length >>= 1
	}
	rp.ContentLength = log2length
	rp.Url = new(UrlProfile)
	rp.Url.Profile(req.URL.Path)
	rp.Qs = new(QueryProfile)
	rp.Qs.Profile(req.URL.Query())
	rp.Headers = new(HeadersProfile)
	rp.Headers.Profile(req.Header)
}

func (rp *ReqProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Method: %v\n", rp.Method))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Proto: %v\n", rp.Proto))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ClientIp: %s\n", rp.ClientIp.String()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ContentLength: %d\n", int(math.Pow(2, float64(rp.ContentLength)))))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Url: %s", rp.Url.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Qs: %s", rp.Qs.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Headers: %s", rp.Headers.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *RespConfig) Normalize() {
	config.Headers.Normalize()
}

func (config *ReqConfig) Normalize() {
	config.Qs.Normalize()
	config.Headers.Normalize()
	config.Url.Normalize()
}

func (config *RespConfig) Decide(rp *RespProfile) string {
	ret := config.Headers.Decide(rp.Headers)
	if ret == "" {
		return ret
	}
	return fmt.Sprintf("HttpResponse: %s", ret)
}

func (config *ReqConfig) Decide(rp *ReqProfile) string {
	var ret string
	ret = config.Url.Decide(rp.Url)
	if ret != "" {
		return fmt.Sprintf("Url: %s", ret)
	}
	ret = config.Qs.Decide(rp.Qs)
	if ret != "" {
		return fmt.Sprintf("QueryString: %s", ret)
	}
	ret = config.Headers.Decide(rp.Headers)
	if ret != "" {
		return fmt.Sprintf("Headers: %s", ret)
	}
	if (rp.ClientIp != nil) && !rp.ClientIp.IsUnspecified() && !rp.ClientIp.IsLoopback() && !rp.ClientIp.IsPrivate() {
		ret = config.ClientIp.Decide(IpSetFromIp(rp.ClientIp))
		if ret != "" {
			return fmt.Sprintf("ClientIp: %s", ret)
		}
	}
	if (rp.HopIp != nil) && !rp.HopIp.IsUnspecified() && !rp.HopIp.IsLoopback() && !rp.HopIp.IsPrivate() {
		ret = config.HopIp.Decide(IpSetFromIp(rp.HopIp))
		if ret != "" {
			return fmt.Sprintf("HopIp: %s", ret)
		}
	}

	methodSet := make(Set)
	methodSet[rp.Method] = true
	ret = config.Method.Decide(methodSet)
	if ret != "" {
		return fmt.Sprintf("Method: %s", ret)
	}
	protoSet := make(Set)
	protoSet[rp.Proto] = true
	ret = config.Proto.Decide(protoSet)
	if ret != "" {
		return fmt.Sprintf("Proto: %s", ret)
	}
	ret = config.ContentLength.Decide(rp.ContentLength)
	if ret != "" {
		return fmt.Sprintf("ContentLength: %s", ret)
	}
	return ""
}

func (config *RespConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Headers: %s", config.Headers.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *ReqConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Method: %v\n", config.Method))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Proto: %v\n", config.Proto))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ClientIp: %v\n", config.ClientIp))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ContentLength: %s\n", config.ContentLength.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Url: %s", config.Url.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Qs: %s", config.Qs.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Headers: %s", config.Headers.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

// Allow typical values - use for development but not in production
func (config *RespConfig) AddTypicalVal() {
	config.Headers.AddTypicalVal()
}

// Allow typical values - use for development but not in production
func (config *ReqConfig) AddTypicalVal() {
	config.Headers.AddTypicalVal()
	config.Url.AddTypicalVal()
	config.Qs.AddTypicalVal()
}
