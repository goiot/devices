// Package ssd1306 contains an SSD1306 OLED driver.
package ssd1306

import (
	"image"

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
type OLED struct{}

// Open opens an SSD1306 OLED display. Once not in use, it needs to
// be close by calling Close.
// The default width is 128, height is 64 if zero values are given.
func Open(o driver.Opener, width, height int) (*OLED, error) {
	panic("not implemented")
}

// DrawByte draws a byte on the OLED display.
func (o *OLED) DrawByte(x, y, int, v byte) error {
	panic("not implemented")
}

// DrawImage draws an image on the OLED display starting from x, y.
func (o *OLED) DrawImage(x, y int, img image.Image) error {
	panic("not implemented")
}

// StartScroll starts scrolling in the horizontal direction starting from
// startY column to endY column.
func (o *OLED) StartScroll(startY, endY int) error {
	panic("not implemented")
}

// StopStrolls stops the scrolling on the display.
func (o *OLED) StopScroll() error {
	panic("not implemented")
}

// Close closes the display.
func (o *OLED) Close() error {
	panic("not implemented")
}
