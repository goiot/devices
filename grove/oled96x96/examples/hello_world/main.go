package main

import (
	"time"

	"github.com/goiot/drivers/i2c/grove/oled96x96"
	"golang.org/x/exp/io/i2c"
)

func main() {
	bus := &i2c.Devfs{
		Dev:  "/dev/i2c-1",
		Addr: oled96x96.Address,
	}

	display, err := oled96x96.New(bus)
	if err != nil {
		panic(err)
	}
	defer display.Close()

	display.Write("Hello World!")

	time.Sleep(5 * time.Second)
	display.PositionCursor(0, 0)
	display.Write("Ciao World!")
	time.Sleep(1 * time.Second)
	display.Inverse()
	time.Sleep(1 * time.Second)
	display.Normal()
	time.Sleep(1 * time.Second)

	display.Off()
}
