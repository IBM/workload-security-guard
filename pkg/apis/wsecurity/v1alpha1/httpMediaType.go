package v1alpha1

import (
	"fmt"
	"mime"
)

//////////////////// MediaTypeProfile ////////////////

// Exposes ValueProfile interface
// TypeToken include rfc7231 type "/" subtype
type MediaTypeProfile struct {
	TypeTokens SetProfile    `json:"type"`   // "text/html"
	Params     KeyValProfile `json:"params"` // {"charset": "utf-8"}
}

func (profile *MediaTypeProfile) Profile(args ...interface{}) {

	str := args[0].(string)
	if mediaType, params, err := mime.ParseMediaType(str); err == nil && mediaType != "" {
		profile.TypeTokens.Profile(mediaType)
		profile.Params.Profile(params)
		return
	}
	// For clients that fail to send media type
	profile.TypeTokens.Profile("none")
	profile.Params.Profile(nil)
}

//////////////////// MediaTypePile ////////////////

// Exposes ValuePile interface
type MediaTypePile struct {
	TypeTokens SetPile    `json:"type"`
	Params     KeyValPile `json:"params"`
}

func (pile *MediaTypePile) Add(valProfile ValueProfile) {
	profile := valProfile.(*MediaTypeProfile)

	pile.TypeTokens.Add(&profile.TypeTokens)
	pile.Params.Add(&profile.Params)
}

func (pile *MediaTypePile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*MediaTypePile)

	pile.TypeTokens.Merge(&otherPile.TypeTokens)
	pile.Params.Merge(&otherPile.Params)
}

func (pile *MediaTypePile) Clear() {

	pile.TypeTokens.Clear()
	pile.Params.Clear()
}

//////////////////// MediaTypeConfig ////////////////

// Exposes ValueConfig interface
type MediaTypeConfig struct {
	TypeTokens SetConfig    `json:"type"`
	Params     KeyValConfig `json:"params"`
}

func (config *MediaTypeConfig) Learn(valPile ValuePile) {
	pile := valPile.(*MediaTypePile)

	config.TypeTokens.Learn(&pile.TypeTokens)
	config.Params.Learn(&pile.Params)
}

func (config *MediaTypeConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*MediaTypeConfig)

	config.TypeTokens.Fuse(&otherConfig.TypeTokens)
	config.Params.Fuse(&otherConfig.Params)
}

func (config *MediaTypeConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*MediaTypeProfile)

	if str := config.TypeTokens.Decide(&profile.TypeTokens); str != "" {
		return fmt.Sprintf("Type: %s", str)
	}
	if str := config.Params.Decide(&profile.Params); str != "" {
		return fmt.Sprintf("Params: %s", str)
	}
	return ""
}
