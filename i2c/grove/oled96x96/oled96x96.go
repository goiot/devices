package oled96x96

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c/driver"
)

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

var (

	// buffer sent to indicate the following data belongs to a command
	cmdCmd = []byte{0x80}
    dataCmd = []byte{0x40}
)

const (
	VerticalModeFlag   = 01
	HorizontalModeFlag = 02

	// Address is the i2c address of the device
	Address  = 0x3c

	lockUnlockCmd    byte = 0xFD // takes a 2nd arg byte
	startLineCmd     byte = 0xA1 // takes a 2nd arg byte
	displayOffCmd    byte = 0xAE
	displayOnCmd     byte = 0xAF
	displayOffsetCmd byte = 0xA2 // takes a 2nd arg byte
	setColAddrCmd    byte = 0x15 // takes 3 arg bytes

	normalDisplayCmd   byte = 0xA4
	InverseDisplayCmd       = 0xA7
	ActivateScrollCmd       = 0x2F
	DectivateScrollCmd      = 0x2E
	contrastLevelCmd        = 0x81
)

// Oled96x96 represents the Grove Oled 96x96 display.
type Oled96x96 struct {
	Conn  driver.Conn
	grayH byte
	grayL byte
}

// New connects to the passed driver, connects and sets it up.
func New(o driver.Opener) (*Oled96x96, error) {
	// TODO(mattetti): switch to `o.Open(Address)` when the exp/io API updated.
	conn, err := o.Open()
	if err != nil {
		return nil, err
	}

	display := &Oled96x96{Conn: conn}

	// Unlock OLED driver IC MCU interface from entering command. i.e: Accept commands
	// Note: locking/unlocking could be exposed to developers later on if needed.
	if err := display.sendCmd(lockUnlockCmd, 0x12); err != nil {
		return display, fmt.Errorf("unlocking OLED interface failed - %v", err)
	}

	if err := display.Off(); err != nil {
		return display, fmt.Errorf("turning off the OLED disaply failed - %v", err)
	}

	// set multiplex ratio, in out case the display is a 96x96 (0x5f = 96)
	if err := display.sendCmd(0xA8, 0x5F); err != nil {
		return display, fmt.Errorf("setting mux ratio failed - %v", err)
	}

	// set the start line
	if err := display.sendCmd(startLineCmd, 0x00); err != nil {
		return display, fmt.Errorf("setting the start line failed - %v", err)
	}

	// set the display offset
	if err := display.sendCmd(displayOffsetCmd, 0x60); err != nil {
		return display, fmt.Errorf("setting the display offset failed - %v", err)
	}

	// set the remap
	if err := display.sendCmd(0xA0, 0x46); err != nil {
		return display, fmt.Errorf("setting the remap failed - %v", err)
	}

	// set the VDD regulator
	if err := display.sendCmd(0xAB, 0x01); err != nil {
		return display, fmt.Errorf("setting vdd internal failed - %v", err)
	}

	if err := display.ContrastLevel(100); err != nil {
		return display, fmt.Errorf("setting contrast failed - %v", err)
	}

	// set the phase length
	if err := display.sendCmd(0xB1, 0x51); err != nil {
		return display, fmt.Errorf("setting the phase length failed - %v", err)
	}

	// set front Clock Divider / Oscillator Frequency
	if err := display.sendCmd(0xB3, 0x01); err != nil {
		return display, fmt.Errorf("setting the display clock failed - %v", err)
	}

	// set linear gray scale table
	if err := display.sendCmd(0xB9); err != nil {
		return display, fmt.Errorf("setting the gray scale table failed - %v", err)
	}

	// set pre charge voltage
	if err := display.sendCmd(0xBC, 0x08); err != nil {
		return display, fmt.Errorf("setting pre charge voltage failed - %v", err)
	}

	if err := display.sendCmd(0xBE, 0x07); err != nil {
		return display, fmt.Errorf("setting VCOMH failed - %v", err)
	}

	if err := display.sendCmd(0xB6, 0x01); err != nil {
		return display, fmt.Errorf("setting second pre-charge period failed - %v", err)
	}

	if err := display.sendCmd(0xD5, 0x62); err != nil {
		return display, fmt.Errorf("enabling second precharge and internal vsl failed - %v", err)
	}

	if err := display.Normal(); err != nil {
		return display, fmt.Errorf("setting the display to normal failed - %v", err)
	}

	if err := display.DisableScroll(); err != nil {
		return display, fmt.Errorf("disabling scrolling failed - %v", err)
	}

	if err := display.On(); err != nil {
		return display, fmt.Errorf("turning on the display failed - %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Row Address
	// This triple byte command specifies row start address and end address of the display data RAM. This
	// command also sets the row address pointer to row start address. This pointer is used to define the current
	// read/write row address in graphic display data RAM
	// start at 0 end at 95
	display.sendCmd(0x75, 0x00, 0x5f)

	// Column Address
	// This triple byte command specifies column start address and end address of the display data RAM. This
	// command also sets the column address pointer to column start address. This pointer is used to define the
	// current read/write column address in graphic display data RAM. If horizontal address increment mode is
	// enabled by command A0h, after finishing read/write one column data, it is incremented automatically to the
	// next column address. Whenever the column address pointer finishes accessing the end column address, it is
	// reset back to start column address and the row address is incremented to the next row.
	//
	// Start from 8th Column of driver IC. This is 0th Column for OLED
	// End at  (8 + 47)th column. Each Column has 2 pixels(segments)
	//
	display.sendCmd(setColAddrCmd, 0x08, 0x37) // Set Column Address

	// Init gray level for text. Default:Brightest White
	display.grayH = 0xF0
	display.grayL = 0x0F

	return display, nil
}

// Close takes care of cleaning things up.
func (o *Oled96x96) Close() error {
	return o.Conn.Close()
}

// Off turns the OLED panel display OFF
func (o *Oled96x96) Off() error {
	return o.sendCmd(displayOffCmd)
}

// On turns the OLED panel display ON
func (o *Oled96x96) On() error {
	return o.sendCmd(displayOnCmd)
}

// Clear clears the whole screen. Should be used before starting a fresh start or after scroll deactivation.
// This function also sets the cursor to top left corner.
func (o *Oled96x96) Clear() {
    for i:=0;i<48;i++ {
        for j=0;j<96;j++) {
            o.sendData(0x00)
        }
    }
}

// Normal sets the display in mormal mode (colors aren't inversed)
func (o *Oled96x96) Normal() error {
	return o.sendCmd(normalDisplayCmd)
}

// Inverse sets the display to inverse mode (colors are inversed)
func (o *Oled96x96) Inverse() {}

// ContrastLevel sets the contrast ratio of OLED display.
// The level can be any number between 0 - 255.
func (o *Oled96x96) ContrastLevel(level int) error {
	if level < 0 || level > 255 {
		return fmt.Errorf("invalid contrast level: %d, should be between 0-255", level)
	}
	return o.sendCmd(contrastLevelCmd, byte(level))
}

// HorizontalMode configures the display to horizontal addressing mode.
func (o *Oled96x96) HorizontalMode() {}

// VerticalMode configures the display to vertical addressing mode.
// The display must be set to vertical mode before printing text.
func (o *Oled96x96) VerticalMode() {}

// PositionCursor sets the text's position (cursor) to Xth Text Row, Yth Text Column.
// The 96x96 OLED is divided into 12 rows and 12 Columns of text.
// These text row and columns should not be confused with the OLED's Row and Column.
func (o *Oled96x96) PositionCursor(row, col int) error {
	// start at 8
	startCol := 0x08 + byte(col*4)
	if err := o.sendCmd(setColAddrCmd, startCol, 0x37); err != nil {
		return fmt.Errorf("failed to set the column - %v", err)
	}

	// Row Address
	if err := o.sendCmd(0x75, byte(row*8), 0x07+(byte(row*8))); err != nil {
		return fmt.Errorf("failed to set the row - %v", err)
	}

	return nil
}

// Write prints the content of the passed text at the cursor's.
func (o *Oled96x96) Write(txt string) error { return nil }

// DrawBitmap displays a binary bitmap on the OLED matrix.
// The data is provided through a slice holding bitmap.
func (o *Oled96x96) DrawBitmap(bitmap []byte) error { return nil }

// HorizontalScrollProperties defines the scrolling behavior.
// StartRow must be in the 0-127 range
// EndRow must be in the 0-127 range and greater than StartRow
// StartColumn must be between 0 and 63.
// EndColumn must be in the 0 and 63 range and greater than StartColumn
func (o *Oled96x96) HorizontalScrollProperties(
	direction ScrollDirection,
	startRow int,
	endRow int,
	startColumn int,
	endColumn int,
	scrollSpeed ScrollSpeed) error {
	return nil
}

// EnableScroll enables and starts scrolling
func (o *Oled96x96) EnableScroll() error {
	return nil
}

// DisableScroll disables and stops scrolling
func (o *Oled96x96) DisableScroll() error {
	return nil
}

// sendCmd sends the passed data preluded by the command byte
func (o *Oled96x96) sendCmd(buf ...byte) error {
	data := append(cmdCmd, buf...)
	return o.Conn.Write(data)
}

// sendData does what you expect it does and maybe even more
func (o *Oled96x96) sendData(buf ...byte) error {
    data := append(dataCmd, buf...)
    return o.Conn.Write(data)
}