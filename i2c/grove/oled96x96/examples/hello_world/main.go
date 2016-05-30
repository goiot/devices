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
	display.Clear()

	display.Write("this is a test")

	time.Sleep(15 * time.Second)
	display.Off()
}
