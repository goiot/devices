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

type Device struct {
	Device *spi.Device

	vals []RGBA

	beginFrame []byte
	endFrame   []byte
}

func New(o driver.Opener, n int) (*Device, error) {
	dev, err := spi.Open(o)
	if err != nil {
		return nil, err
	}
	endFrame := make([]byte, 4*n/2)
	for i := 0; i < len(endFrame); i++ {
		endFrame[i] = 0xff
	}
	return &Device{
		Device:     dev,
		vals:       make([]RGBA, n),
		beginFrame: []byte{0x00, 0x00, 0x00, 0x00},
		endFrame:   endFrame,
	}, nil
}

func (d *Device) SetRGBA(i int, v RGBA) {
	d.vals[i] = v
}

func (d *Device) Display() error {
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

func (d *Device) Close() error {
	return d.Device.Close()
}
