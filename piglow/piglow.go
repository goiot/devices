// Package piglow implements a driver for the Pimoroni PiGlow.
package piglow

import (
	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

// PiGlow represents a PiGlow device
type PiGlow struct {
	conn *i2c.Device
}

// Open opens a new PiGlow. A PiGlow must be closed if no longer in use.
// If the PiGlow has not been powered down since last use, it will be opened
// with it's last programmed state.
func Open(o driver.Opener) (*PiGlow, error) {
	conn, err := i2c.Open(o)
	if err != nil {
		return nil, err
	}

	return &PiGlow{
		conn: conn,
	}, nil
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

	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}

	return nil
}

// LED sets the led 1-18 to the level 0-255.
func (p *PiGlow) LED(led, level int) error {
	if err := p.conn.Write([]byte{byte(led), byte(level)}); err != nil {
		return err
	}

	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}

	return nil
}

// All sets all the LEDs to the level 0-255.
func (p *PiGlow) All(level int) error {
	for i := 1; i <= 18; i++ {
		if err := p.conn.Write([]byte{byte(i), byte(level)}); err != nil {
			return err
		}
	}

	if err := p.conn.Write([]byte{0x16, 0xFF}); err != nil {
		return err
	}

	return nil
}

// Reset resets the internal registers of the PiGlow.
func (p *PiGlow) Reset() error {
	if err := p.conn.Write([]byte{0x17, 0xFF}); err != nil {
		return err
	}

	if err := p.conn.Write([]byte{0x00, 0x01}); err != nil {
		return err
	}

	if err := p.conn.Write([]byte{0x13, 0xFF}); err != nil {
		return err
	}

	if err := p.conn.Write([]byte{0x14, 0xFF}); err != nil {
		return err
	}

	if err := p.conn.Write([]byte{0x15, 0xFF}); err != nil {
		return err
	}

	return nil
}
