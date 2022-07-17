package v1

import (
	"bytes"
	"fmt"
	"math"
	"mime"
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
	ClientIp      IpPile  `json:"cip"`           // 192.168.32.1
	HopIp         IpPile  `json:"hopip"`         // 1.2.3.4
	Method        Set     `json:"method"`        // GET
	Proto         Set     `json:"proto"`         // "HTTP/1.1"
	MediaType     Set     `json:"mediatype"`     // "text/html"
	ContentLength []uint8 `json:"contentlength"` // 0

	Url     UrlPile     `json:"url"`
	Qs      QueryPile   `json:"qs"`
	Headers HeadersPile `json:"headers"`
}

type ReqConfig struct {
	clientIp      CidrSet
	hopIp         CidrSet
	method        Set
	proto         Set
	mediaType     Set
	ClientIp      []string      `json:"cip"`           // subnets for external IPs (normally empty)
	HopIp         []string      `json:"hopip"`         // subnets for external IPs
	Method        []string      `json:"method"`        // GET
	Proto         []string      `json:"proto"`         // "HTTP/1.1"
	MediaType     []string      `json:"mediatype"`     // "text/html"
	ContentLength U8MinmaxSlice `json:"contentlength"` // 0
	Url           UrlConfig     `json:"url"`
	Qs            QueryConfig   `json:"qs"`
	Headers       HeadersConfig `json:"headers"`
}

type ReqProfile struct {
	ClientIp      *IpSet          `json:"cip"`           // 192.168.32.1
	HopIp         *IpSet          `json:"hopip"`         // 1.2.3.4
	Method        string          `json:"method"`        // GET
	Proto         string          `json:"proto"`         // "HTTP/1.1"
	MediaType     string          `json:"mediatype"`     // "text/html"
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
func (p *UrlPile) Append(a *UrlPile) {
	p.Segments = append(p.Segments, a.Segments...)
	p.Val.Append(a.Val)
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
	//fmt.Printf("Path %s, segments %v, len %d\n", path, segments, numSegments)
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

func (u *UrlPile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Val: %s", u.Val.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Segments: %v", u.Segments))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *UrlConfig) Normalize() {
	config.Val.Normalize()
	config.Segments = append(config.Segments, U8Minmax{0, 0})
}

func (config *UrlConfig) Learn(p *UrlPile) {
	config.Segments.Learn(p.Segments)
	config.Val.Learn(p.Val)
}

func (config *UrlConfig) Merge(m *UrlConfig) {
	config.Segments.Merge(m.Segments)
	config.Val.Merge(&m.Val)
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
	config.Val.Spaces = make([]U8Minmax, 1)
	config.Val.Unicodes = make([]U8Minmax, 1)
	config.Val.NonReadables = make([]U8Minmax, 1)
	config.Val.Letters = make([]U8Minmax, 1)
	config.Val.Digits = make([]U8Minmax, 1)
	config.Val.Sequences = make([]U8Minmax, 1)
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
	config.Segments = make([]U8Minmax, 1)
	config.Segments[0].Max = 8
}

func (p *QueryPile) Add(q *QueryProfile) {
	p.Kv.Add(q.Kv)
}

func (p *QueryPile) Clear() {
	p.Kv = new(KeyValPile)
	p.Kv.Clear()
}

func (p *QueryPile) Append(a *QueryPile) {
	p.Kv.Append(a.Kv)
}

func (q *QueryProfile) Profile(m map[string][]string) {
	q.Kv = new(KeyValProfile)
	q.Kv.Profile(m, nil)
}

func (config *QueryConfig) Normalize() {
	config.Kv.Normalize()
}

func (config *QueryConfig) Learn(p *QueryPile) {
	config.Kv.Learn(p.Kv)
}

