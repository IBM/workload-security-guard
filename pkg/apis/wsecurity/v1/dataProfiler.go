package v1

import (
	"bytes"
	"fmt"
	"strings"
)

type SimpleValConfig struct {
	Flags        uint32        `json:"flags"`
	NonReadables U8MinmaxSlice `json:"nonreadables"`
	Spaces       U8MinmaxSlice `json:"spaces"`
	Unicodes     U8MinmaxSlice `json:"unicodes"`
	Digits       U8MinmaxSlice `json:"digits"`
	Letters      U8MinmaxSlice `json:"letters"`
	SpecialChars U8MinmaxSlice `json:"schars"`
	Sequences    U8MinmaxSlice `json:"sequences"`
	//Words        U8MinmaxSlice `json:"words"`
	//Numbers      U8MinmaxSlice `json:"numbers"`
	UnicodeFlags Uint32Slice `json:"unicodeFlags"` //[]uint32
	Mandatory    bool        `json:"mandatory"`
}

type SimpleValPile struct {
	Flags        uint32
	NonReadables []uint8
	Spaces       []uint8
	Unicodes     []uint8
	Digits       []uint8
	Letters      []uint8
	SpecialChars []uint8
	//Words        []uint8
	//Numbers      []uint8
	Sequences    []uint8
	UnicodeFlags Uint32Slice //[]uint32
}

type SimpleValProfile struct {
	Flags        uint32
	NonReadables uint8
	Spaces       uint8
	Unicodes     uint8
	Digits       uint8
	Letters      uint8
	SpecialChars uint8
	//Words        uint8
	//Numbers      uint8
	Sequences    uint8
	UnicodeFlags Uint32Slice //[]uint32
}

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
const ( // Slots for Ascii 0-127
	SpaceSlot           = iota // 32
	ExclamationSlot            // 33
	DoubleQouteSlot            // 34
	NumberSlot                 // 35
	DollarSlot                 // 36
	PrecentSlot                // 37
	AmpersandSlot              // 38
	SingleQouteSlot            // 39
	RoundBrecketSlot           // 40, 41
	AsteriskSlot               // 42
	PlusSlot                   // 43 (10)
	CommaSlot                  // 44
	MinusSlot                  // 45
	PeriodSlot                 // 46
	SlashSlot                  // 47
	ColonSlot                  // 58 (15)
	SemiSlot                   // 59
	LtGtSlot                   // 60, 62
	EqualSlot                  //61
	QuestionSlot               // 63
	AtSlot                     // 64 (20)
	BackslashSlot              // 92 (21)
	SquareBrecketSlot          // 91, 93
	PowerSlot                  // 94
	UnderscoreSlot             // 95
	AccentSlot                 // 96
	PipeSlot                   // 124 (26)
	CurlyBrecketSlot           // 123, 125
	HomeSlot                   // 126
	NonReadableCharSlot        // 0-31, 127 (29)
	CommentsSlot
	HexSlot // (31)
)

/*
const ( // Slots for any code
	LASTSLOT__
)
*/
/*
const (
	TotalCounter = iota
	LetterCounter
	DigitCounter
	SpecialCharCounter
	WordCounter
	NumberCounter
	SpareCounter1__
	SpareCounter2__
)

var CounterName = map[int]string{
	TotalCounter:       "TotalCounter",
	LetterCounter:      "LetterCounter",
	DigitCounter:       "DigitCounter",
	SpecialCharCounter: "SpecialCharCounter",
	WordCounter:        "WordCounter",
	NumberCounter:      "NumberCounter",
	SpareCounter1__:    "<UnusedCounter>",
	SpareCounter2__:    "<UnusedCounter>",
}
*/
const ( // sequence types
	seqNone = iota
	seqLetter
	seqDigit
	seqUnicode
	seqSpace
	seqSpecialChar
	seqNonReadable
)

var FlagName = map[int]string{
	SpaceSlot:           "Space",
	ExclamationSlot:     "Exclamation",
	DoubleQouteSlot:     "DoubleQoute",
	NumberSlot:          "NumberSign",
	DollarSlot:          "DollarSign",
	PrecentSlot:         "PrecentSign",
	SingleQouteSlot:     "SingleQoute",
	RoundBrecketSlot:    "RoundBrecket",
	AsteriskSlot:        "MultiplySign",
	PlusSlot:            "PlusSign",
	AtSlot:              "CommentSign",
	MinusSlot:           "MinusSign",
	PeriodSlot:          "DotSign",
	SlashSlot:           "DivideSign",
	ColonSlot:           "ColonSign",
	SemiSlot:            "SemicolonSign",
	LtGtSlot:            "Less/GreaterThanSign",
	EqualSlot:           "EqualSign",
	QuestionSlot:        "QuestionMark",
	CommaSlot:           "CommaSign",
	SquareBrecketSlot:   "SquareBrecket",
	BackslashSlot:       "ReverseDivideSign",
	PowerSlot:           "PowerSign",
	UnderscoreSlot:      "UnderscoreSign",
	AccentSlot:          "AccentSign",
	CurlyBrecketSlot:    "CurlyBrecket",
	PipeSlot:            "PipeSign",
	NonReadableCharSlot: "NonReadableChar",
	CommentsSlot:        "CommentsCombination",
	HexSlot:             "HexCombination",
}

