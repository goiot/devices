package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goiot/devices/piglow"
	"golang.org/x/exp/io/i2c"
)

func main() {
	p, err := piglow.Open(&i2c.Devfs{Dev: "/dev/i2c-1", Addr: piglow.Address})
	if err != nil {
		panic(err)
	}

	// catch signals and terminate the app
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	if err := p.Setup(); err != nil {
		panic(err)
	}

	for {
		select {
		case <-sigc:
			p.Shutdown()
			p.Close()
			return
		default:
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
	}
}