func (config *QueryConfig) Merge(m *QueryConfig) {
	config.Kv.Merge(&m.Kv)
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

func (q *QueryPile) Marshal(depth int) string {
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
	p.Kv = new(KeyValPile)
	p.Kv.Clear()
}

func (p *HeadersPile) Append(a *HeadersPile) {
	p.Kv.Append(a.Kv)
}

var excpetionHeaders = map[string]bool{"Content-Type": true}

func (h *HeadersProfile) Profile(m map[string][]string) {
	h.Kv = new(KeyValProfile)
	h.Kv.Profile(m, excpetionHeaders)
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

func (h *HeadersPile) Marshal(depth int) string {
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

func (config *HeadersConfig) Learn(p *HeadersPile) {
	config.Kv.Learn(p.Kv)
}

func (config *HeadersConfig) Merge(m *HeadersConfig) {
	config.Kv.Merge(&m.Kv)
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

func (p *RespPile) Append(a *RespPile) {
	p.Headers.Append(&a.Headers)
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

func (rp *RespPile) Marshal(depth int) string {
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
	fmt.Println("Add ReqProfile")
	p.ClientIp.Add(rp.ClientIp)
	p.HopIp.Add(rp.HopIp)
	p.Method.Add(rp.Method)
	p.Proto.Add(rp.Proto)
	if rp.MediaType != "" {
		p.MediaType.Add(rp.MediaType)
	}
	p.ContentLength = append(p.ContentLength, rp.ContentLength)
	p.Url.Add(rp.Url)
	p.Qs.Add(rp.Qs)
	p.Headers.Add(rp.Headers)
}

func (p *ReqPile) Clear() {
	p.ClientIp.Clear()
	p.Method.Clear()
	p.Proto.Clear()
	p.MediaType.Clear()
	p.ContentLength = make([]uint8, 0, 1)
	p.Url.Clear()
	p.Qs.Clear()
	p.Headers.Clear()
}

func (p *ReqPile) Append(a *ReqPile) {
	fmt.Println("Append ReqPile")

	p.ClientIp.Append(&a.ClientIp)
	p.Method.Append(&a.Method)
	p.Proto.Append(&a.Proto)
	p.MediaType.Append(&a.MediaType)
	p.ContentLength = append(p.ContentLength, a.ContentLength...)
	p.Url.Append(&a.Url)
	p.Qs.Append(&a.Qs)
	p.Headers.Append(&a.Headers)
}

func (rp *ReqProfile) Profile(req *http.Request, cip net.IP) {
	var forwarded []string
	var ok bool
	var hopipstr string
	var hopip net.IP
	var mediatype string
	var params map[string]string
	var err error

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
		rp.HopIp = IpSetFromIp(hopip)
	}
	//fmt.Printf("HOP-IP %v %s %s %s \n", hopip, hopipstr, req.Header["X-Forwarded-For"], req.Header["Forwarded"])
	rp.ClientIp = IpSetFromIp(cip)
	rp.HopIp = IpSetFromIp(hopip)

	rp.Method = req.Method
	rp.Proto = req.Proto

	log2length := uint8(0)
	length := req.ContentLength
	mediatype = ""
	if length > 0 {
		for length > 0 {
			log2length++
			length >>= 1
		}
		mediatype, params, err = mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			fmt.Printf("err ParseMediaType %s Content-Type %s (%v) \n", err.Error(), req.Header.Get("Content-Type"), req.Header["Content-Type"])
			mediatype = "X-CAN-NOT-PARSE-MEDIA-TYPE"
		}
		_ = params // TBD  - what should we do with params?
		rp.MediaType = mediatype
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
	description.WriteString(fmt.Sprintf("  MediaType: %v\n", rp.MediaType))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ClientIp: %v\n", rp.ClientIp))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  HopIp: %v\n", rp.HopIp))
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

func (rp *ReqPile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Method: %v\n", rp.Method))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Proto: %v\n", rp.Proto))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  MediaType: %v\n", rp.MediaType))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ClientIp: %v\n", rp.ClientIp))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  HopIp: %v\n", rp.HopIp))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ContentLength: %v\n", rp.ContentLength))
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
func (config *RespConfig) Reconcile() {
}

func (config *ReqConfig) Reconcile() {
	fmt.Printf("Reconcile ReqConfig (%d) %v \n", len(config.ClientIp), config.ClientIp)
	config.clientIp = GetCidrsFromList(config.ClientIp)
	config.hopIp = GetCidrsFromList(config.HopIp)
	AddToSetFromList(config.Method, &config.method)
	AddToSetFromList(config.Proto, &config.proto)
	AddToSetFromList(config.MediaType, &config.mediaType)
}

func (config *ReqConfig) Learn(p *ReqPile) {
	fmt.Println("Learn ReqConfig")
	config.clientIp = GetCidrsFromIpList(p.ClientIp.List)
	config.hopIp = GetCidrsFromIpList(p.HopIp.List)
	config.method.Append(&p.Method)
	config.proto.Append(&p.Proto)
	config.mediaType.Append(&p.MediaType)
	config.ContentLength.Learn(p.ContentLength)
	config.Headers.Learn(&p.Headers)
	config.Qs.Learn(&p.Qs)
	config.Url.Learn(&p.Url)
	config.Method = config.method.List
	config.Proto = config.proto.List
	config.MediaType = config.mediaType.List
	config.ClientIp = config.clientIp.Strings()
	config.HopIp = config.hopIp.Strings()
}

func (config *ReqConfig) Merge(mc *ReqConfig) {
	fmt.Println("Merge ReqConfig")
	config.clientIp.Merge(mc.clientIp)
	config.hopIp.Merge(mc.hopIp)
	config.method.Append(&mc.method)
	config.proto.Append(&mc.proto)
	config.ContentLength.Merge(mc.ContentLength)
	config.Headers.Merge(&mc.Headers)
	config.Qs.Merge(&mc.Qs)
	config.Url.Merge(&mc.Url)
	config.Method = config.method.List
	config.Proto = config.proto.List
	config.MediaType = config.mediaType.List
	config.ClientIp = config.clientIp.Strings()
	config.HopIp = config.hopIp.Strings()
}

func (config *RespConfig) Learn(p *RespPile) {
	config.Headers.Learn(&p.Headers)
}

func (config *RespConfig) Merge(mc *RespConfig) {
	config.Headers.Merge(&mc.Headers)
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
	ret = config.clientIp.Decide(rp.ClientIp)
	if ret != "" {
		return fmt.Sprintf("ClientIp: %s", ret)
	}

	ret = config.hopIp.Decide(rp.HopIp)
	if ret != "" {
		return fmt.Sprintf("HopIp: %s", ret)
	}

	//methodSet := make(Set)
	//methodSet[rp.Method] = true
	ret = config.method.Decide(rp.Method)
	if ret != "" {
		return fmt.Sprintf("Method: %s", ret)
	}
	//protoSet := make(Set)
	//protoSet[rp.Proto] = true
	ret = config.proto.Decide(rp.Proto)
	if ret != "" {
		return fmt.Sprintf("Proto: %s", ret)
	}
	if rp.MediaType != "" {
		ret = config.mediaType.Decide(rp.MediaType)
		if ret != "" {
			return fmt.Sprintf("MediaType: %s", ret)
		}
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
	description.WriteString(fmt.Sprintf("  MediaType: %v\n", config.MediaType))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ClientIp: %v\n", config.ClientIp))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  HopIp: %v\n", config.HopIp))
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
