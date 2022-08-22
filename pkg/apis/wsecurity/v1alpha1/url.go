package v1alpha1

import (
	"fmt"
	"strings"
)

//////////////////// UrlProfile ////////////////

// Exposes ValueProfile interface
type UrlProfile struct {
	Val      SimpleValProfile `json:"val"`
	Segments CountProfile     `json:"segments"`
}

func (profile *UrlProfile) Profile(args ...interface{}) {
	path := args[0].(string)
	segments := strings.Split(path, "/")
	numSegments := len(segments)
	if (numSegments > 0) && segments[0] == "" {
		segments = segments[1:]
		numSegments--
	}
	if (numSegments > 0) && segments[numSegments-1] == "" {
		numSegments--
		segments = segments[:numSegments]

	}
	cleanPath := strings.Join(segments, "")
	//profile.Val = new(SimpleValProfile)
	profile.Val.Profile(cleanPath)

	if numSegments > 0xFF {
		numSegments = 0xFF
	}
	profile.Segments.Profile(uint8(numSegments))
}

//////////////////// UrlPile ////////////////

// Exposes ValuePile interface
type UrlPile struct {
	Val      SimpleValPile `json:"val"`
	Segments CountPile     `json:"segments"`
}

func (pile *UrlPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*UrlProfile)
	pile.Segments.Add(&profile.Segments)
	pile.Val.Add(&profile.Val)
}

func (pile *UrlPile) Clear() {
	pile.Segments.Clear()
	pile.Val.Clear()
}
func (pile *UrlPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*UrlPile)
	pile.Segments.Merge(&otherPile.Segments)
	pile.Val.Merge(&otherPile.Val)
}

//////////////////// UrlConfig ////////////////

// Exposes ValueConfig interface
type UrlConfig struct {
	Val      SimpleValConfig `json:"val"`
	Segments CountConfig     `json:"segments"`
}

func (config *UrlConfig) Learn(valPile ValuePile) {
	pile := valPile.(*UrlPile)
	config.Segments.Learn(&pile.Segments)
	config.Val.Learn(&pile.Val)
}

func (config *UrlConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*UrlConfig)
	config.Segments.Fuse(&otherConfig.Segments)
	config.Val.Fuse(&otherConfig.Val)
}

func (config *UrlConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*UrlProfile)
	if str := config.Segments.Decide(&profile.Segments); str != "" {
		return fmt.Sprintf("Segmengs: %s", str)
	}

	if str := config.Val.Decide(&profile.Val); str != "" {
		return fmt.Sprintf("KeyVal: %s", str)
	}
	return ""
}
