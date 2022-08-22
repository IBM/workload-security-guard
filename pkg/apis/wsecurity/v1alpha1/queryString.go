package v1alpha1

import (
	"fmt"
	"net/url"
)

//////////////////// QueryProfile ////////////////

// Exposes ValueProfile interface
type QueryProfile struct {
	Kv KeyValProfile `json:"kv"`
}

func (profile *QueryProfile) Profile(args ...interface{}) {
	values := (map[string][]string)(args[0].(url.Values))

	profile.Kv.Profile(values, nil)
}

//////////////////// QueryPile ////////////////

// Exposes ValuePile interface
type QueryPile struct {
	Kv *KeyValPile `json:"kv"`
}

func (pile *QueryPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*QueryProfile)
	if pile.Kv == nil {
		pile.Kv = new(KeyValPile)
	}
	pile.Kv.Add(&profile.Kv)
}

func (pile *QueryPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*QueryPile)
	if pile.Kv == nil {
		pile.Kv = new(KeyValPile)
	}
	pile.Kv.Merge(otherPile.Kv)
}

func (pile *QueryPile) Clear() {
	pile.Kv = new(KeyValPile)
	if pile.Kv != nil {
		pile.Kv.Clear()
	}
}

//////////////////// QueryConfig ////////////////

// Exposes ValueConfig interface
type QueryConfig struct {
	Kv KeyValConfig `json:"kv"`
}

func (config *QueryConfig) Learn(valPile ValuePile) {
	pile := valPile.(*QueryPile)
	config.Kv.Learn(pile.Kv)
}

func (config *QueryConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*QueryConfig)
	config.Kv.Fuse(&otherConfig.Kv)
}

func (config *QueryConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*QueryProfile)
	str := config.Kv.Decide(&profile.Kv)
	if str == "" {
		return str
	}
	return fmt.Sprintf("KeyVal: %s", str)
}
