package oled96x96

func ExampleHelloWorld() {
	bus := &i2c.Devfs{
		// change the following value if you use another bus
		Dev:  "/dev/i2c-1",
		Addr: Address,
	}

	display, err := New(bus)
	if err != nil {
		panic(err)
	}
	defer display.Close()

	display.Write("Hello World!")
}
