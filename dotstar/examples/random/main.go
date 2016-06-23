// Package main contains a program that drives a dotstar LED strip.
// This example was written to run on a Raspberry Pi
// Make sure SPI is enabled on in raspi-config
// Connect Data to GPIO10, Clock to GPIO11 (and of course power (5v) and ground)
// Your SPI address might be different, `ls` to /dev to see what addresses are available.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goiot/devices/dotstar"
	"golang.org/x/exp/io/spi"
)

// n is the number of LEDs on the strip.
const n = 8

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func main() {
	d, err := dotstar.Open(&spi.Devfs{Dev: "/dev/spidev0.0", Mode: spi.Mode3}, n)
	if err != nil {
		panic(err)
	}

	// random blinking speed
	speed := time.Duration(random.Intn(700) + 100)
	ticker := time.NewTicker(time.Millisecond * speed)

	// catch signals and terminate the app
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// monitor for signals in the background
	go func() {
		s := <-sigc
		ticker.Stop()
		fmt.Println("\nreceived signal:", s)
		time.Sleep(speed)
		// turn off the LEDs
		for i := 0; i < n; i++ {
			d.SetRGBA(i, dotstar.RGBA{1, 1, 1, 0})
		}
		d.Draw()
		time.Sleep(400 * time.Millisecond)
		d.Close()
		os.Exit(0)
	}()

	for _ = range ticker.C {
		// every x milliseconds, change the colors of each LED
		randLedColors(d)
	}
}

func randLedColors(d *dotstar.LEDs) {
	for i := 0; i < n; i++ {
		d.SetRGBA(i, dotstar.RGBA{
			R: randByte(),
			G: randByte(),
			B: randByte(),
			A: randByte(),
		})
	}

	if err := d.Draw(); err != nil {
		panic(err)
	}
}

func randByte() byte {
	return byte(random.Intn(256))
}
