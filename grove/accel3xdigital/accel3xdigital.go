// Package accel3xdigital allows developers to read x,y,z and acceleratation.
// Currently this library doesn't support any of the interrupts:  Front/Back Interrupt,
// Up/Down/Left/Right Interrupt, Tap Detection Interrupt, GINT (real-time motion tracking), Shake on X-axis, Shake on Y-axis, and
// Shake on Z-axis.
package accel3xdigital

import (
	"errors"
	"fmt"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

var (
	// ErrNotReady warns the user that the device isn't not (yet) ready
	ErrNotReady = errors.New("device is not ready")
)

// Accel3xDigital reresents the Grove 3-Axis Digital Accelerometer(±1.5g)
type Accel3xDigital struct {
	Device *i2c.Device
	State  *State
	// TapEnabled lets devs know if the feature is on or off (default)
	TapEnabled bool
}

// Open connects to the passed driver and sets things up.
// At this point the sensor will sample 32 times a second and will store its information in registries you can read from
// by calling update.
// Note that by default the tap detection is not on. You need to enable this feature manually.
func Open(o driver.Opener) (*Accel3xDigital, error) {
	device, err := i2c.Open(o)
	if err != nil {
		return nil, err
	}

	accel := &Accel3xDigital{Device: device, State: &State{}}

	if err := accel.ChangeMode(StandBy); err != nil {
		return accel, err
	}

	if err := accel.Device.Write([]byte{accelSr, accelAutoSleep32}); err != nil {
		return accel, err
	}

	if err := accel.Device.Write([]byte{accelMode, accelActive}); err != nil {
		return accel, err
	}

	return accel, nil
}

// ChangeMode allows developers to switch between standy and active (default).
func (a *Accel3xDigital) ChangeMode(m Mode) (err error) {
	err = a.Device.Write([]byte{accelMode, byte(m)})
	if err != nil {
		err = fmt.Errorf("failed to change mode - %v", err)
	}
	return err
}

// Enable tap enables checking for taps by increasing the sample rate
func (a *Accel3xDigital) EnableTap() error {
	if err := a.ChangeMode(StandBy); err != nil {
		return err
	}

	if err := a.Device.Write([]byte{accelSr, accelAutoSleep120}); err != nil {
		return err
	}

	// set tap detection sensitivity (how many samples to check for)
	// we are setting the threashold at 80 samples knowing that we are sampling at 120hz
	if err := a.Device.Write([]byte{accelPd, 80}); err != nil {
		return err
	}

	if err := a.Device.Write([]byte{accelMode, accelActive}); err != nil {
		return err
	}
	a.TapEnabled = true
	return nil
}

// SetTapSensitivity sets the debounce filtering requirement to n adjacent tap detection tests to
// be the same to trigger a tap event.
// Note that the sampling rate is 120 samples/second when tap is enabled.
func (a *Accel3xDigital) SetTapSensitivity(n uint8) error {
	if !a.TapEnabled {
		if err := a.EnableTap(); err != nil {
			return err
		}
	}

	if err := a.ChangeMode(StandBy); err != nil {
		return err
	}

	if err := a.Device.Write([]byte{accelPd, byte(n)}); err != nil {
		return err
	}

	if err := a.Device.Write([]byte{accelMode, accelActive}); err != nil {
		return err
	}

	return nil
}

// Update reads the sensor and update the state in memory
func (a *Accel3xDigital) Update() error {
	err := a.Device.Read(stateBuff)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		val := stateBuff[i]
		if ((val >> 6) & 0x01) == 1 {
			return ErrNotReady
		}
	}

	a.State.X = float64((int8(stateBuff[0]) << 2)) / 4.0
	a.State.Y = float64((int8(stateBuff[1]) << 2)) / 4.0
	a.State.Z = float64((int8(stateBuff[2]) << 2)) / 4.0

	tilt := stateBuff[3]
	a.State.Front = (tilt & (1 << 0)) > 0
	a.State.Back = (tilt & (1 << 1)) > 0
	a.State.Tapped = (tilt & (1 << 5)) > 0
	a.State.Alert = (tilt & (1 << 6)) > 0
	a.State.Shaken = (tilt & (1 << 7)) > 0

	masks := [3]bool{
		tilt&(1<<2) > 0,
		tilt&(1<<3) > 0,
		tilt&(1<<4) > 0,
	}

	switch masks {
	case downMask:
		a.State.Position = Down
	case upMask:
		a.State.Position = Up
	case leftMask:
		a.State.Position = Left
	case rightMask:
		a.State.Position = Right
	default:
		a.State.Position = Unknown
	}

	// report race conditions
	if a.State.Alert {
		return errors.New("error reading state, try again")
	}

	return nil
}

// Close puts the device on standby
func (a *Accel3xDigital) Close() error {
	a.ChangeMode(StandBy)
	return a.Device.Close()
}

// State contains the last read state of the device
type State struct {
	// Front is true if the equipment is lying on its front
	Front bool
	// Back is true if the equipment is lying on its back
	Back     bool
	Position Position
	// Tapped reports that the device was tapped, after reading the data once, the flag is cleared.
	Tapped bool
	// Shaken reports that the device was shaken, after reading the data once, the flag is cleared.
	Shaken bool
	// Alert can be triggered if there was a race condition when we tried to read (the sensor was updating the data at the same time)
	// If you get an alert, ignore the data and read again.
	Alert bool
	// X axis value
	X float64
	// Y axis value
	Y float64
	// Z axis value
	Z float64
}

// Acceleration returns the acceleration (g) for each axis (x,y,z)
func (s *State) Acceleration() (float64, float64, float64) {
	return s.X / 21.0, s.Y / 21.0, s.Y / 21.0
}

// String implements the stringer interface
func (s *State) String() string {
	xg, yg, zg := s.Acceleration()
	orientation := "unknown"
	if s.Front {
		orientation = "front facing"
	} else if s.Back {
		orientation = "back facing"
	}
	return fmt.Sprintf(`Current State:
  X: %2.f, Y: %2.f, Z: %2.f
Acceleration:
  X: %2.fg, Y: %2.fg, Z: %2.fg
Orientation: %s
Position: %s
Was Shaken?: %t,
Was tapped?: %t
`, s.X, s.Y, s.Z,
		xg, yg, zg,
		orientation,
		s.Position,
		s.Shaken,
		s.Tapped)
}
