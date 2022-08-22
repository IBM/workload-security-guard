package v1alpha1

import (
	"fmt"
	"strings"
)

//////////////////// KeyValProfile ////////////////

// Exposes ValueProfile interface
type KeyValProfile map[string]*SimpleValProfile

// Profile a generic map of key vals where we expect:
func (profile *KeyValProfile) Profile(args ...interface{}) {
	switch keyValMap := args[0].(type) {
	case nil:
		*profile = nil
	case map[string]string:
		if len(keyValMap) == 0 { // no keys
			*profile = nil
			return
		}
		*profile = make(map[string]*SimpleValProfile, len(keyValMap))
		for k, v := range keyValMap {
			// Profile the concatenated value
			(*profile)[k] = new(SimpleValProfile)
			(*profile)[k].Profile(v)
		}
	case map[string][]string:
		if len(keyValMap) == 0 { // no keys
			*profile = nil
			return
		}
		*profile = make(map[string]*SimpleValProfile, len(keyValMap))
		for k, v := range keyValMap {
			// Concatenate all strings into one value
			// Appropriate for evaluating []string where order should be also preserved
			val := strings.Join(v, " ")

			// Profile the concatenated value
			(*profile)[k] = new(SimpleValProfile)
			(*profile)[k].Profile(val)
		}
	default:
		panic("Unsupported type in KeyValProfile")
	}

}

//////////////////// KeyValPile ////////////////

// Exposes ValuePile interface
type KeyValPile map[string]*SimpleValPile

func (pile *KeyValPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*KeyValProfile)

	if *pile == nil {
		*pile = make(map[string]*SimpleValPile, 16)
	}
	for key, v := range *profile {
		svp, exists := (*pile)[key]
		if !exists {
			svp = new(SimpleValPile)
			(*pile)[key] = svp
		}
		svp.Add(v)
	}
}

func (pile *KeyValPile) Clear() {
	*pile = nil
}

func (pile *KeyValPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*KeyValPile)

	if otherPile == nil {
		return
	}
	if *pile == nil {
		*pile = *otherPile
		return
	}
	for key, val := range *otherPile {
		if myVal, exists := (*pile)[key]; exists {
			myVal.Merge(val)
		} else {
			(*pile)[key] = val
		}
	}
}

//////////////////// KeyValConfig ////////////////

// Exposes ValueConfig interface
type KeyValConfig struct {
	Vals          map[string]*SimpleValConfig `json:"vals"`          // Profile the value of whitelisted keys
	OtherVals     *SimpleValConfig            `json:"otherVals"`     // Profile the values of other keys
	OtherKeynames *SimpleValConfig            `json:"otherKeynames"` // Profile the keynames of other keys
}

func (config *KeyValConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*KeyValProfile)

	if profile == nil {
		return ""
	}

	// For each key-val, decide
	for k, v := range *profile {
		// Decide based on a known keys
		if config.Vals != nil && config.Vals[k] != nil {
			if ret := config.Vals[k].Decide(v); ret != "" {
				return fmt.Sprintf("Known Key %s: %s", k, ret)
			}
			continue
		}
		// Decide based on unknown key...
		if config.OtherKeynames == nil || config.OtherVals == nil {
			return fmt.Sprintf("Key %s is not known", k)
		}
		// Cosnider the keyname
		var keyname SimpleValProfile
		keyname.Profile(k)
		if ret := config.OtherKeynames.Decide(&keyname); ret != "" {
			return fmt.Sprintf("Other keyname %s: %s", k, ret)
		}
		// Cosnider the key value
		if ret := config.OtherVals.Decide(v); ret != "" {
			return fmt.Sprintf("Other keyname %s: %s", k, ret)
		}
		continue
	}
	return ""
}

// Learn implementation currently is not optimized for a large number of keys
// When the number of keys grow, Learn may reduce the number of known keys by
// aggregating all known keys which have common low security fingerprint into
// OtherKeynames and OtherVals
// TBD - left for future
func (config *KeyValConfig) Learn(valPile ValuePile) {
	pile := valPile.(*KeyValPile)

	config.OtherVals = nil
	config.OtherKeynames = nil

	if pile == nil {
		config.Vals = nil
		return
	}

	// learn known keys
	config.Vals = make(map[string]*SimpleValConfig, len(*pile))
	for k, v := range *pile {
		svc := new(SimpleValConfig)
		svc.Learn(v)
		config.Vals[k] = svc
	}
}

func (config *KeyValConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*KeyValConfig)

	if otherConfig == nil {
		return
	}
	if config.Vals == nil {
		config.Vals = otherConfig.Vals
	} else {
		// fuse known keys
		for k, v := range otherConfig.Vals {
			svc, exists := config.Vals[k]
			if exists {
				svc.Fuse(v)
			} else {
				config.Vals[k] = v
			}
		}
	}

	// fuse keynames of unknown keys
	if config.OtherKeynames == nil {
		config.OtherKeynames = otherConfig.OtherKeynames
	} else {
		config.OtherKeynames.Fuse(otherConfig.OtherKeynames)
	}

	// fuse key values of unknown keys
	if config.OtherVals == nil {
		config.OtherVals = otherConfig.OtherVals
	} else {
		config.OtherVals.Fuse(otherConfig.OtherVals)
	}
}
