package v1

import (
	"fmt"
	"testing"
	"time"
)

const str1 = `What is Lorem Ipsum?
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.

Why do we use it?
`

var flags1 uint32 = uint32(0b100000000010000010100010000001)

func TestDescriptionSimpleVals(t *testing.T) {
	t.Run("LoremIpsum", func(t *testing.T) {
		fmt.Printf("TestDescriptionSimpleVals: %s\n", NameFlags(0b101011111001001101110000001))
	})

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
		svc.AddValExample(str)
		svp.Profile(str)
		svp.Describe()
		if ret := svc.Decide(&svp); ret != "" {
			t.Errorf("ProfileSimpleVal() Decided based on example returned %s", ret)
		}
		svpLonger.Profile(strLonger)
		svpShorter.Profile(strShorter)
		svpSpecial.Profile(strSpecial)
		svpUnicode.Profile(strUnicode)

		fmt.Printf("svc %s\n", svc.Describe())
		fmt.Printf("svpLonger %s\n", svpLonger.Describe())
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
		svc.AddValExample(str)
		svc.AddValExample(strLonger)
		svc.AddValExample(strShorter)
		svc.AddValExample(strSpecial)
		svc.AddValExample(strUnicode)
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

	})
}

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
		if svp.Flags != flags1 {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags1)
		}
		if svp.Spaces != 97 {
			t.Errorf("ProfileSimpleVal() Spaces = %v, want %v", svp.Spaces, 97)
		}
		if svp.NonReadables != 4 {
			t.Errorf("ProfileSimpleVal() NonReadables = %v, want %v", svp.NonReadables, 4)
		}
		if svp.Unicodes != 0 {
			t.Errorf("ProfileSimpleVal() Unicodes = %v, want %v", svp.Unicodes, 0)
		}
		if svp.Letters != 255 {
			t.Errorf("ProfileSimpleVal() Letters = %v, want %v", svp.Letters, 255)
		}
		if svp.Digits != 8 {
			t.Errorf("ProfileSimpleVal() Digits = %v, want %v", svp.Digits, 8)
		}
		if svp.SpecialChars != 11 {
			t.Errorf("ProfileSimpleVal() SpecialChars = %v, want %v", svp.SpecialChars, 11)
		}
		if svp.Sequences != 214 {
			t.Errorf("ProfileSimpleVal() Sequences = %v, want %v", svp.Sequences, 214)
		}
		//if svp.Words != 255 {
		//	t.Errorf("ProfileSimpleVal() Words = %v, want %v", svp.Words, 101)
		//}
		//if svp.Numbers != 255 {
		//	t.Errorf("ProfileSimpleVal() Numbers = %v, want %v", svp.Numbers, 2)
		//}
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
			var flags uint32
			// Slots and counters for AsciiDaya:
			// 0-31 (32) nonReadableRCharCounter
			// 32-47 (16) slots 0-15 respectivly
			// 48-57 (10) digitCounter
			// 58-64 (6) slots 16-22
			// 65-90 (26) smallLetterCounter
			// 91-96 (6) slots 23-28
			// 97-122 (26) capitalLetterCounter
			// 123-126 (4) slots 29-32
			// 127 (1) nonReadableRCharCounter
			// Slots:
			// <SPACE> ! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~
			//    0    1 2 3 4 5 6 7 8 8 9 0 1 2 3 4 5 6 7 8 7 9 0 1 2 1 3 4 5 6 7 6 9 0 1 2
			var notSpecial bool

			if c < 32 { //0-31
				flags |= 1 << NonReadableCharSlot
			} else {
				switch c {
				case 32: // SPACE
					flags |= 1
				case 33: // !
					flags |= 1 << 1
				case 34: // "
					flags |= 1 << 2
				case 35: // #
					flags |= 1 << 3
				case 36: // $
					flags |= 1 << 4
				case 37: // %
					flags |= 1 << 5
				case 38: // &
					flags |= 1 << 6
				case 39: // '
					flags |= 1 << 7
				case 40: // (
					flags |= 1 << 8
				case 41: // )
					flags |= 1 << 8
				case 42: // *
					flags |= 1 << 9
				case 43: // +
					flags |= 1 << 10
				case 44: // ,
					flags |= 1 << 11
				case 45: // -
					flags |= 1 << 12
				case 46: // .
					flags |= 1 << 13
				case 47: // /
					flags |= 1 << 14
				case 58: // :
					flags |= 1 << 15
				case 59: // ;
					flags |= 1 << 16
				case 60: // <
					flags |= 1 << 17
				case 61: // =
					flags |= 1 << 18
				case 62: // >
					flags |= 1 << 17
				case 63: // ?
					flags |= 1 << 19
				case 64: // @
					flags |= 1 << 20
				case 91: // [
					flags |= 1 << 22
				case 92: // \
					flags |= 1 << 21
				case 93: // ]
					flags |= 1 << 22
				case 94: // ^
					flags |= 1 << 23
				case 95: // _
					flags |= 1 << 24
				case 96: // `
					flags |= 1 << 25
				case 122: // {
					flags |= 1 << 27
				case 124: // |
					flags |= 1 << 26
				case 125: // }
					flags |= 1 << 27
				case 126: // ~
					flags |= 1 << 28
				default:
					notSpecial = true
				}
			}
			if !notSpecial {
				specialCharCounter++
			}
			if svp.Flags != flags {
				t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
			}
			if svp.Spaces != 0 {
				t.Errorf("ProfileSimpleVal() Spaces = %d instead of 0", svp.Spaces)
			}
			if svp.Letters != uint8(letterCounter) {
				t.Errorf("ProfileSimpleVal() Letters = %d instead of %d", svp.Letters, letterCounter)
			}
			if svp.Digits != uint8(digitCounter) {
				t.Errorf("ProfileSimpleVal() Digits = %d instead of %d", svp.Digits, digitCounter)
			}
			if svp.SpecialChars != uint8(specialCharCounter) {
				t.Errorf("ProfileSimpleVal() %d SpecialChars = %d instead of %d %v", c, svp.SpecialChars, specialCharCounter, svp)
			}
		})
	}
	var flags uint32
	var str, name string
	flags = 0x1<<CommentsSlot | 0x1<<SlashSlot | 0x1<<AsteriskSlot
	str = "/*"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		var svp SimpleValProfile
		//svc := new(SimpleValConfig)

		svp.Profile(str)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	flags = 0x1<<CommentsSlot | 0x1<<SlashSlot | 0x1<<AsteriskSlot
	str = "*/"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	flags = 0x1 << HexSlot //| 0x1<<DigitSlot | 0x1<<LetterSlot
	str = "0x"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
	})
	flags = 0x1 << HexSlot //| 0x1<<DigitSlot | 0x1<<LetterSlot
	str = "0X"
	name = fmt.Sprintf("testString %s", str)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.UnicodeFlags != nil {
			t.Errorf("ProfileSimpleVal() expected no UnicodeFlags!")
		}
	})
	const str2 = "日本語"
	//flags = 0x1 << UnicodeCharSlot
	unicode := []uint32{0, 0, 0, 9216, 1048576}
	name = fmt.Sprintf("testString %s", str2)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str2)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.Spaces != uint8(len([]rune(str2))) {
			t.Errorf("ProfileSimpleVal() Spaces = %d, want %d", svp.Spaces, uint8(len([]rune(str2))))
		}
		if svp.UnicodeFlags == nil || len(svp.UnicodeFlags) != len(unicode) || svp.UnicodeFlags[len(svp.UnicodeFlags)-1] != unicode[len(unicode)-1] {
			t.Errorf("ProfileSimpleVal() UnicodeFlags = %v", svp.UnicodeFlags)
		}
	})
	const str3 = "{!}"
	flags = SetFlags([]int{ExclamationSlot, CurlyBrecketSlot})
	name = fmt.Sprintf("testString %s", str3)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str3)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %s, want %s", svp.NameFlags(), NameFlags(flags))
		}
		if svp.Spaces != uint8(len([]rune(str3))) {
			t.Errorf("ProfileSimpleVal() Spaces = %d, want %d", svp.Spaces, uint8(len([]rune(str3))))
		}
		if svp.SpecialChars != uint8(len([]rune(str3))) {
			t.Errorf("ProfileSimpleVal() SpecialChars = %d, want %d %v", svp.SpecialChars, uint8(len([]rune(str3))), svp)
		}
		if svp.UnicodeFlags != nil {
			t.Errorf("ProfileSimpleVal() expected no UnicodeFlags!")
		}

	})
	const str4 = "123"
	flags = 0
	name = fmt.Sprintf("testString %s", str4)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str4)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.Spaces != uint8(len([]rune(str4))) {
			t.Errorf("ProfileSimpleVal() Spaces = %d, want %d", svp.Spaces, uint8(len([]rune(str4))))
		}
		if svp.Digits != uint8(len([]rune(str4))) {
			t.Errorf("ProfileSimpleVal() Digits = %d, want %d", svp.Digits, uint8(len([]rune(str4))))
		}
		if svp.UnicodeFlags != nil {
			t.Errorf("ProfileSimpleVal() expected no UnicodeFlags!")
		}

	})
	const str5 = "aBc"
	flags = 0
	name = fmt.Sprintf("testString %s", str5)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str5)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.Spaces != uint8(len([]rune(str5))) {
			t.Errorf("ProfileSimpleVal() Spaces = %d, want %d", svp.Spaces, uint8(len([]rune(str5))))
		}
		if svp.Letters != uint8(len([]rune(str5))) {
			t.Errorf("ProfileSimpleVal() Letters = %d, want %d", svp.Letters, uint8(len([]rune(str5))))
		}
		if svp.UnicodeFlags != nil {
			t.Errorf("ProfileSimpleVal() expected no UnicodeFlags!")
		}

	})
	var str6 string = string([]rune{rune(200), rune(201), rune(202)})
	//flags = 0x1 << UnicodeCharSlot
	unicode = []uint32{1}
	name = fmt.Sprintf("testString %s", str6)
	t.Run(name, func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str6)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.Spaces != uint8(len([]rune(str6))) {
			t.Errorf("ProfileSimpleVal() Spaces = %d, want %d", svp.Spaces, uint8(len([]rune(str6))))
		}
		if svp.UnicodeFlags == nil || len(svp.UnicodeFlags) != 1 || svp.UnicodeFlags[0] != 1 {
			t.Errorf("ProfileSimpleVal() expected UnicodeFlags %v received %v!", unicode, svp.UnicodeFlags)
		}

	})
	byteSlice := make([]rune, 255*2*26)
	for j := 0; j < 26; j++ {
		for i := rune(0); i < 255; i++ {
			byteSlice[int(j)*255*2+2*int(i)] = i
			byteSlice[int(j)*255*2+2*int(i)+1] = 32
		}
	}
	str7 := string(byteSlice)
	flags = 0xFFFFFFFF
	unicode = []uint32{1}
	t.Run("testString allletters", func(t *testing.T) {
		svp := new(SimpleValProfile)
		svp.Profile(str7)
		if svp.Flags != flags {
			t.Errorf("ProfileSimpleVal() Flags = %b, want %b", svp.Flags, flags)
		}
		if svp.Spaces != uint8(len([]rune(str6))) {
			t.Errorf("ProfileSimpleVal() Spaces = %d, want %d", svp.Spaces, uint8(len([]rune(str6))))
		}
		if svp.Sequences != uint8(len([]rune(str6))) {
			t.Errorf("ProfileSimpleVal() Sequences = %d, want %d", svp.Sequences, uint8(len([]rune(str6))))
		}
		//if svp.Words != uint8(len([]rune(str6))) {
		//	t.Errorf("ProfileSimpleVal() Words = %d, want %d", svp.Words, uint8(len([]rune(str6))))
		//}
		//if svp.Digits != uint8(len([]rune(str6))) {
		//	t.Errorf("ProfileSimpleVal() Digits = %d, want %d", svp.Digits, uint8(len([]rune(str6))))
		//}
		if svp.SpecialChars != uint8(len([]rune(str6))) {
			t.Errorf("ProfileSimpleVal() SpecialChars = %d, want %d", svp.SpecialChars, uint8(len([]rune(str6))))
		}
		if svp.Sequences != uint8(len([]rune(str6))) {
			t.Errorf("ProfileSimpleVal() Sequences = %d, want %d", svp.Sequences, uint8(len([]rune(str6))))
		}
		//if svp.Words != uint8(len([]rune(str6))) {
		//	t.Errorf("ProfileSimpleVal() Words = %d, want %d", svp.Words, uint8(len([]rune(str6))))
		//}
		//if svp.Numbers != uint8(len([]rune(str6))) {
		//	t.Errorf("ProfileSimpleVal() Numbers = %d, want %d", svp.Numbers, uint8(len([]rune(str6))))
		//}
		if svp.UnicodeFlags == nil || len(svp.UnicodeFlags) != 1 || svp.UnicodeFlags[0] != 1 {
			t.Errorf("ProfileSimpleVal() expected UnicodeFlags %v received %v!", unicode, svp.UnicodeFlags)
		}

	})

}
