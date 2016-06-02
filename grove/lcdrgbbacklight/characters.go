package lcdrgbbacklight

// CustomLCDChars is a map of CGRAM characters that can be loaded
// into a LCD screen to display custom characters. Some LCD screens such
// as the Grove screen (jhd1313m1) isn't loaded with latin 1 characters.
// It's up to the developer to load the set up to 8 custom characters and
// update the input text so the character is swapped by a byte reflecting
// the position of the custom character to use.
// See SetCustomChar
var CustomLCDChars = map[string][8]byte{
	"é":       [8]byte{130, 132, 142, 145, 159, 144, 142, 128},
	"è":       [8]byte{136, 132, 142, 145, 159, 144, 142, 128},
	"ê":       [8]byte{132, 138, 142, 145, 159, 144, 142, 128},
	"à":       [8]byte{136, 134, 128, 142, 145, 147, 141, 128},
	"â":       [8]byte{132, 138, 128, 142, 145, 147, 141, 128},
	"á":       [8]byte{2, 4, 14, 1, 15, 17, 15, 0},
	"î":       [8]byte{132, 138, 128, 140, 132, 132, 142, 128},
	"í":       [8]byte{2, 4, 12, 4, 4, 4, 14, 0},
	"û":       [8]byte{132, 138, 128, 145, 145, 147, 141, 128},
	"ù":       [8]byte{136, 134, 128, 145, 145, 147, 141, 128},
	"ñ":       [8]byte{14, 0, 22, 25, 17, 17, 17, 0},
	"ó":       [8]byte{2, 4, 14, 17, 17, 17, 14, 0},
	"heart":   [8]byte{0, 10, 31, 31, 31, 14, 4, 0},
	"smiley":  [8]byte{0, 0, 10, 0, 0, 17, 14, 0},
	"frowney": [8]byte{0, 0, 10, 0, 0, 0, 14, 17},
}
