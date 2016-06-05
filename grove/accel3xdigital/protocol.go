package accel3xdigital

// Mode The sensor has three power modes: Off Mode, Standby Mode, and Active Mode to offer the customer different power
// consumption options. The sensor is only capable of running in one of these modes at a time.
type Mode byte

var (
	downMask  = [3]bool{true, false, true}
	upMask    = [3]bool{false, true, true}
	leftMask  = [3]bool{true, false, false}
	rightMask = [3]bool{false, true, false}

	stateBuff = make([]byte, 4)
)

const (
	// Standby Mode is ideal for battery operated products. When Standby Mode is active the device outputs are turned off
	// providing a significant reduction in operating current. When the device is in Standby Mode the current will be reduced to
	// approximately 3 µA. Standby Mode is entered as soon as both analog and digital power supplies are up.
	// In this mode, the device can read and write to the registers with I2C, but no new measurements can be taken.
	StandBy = Mode(accelStandBy)
	// Active Mode, continuous measurement on all three axes is enabled. In addition, the user can choose to enable:
	// Shake Detection, Tap Detection, Orientation Detection, and/or Auto-Wake/Sleep Feature and in this mode the digital analysis for
	// any of these functions is done.
	Active = Mode(accelActive)
)

const (
	// Address is the i2c address of this sensor
	Address = 0x4c

	accelX    = 0x00
	accelY    = 0x01
	accelZ    = 0x02
	accelTilt = 0x03

	accelSrst  = 0x04
	accelSpcnt = 0x05
	accelIntsu = 0x06

	accelMode    = 0x07
	accelStandBy = 0x00
	accelActive  = 0x01

	// sample rate
	accelSr           = 0x08
	accelAutoSleep120 = 0x00
	accelAutoSleep64  = 0x01
	// 32 samples per second (default)
	accelAutoSleep32 = 0x02
	accelAutoSleep16 = 0x03
	accelAutoSleep8  = 0x04
	accelAutoSleep4  = 0x05
	accelAutoSleep2  = 0x06
	accelAutoSleep1  = 0x07

	accelPdet = 0x09
	accelPd   = 0x0A
)

// Position indicates the position of the sensor/device
type Position int

const (
	// Unknown condition of up or down or left or right
	Unknown Position = iota
	// Left is true if in landscape mode to the left
	Left
	// Right is true if in landscape mode to the right
	Right
	// Down is true if standing vertically in inverted orientation
	Down
	// Up is true if standing vertically in normal orientation
	Up
)

func (p Position) String() string {
	switch p {
	case Right:
		return "right"
	case Left:
		return "left"
	case Down:
		return "down"
	case Up:
		return "up"
	default:
		return "unkown"
	}
}
