// Package dotstar implements a driver for the dotstar LEDs.
package dotstar

import (
	"golang.org/x/exp/io/spi"
	"golang.org/x/exp/io/spi/driver"
)

type RGBA struct {
	R byte
	G byte
	B byte
	A byte
}

// LEDs represent a strip of dotstar LEDs.
type LEDs struct {
	// Device is the underlying SPI bus that is used to
	// communicate the LED strip.
	Device *spi.Device

	vals []RGBA
}

// New creates a new LED strip with n dotstar LEDs. An LED strip
// must be closed if no longer in use.
func New(o driver.Opener, n int) (*LEDs, error) {
	dev, err := spi.Open(o)
	if err != nil {
		return nil, err
	}

	return &LEDs{
		Device: dev,
		vals:   make([]RGBA, n),
	}, nil
}

// SetRGBA sets the ith LED to the given RGBA value.
func (d *LEDs) SetRGBA(i int, v RGBA) {
	d.vals[i] = v
}

// Display displays the RGBA values set on the actual LED strip.
func (d *LEDs) Display() error {
	// TODO(jbd): dotstar allows other RGBA allignments, support those layouts.
	n := len(d.vals)
	tx := make([]byte, 1+(4*n)+(n/2))
	tx[0] = 0x00
	tx[1] = 0x00
	tx[2] = 0x00
	tx[3] = 0x00

	for i, c := range d.vals {
		j := (i + 1) * 4
		tx[j] = c.A
		tx[j+1] = c.B
		tx[j+2] = c.G
		tx[j+3] = c.R
	}

	// end frame with at least n/2 0xff vals
	for i := 4*n + 1; i < len(tx); i++ {
		tx[i] = 0xff
	}

	return d.Device.Tx(tx, nil)
}

// Close frees the underlying resources. It must be called once
// the LED strip is no longer in use.
func (d *LEDs) Close() error {
	return d.Device.Close()
}
