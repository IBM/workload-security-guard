package v1alpha1

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const str1 = `What is Lorem Ipsum?
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was xxx in the 1960s with the release of xxx sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like xxx PageMaker including versions of Lorem Ipsum.

Why do we use it?
`

var flags1 uint32 = uint32(0b1000001010001000000)

func setFlags(slots []int) (f uint32) {
	for _, slot := range slots {
		f = f | (0x1 << slot)
	}
	return
}

func TestConfigSimpleVals(t *testing.T) {
	var svp SimpleValProfile
	var svpLonger SimpleValProfile
	var svpShorter SimpleValProfile
	var svpSpecial SimpleValProfile
	var svpUnicode SimpleValProfile
	var svc SimpleValConfig
	str := "A23/*מכבי*/"
	strLonger := "A23/*מכבימכבי*/"
	strShorter := "A23/*מי*/"
	strSpecial := "A23--מכבי*/"
	strUnicode := "A23/*世界世界*/"
	t.Run("LoremIpsum", func(t *testing.T) {
		//svc.AddValExample(str)
		svp.Profile(str)
		//svp.Describe()
		/*
			if ret := svc.Decide(&svp); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
		*/svpLonger.Profile(strLonger)
		svpShorter.Profile(strShorter)
		svpSpecial.Profile(strSpecial)
		svpUnicode.Profile(strUnicode)

		fmt.Printf("svc %v\n", svc)
		fmt.Printf("svpLonger %c\n", svpLonger)
		/*
			if ret := svc.Decide(&svpLonger); ret == "" {
				t.Errorf("ProfileSimpleVal() Decided ok but expected a reject")
			}
			if ret := svc.Decide(&svpShorter); ret == "" {
				t.Errorf("ProfileSimpleVal() Decided ok but expected a reject")
			}
			if ret := svc.Decide(&svpSpecial); ret == "" {
				t.Errorf("ProfileSimpleVal() Decided ok but expected a reject")
			}
			if ret := svc.Decide(&svpUnicode); ret == "" {
				t.Errorf("ProfileSimpleVal() Decided ok but expected a reject")
			}
		*/
		//svc.AddValExample(str)
		//svc.AddValExample(strLonger)
		//svc.AddValExample(strShorter)
		//svc.AddValExample(strSpecial)
		//svc.AddValExample(strUnicode)
		/*
			if ret := svc.Decide(&svp); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
			if ret := svc.Decide(&svpLonger); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
			if ret := svc.Decide(&svpShorter); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
			if ret := svc.Decide(&svpSpecial); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
			if ret := svc.Decide(&svpUnicode); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
		*/
	})

}

