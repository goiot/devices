// Package dotstar implements a driver for the dotstar LEDs.
package dotstar

import (
	"golang.org/x/exp/io/spi"
	"golang.org/x/exp/io/spi/driver"
)

type RGBA struct {
	R int
	G int
	B int
	A int
}

type Device struct {
	Device *spi.Device

	vals []RGBA
}

func New(o driver.Opener, n int) (*Device, error) {
	dev, err := spi.Open(o)
	if err != nil {
		return nil, err
	}
	return &Device{
		Device: dev,
		vals:   make([]RGBA, n),
	}, nil
}

func (d *Device) SetRGBA(i int, v RGBA) {
	d.vals[i] = v
}

func (d *Device) Display() error {
	panic("not yet implemented")
}

func (d *Device) Close() error {
	return d.Device.Close()
}
