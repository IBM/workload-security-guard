package spec

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const str1 = `What is Lorem Ipsum?
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.

Why do we use it?
It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).

`

var basicCounters1 [8]uint8 = [8]uint8{255, 255, 8, 29, 200, 0, 0, 0}
var flags1 uint64 = uint64(60131668865)

func TestProfileSimpleVals(t *testing.T) {
	var svp *SimpleValProfile
	//svc := new(SimpleValConfig)
	t.Run("LoremIpsum", func(t *testing.T) {
		//fmt.Printf("tt.args.str %s", tt.args.str)
		start := time.Now()

		for i := 0; i < 1; i++ {
			svp = new(SimpleValProfile)
			svp.Profile(str1)
		}
		elapsed := time.Since(start)
		fmt.Printf("Time is %s", elapsed)
		//t.Errorf("svp %v", svp)
		if !reflect.DeepEqual(svp.Flags, flags1) {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags1)
		}

		if !reflect.DeepEqual(svp.BasicCounters, basicCounters1) {
			t.Errorf("ProfileSimpleVal() BasicCounters = %v, want %v", svp.BasicCounters, basicCounters1)
		}

	})
	for c := 0; c < 257; c++ {
		r := rune(c)
		str := string(r)
		name := fmt.Sprintf("testString %s", str)
		t.Run(name, func(t *testing.T) {
			digitCounter := 0
			letterCounter := 0
			specialCharCounter := 0

			svp := new(SimpleValProfile)
			svp.Profile(string(r))
			var flags uint64
			if c < 32 { //0-31
				flags |= 1 << nonReadableCharCounterSlot
			} else if c < 33 { //32 space
				flags |= 1
			} else if c < 48 { //32-47
				slot := uint(c - 32)
				flags |= 1 << slot
				specialCharCounter++
			} else if c < 58 { //48-57
				digitCounter++
				//	flags |= 1 << DigitSlot
			} else if c < 65 { //58-64
				slot := uint(c - 58 + 16)
				flags |= 1 << slot
				specialCharCounter++
			} else if c < 91 { //65-90
				letterCounter++
				//	flags |= 1 << LetterSlot
			} else if c < 97 { //91-96
				slot := uint(c - 91 + 23)
				flags |= 1 << slot
				specialCharCounter++
			} else if c < 123 { //97-122
				letterCounter++
				//	flags |= 1 << LetterSlot
			} else if c < 127 { //123-126
				slot := uint(c - 123 + 29)
				flags |= 1 << slot
				specialCharCounter++
			} else if c < 128 { //127
				flags |= 1 << nonReadableCharCounterSlot
			} else if c < 256 { //128-255
				flags |= 1 << UpperAsciiSlot
			} else { //unicode
				flags |= 1 << UnicodeCharSlot
			}
			if !reflect.DeepEqual(svp.Flags, flags) {
				t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
			}
			if svp.BasicCounters[BasicTotalCounter] != 1 {
				t.Errorf("ProfileSimpleVal() BasicTotalCounter = %d instead of 1", svp.BasicCounters[BasicTotalCounter])
			}
			if svp.BasicCounters[BasicLetterCounter] != uint8(letterCounter) {
				t.Errorf("ProfileSimpleVal() BasicLetterCounter = %d instead of %d", svp.BasicCounters[BasicLetterCounter], letterCounter)
			}
			if svp.BasicCounters[BasicDigitCounter] != uint8(digitCounter) {
				t.Errorf("ProfileSimpleVal() BasicDigitCounter = %d instead of %d", svp.BasicCounters[BasicDigitCounter], digitCounter)
			}
			if svp.BasicCounters[BasicSpecialCounter] != uint8(specialCharCounter) {
				t.Errorf("ProfileSimpleVal() %d BasicSpecialCounter = %d instead of %d %v", c, svp.BasicCounters[BasicSpecialCounter], specialCharCounter, svp)
			}
		})
	}
	var flags uint64
	var str, name string
	flags = 0x1<<SlashAsteriskCommentSlot | 0x1<<DivSlot | 0x1<<MultSlot
	str = "/*"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		var svp SimpleValProfile
		//svc := new(SimpleValConfig)

		svp.Profile(str)
		if !reflect.DeepEqual(svp.Flags, flags) {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	flags = 0x1<<SlashAsteriskCommentSlot | 0x1<<DivSlot | 0x1<<MultSlot
	str = "*/"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		if !reflect.DeepEqual(svp.Flags, flags) {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	flags = 0x1 << HexSlot //| 0x1<<DigitSlot | 0x1<<LetterSlot
	str = "0x"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		if !reflect.DeepEqual(svp.Flags, flags) {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	flags = 0x1 << HexSlot //| 0x1<<DigitSlot | 0x1<<LetterSlot
	str = "0X"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		if !reflect.DeepEqual(svp.Flags, flags) {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	const str2 = "日本語"
	flags = 0x1 << UnicodeCharSlot
	name = fmt.Sprintf("testString %s", str2)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str2)
		if !reflect.DeepEqual(svp.Flags, flags) {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.BasicCounters[BasicTotalCounter] != uint8(len([]rune(str2))) {
			t.Errorf("ProfileSimpleVal() BasicTotalCounter = %d, want %d", svp.BasicCounters[BasicTotalCounter], uint8(len([]rune(str2))))
		}
	})
	const str3 = "{!}"
	flags = 0x1 << UpperAsciiSlot
	name = fmt.Sprintf("testString %s", str3)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str3)
		if svp.BasicCounters[BasicTotalCounter] != uint8(len([]rune(str3))) {
			t.Errorf("ProfileSimpleVal() BasicTotalCounter = %d, want %d", svp.BasicCounters[BasicTotalCounter], uint8(len([]rune(str3))))
		}
		if svp.BasicCounters[BasicSpecialCounter] != uint8(len([]rune(str3))) {
			t.Errorf("ProfileSimpleVal() BasicSpecialCounter = %d, want %d %v", svp.BasicCounters[BasicSpecialCounter], uint8(len([]rune(str3))), svp)
		}
	})
	const str4 = "123"
	flags = 0x1 << UnicodeCharSlot
	name = fmt.Sprintf("testString %s", str4)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str4)
		if svp.BasicCounters[BasicTotalCounter] != uint8(len([]rune(str4))) {
			t.Errorf("ProfileSimpleVal() BasicTotalCounter = %d, want %d", svp.BasicCounters[BasicTotalCounter], uint8(len([]rune(str4))))
		}
		if svp.BasicCounters[BasicDigitCounter] != uint8(len([]rune(str4))) {
			t.Errorf("ProfileSimpleVal() BasicDigitCounter = %d, want %d", svp.BasicCounters[BasicDigitCounter], uint8(len([]rune(str4))))
		}
	})
	const str5 = "aBc"
	flags = 0x1 << UnicodeCharSlot
	name = fmt.Sprintf("testString %s", str5)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str5)
		if svp.BasicCounters[BasicTotalCounter] != uint8(len([]rune(str5))) {
			t.Errorf("ProfileSimpleVal() BasicTotalCounter = %d, want %d", svp.BasicCounters[BasicTotalCounter], uint8(len([]rune(str5))))
		}
		if svp.BasicCounters[BasicLetterCounter] != uint8(len([]rune(str5))) {
			t.Errorf("ProfileSimpleVal() BasicLetterCounter = %d, want %d", svp.BasicCounters[BasicLetterCounter], uint8(len([]rune(str5))))
		}
	})
	var str6 string = string([]rune{rune(200), rune(201), rune(202)})
	flags = 0x1 << UnicodeCharSlot
	name = fmt.Sprintf("testString %s", str6)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str6)
		if svp.BasicCounters[BasicTotalCounter] != uint8(len([]rune(str6))) {
			t.Errorf("ProfileSimpleVal() BasicTotalCounter = %d, want %d", svp.BasicCounters[BasicTotalCounter], uint8(len([]rune(str6))))
		}
	})
}
