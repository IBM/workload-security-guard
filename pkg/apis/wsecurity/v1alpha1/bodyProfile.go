package v1alpha1

import "fmt"

//////////////////// BodyProfile ////////////////

// Exposes ValueProfile interface
type BodyProfile struct {
	Unstructured *SimpleValProfile  `json:"unstructured"`
	Structured   *StructuredProfile `json:"structured"`
}

func (profile *BodyProfile) Profile(args ...interface{}) {
	//*profile = AsciiFlagsProfile(args[0].(uint32))
}

/*
// Future: Missing implementation
func (profile *BodyProfile) String(depth int) string {
	return "Missing Implementation"
}
*/
//////////////////// BodyPile ////////////////

// Exposes ValuePile interface
type BodyPile struct {
	Unstructured *SimpleValPile  `json:"unstructured"`
	Structured   *StructuredPile `json:"structured"`
}

// Future: Missing implementation
func (pile *BodyPile) Add(valProfile ValueProfile) {
	//profile := valProfile.(*BodyProfile)

}

// Future: Missing implementation
func (pile *BodyPile) Clear() {

}

// Future: Missing implementation
func (pile *BodyPile) Merge(otherValPile ValuePile) {
	//otherPile := otherValPile.(*BodyPile)

}

// Future: Missing implementation
//func (pile *BodyPile) String(depth int) string {
//	return "Missing Implementation"
//}

//////////////////// BodyConfig ////////////////

// Exposes ValueConfig interface
type BodyConfig struct {
	Unstructured *SimpleValConfig  `json:"unstructured"`
	Structured   *StructuredConfig `json:"structured"`
}

func (config *BodyConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*BodyProfile)

	if profile.Structured != nil {
		if config.Structured != nil {
			str := config.Structured.Decide(profile.Structured)
			if str != "" {
				return fmt.Sprintf("Body %s", str)
			}
		} else {
			return "Structured Body not allowed"
		}
	}
	if profile.Unstructured != nil {
		if config.Unstructured != nil {
			str := config.Unstructured.Decide(profile.Unstructured)
			if str != "" {
				return fmt.Sprintf("Body %s", str)
			}
		} else {
			return "Unstructured Body not allowed"
		}
	}
	return ""
}

// Future: Missing implementation
func (config *BodyConfig) Learn(valPile ValuePile) {
	//pile := valPile.(*BodyPile)

}

// Future: Missing implementation
func (config *BodyConfig) Fuse(otherValConfig ValueConfig) {
	//otherConfig := otherValConfig.(*BodyConfig)

}