func SetFlags(slots []int) (f uint32) {
	for _, slot := range slots {
		f = f | (0x1 << slot)
	}
	return
}
func NameFlags(f uint32) string {
	var ret bytes.Buffer
	mask := uint32(0x1)
	for i := 0; i < 32; i++ {
		if (f & mask) != 0 {
			ret.WriteString(FlagName[i])
			ret.WriteString(" ")
			f = f ^ mask
		}
		mask = mask << 1
	}
	if f != 0 {
		ret.WriteString("<UnnamedFlags>")
	}
	return ret.String()
}

func NewSimpleValConfig(spaces, unicodes, nonreadables, letters, digits, specialChars, sequences uint8) *SimpleValConfig {
	svc := new(SimpleValConfig)
	svc.Spaces = make([]U8Minmax, 0, 1)
	svc.Unicodes = make([]U8Minmax, 0, 1)
	svc.NonReadables = make([]U8Minmax, 0, 1)
	svc.Letters = make([]U8Minmax, 0, 1)
	svc.Digits = make([]U8Minmax, 0, 1)
	svc.SpecialChars = make([]U8Minmax, 0, 1)
	svc.Sequences = make([]U8Minmax, 0, 1)
	//svc.Words = make([]U8Minmax, 0, 1)
	//svc.Numbers = make([]U8Minmax, 0, 1)

	svc.Spaces[0].Max = spaces
	svc.NonReadables[0].Max = nonreadables
	svc.Unicodes[0].Max = unicodes
	svc.Letters[0].Max = letters
	svc.Digits[0].Max = digits
	svc.SpecialChars[0].Max = specialChars
	svc.Sequences[0].Max = sequences
	//svc.Words[0].Max = words
	//svc.Numbers[0].Max = numbers
	return svc
}

func (svp *SimpleValProfile) NameFlags() string {
	return NameFlags(svp.Flags)
}

/*
func convert64To32(v uint64) (vL uint32, vH uint32) {
	vL = uint32(v)
	vH = uint32(v >> 32)
	return
}

func convert32To64(vL uint32, vH uint32) (v uint64) {
	v = (uint64(vH) << 32) & uint64(vL)
	return
}
*/
func (config *SimpleValConfig) Normalize() {
	config.Digits = append(config.Digits, U8Minmax{0, 0})
	config.Spaces = append(config.Spaces, U8Minmax{0, 0})
	config.Unicodes = append(config.Unicodes, U8Minmax{0, 0})
	config.NonReadables = append(config.NonReadables, U8Minmax{0, 0})
	config.Letters = append(config.Letters, U8Minmax{0, 0})
	config.SpecialChars = append(config.SpecialChars, U8Minmax{0, 0})
	config.Sequences = append(config.Sequences, U8Minmax{0, 0})
	//config.Words = append(config.Words, U8Minmax{0, 0})
	//config.Numbers = append(config.Numbers, U8Minmax{0, 0})
}

