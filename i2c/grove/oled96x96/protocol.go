package oled96x96

// ScrollDirection is the type determining the scrolling direction of text
type ScrollDirection byte

// ScrollSpeed is the type determining the speed of scrolling
type ScrollSpeed byte

var (
	ScrollLeft  ScrollDirection = 0x00
	ScrollRight ScrollDirection = 0x01

	Scroll2Frames   ScrollSpeed = 0x7
	Scroll3Frames   ScrollSpeed = 0x4
	Scroll4Frames   ScrollSpeed = 0x5
	Scroll5Frames   ScrollSpeed = 0x0
	Scroll25Frames  ScrollSpeed = 0x6
	Scroll64Frames  ScrollSpeed = 0x1
	Scroll128Frames ScrollSpeed = 0x2
	Scroll256Frames ScrollSpeed = 0x3
)

var (
	// buffer sent to indicate the following data belongs to a command
	cmdCmdBuf  = []byte{cmdCmd, 0x0}
	dataCmdBuf = []byte{dataCmd, 0x0}
)

const (
	VerticalModeFlag   = 01
	HorizontalModeFlag = 02

	// Address is the i2c address of the device
	Address = 0x3c

	cmdCmd           byte = 0x80
	dataCmd          byte = 0x40
	lockUnlockCmd    byte = 0xFD // takes a 2nd arg byte
	startLineCmd     byte = 0xA1 // takes a 2nd arg byte
	displayOffCmd    byte = 0xAE
	displayOnCmd     byte = 0xAF
	displayOffsetCmd byte = 0xA2 // takes a 2nd arg byte
	setColAddrCmd    byte = 0x15 // takes 3 arg bytes

	normalDisplayCmd   byte = 0xA4
	inverseDisplayCmd  byte = 0xA7
	ActivateScrollCmd       = 0x2F
	DectivateScrollCmd      = 0x2E
	contrastLevelCmd        = 0x81
)

// sendCmd sends the passed data preluded by the command byte
func (o *Oled96x96) sendCmd(buf ...byte) error {
	for _, b := range buf {
		cmdCmdBuf[1] = b
		if err := o.Conn.Write(cmdCmdBuf); err != nil {
			return err
		}
	}
	return nil
}

// sendData does what you expect it does and maybe even more
func (o *Oled96x96) sendData(buf ...byte) error {
	for _, b := range buf {
		dataCmdBuf[1] = b
		if err := o.Conn.Write(dataCmdBuf); err != nil {
			return err
		}
	}
	return nil
}
