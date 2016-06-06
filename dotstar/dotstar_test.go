package dotstar_test

import (
	"github.com/goiot/devices/dotstar"
	"golang.org/x/exp/io/spi"
)

func Example() {
	d, err := dotstar.Open( // a strip with 5 LEDs
		&spi.Devfs{Dev: "/dev/spi0.1", Mode: spi.Mode3}, 5)
	if err != nil {
		panic(err)
	}

	d.SetRGBA(0, dotstar.RGBA{255, 0, 0, 16}) // Red
	d.SetRGBA(1, dotstar.RGBA{0, 255, 0, 16}) // Green
	d.SetRGBA(2, dotstar.RGBA{0, 0, 255, 16}) // Blue
	d.SetRGBA(3, dotstar.RGBA{255, 0, 0, 8})  // Half dim red
	d.SetRGBA(4, dotstar.RGBA{0, 0, 255, 8})  // Half dim blue

	if err := d.Display(); err != nil {
		panic(err)
	}
}