func (config *SimpleValConfig) Decide(svp *SimpleValProfile) string {
	//flagsL, flagsH := convert64To32(svp.Flags)
	//flagsL = flagsL & ^config.FlagsL
	//flagsH = flagsH & ^config.FlagsH
	//flags := convert32To64(flagsL, flagsH)
	flags := svp.Flags & ^config.Flags
	//if (flagsL != 0) || (flagsH != 0) {
	//return fmt.Sprintf("Unexpected FlagsL %s (%x) in Value", NameFlags(flags), flags)
	//}
	if flags != 0 {
		return fmt.Sprintf("Unexpected Flags %s (%x) in Value", NameFlags(flags), flags)
	}
	if ret := config.UnicodeFlags.Decide(svp.UnicodeFlags); ret != "" {
		return ret
	}
	if ret := config.Spaces.Decide(svp.Spaces); ret != "" {
		return fmt.Sprintf("Spaces: %s", ret)
	}
	if ret := config.Unicodes.Decide(svp.Unicodes); ret != "" {
		return fmt.Sprintf("Unicodes: %s", ret)
	}
	if ret := config.NonReadables.Decide(svp.NonReadables); ret != "" {
		return fmt.Sprintf("NonReadables: %s", ret)
	}
	if ret := config.Digits.Decide(svp.Digits); ret != "" {
		return fmt.Sprintf("Digits: %s", ret)
	}
	if ret := config.Letters.Decide(svp.Letters); ret != "" {
		return fmt.Sprintf("Letters: %s", ret)
	}
	if ret := config.SpecialChars.Decide(svp.SpecialChars); ret != "" {
		return fmt.Sprintf("SpecialChars: %s", ret)
	}
	if ret := config.Sequences.Decide(svp.Sequences); ret != "" {
		return fmt.Sprintf("Sequences: %s", ret)
	}
	//if ret := config.Words.Decide(svp.Words); ret != "" {
	//	return fmt.Sprintf("Words: %s", ret)
	//}
	//if ret := config.Numbers.Decide(svp.Numbers); ret != "" {
	//	return fmt.Sprintf("Numbers: %s", ret)
	//}
	return ""
}

func (p *SimpleValPile) Add(svp *SimpleValProfile) {
	p.Digits = append(p.Digits, svp.Digits)
	p.Letters = append(p.Letters, svp.Letters)
	p.Sequences = append(p.Sequences, svp.Sequences)
	//p.Numbers = append(p.Numbers, svp.Numbers)
	//p.Words = append(p.Words, svp.Words)
	p.Spaces = append(p.Spaces, svp.Spaces)
	p.Unicodes = append(p.Unicodes, svp.Unicodes)
	p.NonReadables = append(p.NonReadables, svp.NonReadables)
	p.SpecialChars = append(p.SpecialChars, svp.SpecialChars)
	p.Flags |= svp.Flags
	for i := 0; i < len(p.UnicodeFlags); i++ {
		p.UnicodeFlags[i] |= svp.UnicodeFlags[i]
	}
	for i := len(p.UnicodeFlags); i < len(svp.UnicodeFlags); i++ {
		p.UnicodeFlags = append(p.UnicodeFlags, svp.UnicodeFlags[i])
	}
}

