package accel3xdigital_test

import (
	"fmt"
	"time"

	"github.com/goiot/devices/grove/accel3xdigital"
	"golang.org/x/exp/io/i2c"
)

func Example() {
	accel, err := accel3xdigital.Open(&i2c.Devfs{
		Dev:  "/dev/i2c-1",
		Addr: accel3xdigital.Address,
	})
	if err != nil {
		panic(err)
	}

	defer accel.Close()

	for i := 0; i < 20; i++ {
		accel.Update()
		fmt.Println(accel.State)
		time.Sleep(500 * time.Millisecond)
	}
}
