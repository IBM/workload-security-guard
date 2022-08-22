package v1alpha1

import (
	"fmt"
	"net/http"
)

// FUTURE Add support for Trailers, review Resp to find more security data
//////////////////// RespProfile ////////////////

// Exposes ValueProfile interface
type RespProfile struct {
	Headers       HeadersProfile   `json:"headers"`
	MediaType     MediaTypeProfile `json:"mediatype"`     // "text/html"
	ContentLength CountProfile     `json:"contentlength"` // 0
}

func (profile *RespProfile) Profile(args ...interface{}) {
	resp := args[0].(*http.Response)

	profile.Headers.Profile(resp.Header)

	profile.MediaType.Profile(resp.Header.Get("Content-Type"))

	length := resp.ContentLength
	if length > 0 {
		var log2length uint8

		for length > 0 {
			log2length++
			length >>= 1
		}
		profile.ContentLength = CountProfile(log2length)
	}

}

//////////////////// RespPile ////////////////

// Exposes ValuePile interface
type RespPile struct {
	Headers       HeadersPile   `json:"headers"`
	MediaType     MediaTypePile `json:"mediatype"`
	ContentLength CountPile     `json:"contentlength"`
}

func (pile *RespPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*RespProfile)

	pile.Headers.Add(&profile.Headers)
	pile.MediaType.Add(&profile.MediaType)
	pile.ContentLength.Add(&profile.ContentLength)
}

func (pile *RespPile) Clear() {
	pile.Headers.Clear()
	pile.MediaType.Clear()
	pile.ContentLength.Clear()
}

func (pile *RespPile) Merge(otherValProfile ValuePile) {
	otherPile := otherValProfile.(*RespPile)

	pile.Headers.Merge(&otherPile.Headers)
	pile.MediaType.Merge(&otherPile.MediaType)
	pile.ContentLength.Merge(&otherPile.ContentLength)
}

//////////////////// RespConfig ////////////////

// Exposes ValueConfig interface
type RespConfig struct {
	Headers       HeadersConfig   `json:"headers"`
	MediaType     MediaTypeConfig `json:"mediatype"`
	ContentLength CountConfig     `json:"contentlength"`
}

func (config *RespConfig) Learn(valPile ValuePile) {
	pile := valPile.(*RespPile)

	config.Headers.Learn(&pile.Headers)
	config.MediaType.Learn(&pile.MediaType)
	config.ContentLength.Learn(&pile.ContentLength)
}

func (config *RespConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*RespConfig)

	config.Headers.Fuse(&otherConfig.Headers)
	config.MediaType.Fuse(&otherConfig.MediaType)
	config.ContentLength.Fuse(&otherConfig.ContentLength)
}

func (config *RespConfig) Decide(valProfile ValueProfile) string {
	profile := valProfile.(*RespProfile)

	if ret := config.Headers.Decide(&profile.Headers); ret != "" {
		return fmt.Sprintf("Headers: %s", ret)
	}
	if ret := config.MediaType.Decide(&profile.MediaType); ret != "" {
		return fmt.Sprintf("Media Type: %s", ret)
	}
	if ret := config.ContentLength.Decide(&profile.ContentLength); ret != "" {
		return fmt.Sprintf("Content Length: %s", ret)
	}
	return ""
}
