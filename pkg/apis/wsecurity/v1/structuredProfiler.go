package v1

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

/*
type ConfigInterface interface {
	Learn(p *PileInterface)
	Merge(m *ConfigInterface)
	Normalize()
	Decide(kvp *ProfileInterface) string
	Marshal(depth int) string
	Describe() string

	Clear()
}

type PileInterface interface {
	Append(*PileInterface)
	Marshal(depth int) string
	Pile(*ProfileInterface)
	Add(*ProfileInterface)
	Clear()
}

type ProfileInterface interface {
	Marshal(depth int) string
	DeepCopyInto(*ProfileInterface)
}
*/
type StructuredConfig struct {
	Kind string                       `json:"kind"` // boolean, number, string, skip, array, object
	Val  *SimpleValConfig             `json:"val"`  // used for: array, boolean, number, string items
	Kv   map[string]*StructuredConfig `json:"kv"`   // used for: object items
}

type StructuredPile struct {
	Kind string                     `json:"kind"` // bool, float64, string, array, map
	Val  *SimpleValPile             `json:"val"`  // used for: array, boolean, number, string items
	Kv   map[string]*StructuredPile `json:"kv"`   // used for: object items
}

func (config *StructuredConfig) Learn(p *StructuredPile) {
	//??? config.Kind.Learn(p.)
	config.Val.Learn(p.Val)
	//config.Kv.Lean(p.Kv)
}

func (config *StructuredConfig) Merge(m *StructuredConfig) {
	//??? config.Kind.Learn(p.)
	config.Val.Merge(m.Val)
	//config.Kv.Merge(m.Kv)
}

//  JsonProfile struct - maintain the profile of a json with some structure
//	Data Types: The default Golang data types for decoding and encoding JSON are as follows:
//		bool for JSON booleans.
//		float64 for JSON numbers.
//		string for JSON strings.
//		nil for JSON null.
//		array as JSON array.
//		map or struct as JSON Object.
type StructuredProfile struct {
	Kind string                        `json:"kind"` // bool, float64, string, array, map
	Vals []SimpleValProfile            `json:"vals"` // used for: array, boolean, number, string items
	Kv   map[string]*StructuredProfile `json:"kv"`   // used for: object items

}

/*
func process(o interface{}) {
	if reflect.ValueOf(o).Kind() == reflect.Map {
		fmt.Printf("%v is a map!\n", o)
		for k, v := range o.(map[string]interface{}) {

			fmt.Printf("%s: ", k)
			process(v)
		}
		return
	}
	if reflect.ValueOf(o).Kind() == reflect.Slice {
		fmt.Printf("%v is a slice!\n", o)
		for i, v := range o.([]interface{}) {
			fmt.Printf("%d: ", i)
			process(v)
		}
		return
	}
	fmt.Printf("%v is just a value\n", o)
}
*/

func (config *StructuredConfig) Decide(jp *StructuredProfile) string {
	if config.Kind != jp.Kind {
		if config.Kind == "skip" {
			return ""
		} else {
			return fmt.Sprintf("Structured -  kind mismatch allowed %s has %s", config.Kind, jp.Kind)
		}
	}
	switch config.Kind {
	case "object":
		//fmt.Printf(">> Object profile %v\n", jp.Kv)
		//fmt.Printf(">> Object config %v\n", config.Kv)
		for jpk, jpv := range jp.Kv {
			if config.Kv[jpk] == nil {
				return fmt.Sprintf("Structured - key not allowed: %s", jpk)
			} else {
				str := config.Kv[jpk].Decide(jpv)
				if str != "" {
					return fmt.Sprintf("Structured - Key: %s - %s", jpk, str)
				}
			}
		}
	case "list":
		for _, jpv := range jp.Vals {
			str := config.Val.Decide(&jpv)
			if str != "" {
				return fmt.Sprintf("Structured - array val: %s", str)
			}
		}
	default: // number, string, boolean
		str := config.Val.Decide(&jp.Vals[0])
		if str != "" {
			return fmt.Sprintf("Structured - Val: %s", str)
		}
	}
	return ""
}

