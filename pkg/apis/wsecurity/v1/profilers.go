package v1

import (
	"bytes"
	"fmt"
)

type U8Minmax struct {
	Min uint8 `json:"min"`
	Max uint8 `json:"max"`
}

type Set struct {
	List []string
	m    map[string]bool
}

type U8MinmaxSlice []U8Minmax
type Uint32Slice []uint32

func (p *Set) Add(v string) {
	if !p.m[v] {
		p.m[v] = true
		p.List = append(p.List, v)

	}
}

func (p *Set) Clear() {
	p.m = make(map[string]bool)
	p.List = make([]string, 0, 1)
}

func (p *Set) Append(a *Set) {
	if p.m == nil {
		p.m = make(map[string]bool, len(a.List))
	}
	if p.List == nil {
		p.List = make([]string, len(a.List))
	}
	for _, v := range a.List {
		if !p.m[v] {
			p.m[v] = true
			p.List = append(p.List, v)
		}
	}
}

func (p Set) Decide(s string) string {
	_, exists := p.m[s]
	if !exists {
		return fmt.Sprintf("Unexpected key %s in Set", s)
	}
	return ""
}

func AddToSetFromList(list []string, set *Set) {
	if set.m == nil {
		set.m = make(map[string]bool, len(list))
	}
	if set.List == nil {
		set.List = make([]string, len(list))
	}
	for _, v := range list {
		if !set.m[v] {
			set.m[v] = true
			set.List = append(set.List, v)
		}
	}
}

func AddToListFromSet(set *Set, list *[]string) {
	*list = append(*list, set.List...)
}
func (base Uint32Slice) Add(val Uint32Slice) Uint32Slice {
	if missing := len(val) - len(base); missing > 0 {
		// Dynamically allocate as many blockElements as needed for this config
		base = append(base, make([]uint32, missing)...)
	}
	for i, v := range val {
		base[i] = base[i] | v
	}
	return base
}

func (base Uint32Slice) Decide(val Uint32Slice) string {
	for i, v := range val {
		if v == 0 {
			continue
		}
		if i < len(base) && (v & ^base[i]) == 0 {
			continue
		}
		return fmt.Sprintf("Unexpected Unicode Flag %x on Element %d", v, i)
	}
	return ""
}

func (uint64Slice Uint32Slice) Describe() string {
	if uint64Slice != nil {
		var description bytes.Buffer
		description.WriteString("Unicode slice: ")
		description.WriteString(fmt.Sprintf("%x", uint64Slice[0]))
		for i := 1; i < len(uint64Slice); i++ {
			description.WriteString(fmt.Sprintf(", %x", uint64Slice[i]))
		}
		return description.String()
	}
	return ""
}

func (uint64Slice Uint32Slice) Marshal() string {
	if uint64Slice == nil {
		return "null"
	}
	if len(uint64Slice) == 0 {
		return "[]"
	}
	var description bytes.Buffer
	description.WriteString(fmt.Sprintf("[%x", uint64Slice[0]))
	for i := 1; i < len(uint64Slice); i++ {
		description.WriteString(fmt.Sprintf(",%x", uint64Slice[i]))
	}
	description.WriteString("]")
	return description.String()

}

func (mms U8MinmaxSlice) Decide(v uint8) string {
	if v == 0 {
		return ""
	}
	// v>0
	if len(mms) == 0 {
		return fmt.Sprintf("Value %d Not Allowed!", v)
	}

	for j := 0; j < len(mms); j++ {
		if v < mms[j].Min {
			break
		}
		if v <= mms[j].Max { // found ok interval
			return ""
		}
	}
	return fmt.Sprintf("Counter out of Range: %d", v)
}

func (mms U8MinmaxSlice) AddValExample(v uint8) U8MinmaxSlice {
	if len(mms) == 0 {
		mms = append(mms, U8Minmax{v, v})
	} else {
		if mms[0].Min > v {
			mms[0].Min = v
		}
		if mms[0].Max < v {
			mms[0].Max = v
		}
	}
	return mms
}

func (mms *U8MinmaxSlice) Learn(list []uint8) {
	min := uint8(0)
	max := uint8(0)
	if len(list) >= 0 {
		min = list[0]
		max = list[0]
	}
	for _, v := range list {
		if min > v {
			min = v
		}
		if max < v {
			max = v
		}
	}
	*mms = append(*mms, U8Minmax{min, max})
}

func (mms *U8Minmax) Merge(m *U8Minmax) {
	if mms.Min > m.Min {
		mms.Min = m.Min
	}
	if mms.Max < m.Max {
		mms.Max = m.Max
	}
}

func (mms *U8MinmaxSlice) Merge(m U8MinmaxSlice) {
	//(*mms)[0].Merge(&(m[0]))

	var found bool
	for _, v := range m {
		for _, mm := range *mms {
			if mm.Min < v.Min {
				if mm.Max > v.Max {
					found = true
					break
				}
				// mm.Max < v.Max
				if mm.Max >= v.Min {
					mm.Merge(&v)
					found = true
					break
				}
				// mm.Max < v.Min
			}
			// mm.Min > v.Min
			if mm.Min > v.Max {
				continue
			}
			mm.Merge(&v)
			found = true
			break
		}
		if !found {
			*mms = append(*mms, v)
		}
	}
}

func (mms U8MinmaxSlice) Describe() string {
	if len(mms) == 0 {
		return "No Limit"
	}
	var description bytes.Buffer
	description.WriteString(fmt.Sprintf("%d-%d", mms[0].Min, mms[0].Max))
	for j := 1; j < len(mms); j++ {
		description.WriteString(fmt.Sprintf(", %d-%d", mms[j].Min, mms[j].Max))
	}
	return description.String()
}

func (mms U8MinmaxSlice) Marshal() string {
	if len(mms) == 0 {
		return "null"
	}
	var description bytes.Buffer
	description.WriteString(fmt.Sprintf("[{Min:%d,Max: %d", mms[0].Min, mms[0].Max))
	for j := 1; j < len(mms); j++ {
		description.WriteString(fmt.Sprintf("}, {Min:%d,Max: %d", mms[j].Min, mms[j].Max))
	}
	description.WriteString("}]")
	return description.String()
}
