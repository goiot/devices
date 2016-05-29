package oled96x96

// ScrollDirection is the type determining the scrolling direction of text
type ScrollDirection byte

// ScrollSpeed is the type determining the speed of scrolling
type ScrollSpeed byte

var (
	ScrollLeft  ScrollDirection = 0x00
	ScrollRight ScrollDirection = 0x01

	Scroll2Frames   ScrollSpeed = 0x7
	Scroll3Frames   ScrollSpeed = 0x4
	Scroll4Frames   ScrollSpeed = 0x5
	Scroll5Frames   ScrollSpeed = 0x0
	Scroll25Frames  ScrollSpeed = 0x6
	Scroll64Frames  ScrollSpeed = 0x1
	Scroll128Frames ScrollSpeed = 0x2
	Scroll256Frames ScrollSpeed = 0x3
)

const (
	VerticalModeFlag   = 01
	HorizontalModeFlag = 02

	Address     = 0x3c
	CommandMode = 0x80
	DataMode    = 0x40

	DisplayOffCmd = 0xAE
	DisplayOnCmd  = 0xAF

	NormalDisplayCmd    = 0xA4
	InverseDisplayCmd   = 0xA7
	ActivateScrollCmd   = 0x2F
	DectivateScrollCmd  = 0x2E
	SetContrastLevelCmd = 0x81
)

// Clear clears the whole screen. Should be used before starting a fresh start or after scroll deactivation.
// This function also sets the cursor to top left corner.
func Clear() {}

// NormalDisplay sets the display in mormal mode (colors aren't inversed)
func NormalDisplay() {}

// InverseDisplay sets the display to inverse mode (colors are inversed)
func InverseDisplay() {}

// SetContrastLevel sets the contrast ratio of OLED display.
// The level can be any number between 0 - 255.
func SetContrastLevel(level uint8) {}

// HorizontalMode configures the display to horizontal addressing mode.
func HorizontalMode() {}

// VerticalMode configures the display to vertical addressing mode.
// The display must be set to vertical mode before printing text.
func VerticalMode() {}

// PositionCursor sets the text's position (cursor) to Xth Text Row, Yth Text Column.
// The 96x96 OLED is divided into 12 rows and 12 Columns of text.
// These text row and columns should not be confused with the OLED's Row and Column.
func PositionCursor(row, col int) error { return nil }

// Write prints the content of the passed text at the cursor's.
func Write(txt string) error { return nil }

// DrawBitmap displays a binary bitmap on the OLED matrix.
// The data is provided through a slice holding bitmap.
func DrawBitmap(bitmap []byte) error { return nil }

// HorizontalScrollProperties defines the scrolling behavior.
// StartRow must be in the 0-127 range
// EndRow must be in the 0-127 range and greater than StartRow
// StartColumn must be between 0 and 63.
// EndColumn must be in the 0 and 63 range and greater than StartColumn
func HorizontalScrollProperties(
	direction ScrollDirection,
	startRow int,
	endRow int,
	startColumn int,
	endColumn int,
	scrollSpeed ScrollSpeed) error {
	return nil
}

func EnableScroll() {}

func DisableScroll() {}
