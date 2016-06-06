package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goiot/devices/grove/accel3xdigital"
	"golang.org/x/exp/io/i2c"
)

func main() {
	accel, err := accel3xdigital.New(&i2c.Devfs{
		Dev:  "/dev/i2c-1",
		Addr: accel3xdigital.Address,
	})
	if err != nil {
		panic(err)
	}

	// Enabling tap switches from 32 samples/second to 120
	err = accel.EnableTap()
	if err != nil {
		panic(err)
	}

	// channel to push to if we want to exit in a clean way
	quitQ := make(chan bool)

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
		fmt.Println("\nreceived signal:", s)
		quitQ <- true
	}()

	for {
		select {
		case <-quitQ:
			accel.Close()
			fmt.Println("Ciao! :)")
			os.Exit(0)
		case <-time.After(500 * time.Millisecond):
			accel.Update()
			fmt.Println(accel.State)
		}
	}

}
