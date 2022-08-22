package v1alpha1

import (
	"fmt"
	"net/http"
)

var exceptionHeaders = map[string]bool{"Content-Type": true}

//////////////////// HeadersProfile ////////////////

// Exposes ValueProfile interface
type HeadersProfile struct {
	Kv KeyValProfile `json:"kv"`
}

func (profile *HeadersProfile) Profile(args ...interface{}) {
	headers := (map[string][]string)(args[0].(http.Header))
	//profile.Kv = new(KeyValProfile)
	profile.Kv.Profile(headers, exceptionHeaders)
}

//////////////////// HeadersPile ////////////////

// Exposes ValuePile interface
type HeadersPile struct {
	Kv *KeyValPile `json:"kv"`
}

func (pile *HeadersPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*HeadersProfile)
	if pile.Kv == nil {
		pile.Kv = new(KeyValPile)
	}
	pile.Kv.Add(&profile.Kv)
}

func (pile *HeadersPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*HeadersPile)
	if pile.Kv == nil {
		pile.Kv = new(KeyValPile)
	}
	pile.Kv.Merge(otherPile.Kv)
}

func (pile *HeadersPile) Clear() {
	pile.Kv = new(KeyValPile)
	if pile.Kv != nil {
		pile.Kv.Clear()
	}
}

//////////////////// HeadersConfig ////////////////

// Exposes ValueConfig interface
type HeadersConfig struct {
	Kv KeyValConfig `json:"kv"`
}

func (config *HeadersConfig) Learn(valPile ValuePile) {
	pile := valPile.(*HeadersPile)
	config.Kv.Learn(pile.Kv)
}

func (config *HeadersConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*HeadersConfig)
	config.Kv.Fuse(&otherConfig.Kv)
}

func (config *HeadersConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*HeadersProfile)
	str := config.Kv.Decide(&profile.Kv)
	if str == "" {
		return str
	}
	return fmt.Sprintf("KeyVal: %s", str)
}