func TestProfileSimpleVals(t *testing.T) {
	var svp *SimpleValProfile
	var svc *SimpleValConfig

	unicode := []uint32{}

	//svc := new(SimpleValConfig)
	t.Run("LoremIpsum", func(t *testing.T) {
		//fmt.Printf("tt.args.str %s", tt.args.str)
		start := time.Now()

		svp = new(SimpleValProfile)
		svc = new(SimpleValConfig)
		svPile := new(SimpleValPile)

		svp.Profile(str1)
		Confirm(t, svp, 255, 8, 11, 97, 4, flags1, 0, unicode)

		svPile.Add(svp)
		//svPile.String(0)
		svPile.Merge(svPile)
		if ret := svc.Decide(svp); ret == "" {
			t.Errorf("ProfileSimpleVal() Decided ok but expected a reject")
		}

		//svc.AddValExample(str1)
		/*
			if ret := svc.Decide(svp); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
			if ret := svc.Decide(svp); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
		*/

		elapsed := time.Since(start)
		fmt.Printf("Time is %s", elapsed)
		//svc.Normalize()
		/*
			if ret := svc.Decide(svp); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}
			svc.Learn(svPile)
			if ret := svc.Decide(svp); ret != "" {
				t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
			}

			//svc = NewSimpleValConfig(0, 0, 0, 32, 32, 0, 16)
			if ret := svc.Decide(svp); ret == "" {
				t.Errorf("ProfileSimpleVal() Decided ok but expected a reject")
			}
		*/
	})

	var flags uint32
	var str, name string
	flags = 0x1<<CommentsSlot | 0x1<<SlashSlot | 0x1<<AsteriskSlot
	str = "/*"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp = new(SimpleValProfile)
		svp.Profile(str)
		Confirm(t, svp, 0, 0, 2, 0, 0, flags, 0, unicode)
	})
	flags = 0x1<<CommentsSlot | 0x1<<SlashSlot | 0x1<<AsteriskSlot
	str = "*/"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp = new(SimpleValProfile)
		svp.Profile(str)
		Confirm(t, svp, 0, 0, 2, 0, 0, flags, 0, unicode)
	})
	flags = 0x1 << HexSlot //| 0x1<<DigitSlot | 0x1<<LetterSlot
	str = "0x"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		Confirm(t, svp, 1, 1, 0, 0, 0, flags, 0, unicode)
	})
	flags = 0x1 << HexSlot //| 0x1<<DigitSlot | 0x1<<LetterSlot
	str = "0X"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		Confirm(t, svp, 1, 1, 0, 0, 0, flags, 0, unicode)
	})

	const str2 = "日本語"
	flags = 0
	unicode = []uint32{0, 0, 0, 0, 0, 0, 9216, 0, 1048576}
	name = fmt.Sprintf("testString %s", str2)

	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str2)
		Confirm(t, svp, 0, 0, 0, 0, 0, 0, 3, unicode)
	})
	const str3 = "{!}"
	flags = setFlags([]int{ExclamationSlot, CurlyBracketSlot})
	unicode = []uint32{}
	name = fmt.Sprintf("testString %s", str3)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str3)
		Confirm(t, svp, 0, 0, 3, 0, 0, flags, 0, unicode)
	})
	const str4 = "123"
	flags = 0
	name = fmt.Sprintf("testString %s", str4)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str4)
		Confirm(t, svp, 0, 3, 0, 0, 0, flags, 0, unicode)
	})
	const str5 = "aBc"
	flags = 0
	name = fmt.Sprintf("testString %s", str5)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str5)
		Confirm(t, svp, 3, 0, 0, 0, 0, flags, 0, unicode)
	})
	var str6 string = string([]rune{rune(200), rune(201), rune(202)})
	flags = 0
	unicode = []uint32{1}
	name = fmt.Sprintf("testString %s", str6)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str6)
		Confirm(t, svp, 0, 0, 0, 0, 0, flags, 3, unicode)
	})
	byteSlice := make([]rune, 255*2*26)
	for j := 0; j < 26; j++ {
		for i := rune(0); i < 255; i++ {
			byteSlice[int(j)*255*2+2*int(i)] = i
			byteSlice[int(j)*255*2+2*int(i)+1] = 32
		}
	}
	str7 := string(byteSlice)
	flags = 0x0FFFFFFF
	unicode = []uint32{1}
	t.Run("testString all letters", func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str7)
		Confirm(t, svp, 255, 255, 255, 255, 255, flags, 255, unicode)
	})
}

func ConfirmMulti(t *testing.T, svp *SimpleValProfile, letters uint8, digits uint8, specialChars uint8, spaces uint8,
	nonReadables uint8, flags uint32, unicodes uint8, unicodeFlags []uint32, multi uint8) {
	Confirm(t, svp, letters*multi, digits*multi, specialChars*multi, spaces*multi, nonReadables*multi, flags, unicodes*multi, unicodeFlags)
}

func Confirm(t *testing.T, svp *SimpleValProfile, letters uint8, digits uint8, specialChars uint8, spaces uint8,
	nonReadables uint8, flags uint32, unicodes uint8, unicodeFlags []uint32) {
	if uint32(svp.Flags) != flags {
		t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
	}
	if uint8(svp.Spaces) != spaces {
		t.Errorf("ProfileSimpleVal() Spaces = %d instead of %d", svp.Spaces, spaces)
	}
	if uint8(svp.Letters) != letters {
		t.Errorf("ProfileSimpleVal() Letters = %d instead of %d", svp.Letters, letters)
	}
	if uint8(svp.Digits) != digits {
		t.Errorf("ProfileSimpleVal() Digits = %d instead of %d", svp.Digits, digits)
	}
	if uint8(svp.SpecialChars) != specialChars {
		t.Errorf("ProfileSimpleVal() SpecialChars = %d instead of %d", svp.SpecialChars, specialChars)
	}
	if uint8(svp.NonReadables) != nonReadables {
		t.Errorf("ProfileSimpleVal() NonReadables = %d instead of %d", svp.NonReadables, nonReadables)
	}
	if uint8(svp.Unicodes) != unicodes {
		t.Errorf("ProfileSimpleVal() Unicodes = %d instead of %d", svp.Unicodes, 0)
	}
	if len(svp.UnicodeFlags) != len(unicodeFlags) {
		t.Errorf("ProfileSimpleVal() UnicodeFlags = %d instead of %d", len(svp.UnicodeFlags), len(unicodeFlags))
	}
	for i := 0; i < len(unicodeFlags); i++ {
		if svp.UnicodeFlags[i] != unicodeFlags[i] {
			t.Errorf("ProfileSimpleVal() UnicodeFlags[%d] = %d instead of %d", i, svp.UnicodeFlags[i], unicodeFlags[i])
		}
	}
}

