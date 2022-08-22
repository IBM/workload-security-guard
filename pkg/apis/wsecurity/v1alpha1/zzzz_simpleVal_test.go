package v1alpha1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSimpleValProfile_Profile(t *testing.T) {
	nonReadables := []byte{0x7, 0x11, 0x10, 0x0, 0x7F}

	type fields struct {
		Flags        uint32
		NonReadables uint8
		Spaces       uint8
		Unicodes     uint8
		Digits       uint8
		Letters      uint8
		SpecialChars uint8
		Sequences    uint8
		UnicodeFlags []uint32
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{{
		name: "letters",
		fields: fields{
			Flags:        0xc,
			NonReadables: 0,
			Spaces:       0,
			Unicodes:     0,
			Digits:       0,
			Letters:      3,
			SpecialChars: 2,
			Sequences:    2,
			UnicodeFlags: nil,
		},
		args: args{
			args: []interface{}{"aac$#"},
		},
	}, {
		name: "digits and spaces",
		fields: fields{
			Flags:        0,
			NonReadables: 0,
			Spaces:       5,
			Unicodes:     0,
			Digits:       6,
			Letters:      0,
			SpecialChars: 0,
			Sequences:    5,
			UnicodeFlags: nil,
		},
		args: args{
			args: []interface{}{" 12 1234   "},
		},
	}, {
		name: "Nonreadables",
		fields: fields{
			Flags:        0,
			NonReadables: 5,
			Spaces:       0,
			Unicodes:     0,
			Digits:       0,
			Letters:      0,
			SpecialChars: 0,
			Sequences:    1,
			UnicodeFlags: nil,
		},
		args: args{
			args: []interface{}{string(nonReadables)},
		},
	}, {
		name: "Lorem Ipsum",
		fields: fields{
			Flags:        0x41440,
			NonReadables: 16,
			Spaces:       97,
			Unicodes:     0,
			Digits:       8,
			Letters:      255,
			SpecialChars: 11,
			Sequences:    214,
			UnicodeFlags: nil,
		},
		args: args{
			args: []interface{}{loremIpsum},
		},
	}, {
		name: "hebrew1",
		fields: fields{
			Flags:        0x40002100,
			NonReadables: 0,
			Spaces:       0,
			Unicodes:     4,
			Digits:       2,
			Letters:      1,
			SpecialChars: 4,
			Sequences:    5,
			UnicodeFlags: []uint32{1024},
		},
		args: args{
			args: []interface{}{hebrew1},
		},
	}, {
		name: "hebrew2",
		fields: fields{
			Flags:        0x40002100,
			NonReadables: 0,
			Spaces:       0,
			Unicodes:     8,
			Digits:       2,
			Letters:      1,
			SpecialChars: 4,
			Sequences:    5,
			UnicodeFlags: []uint32{1024},
		},
		args: args{
			args: []interface{}{hebrew2},
		},
	}, {
		name: "Chineese",
		fields: fields{
			Flags:        0x40002100,
			NonReadables: 0,
			Spaces:       0,
			Unicodes:     4,
			Digits:       2,
			Letters:      1,
			SpecialChars: 4,
			Sequences:    5,
			UnicodeFlags: []uint32{0, 0, 0, 0, 134217728, 0, 0, 512},
		},
		args: args{
			args: []interface{}{chineese},
		},
	}}
	for _, tt := range tests {
		var args []interface{}
		var profiles []ValueProfile
		var piles []ValuePile
		var configs []ValueConfig
		for i := 0; i < 10; i++ {
			profiles = append(profiles, new(SimpleValProfile))
			piles = append(piles, new(SimpleValPile))
			configs = append(configs, new(SimpleValConfig))
		}
		args = append(args, tt.args.args)

		ValueTests_Test(t, profiles, piles, configs, args...)

		t.Run(tt.name, func(t *testing.T) {
			profileResult := &SimpleValProfile{
				Flags:        AsciiFlagsProfile(tt.fields.Flags),
				NonReadables: CountProfile(tt.fields.NonReadables),
				Spaces:       CountProfile(tt.fields.Spaces),
				Unicodes:     CountProfile(tt.fields.Unicodes),
				Digits:       CountProfile(tt.fields.Digits),
				Letters:      CountProfile(tt.fields.Letters),
				SpecialChars: CountProfile(tt.fields.SpecialChars),
				Sequences:    CountProfile(tt.fields.Sequences),
				UnicodeFlags: FlagSliceProfile(tt.fields.UnicodeFlags),
			}
			if !reflect.DeepEqual(profileResult, profiles[0]) {
				t.Errorf("Profile() want %v got %v\n%v\n%v", profileResult, profiles[0], profileResult, profiles[0])
			}
		})
	}
}

func TestSimpleValProfile_Minimal(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		var profile1, profile2, profile3 SimpleValProfile
		var pile1, pile2, pile3 SimpleValPile
		var config1, config2 SimpleValConfig
		profile1.Profile("abc")
		profile2.Profile("123abc")
		profile3.Profile("abcd")
		pile1.Add(&profile1)
		pile2.Add(&profile3)
		pile2.Add(&profile2)
		pile3.Add(&profile2)
		pile1.Merge(&pile3)
		config1.Learn(&pile1)
		config2.Learn(&pile2)
		config1.Fuse(&config2)
		if config1.Digits[0].Min != 0 || config1.Digits[0].Max != 3 ||
			config1.Letters[0].Min != 3 || config1.Letters[0].Max != 4 {
			t.Errorf("Config got %v\n", config1)

		}
		var bytes []byte
		bytes, err := json.Marshal(pile1)
		if err != nil {
			t.Errorf("json.Marshal Error %v", err.Error())
		}
		fmt.Println(string(bytes))

	})
}
