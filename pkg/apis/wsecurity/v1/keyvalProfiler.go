package v1

import (
	"bytes"
	"fmt"
	"strings"
)

type KeyValConfig struct {
	Vals map[string]*SimpleValConfig `json:"vals"` // Profile the value of whitelisted keys
	//MinimalSet    map[string]void             `json:"minimalSet,omitempty"`    // Mandatory keys
	//MinimalSet    map[string]struct{} `json:"minimalSet,omitempty"`    // Mandatory keys
	OtherVals     *SimpleValConfig `json:"otherVals"`     // Profile the values of other keys
	OtherKeynames *SimpleValConfig `json:"otherKeynames"` // Profile the keynames of other keys
}

type KeyValPile struct {
	Vals map[string]*SimpleValPile
}

type KeyValProfile struct {
	Vals map[string]*SimpleValProfile
}

//type void struct{}

func (p *KeyValPile) Add(kv *KeyValProfile) {
	for key, kv_profile := range kv.Vals {
		key_pile, exists := p.Vals[key]
		if !exists {
			key_pile = new(SimpleValPile)
			p.Vals[key] = key_pile
		}
		key_pile.Add(kv_profile)
	}
}

func (p *KeyValPile) Clear() {
	p.Vals = make(map[string]*SimpleValPile)
}

func (p *KeyValPile) Append(a *KeyValPile) {
	for key, kv_profile := range a.Vals {
		key_pile, exists := p.Vals[key]
		if !exists {
			key_pile = new(SimpleValPile)
			p.Vals[key] = key_pile
		}
		key_pile.Append(kv_profile)
	}
}

// Profile a generic map of key vals where we expect:
// keys belonging to some contstant list of keys
// vals have some defined charactaristics
func (kvp *KeyValProfile) Profile(m map[string][]string, exceptions map[string]bool) {
	if len(m) == 0 { // no keys
		return
	}
	kvp.Vals = make(map[string]*SimpleValProfile, len(m))
	for k, v := range m {
		if exceptions[k] {
			continue
		}
		//var keyConfig *SimpleValConfig
		//if config.Vals != nil {
		//	keyConfig = config.Vals[k]
		//}
		//if keyConfig == nil {
		//}

		val := strings.Join(v, " ")
		kvp.Vals[k] = new(SimpleValProfile)
		kvp.Vals[k].Profile(val)
	}
}