// Profile generic value where we expect:
// some short combination of chars
// mainly english letters and/or digits (ascii)
// potentially some small content of special chars
// typically no unicode
func (svp *SimpleValProfile) Profile(str string) {
	var flags uint32
	unicodeFlags := []uint32{}
	digitCounter := uint(0)
	letterCounter := uint(0)
	specialCharCounter := uint(0)
	//wordCounter := uint(0)
	sequenceCounter := uint(0)
	//numberCounter := uint(0)
	nonReadableCounter := uint(0)
	spaceCounter := uint(0)
	totalCounter := uint(0)
	unicodeCounter := uint(0)
	var zero, asterisk, slash, minus bool
	//digits := 0
	//letters := 0
	seqType := seqNone
	seqPrevType := seqNone
	for _, c := range str {
		//letter := false
		//digit := false
		totalCounter++
		if c < 'a' { //0-96
			if c < 'A' { // 0-64
				if c < '0' { //0-47
					if c > 32 { //33-47
						seqType = seqSpecialChar
						specialCharCounter++
						slot := uint(c - 32)
						if c >= 41 { //41-47
							slot = slot - 1
						}
						flags |= 0x1 << slot
						if c == '/' {
							if asterisk {
								flags |= 1 << CommentsSlot
							}
						}
						if slash && c == '*' {
							flags |= 1 << CommentsSlot
						}
						if minus && c == '-' {
							flags |= 1 << CommentsSlot
						}
					} else if c < 32 { //0-31
						seqType = seqNonReadable
						nonReadableCounter++
						flags |= 1 << NonReadableCharSlot
					} else { //32 space
						seqType = seqSpace
						spaceCounter++
						flags |= 0x1
					}
				} else if c <= '9' { //48-57  012..9
					seqType = seqDigit
					digitCounter++
					//digit = true
					//digits++
				} else { //58-64
					seqType = seqSpecialChar
					specialCharCounter++
					slot := uint(c - 58 + 15)
					if c > 61 { // 63-64
						slot = slot - 1
						if c == 62 { //62
							slot = slot - 1
						}
					}
					flags |= 0x1 << slot
				}
			} else if c <= 'Z' { //65-90    ABC..Z
				seqType = seqLetter
				letterCounter++
				if zero && c == 'X' {
					flags |= 0x1 << HexSlot
				}
				//letter = true
				//letters++
			} else { //91-96
				seqType = seqSpecialChar
				specialCharCounter++
				slot := uint(c - 91 + 20)
				if c == 91 {
					slot = slot + 2
				}
				flags |= 0x1 << slot
			}
		} else if c <= 'z' { //97-122   abc..z
			seqType = seqLetter
			letterCounter++
			if zero && c == 'x' {
				flags |= 0x1 << HexSlot
			}
			//letter = true
			//letters++
		} else if c < 127 { //123-126
			seqType = seqSpecialChar
			specialCharCounter++
			slot := uint(c - 123 + 25)
			if c == 123 {
				slot = slot + 2
			}
			flags |= 0x1 << slot
		} else if c < 128 { //127
			seqType = seqNonReadable
			nonReadableCounter++
			flags |= 0x1 << NonReadableCharSlot
		} else {
			// Unicode -  128 and onwards

			// Next we use a rought but quick way to profile unicodes using blocks of 128 codes
			// Block 0 is 128-255, block 1 is 256-383...
			// BlockBit represent the bit in a blockElement. Each blockElement carry 64 bits
			seqType = seqUnicode
			unicodeCounter++
			block := (c / 0x80) - 1
			blockBit := int(block & 0x1F)
			blockElement := int(block / 0x20)
			if blockElement >= len(unicodeFlags) {
				// Dynamically allocate as many blockElements as needed for this profile
				unicodeFlags = append(unicodeFlags, make([]uint32, blockElement-len(unicodeFlags)+1)...)
			}
			unicodeFlags[blockElement] |= 0x1 << blockBit
		}
		zero = (c == '0')
		asterisk = (c == '*')
		slash = (c == '/')
		minus = (c == '-')

		if seqType != seqPrevType {
			sequenceCounter++
			seqPrevType = seqType
		}

		//if letters > 0 && !letter {
		//	wordCounter++
		//	letters = 0
		//}
		//if digits > 0 && !digit {
		//	numberCounter++
		///	digits = 0
		//}
	}
	//if letters > 0 {
	//	wordCounter++
	//}
	//if digits > 0 {
	//	numberCounter++
	//}
	if totalCounter > 0xFF {
		totalCounter = 0xFF
		if digitCounter > 0xFF {
			digitCounter = 0xFF
		}
		if letterCounter > 0xFF {
			letterCounter = 0xFF
		}
		if specialCharCounter > 0xFF {
			specialCharCounter = 0xFF
		}
		if unicodeCounter > 0xFF {
			unicodeCounter = 0xFF
		}
		if spaceCounter > 0xFF {
			spaceCounter = 0xFF
		}
		if nonReadableCounter > 0xFF {
			nonReadableCounter = 0xFF
		}
		if sequenceCounter > 0xFF {
			sequenceCounter = 0xFF
		}
		//if numberCounter > 0xFF {
		//	numberCounter = 0xFF
		//}
		//if wordCounter > 0xFF {
		//	wordCounter = 0xFF
		//}
	}

	svp.Spaces = uint8(spaceCounter)
	svp.Unicodes = uint8(unicodeCounter)
	svp.NonReadables = uint8(nonReadableCounter)
	svp.Digits = uint8(digitCounter)
	svp.Letters = uint8(letterCounter)
	svp.SpecialChars = uint8(specialCharCounter)
	svp.Sequences = uint8(sequenceCounter)
	//svp.Words = uint8(wordCounter)
	//svp.Numbers = uint8(numberCounter)

	svp.Flags = flags
	if len(unicodeFlags) > 0 {
		svp.UnicodeFlags = unicodeFlags
	}
	//fmt.Println(svp.Describe())

}

func (svp *SimpleValProfile) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Flags: 0x%x,\n", svp.Flags))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  UnicodeFlags: %s,\n", svp.UnicodeFlags.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Spaces: %d,\n", svp.Spaces))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Unicodes: %d,\n", svp.Unicodes))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  NonReadables: %d,\n", svp.NonReadables))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Letters: %d,\n", svp.Letters))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Digits: %d,\n", svp.Digits))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  SpecialChars: %d,\n", svp.SpecialChars))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Sequences: %d,\n", svp.Sequences))
	description.WriteString(shift)
	//description.WriteString(fmt.Sprintf("  Words: %d,\n", svp.Words))
	//description.WriteString(shift)
	//description.WriteString(fmt.Sprintf("  Numbers: %d,\n", svp.Numbers))
	//description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (svp *SimpleValProfile) Describe() string {
	var description bytes.Buffer
	description.WriteString("Flags: ")
	description.WriteString(svp.NameFlags())
	description.WriteString(svp.UnicodeFlags.Describe())
	description.WriteString(fmt.Sprintf("Spaces: %d", svp.Spaces))
	description.WriteString(fmt.Sprintf("Unicodes: %d", svp.Unicodes))
	description.WriteString(fmt.Sprintf("NonReadables: %d", svp.NonReadables))
	description.WriteString(fmt.Sprintf("Letters: %d", svp.Letters))
	description.WriteString(fmt.Sprintf("Digits: %d", svp.Digits))
	description.WriteString(fmt.Sprintf("SpecialChars: %d", svp.SpecialChars))
	description.WriteString(fmt.Sprintf("Sequences: %d", svp.Sequences))
	//description.WriteString(fmt.Sprintf("Words: %d", svp.Words))
	//description.WriteString(fmt.Sprintf("Numbers: %d", svp.Numbers))

	return description.String()
}

