// Package lcdrgbbacklight implements a driver for the Grove LCD RGB Backlight display.
package lcdrgbbacklight

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

// LCDRGBBacklight is a driver for the Jhd1313m1 LCD display which has two i2c addresses,
// one belongs to a controller and the other controls solely the backlight.
type LCDRGBBacklight struct {
	LCD *i2c.Device
	RGB *i2c.Device
}

// Open connects to the lcd and rgb openers, connects and sets up.
func Open(o driver.Opener) (*LCDRGBBacklight, error) {
	lcdD, err := i2c.Open(o, lcdAddr)
	if err != nil {
		return nil, fmt.Errorf("LCD driver failed to connect - %v", err)
	}

	rgbD, err := i2c.Open(o, rgbAddr)
	if err != nil {
		return nil, fmt.Errorf("RGB driver failed to connect - %v", err)
	}

	display := &LCDRGBBacklight{
		LCD: lcdD,
		RGB: rgbD,
	}

	time.Sleep(50 * time.Millisecond)
	if err := lcdD.Write([]byte{lcdCmd, lcdFunctionSet | lcd2Line}); err != nil {
		return nil, fmt.Errorf("LCD failed to initialize (part 1) - %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if err := lcdD.Write([]byte{lcdCmd, lcdDisplayControl | lcdDisplayOn}); err != nil {
		return nil, fmt.Errorf("LCD failed to initialize (part 2) - %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if err := display.Clear(); err != nil {
		return nil, fmt.Errorf("display failed to clear - %v", err)
	}

	if err := lcdD.Write([]byte{lcdCmd, lcdEntryModeSet | lcdEntryLeft | lcdEntryShiftDecrement}); err != nil {
		return nil, fmt.Errorf("failed to initialize (part 3) - %v", err)
	}

	// registry 0
	if err := display.setReg(0x0, 0x1); err != nil {
		return nil, fmt.Errorf("failed to set registry 0 - %v", err)
	}

	// registry 1
	if err := display.setReg(0x1, 0x0); err != nil {
		return nil, fmt.Errorf("failed to set registry 1 - %v", err)
	}

	// registry 8
	if err := display.setReg(0x08, 0xAA); err != nil {
		return nil, fmt.Errorf("failed to set registry 8 - %v", err)
	}

	if err := display.SetRGB(255, 255, 255); err != nil {
		return nil, fmt.Errorf("failed to set registry 8 - %v", err)
	}

	return display, nil
}

// SetRGB sets the Red Green Blue value of backlit.
func (d *LCDRGBBacklight) SetRGB(r, g, b int) error {
	if err := d.setReg(regRed, byte(r)); err != nil {
		return err
	}
	if err := d.setReg(regGreen, byte(g)); err != nil {
		return err
	}
	return d.setReg(regBlue, byte(b))
}

// Clear clears the text on the lCD display.
func (d *LCDRGBBacklight) Clear() error {
	err := d.command([]byte{lcdClearDisplay})
	return err
}

// Home sets the cursor to the origin position on the display.
func (d *LCDRGBBacklight) Home() error {
	err := d.command([]byte{lcdReturnHome})
	// This wait fixes a race condition when calling home and clear back to back.
	time.Sleep(2 * time.Millisecond)
	return err
}

// Write displays the passed message on the screen.
func (d *LCDRGBBacklight) Write(message string) error {
	// This wait fixes an odd bug where the clear function doesn't always work properly.
	time.Sleep(1 * time.Millisecond)
	for _, val := range message {
		if val == '\n' {
			if err := d.SetPosition(16); err != nil {
				return err
			}
			continue
		}
		if err := d.LCD.Write([]byte{lcdData, byte(val)}); err != nil {
			return err
		}
	}
	return nil
}

// SetPosition sets the cursor and the data display to pos.
// 0..15 are the positions in the first display line.
// 16..32 are the positions in the second display line.
func (d *LCDRGBBacklight) SetPosition(pos int) (err error) {
	if pos < 0 || pos > 31 {
		err = ErrInvalidPosition
		return
	}
	offset := byte(pos)
	if pos >= 16 {
		offset -= 16
		offset |= lcd2ndLineOffset
	}
	err = d.command([]byte{lcdSetDdramAddr | offset})
	return
}

// Scroll scrolls the text on the display
func (d *LCDRGBBacklight) Scroll(leftToRight bool) error {
	if leftToRight {
		return d.command([]byte{lcdCursorShift | lcdDisplayMove | lcdMoveLeft})
	}

	return d.command([]byte{lcdCursorShift | lcdDisplayMove | lcdMoveRight})
}

// CustomChar sets one of the 8 CGRAM locations with a custom character.
// The custom character can be used by writing a byte of value 0 to 7.
// When you are using LCD as 5x8 dots in function set then you can define a total of 8 user defined patterns
// (1 Byte for each row and 8 rows for each pattern).
// Use http://www.8051projects.net/lcd-interfacing/lcd-custom-character.php to create your own
// characters.
// To use a custom character, write byte value of the custom character position as a string after
// having setup the custom character.
func (d *LCDRGBBacklight) SetCustomChar(pos int, charMap [8]byte) error {
	if pos > 7 {
		return fmt.Errorf("can't set a custom character at a position greater than 7")
	}
	location := uint8(pos)
	if err := d.command([]byte{lcdSetCgramAddr | (location << 3)}); err != nil {
		return err
	}

	return d.LCD.Write(append([]byte{lcdData}, charMap[:]...))
}

// Close cleans up the connections
func (d *LCDRGBBacklight) Close() error {
	d.Clear()
	d.SetRGB(0, 0, 0)
	if err := d.LCD.Close(); err != nil {
		return err
	}
	if err := d.RGB.Close(); err != nil {
		return err
	}
	return nil
}
