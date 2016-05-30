package oled96x96_test

import "github.com/goiot/drivers/i2c/grove/oled96x96"

func ExampleHelloWorld() {
	bus := &i2c.Devfs{
		// change the following value if you use another bus
		Dev:  "/dev/i2c-1",
		Addr: oled96x96.Address,
	}

	display, err := oled96x96.New(bus)
	if err != nil {
		panic(err)
	}
	defer display.Close()

	display.Write("Hello World!")
}
