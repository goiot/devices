package oled96x96_test

import (
	"github.com/goiot/devices/grove/oled96x96"
	"golang.org/x/exp/io/i2c"
)

func Example() {
	bus := &i2c.Devfs{
		// change the following value if you use another bus
		Dev:  "/dev/i2c-1",
		Addr: oled96x96.Address,
	}

	display, err := oled96x96.Open(bus)
	if err != nil {
		panic(err)
	}
	defer display.Close()

	display.Write("Hello World!")
}
