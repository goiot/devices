// Package ssd1306 contains an SSD1306 OLED driver.
package ssd1306

import (
	"fmt"
	"image"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

const (
	ssd1306_LCDWIDTH  = 128
	ssd1306_LCDHEIGHT = 64

	Address = 0x3C // Address is the I2C address of the device.

	// On or off registers.
	ssd1306_DISPLAY_ON  = 0xAF
	ssd1306_DISPLAY_OFF = 0xAE

	// Scrolling registers.
	ssd1306_ACTIVATE_SCROLL                      = 0x2F
	ssd1306_DEACTIVATE_SCROLL                    = 0x2E
	ssd1306_SET_VERTICAL_SCROLL_AREA             = 0xA3
	ssd1306_RIGHT_HORIZONTAL_SCROLL              = 0x26
	ssd1306_LEFT_HORIZONTAL_SCROLL               = 0x27
	ssd1306_VERTICAL_AND_RIGHT_HORIZONTAL_SCROLL = 0x29
	ssd1306_VERTICAL_AND_LEFT_HORIZONTAL_SCROLL  = 0x2A
)

// OLED represents an SSD1306 OLED display.
type OLED struct {
	dev *i2c.Device

	w   int    // width of the display
	h   int    // height of the display
	buf []byte // each pixel is represented by a bit
}

var initSeq = []byte{
	0xae,
	0x00,
	0x10,
	0x40,
	0x81,
	0xcf,
	0xa1,
	0xc8,
	0xa6,
	0xa4,
	0xa8,
	0x3f,
	0xd3,
	0x00,
	0xd5,
	0x80,
	0xd9,
	0xf1,
	0xda,
	0x12,
	0xdb,
	0x40,
	0x20,
	0x00,
	0x8d,
	0x14,
	0xa5,
	0xaf,
}

// Open opens an SSD1306 OLED display. Once not in use, it needs to
// be close by calling Close.
// The default width is 128, height is 64 if zero values are given.
func Open(o driver.Opener, width, height int) (*OLED, error) {
	if width == 0 {
		width = ssd1306_LCDWIDTH
	}
	if height == 0 {
		height = ssd1306_LCDHEIGHT
	}
	dev, err := i2c.Open(o)
	if err != nil {
		return nil, err
	}
	if err := dev.Write(initSeq); err != nil {
		dev.Close()
		return nil, err
	}
	buf := make([]byte, width*(height/8)+1)
	buf[0] = 0x40 // start frame of pixel data
	return &OLED{dev: dev, w: width, h: height, buf: buf}, nil
}

// On turns on the display if it is off.
func (o *OLED) On() error {
	return o.dev.Write([]byte{ssd1306_DISPLAY_ON})
}

// Off turns off the display if it is on.
func (o *OLED) Off() error {
	return o.dev.Write([]byte{ssd1306_DISPLAY_OFF})
}

// Clear clears the entire display.
func (o *OLED) Clear() error {
	for i := 0; i < len(o.buf); i++ {
		o.buf[0] = 0
	}
	return o.drawBuf(0, o.buf)
}

// DrawPixel draws a pixel on the OLED display.
// v can either be 0 or 1.
func (o *OLED) DrawPixel(x, y int, v byte) error {
	if x >= o.w || y >= o.h {
		return fmt.Errorf("(x=%v, y=%v) is out of bounds on this %vx%v display", x, y, o.w, o.h)
	}
	if v > 1 {
		return fmt.Errorf("value needs to be either 0 or 1; given %v", v)
	}
	i := x + (y/8)*o.w + 1
	o.buf[i] = o.buf[i] | v<<uint((y%8))
	return o.drawBuf(x, o.buf)
}

// DrawImage draws an image on the OLED display starting from x, y.
func (o *OLED) DrawImage(x, y int, img image.Image) error {
	panic("not implemented")
}

func (o *OLED) drawBuf(x int, buf []byte) error {
	if err := o.dev.Write([]byte{0xa4}); err != nil { // the write mode
		return err
	}
	if err := o.dev.Write([]byte{0x40 | byte(x)}); err != nil { //set the start line to x
		return err
	}
	return o.dev.Write(buf)
}

// StartScroll starts scrolling in the horizontal direction starting from
// startY column to endY column.
func (o *OLED) EnableScroll(startY, endY int) error {
	panic("not implemented")
}

// StopStrolls stops the scrolling on the display.
func (o *OLED) DisableScroll() error {
	panic("not implemented")
}

// Width returns the display width.
func (o *OLED) Width() int { return o.w }

// Height returns the display height.
func (o *OLED) Height() int { return o.h }

// Close closes the display.
func (o *OLED) Close() error {
	return o.dev.Close()
}
