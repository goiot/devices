package main

import (
	"github.com/goiot/devices/ssd1306"
	"golang.org/x/exp/io/i2c"
)

func main() {
	d, err := ssd1306.Open(&i2c.Devfs{Dev: "/dev/i2c-1", Addr: ssd1306.Address})
	if err != nil {
		panic(err)
	}
	defer d.Close()

	// clear the display before putting on anything
	if err := d.Clear(); err != nil {
		panic(err)
	}

	for i := 0; i < d.Width(); i += 2 {
		for j := 0; j < d.Height(); j += 2 {
			d.SetPixel(i, j, 0x1)
		}
	}

	// display the pattern on the display
	if err := d.Draw(); err != nil {
		panic(err)
	}
}