func (kvp *KeyValProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")

	if len(kvp.Vals) > 0 {
		description.WriteString(shift)
		description.WriteString("  Vals: {\n")
		for k, v := range kvp.Vals {
			description.WriteString(shift)
			description.WriteString(fmt.Sprintf("  , %s: %s", k, v.Marshal(depth+2)))
		}
		description.WriteString(shift)
		description.WriteString("  }\n")
	}

	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (kvp *KeyValPile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")

	if len(kvp.Vals) > 0 {
		description.WriteString(shift)
		description.WriteString("  Vals: {\n")
		for k, v := range kvp.Vals {
			description.WriteString(shift)
			description.WriteString(fmt.Sprintf("  , %s: %s", k, v.Marshal(depth+2)))
		}
		description.WriteString(shift)
		description.WriteString("  }\n")
	}

	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (config *KeyValConfig) Learn(p *KeyValPile) {
	config.Vals = make(map[string]*SimpleValConfig)
	for k, v := range p.Vals {
		svc := new(SimpleValConfig)
		svc.Learn(v)
		config.Vals[k] = svc
	}
}

func (config *KeyValConfig) Merge(m *KeyValConfig) {
	for mk, mv := range m.Vals {
		v, exists := config.Vals[mk]
		if exists {
			v.Merge(mv)
		} else {
			config.Vals[mk] = mv
		}
	}
}

func (config *KeyValConfig) Normalize() {
	//config.MinimalSet = make(map[string]void)
	//config.MinimalSet = make(map[string]struct{})
	if config.OtherVals == nil {
		config.OtherVals = new(SimpleValConfig)
	}
	if config.OtherKeynames == nil {
		config.OtherKeynames = new(SimpleValConfig)
	}
	config.OtherVals.Normalize()
	config.OtherKeynames.Normalize()
	for k, v := range config.Vals {
		if v == nil {
			v = new(SimpleValConfig)
			config.Vals[k] = v
		}
		v.Normalize()
	}
}

func (config *KeyValConfig) Decide(kvp *KeyValProfile) string {
	//if config == nil || !config.Enable {
	//	return ""
	//}

	// Duplicate minimalSet map
	//var required void
	var required struct{}
	//minimalSet := make(map[string]void, len(config.MinimalSet))
	minimalSet := make(map[string]struct{}, len(config.Vals))

	for k, v := range config.Vals {
		if v.Mandatory {
			minimalSet[k] = required
		}
	}

	// For each key-val, decide! and remove from minimalSet
	if kvp.Vals != nil {
		for k, v := range kvp.Vals {
			delete(minimalSet, k) // Remove from minimalSet
			// Decide based on a known key
			if config.Vals != nil && config.Vals[k] != nil {
				if ret := config.Vals[k].Decide(v); ret != "" {
					return fmt.Sprintf("Known Key %s: %s", k, ret)
				}
				continue
			}
			// Not a known key...
			if config.OtherKeynames == nil || config.OtherVals == nil {
				return fmt.Sprintf("Key %s is not known", k)
			}
			// Decide keyname of not known key
			var keynames SimpleValProfile
			keynames.Profile(k)
			if ret := config.OtherKeynames.Decide(&keynames); ret != "" {
				return fmt.Sprintf("Other keyname %s: %s", k, ret)
			}
			// Decide val of not known key
			if ret := config.OtherVals.Decide(v); ret != "" {
				return fmt.Sprintf("Other keyname %s: %s", k, ret)
			}
			continue
		}
	}
	// Once we oked all keys, check if there are missing mandatory keys
	if len(minimalSet) > 0 {
		keys := make([]string, len(minimalSet))
		for k := range minimalSet {
			keys = append(keys, k)
		}
		return fmt.Sprintf("KeyVal missing mandatory keys %s", strings.Join(keys, ", "))
	}
	return ""
}

// Allow a list of specific keys and an example of their values
// Can be called multiple times to add keys or to add examples for values
// Use this when the keynames are known in advance
// Call multiple times to show the entire range of values per key
// For keys not known in advance, use WhitelistByExample() instead
func (config *KeyValConfig) WhitelistKnownKeys(m map[string]string) {
	if config.Vals == nil {
		config.Vals = make(map[string]*SimpleValConfig, len(m))
	}
	for k, v := range m {
		if config.Vals[k] == nil {
			config.Vals[k] = new(SimpleValConfig)
		}
		config.Vals[k].AddValExample(v)
	}
}

// Define which of the known keynames is mandatory (if any)
// Must call WhitelistKnownKeys before setting keys as Mandatory
func (config *KeyValConfig) SetMandatoryKeys(minimalSet []string) {
	if config.Vals == nil {
		panic("Keys should be set with WhitelistKnownKeys before becoming Mandatory")
	}

	//if config.MinimalSet == nil {
	//config.MinimalSet = make(map[string]void, len(minimalSet))
	//	config.MinimalSet = make(map[string]struct{}, len(minimalSet))
	//}

	//var required void
	//var required struct{}

	for _, k := range minimalSet {
		if _, exists := config.Vals[k]; !exists {
			panic(fmt.Sprintf("Key \"%s\" should be set with WhitelistKnownKeys before becoming Mandatory", k))
		}
		config.Vals[k].Mandatory = true
	}
}

// Allow keynames and their values based on examples
// Can be called multiple times to add examples for keynames or values
// Use this when the keynames are not known in advance
// Call multiple times to show the entire range of keynames and values
// When keys are known in advance, use WhitelistKnownKeys() instead
func (config *KeyValConfig) WhitelistByExample(k string, v string) {
	if config.OtherKeynames == nil {
		config.OtherKeynames = new(SimpleValConfig)

	}
	config.OtherKeynames.AddValExample(k)

	if config.OtherVals == nil {
		config.OtherVals = new(SimpleValConfig)

	}
	config.OtherVals.AddValExample(v)
}

func (config *KeyValConfig) Describe() string {
	var description bytes.Buffer

	if config.Vals != nil {
		for k, v := range config.Vals {
			if v.Mandatory {
				description.WriteString(" | MandatoryKey: ")
			} else {
				description.WriteString(" | OptionalKey: ")
			}
			description.WriteString(k)
			description.WriteString(" => ")
			description.WriteString(v.Describe())
		}
	}

	if config.OtherVals != nil {
		description.WriteString(" | OtherVals: ")
		description.WriteString(config.OtherVals.Describe())
	}
	if config.OtherKeynames != nil {
		description.WriteString(" | OtherKeynames: ")
		description.WriteString(config.OtherKeynames.Describe())
	}

	return description.String()
}

func (config *KeyValConfig) Marshal(depth int) string {
	var description bytes.Buffer
	var started bool
	started = false
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")

	if len(config.Vals) > 0 {
		description.WriteString(shift)
		description.WriteString("  Vals: {\n")
		for k, v := range config.Vals {
			description.WriteString(shift)
			description.WriteString(fmt.Sprintf("  , %s: %s", k, v.Marshal(depth+1)))
		}
		description.WriteString(shift)
		description.WriteString("  }\n")
		started = true
	}

	if config.OtherKeynames != nil {
		description.WriteString(shift)
		if started {
			description.WriteString(", ")
		} else {
			description.WriteString("  ")
		}
		description.WriteString(fmt.Sprintf("OtherKeynames: %s\n", config.OtherKeynames.Marshal(depth+1)))
		started = true
	}
	if config.OtherVals != nil {
		description.WriteString(shift)
		if started {
			description.WriteString(", ")
		} else {
			description.WriteString("  ")
		}
		description.WriteString(fmt.Sprintf("OtherVals: %s\n", config.OtherVals.Marshal(depth+1)))
	}
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}
