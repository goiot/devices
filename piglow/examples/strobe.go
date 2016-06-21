package main

import (
	"time"

	"github.com/goiot/devices/piglow"

	"golang.org/x/exp/io/i2c"
)

func main() {
	piglow, err := piglow.Open(
		&i2c.Devfs{
			Dev:  "/dev/i2c-1",
			Addr: 0x54,
		},
	)
	if err != nil {
		panic(err)
	}

	defer piglow.Close()

	if err := piglow.Reset(); err != nil {
		panic(err)
	}

	for {
		time.Sleep(50 * time.Millisecond)

		for i := 1; i <= 18; i++ {
			if err := piglow.LED(i, 1); err != nil {
				panic(err)
			}
			time.Sleep(10 * time.Millisecond)
		}

		time.Sleep(50 * time.Millisecond)

		for i := 18; i > 0; i-- {
			if err := piglow.LED(i, 0); err != nil {
				panic(err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
