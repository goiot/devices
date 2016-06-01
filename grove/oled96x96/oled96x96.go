// Package oled96x96 implements a driver for the Grove OLED grayscale 96x96 display.
package oled96x96

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

// OLED96x96 represents the Grove Oled 96x96 display.
type OLED96x96 struct {
	Device *i2c.Device
	Font   Font
	grayH  byte
	grayL  byte
}

// New connects to the passed driver, connects and sets it up.
func New(o driver.Opener) (*OLED96x96, error) {
	// TODO(mattetti): switch to `o.Open(Address)` when the exp/io API updated.
	device, err := i2c.Open(o)
	if err != nil {
		return nil, err
	}

	display := &OLED96x96{Device: device, Font: DefaultFont()}

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

	// set the remap to vertical mode
	if err := display.VerticalMode(); err != nil {
		return display, fmt.Errorf("setting vertical mode failed - %v", err)
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
	display.sendCmd(setColAddrCmd, 0x08, 0x37)

	display.Clear()

	// Init gray level for text.
	display.grayH = 0xF0
	display.grayL = 0x0F

	if err := display.Normal(); err != nil {
		return display, fmt.Errorf("setting the display to normal failed - %v", err)
	}

	if err := display.VerticalMode(); err != nil {
		return display, fmt.Errorf("setting vertical mode failed - %v", err)
	}

	if err := display.PositionCursor(0, 0); err != nil {
		return display, fmt.Errorf("resetting cursor's position failed - %v", err)
	}

	return display, nil
}

// Close takes care of cleaning things up.
func (o *OLED96x96) Close() error {
	return o.Device.Close()
}

// Off turns the OLED panel display OFF
func (o *OLED96x96) Off() error {
	return o.sendCmd(displayOffCmd)
}

// On turns the OLED panel display ON
func (o *OLED96x96) On() error {
	return o.sendCmd(displayOnCmd)
}

// Clear clears the whole screen. Should be used before starting a fresh start or after scroll deactivation.
// This function also sets the cursor to top left corner.
func (o *OLED96x96) Clear() error {
	// 48*96 = 4608
	nullData := make([]byte, 4609)
	nullData[0] = dataCmd
	return o.Device.Write(nullData)
}

// Normal sets the display in mormal mode (colors aren't inversed)
func (o *OLED96x96) Normal() error {
	return o.sendCmd(normalDisplayCmd)
}

// Inverse sets the display to inverse mode (colors are inversed)
func (o *OLED96x96) Inverse() error {
	return o.sendCmd(inverseDisplayCmd)
}

// ContrastLevel sets the contrast ratio of OLED display.
// The level can be any number between 0 - 255.
func (o *OLED96x96) ContrastLevel(level int) error {
	if level < 0 || level > 255 {
		return fmt.Errorf("invalid contrast level: %d, should be between 0-255", level)
	}
	return o.sendCmd(contrastLevelCmd, byte(level))
}

// HorizontalMode configures the display to horizontal addressing mode.
func (o *OLED96x96) HorizontalMode() error {
	// horizontal mode
	if err := o.sendCmd(0xA0, 0x42); err != nil {
		return err
	}
	// row address (0 to 95)
	if err := o.sendCmd(0x75, 0x00, 0x5f); err != nil {
		return err
	}
	// col address
	if err := o.sendCmd(0x15, 0x08, 0x37); err != nil {
		return err
	}
	return nil
}

// VerticalMode configures the display to vertical addressing mode.
// The display must be set to vertical mode before printing text.
func (o *OLED96x96) VerticalMode() error {
	return o.sendCmd(0xA0, 0x46)
}

// PositionCursor sets the text's position (cursor) to Xth Text Row, Yth Text Column.
// The 96x96 OLED is divided into 12 rows and 12 Columns of text.
// These text row and columns should not be confused with the OLED's Row and Column.
func (o *OLED96x96) PositionCursor(row, col int) error {
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
func (o *OLED96x96) Write(txt string) error {
	var c, bit1, bit2 byte
	letterLen := len(o.Font)

	// TODO (mattetti): support unicode
	pushChar := func(r rune) error {
		n := int(r)
		var j uint8
		for i := 0; i < 8; i = i + 2 {
			for j = 0; j < 8; j++ {
				c, bit1, bit2 = 0x00, 0x00, 0x00
				// Character is constructed two pixel at a time using vertical mode from the default 8x8 font
				// Guard to prevent using characters not supported in the used font.
				if n <= letterLen {
					bit1 = (o.Font[n][i] >> j) & 0x01
					bit2 = (o.Font[n][i+1] >> j) & 0x01
				}
				// Each bit is changed to a nibble
				if bit1 > 0 {
					c |= o.grayH
				} else {
					c |= 0x00
				}
				if bit2 > 0 {
					c |= o.grayL
				} else {
					c |= 0x00
				}
				if err := o.sendData(c); err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, b := range txt {
		if err := pushChar(b); err != nil {
			return err
		}
	}

	return nil
}

// DrawBitmap displays a binary bitmap on the OLED matrix.
// The data is provided through a slice holding bitmap.
func (o *OLED96x96) DrawBitmap(bitmap []byte) error { return nil }

// HorizontalScrollProperties defines the scrolling behavior.
// StartRow must be in the 0-127 range
// EndRow must be in the 0-127 range and greater than StartRow
// StartColumn must be between 0 and 63.
// EndColumn must be in the 0 and 63 range and greater than StartColumn
func (o *OLED96x96) HorizontalScrollProperties(
	direction ScrollDirection,
	startRow int,
	endRow int,
	startColumn int,
	endColumn int,
	scrollSpeed ScrollSpeed) error {
	return nil
}

// EnableScroll enables and starts scrolling
func (o *OLED96x96) EnableScroll() error {
	return nil
}

// DisableScroll disables and stops scrolling
func (o *OLED96x96) DisableScroll() error {
	return nil
}
