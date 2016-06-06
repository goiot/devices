// +build linux

// Package main contains a program that drives a dotstar LED strip.
package main

import (
	"github.com/goiot/devices/dotstar"
	"golang.org/x/exp/io/spi"
)

// n is the number of LEDs on the strip.
const n = 100

func main() {
	d, err := dotstar.Open(&spi.Devfs{Dev: "/dev/spi0.1", Mode: spi.Mode3}, n)
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		d.SetRGBA(i, dotstar.RGBA{255, 0, 0, 16}) // Brightest red
	}

	if err := d.Display(); err != nil {
		panic(err)
	}
}
