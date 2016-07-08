package lcdrgbbacklight

import "errors"

const (
	lcdAddr = 0x3E
	rgbAddr = 0x62
)

const (
	regRed   = 0x04
	regGreen = 0x03
	regBlue  = 0x02

	lcdClearDisplay        = 0x01
	lcdReturnHome          = 0x02
	lcdEntryModeSet        = 0x04
	lcdDisplayControl      = 0x08
	lcdCursorShift         = 0x10
	lcdFunctionSet         = 0x20
	lcdSetCgramAddr        = 0x40
	lcdSetDdramAddr        = 0x80
	lcdEntryRight          = 0x00
	lcdEntryLeft           = 0x02
	lcdEntryShiftIncrement = 0x01
	lcdEntryShiftDecrement = 0x00
	lcdDisplayOn           = 0x04
	lcdDisplayOff          = 0x00
	lcdCursorOn            = 0x02
	lcdCursorOff           = 0x00
	lcdBlinkOn             = 0x01
	lcdBlinkOff            = 0x00
	lcdDisplayMove         = 0x08
	lcdCursorMove          = 0x00
	lcdMoveRight           = 0x04
	lcdMoveLeft            = 0x00
	lcd2Line               = 0x08
	lcdCmd                 = 0x80
	lcdData                = 0x40

	lcd2ndLineOffset = 0x40
)

var (
	ErrInvalidPosition = errors.New("Invalid position value")
)

func (d *LCDRGBBacklight) setReg(cmd byte, data byte) error {
	// TODO(mattetti): reuse a buffer instead of reallocating
	if err := d.RGB.Write([]byte{cmd, data}); err != nil {
		return err
	}
	return nil
}

// TODO(mattetti): change the arg to be variadic
func (d *LCDRGBBacklight) command(buf []byte) error {
	return d.LCD.Write(append([]byte{lcdCmd}, buf...))
}
