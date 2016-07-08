package main

import (
	"image"
	"os"

	_ "image/png"

	"github.com/goiot/devices/monochromeoled"
	"golang.org/x/exp/io/i2c"
)

func main() {
	rc, err := os.Open("./golang.png")
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	m, _, err := image.Decode(rc)
	if err != nil {
		panic(err)
	}

	d, err := monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		panic(err)
	}
	defer d.Close()

	// clear the display before putting on anything
	if err := d.Clear(); err != nil {
		panic(err)
	}
	if err := d.SetImage(0, 0, m); err != nil {
		panic(err)
	}
	if err := d.Draw(); err != nil {
		panic(err)
	}
}
