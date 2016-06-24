package main

import (
	"time"

	"github.com/goiot/devices/piglow"
	"golang.org/x/exp/io/i2c"
)

func main() {
	p, err := piglow.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		panic(err)
	}

	defer func() {
		p.Shutdown()
		p.Close()
	}()

	time.Sleep(50 * time.Millisecond)
	for i := 1; i <= 18; i++ {
		if err := p.SetLEDBrightness(i, 1); err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	for i := 18; i > 0; i-- {
		if err := p.SetLEDBrightness(i, 0); err != nil {
			panic(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
