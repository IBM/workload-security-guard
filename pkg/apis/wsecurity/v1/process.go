package v1

import (
	"bytes"
	"fmt"
	"strings"
)

type ProcessPile struct {
	ResponseTime   []uint8 `json:"responsetime"`
	CompletionTime []uint8 `json:"completiontime"`
}

type ProcessConfig struct {
	ResponseTime   U8MinmaxSlice `json:"responsetime"`
	CompletionTime U8MinmaxSlice `json:"completiontime"`
}

type ProcessProfile struct {
	ResponseTime   uint8 `json:"responsetime"`
	CompletionTime uint8 `json:"completiontime"`
}

func (config *ProcessConfig) Normalize() {
	config.ResponseTime = append(config.ResponseTime, U8Minmax{0, 0})
	config.CompletionTime = append(config.CompletionTime, U8Minmax{0, 0})
}

func (config *ProcessConfig) Decide(pp *ProcessProfile) string {
	var ret string
	ret = config.ResponseTime.Decide(pp.ResponseTime)
	if ret == "" {
		ret = config.CompletionTime.Decide(pp.CompletionTime)
		if ret == "" {
			return ret
		}
	}
	return fmt.Sprintf("Processing: %s", ret)
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
