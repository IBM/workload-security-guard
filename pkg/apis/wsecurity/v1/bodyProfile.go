package v1

import "fmt"

type BodyProfile struct {
	Unstructured *SimpleValProfile  `json:"unstructured"`
	Structured   *StructuredProfile `json:"structured"`
}

type BodyConfig struct {
	Unstructured *SimpleValConfig  `json:"unstructured"`
	Structured   *StructuredConfig `json:"structured"`
}

func (config *BodyConfig) Decide(bp *BodyProfile) string {
	if bp.Structured != nil {
		if config.Structured != nil {
			str := config.Structured.Decide(bp.Structured)
			if str != "" {
				return fmt.Sprintf("Body %s", str)
			}
		} else {
			return "Structured Body not allowed"
		}
	}
	if bp.Unstructured != nil {
		if config.Unstructured != nil {
			str := config.Unstructured.Decide(bp.Unstructured)
			if str != "" {
				return fmt.Sprintf("Body %s", str)
			}
		} else {
			return "Unstructured Body not allowed"
		}
	}
	return ""
}
