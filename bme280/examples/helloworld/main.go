package main

import (
	"fmt"

	"github.com/goiot/devices/bme280"
	"golang.org/x/exp/io/i2c"
)

func main() {
	sensor, err := bme280.Open(&i2c.Devfs{Dev: "/dev/i2c-2"})
	if err != nil {
		panic(err)
	}

	err = sensor.Update(bme280.Oversample16)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Temperature (celcius) = %f\n", sensor.State.Temperature)
	fmt.Printf("Pressure (hpa) = %f\n", sensor.State.Pressure)
	fmt.Printf("Humidity (%%RH) = %f\n", sensor.State.Humidity)

}
