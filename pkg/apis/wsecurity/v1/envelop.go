package v1

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type EnvelopPile struct {
	ResponseTime   []uint8 `json:"responsetime"`
	CompletionTime []uint8 `json:"completiontime"`
}

type EnvelopConfig struct {
	ResponseTime   U8MinmaxSlice `json:"responsetime"`
	CompletionTime U8MinmaxSlice `json:"completiontime"`
}

type EnvelopProfile struct {
	ResponseTime   uint8 `json:"responsetime"`
	CompletionTime uint8 `json:"completiontime"`
}

func (p *EnvelopPile) Add(ep *EnvelopProfile) {
	p.CompletionTime = append(p.CompletionTime, ep.CompletionTime)
	p.ResponseTime = append(p.ResponseTime, ep.ResponseTime)
}

func (p *EnvelopPile) Append(a *EnvelopPile) {
	p.CompletionTime = append(p.CompletionTime, a.CompletionTime...)
	p.ResponseTime = append(p.ResponseTime, a.ResponseTime...)
}

func (config *EnvelopConfig) Reconcile() {

}

func (config *EnvelopConfig) Learn(p *EnvelopPile) {
	config.CompletionTime.Learn(p.CompletionTime)
	config.ResponseTime.Learn(p.ResponseTime)
}

func (config *EnvelopConfig) Merge(m *EnvelopConfig) {
	config.CompletionTime.Merge(m.CompletionTime)
	config.ResponseTime.Merge(m.ResponseTime)
}

func (config *EnvelopConfig) Normalize() {
	config.ResponseTime = append(config.ResponseTime, U8Minmax{0, 0})
	config.CompletionTime = append(config.CompletionTime, U8Minmax{0, 0})
	//config.Tcp4Peers = new(CidrSet)
	//config.Udp4Peers = new(CidrSet)
	//config.Udplite4Peers = new(CidrSet)
	//config.Tcp6Peers = new(CidrSet)
	//config.Udp6Peers = new(CidrSet)
	//config.Udplite6Peers = new(CidrSet)

}

func (config *EnvelopConfig) Decide(pp *EnvelopProfile) string {
	var ret string
	ret = config.ResponseTime.Decide(pp.ResponseTime)
	if ret != "" {
		return fmt.Sprintf("ResponseTime: %s", ret)
	}
	ret = config.CompletionTime.Decide(pp.CompletionTime)
	if ret != "" {
		return fmt.Sprintf("CompletionTime: %s", ret)
	}
	return ""
}

func (config *EnvelopConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ResponseTime: %s,\n", config.ResponseTime.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  CompletionTime: %s,\n", config.CompletionTime.Marshal()))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

// Allow typical values - use for development but not in production
func (config *EnvelopConfig) AddTypicalVal() {
	config.ResponseTime = make([]U8Minmax, 1)
	config.ResponseTime[0].Max = 60
	config.CompletionTime = make([]U8Minmax, 1)
	config.CompletionTime[0].Max = 120
}

// Profile timestamps and /proc
func (pp *EnvelopProfile) Profile(reqTime time.Time, respTime time.Time, endTime time.Time) {
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

func (pp *EnvelopProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ResponseTime: %d\n", pp.ResponseTime))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  CompletionTime: %d\n", pp.CompletionTime))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (pp *EnvelopPile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ResponseTime: %v\n", pp.ResponseTime))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  CompletionTime: %v\n", pp.CompletionTime))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}
