package v1alpha1

/*
import (
	"bytes"
	"fmt"
)

type u8MinMax struct {
	Min uint8 `json:"min"`
	Max uint8 `json:"max"`
}

func (mms *u8MinMax) merge(m *u8MinMax) {
	if mms.Min > m.Min {
		mms.Min = m.Min
	}
	if mms.Max < m.Max {
		mms.Max = m.Max
	}
}

// Exposes ValueConfig interface
type U8MinMaxConfig []u8MinMax

func (mms U8MinMaxConfig) Decide(v uint8) string {
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

func (mms *U8MinMaxConfig) Learn(list []uint8) {
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
	*mms = append(*mms, u8MinMax{min, max})
}

func (mms *U8MinMaxConfig) Fuse(m U8MinMaxConfig) {
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
					mm.merge(&v)
					found = true
					break
				}
				// mm.Max < v.Min
			}
			// mm.Min > v.Min
			if mm.Min > v.Max {
				continue
			}
			mm.merge(&v)
			found = true
			break
		}
		if !found {
			*mms = append(*mms, v)
		}
	}
}

func (mms U8MinMaxConfig) String() string {
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

func (mms U8MinMaxConfig) AddValExample(v uint8) U8MinMaxConfig {
	if len(mms) == 0 {
		mms = append(mms, u8MinMax{v, v})
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
*/
