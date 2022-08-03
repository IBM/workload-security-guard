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

type BodyPile struct {
	Unstructured *SimpleValPile  `json:"unstructured"`
	Structured   *StructuredPile `json:"structured"`
}

func (config *BodyConfig) Reconcile() {

}

func (config *BodyConfig) Learn(p *BodyPile) {
	config.Structured.Learn(p.Structured)
	config.Unstructured.Learn(p.Unstructured)
}

func (config *BodyConfig) Merge(m *BodyConfig) {
	config.Structured.Merge(m.Structured)
	config.Unstructured.Merge(m.Unstructured)
}

func (p *BodyPile) Clear() {
	p.Structured = new(StructuredPile)
	p.Unstructured = new(SimpleValPile)
}

func (p *BodyPile) Append(a *BodyPile) {
	p.Structured.Append(a.Structured)
	p.Unstructured.Append(a.Unstructured)
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

func (config *BodyConfig) Marshal(depth int) string {
	return "<ReqBody Config>"
}

func (config *BodyPile) Marshal(depth int) string {
	return "<ReqBody Config>"
}

func (p *BodyPile) Add(h *BodyProfile) {
	//p.Kv.Add(h.Kv)
}