// Allow generic value based on example (whitelisting)
// Call multiple times top present multiple examples
func (config *SimpleValConfig) AddValExample(str string) {
	svp := new(SimpleValProfile)
	svp.Profile(str)
	//flagsL, flagsH := convert64To32(svp.Flags)
	//config.FlagsL |= flagsL
	//config.FlagsH |= flagsH
	config.Flags |= svp.Flags
	config.Spaces = config.Spaces.AddValExample(svp.Spaces)
	config.Unicodes = config.Unicodes.AddValExample(svp.Unicodes)
	config.NonReadables = config.NonReadables.AddValExample(svp.NonReadables)
	config.Digits = config.Digits.AddValExample(svp.Digits)
	config.Letters = config.Letters.AddValExample(svp.Letters)
	config.SpecialChars = config.SpecialChars.AddValExample(svp.SpecialChars)
	config.Sequences = config.Sequences.AddValExample(svp.Sequences)
	//config.Words = config.Words.AddValExample(svp.Words)
	//config.Numbers = config.Numbers.AddValExample(svp.Numbers)
	config.UnicodeFlags = config.UnicodeFlags.Add(svp.UnicodeFlags)
}

func (config *SimpleValConfig) NameFlags() string {
	//flags := convert32To64(config.FlagsL, config.FlagsH)
	flags := config.Flags
	return NameFlags(flags)
}

func (config *SimpleValConfig) Describe() string {
	var description bytes.Buffer
	description.WriteString("Flags: ")
	description.WriteString(config.NameFlags())

	description.WriteString(config.UnicodeFlags.Describe())

	description.WriteString("Spaces: ")
	description.WriteString(config.Spaces.Describe())
	description.WriteString("NonReadables: ")
	description.WriteString(config.NonReadables.Describe())
	description.WriteString("Unicodes: ")
	description.WriteString(config.Unicodes.Describe())
	description.WriteString("Letters: ")
	description.WriteString(config.Letters.Describe())
	description.WriteString("Digits: ")
	description.WriteString(config.Digits.Describe())
	description.WriteString("SpecialChars: ")
	description.WriteString(config.SpecialChars.Describe())
	description.WriteString("Sequences: ")
	description.WriteString(config.Sequences.Describe())
	//description.WriteString("Words: ")
	//description.WriteString(config.Words.Describe())
	//description.WriteString("Numbers: ")
	//description.WriteString(config.Numbers.Describe())
	return description.String()
}

func (config *SimpleValConfig) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	//description.WriteString(fmt.Sprintf("  FlagsL: 0x%x,\n", config.FlagsL))
	//description.WriteString(shift)
	//description.WriteString(fmt.Sprintf("  FlagsH: 0x%x,\n", config.FlagsH))
	//description.WriteString(shift)
	if config.Mandatory {
		description.WriteString("  Mandatory!\n")
		description.WriteString(shift)
	}
	description.WriteString(fmt.Sprintf("  Flags: 0x%x,\n", config.Flags))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  UnicodeFlags: %s,\n", config.UnicodeFlags.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Spaces: %s,\n", config.Spaces.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Unicodes: %s,\n", config.Unicodes.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  NonReadables: %s,\n", config.NonReadables.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Letters: %s,\n", config.Letters.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Digits: %s,\n", config.Digits.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  SpecialChars: %s,\n", config.SpecialChars.Marshal()))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Sequences: %s,\n", config.Sequences.Marshal()))
	description.WriteString(shift)
	//description.WriteString(fmt.Sprintf("  Words: %s,\n", config.Words.Marshal()))
	//description.WriteString(shift)
	//description.WriteString(fmt.Sprintf("  Numbers: %s,\n", config.Numbers.Marshal()))
	//description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}