func TestProfileChars0_127(t *testing.T) {
	var unicodeFlags []uint32
	var c uint8
	for c = 0; c < 128; c++ {
		r := rune(c)
		str := string(r)

		name := fmt.Sprintf("testString %s", str)
		t.Run(name, func(t *testing.T) {
			svp := new(SimpleValProfile)
			svp.Profile(str)

			switch asciiMap[c] {
			case LetterSlot:
				Confirm(t, svp, 1, 0, 0, 0, 0, 0, 0, unicodeFlags)
			case DigitSlot:
				Confirm(t, svp, 0, 1, 0, 0, 0, 0, 0, unicodeFlags)
			case NonReadableSlot:
				Confirm(t, svp, 0, 0, 0, 0, 1, 0, 0, unicodeFlags)
			case SpaceSlot:
				Confirm(t, svp, 0, 0, 0, 1, 0, 0, 0, unicodeFlags)
			default:
				Confirm(t, svp, 0, 0, 1, 0, 0, 1<<asciiMap[c], 0, unicodeFlags)
			}
		})
	}
}

func TestProfileChars128_2560(t *testing.T) {
	var unicodeFlags []uint32
	for c := 128; c <= 25600; c += 128 {
		r := rune(c)
		str := string(r)
		name := fmt.Sprintf("testString %s", str)
		t.Run(name, func(t *testing.T) {
			svp := new(SimpleValProfile)
			svp.Profile(str)

			block := (c / 0x80) - 1
			blockBit := int(block & 0x1F)
			blockElement := int(block / 0x20)

			unicodeFlags = make([]uint32, blockElement+1)
			unicodeFlags[blockElement] |= 0x1 << blockBit

			Confirm(t, svp, 0, 0, 0, 0, 0, 0, 1, unicodeFlags)
		})
	}
}

func MultiCharTest(t *testing.T, str string, num uint8) {
	var unicodeFlags []uint32
	name := fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)

		switch asciiMap[str[0]] {
		case LetterSlot:
			ConfirmMulti(t, svp, 1, 0, 0, 0, 0, 0, 0, unicodeFlags, num)
		case DigitSlot:
			ConfirmMulti(t, svp, 0, 1, 0, 0, 0, 0, 0, unicodeFlags, num)
		case NonReadableSlot:
			ConfirmMulti(t, svp, 0, 0, 0, 0, 1, 0, 0, unicodeFlags, num)
		case SpaceSlot:
			ConfirmMulti(t, svp, 0, 0, 0, 1, 0, 0, 0, unicodeFlags, num)
		default:
			ConfirmMulti(t, svp, 0, 0, 1, 0, 0, 1<<asciiMap[str[0]], 0, unicodeFlags, num)
		}
	})
}

func TestProfileRepeatedChars(t *testing.T) {
	t.Run("Repeated NonReadable", func(t *testing.T) {
		MultiCharTest(t, strings.Repeat(string(rune(4)), 10), 10)
		MultiCharTest(t, strings.Repeat(string(rune(4)), 1000), 255)
	})
	t.Run("Repeated Letter", func(t *testing.T) {
		MultiCharTest(t, strings.Repeat(string(rune(65)), 10), 10)
		MultiCharTest(t, strings.Repeat(string(rune(65)), 1000), 255)
	})
	t.Run("Repeated Digit", func(t *testing.T) {
		MultiCharTest(t, strings.Repeat(string(rune(48)), 10), 10)
		MultiCharTest(t, strings.Repeat(string(rune(48)), 1000), 255)
	})
	t.Run("Repeated Space", func(t *testing.T) {
		MultiCharTest(t, strings.Repeat(string(rune(32)), 10), 10)
		MultiCharTest(t, strings.Repeat(string(rune(32)), 1000), 255)
	})
	t.Run("Repeated SpacialChar", func(t *testing.T) {
		MultiCharTest(t, strings.Repeat(string(rune(33)), 10), 10)
		MultiCharTest(t, strings.Repeat(string(rune(33)), 1000), 255)
	})
	t.Run("Repeated SpacialChar", func(t *testing.T) {
		MultiCharTest(t, strings.Repeat(string(rune(33)), 10), 10)
		MultiCharTest(t, strings.Repeat(string(rune(33)), 1000), 255)
	})

	var unicodeFlags [1]uint32 = [1]uint32{1}
	t.Run("Repeated 10 Unicode", func(t *testing.T) {
		str := strings.Repeat(string(rune(128)), 10)
		svp := new(SimpleValProfile)
		svp.Profile(str)
		ConfirmMulti(t, svp, 0, 0, 0, 0, 0, 0, 1, unicodeFlags[:], 10)
	})

	t.Run("Repeated 1000 Unicode", func(t *testing.T) {
		str := strings.Repeat(string(rune(128)), 1000)
		svp := new(SimpleValProfile)
		svp.Profile(str)
		ConfirmMulti(t, svp, 0, 0, 0, 0, 0, 0, 1, unicodeFlags[:], 255)
	})

}
