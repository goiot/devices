// Package piglow implements a driver for the Pimoroni PiGlow.
package piglow

import (
	"fmt"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

const addr = 0x54 // address is the I2C address of the device.

// PiGlow represents a PiGlow device
type PiGlow struct {
	conn *i2c.Device
}

// Reset resets the internal registers
func (p *PiGlow) Reset() error {
	return p.conn.Write([]byte{0x17, 0xFF})
}

// Shutdown sets the software shutdown mode of the PiGlow
func (p *PiGlow) Shutdown() error {
	return p.conn.Write([]byte{0x00, 0x00})
}

// Enable enables the PiGlow for normal operations
func (p *PiGlow) Enable() error {
	return p.conn.Write([]byte{0x00, 0x01})
}

// Setup enables normal operations, resets the internal registers, and enables
// all LED control registers
func (p *PiGlow) Setup() error {
	if err := p.Reset(); err != nil {
		return err
	}
	if err := p.Enable(); err != nil {
		return err
	}
	if err := p.SetLEDControlRegister(1, 0xFF); err != nil {
		return err
	}
	if err := p.SetLEDControlRegister(2, 0xFF); err != nil {
		return err
	}
	if err := p.SetLEDControlRegister(3, 0xFF); err != nil {
		return err
	}
	return nil
}

// Open opens a new PiGlow. A PiGlow must be closed if no longer in use.
// If the PiGlow has not been powered down since last use, it will be opened
// with it's last programmed state.
func Open(o driver.Opener) (*PiGlow, error) {
	conn, err := i2c.Open(o, addr)
	if err != nil {
		return nil, err
	}
	g := &PiGlow{conn: conn}
	return g, nil
}

// Close frees the underlying resources. It must be called once
// the PiGlow is no longer in use.
func (p *PiGlow) Close() error {
	return p.conn.Close()
}

// Green sets all the green LEDs to the level of 0-255.
func (p *PiGlow) Green(level int) error {
	if err := p.conn.Write([]byte{0x04, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x06, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x0E, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}
	return nil
}

// Blue sets all the blue LEDs to the level of 0-255.
func (p *PiGlow) Blue(level int) error {
	if err := p.conn.Write([]byte{0x05, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x0C, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x0F, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}
	return nil
}

// Yellow sets all the yellow LEDs to the level of 0-255.
func (p *PiGlow) Yellow(level int) error {
	if err := p.conn.Write([]byte{0x03, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x09, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x10, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}
	return nil
}

// Orange sets all the orange LEDs to the level of 0-255.
func (p *PiGlow) Orange(level int) error {
	if err := p.conn.Write([]byte{0x02, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x08, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x11, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}
	return nil
}

// White sets all the white LEDs to the level of 0-255.
func (p *PiGlow) White(level int) error {
	if err := p.conn.Write([]byte{0x0A, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x0B, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x0D, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}
	return nil
}

// Red sets all the red LEDs to the level of 0-255.
func (p *PiGlow) Red(level int) error {
	if err := p.conn.Write([]byte{0x01, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x07, byte(level)}); err != nil {
		return err
	}
	if err := p.conn.Write([]byte{0x12, byte(level)}); err != nil {
		return err
	}
	return p.conn.Write([]byte{0x16, 0xFF})
}

// SetLEDControlRegister sets the control register 1-3 to the bitmask enables.
//   bitmask definition:
//   0 - LED disabled
//   1 - LED enabled
//   LED Control Register 1 - LED channel 1  to 6   bits 0-5
//   LED Control Register 2 - LED channel 7  to 12  bits 0-5
//   LED Control Register 3 - LED channel 13 to 18  bits 0-5
func (p *PiGlow) SetLEDControlRegister(register, enables int) error {
	var address byte

	switch register {
	case 1:
		address = 0x13
	case 2:
		address = 0x14
	case 3:
		address = 0x15
	default:
		return fmt.Errorf("%d is an unknown register", register)
	}

	if err := p.conn.Write([]byte{address, byte(enables)}); err != nil {
		return err
	}
	return p.conn.Write([]byte{0x16, 0xFF})
}

// SetLEDBrightness sets the led 1-18 to the level 0-255.
func (p *PiGlow) SetLEDBrightness(led, level int) error {
	if err := p.conn.Write([]byte{byte(led), byte(level)}); err != nil {
		return err
	}
	return p.conn.Write([]byte{0x16, 0xFF})
}

// SetBrightness sets all the LEDs to the level 0-255.
func (p *PiGlow) SetBrightness(level int) error {
	for i := 1; i <= 18; i++ {
		if err := p.conn.Write([]byte{byte(i), byte(level)}); err != nil {
			return err
		}
	}
	return p.conn.Write([]byte{0x16, 0xFF})
}
