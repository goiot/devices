package piglow_test

import (
	"time"

	"github.com/goiot/devices/piglow"
	"golang.org/x/exp/io/i2c"
)

func Example() {
	p, _ := piglow.Open(&i2c.Devfs{Dev: "/dev/i2c-1", Addr: piglow.Address})

	p.Setup()

	brightness := 0
	for i := 0; i < 10; i++ {
		brightness ^= 1
		p.SetBrightness(brightness)
		time.Sleep(300 * time.Millisecond)
	}
}