// Profile a generic json
func (jp *StructuredProfile) Profile(o interface{}) {
	if o == nil {
		return
	}
	switch reflect.ValueOf(o).Kind() {
	case reflect.Slice:
		jp.Kind = "array"
		//fmt.Printf("JsonProfile Profile is a slice! \n")
		s := o.([]interface{})
		jp.Vals = make([]SimpleValProfile, len(s))

		for i, v := range s {
			switch reflect.ValueOf(v).Kind() {
			case reflect.Map:
				jp.Vals[i].Profile(fmt.Sprintf("%v", v))
			case reflect.Slice:
				jp.Vals[i].Profile(fmt.Sprintf("%v", v))
			case reflect.Float64:
				jp.Vals[i].Profile(fmt.Sprintf("%v", v.(float64)))
			case reflect.Bool:
				jp.Vals[i].Profile(fmt.Sprintf("%t", v.(bool)))
			case reflect.String:
				jp.Vals[i].Profile(v.(string))
			default:
				panic("StructuredProfile.Profile() unknown Kind in Array")
			}
		}
	case reflect.Map:
		jp.Kind = "object"
		//("JsonProfile Profile is a map! \n")
		m := o.(map[string]interface{})
		jp.Kv = make(map[string]*StructuredProfile, len(m))
		for k, v := range m {
			jp.Kv[k] = new(StructuredProfile)
			jp.Kv[k].Profile(v)
		}
	case reflect.Float64:
		jp.Kind = "number"
		jp.Vals = make([]SimpleValProfile, 1)
		jp.Vals[0].Profile(fmt.Sprintf("%f", o.(float64)))
	case reflect.Bool:
		jp.Kind = "boolean"
		jp.Vals = make([]SimpleValProfile, 1)
		jp.Vals[0].Profile(fmt.Sprintf("%t", o.(bool)))
	case reflect.String:
		jp.Kind = "string"
		jp.Vals = make([]SimpleValProfile, 1)
		jp.Vals[0].Profile(o.(string))
	default:
		panic("StructuredProfile.Profile() unknown Kind")
	}
	//fmt.Printf("StructuredProfile Profile jp.Kind %s\n", jp.Kind)
}

// Profile a PostForm
func (jp *StructuredProfile) ProfilePostForm(m map[string][]string) {
	jp.Kind = "object"
	//fmt.Printf("JsonProfile ProfilePostForm\n")
	jp.Kv = make(map[string]*StructuredProfile, len(m))
	for k, sv := range m {
		jp.Kv[k] = new(StructuredProfile)
		jp.Kv[k].Kind = "array"
		jp.Kv[k].Vals = make([]SimpleValProfile, len(sv))
		for i, v := range sv {
			jp.Kv[k].Vals[i].Profile(v)
		}
	}
}

func (jp *StructuredProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	//fmt.Printf("Marshal: starts %d with kind %s\n", depth, jp.Kind)
	switch jp.Kind {
	case "object":
		description.WriteString(shift)
		description.WriteString("{\n")
		for k, v := range jp.Kv {
			description.WriteString(fmt.Sprintf("  , %s: %s", k, v.Marshal(depth+1)))
		}
		description.WriteString(shift)
		description.WriteString("}\n")
	case "array":
		description.WriteString(shift)
		description.WriteString("[\n")
		for _, v := range jp.Vals {
			description.WriteString(fmt.Sprintf("%s ", v.Marshal(depth+1)))
		}
		description.WriteString(shift)
		description.WriteString("]\n")
	case "boolean":
		description.WriteString(jp.Vals[0].Marshal(depth + 1))
	case "number":
		description.WriteString(jp.Vals[0].Marshal(depth + 1))
	case "string":
		description.WriteString(jp.Vals[0].Marshal(depth + 1))
	}

	return description.String()
}

func (p *StructuredPile) Append(a *StructuredPile) {

}
