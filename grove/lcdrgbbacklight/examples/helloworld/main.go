package main

import (
	"time"

	"github.com/goiot/devices/grove/lcdrgbbacklight"
	"golang.org/x/exp/io/i2c"
)

func main() {
	display, err := lcdrgbbacklight.Open(
		&i2c.Devfs{
			Dev:  "/dev/i2c-1",
			Addr: lcdrgbbacklight.LCDAddress,
		},
		&i2c.Devfs{
			Dev:  "/dev/i2c-1",
			Addr: lcdrgbbacklight.RGBAddress,
		},
	)
	if err != nil {
		panic(err)
	}
	defer display.Close()

	display.Write("Hello World!")
	display.SetRGB(0, 255, 0)

	time.Sleep(5 * time.Second)

	display.Clear()
	display.Home()
	display.SetRGB(255, 124, 0)

	display.SetCustomChar(0, lcdrgbbacklight.CustomLCDChars["smiley"])
	display.Write("goodbye\nhave a nice day " + string(byte(0)))

	ticker := time.NewTicker(time.Millisecond * 200)
	go func() {
		for _ = range ticker.C {
			display.Scroll(false)
		}
	}()
	time.Sleep(4 * time.Second)
	ticker.Stop()
}
