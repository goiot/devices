package piglow

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"golang.org/x/exp/io/i2c/driver"
)

type opener struct {
	buf *bytes.Buffer
}

func (o opener) Open(addr int, tenbit bool) (driver.Conn, error) {
	return &conn{buf: o.buf}, nil
}

type conn struct {
	buf *bytes.Buffer
}

func (c conn) Tx(w, r []byte) error {
	if w != nil {
		if _, err := c.buf.Write(w); err != nil {
			return err
		}
	}
	if r != nil {
		if _, err := c.buf.Read(r); err != nil {
			return err
		}
	}
	return nil
}

func (conn) Close() error {
	return nil
}

func openPiGlow(t *testing.T) (*PiGlow, *bytes.Buffer) {
	o := opener{
		buf: bytes.NewBuffer([]byte{}),
	}
	device, err := Open(o)
	if err != nil {
		t.Fatal(err)
	}
	return device, o.buf
}

func assert(t *testing.T, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("got = %v, want = %v", got, want)
	}
}

func TestGreen(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Green(1); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x04, 0x01,
		0x06, 0x01,
		0x0E, 0x01,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestBlue(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Blue(1); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x05, 0x01,
		0x0C, 0x01,
		0x0F, 0x01,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestYellow(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Yellow(1); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x03, 0x01,
		0x09, 0x01,
		0x010, 0x01,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestOrange(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Orange(1); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x02, 0x01,
		0x08, 0x01,
		0x011, 0x01,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestWhite(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.White(1); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x0A, 0x01,
		0x0B, 0x01,
		0x0D, 0x01,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestRed(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Red(1); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x01, 0x01,
		0x07, 0x01,
		0x012, 0x01,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestSetBrightness(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.SetBrightness(10); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x01, 0x0A,
		0x02, 0x0A,
		0x03, 0x0A,
		0x04, 0x0A,
		0x05, 0x0A,
		0x06, 0x0A,
		0x07, 0x0A,
		0x08, 0x0A,
		0x09, 0x0A,
		0x0a, 0x0A,
		0x0b, 0x0A,
		0x0c, 0x0A,
		0x0d, 0x0A,
		0x0e, 0x0A,
		0x0f, 0x0A,
		0x10, 0x0A,
		0x11, 0x0A,
		0x012, 0x0A,
		0x16, 0xFF,
	}
	assert(t, want, buf.Bytes())
}

func TestSetLEDBrightness(t *testing.T) {
	device, buf := openPiGlow(t)

	var states = []struct {
		led   int
		level int
		want  []byte
	}{
		{1, 5, []byte{0x1, 5, 0x16, 0xFF}},
		{2, 10, []byte{0x2, 10, 0x16, 0xFF}},
		{3, 15, []byte{0x3, 15, 0x16, 0xFF}},
		{4, 20, []byte{0x4, 20, 0x16, 0xFF}},
		{5, 25, []byte{0x5, 25, 0x16, 0xFF}},
		{6, 30, []byte{0x6, 30, 0x16, 0xFF}},
		{7, 35, []byte{0x7, 35, 0x16, 0xFF}},
		{8, 40, []byte{0x8, 40, 0x16, 0xFF}},
		{9, 45, []byte{0x9, 45, 0x16, 0xFF}},
		{10, 50, []byte{0x0A, 50, 0x16, 0xFF}},
		{11, 55, []byte{0x0B, 55, 0x16, 0xFF}},
		{12, 60, []byte{0x0C, 60, 0x16, 0xFF}},
		{13, 65, []byte{0x0D, 65, 0x16, 0xFF}},
		{14, 70, []byte{0x0E, 70, 0x16, 0xFF}},
		{15, 75, []byte{0x0F, 75, 0x16, 0xFF}},
		{16, 80, []byte{0x10, 80, 0x16, 0xFF}},
		{17, 85, []byte{0x11, 85, 0x16, 0xFF}},
		{18, 90, []byte{0x12, 90, 0x16, 0xFF}},
	}

	for _, state := range states {
		buf.Reset()
		if err := device.SetLEDBrightness(state.led, state.level); err != nil {
			t.Log(err)
		}
		assert(t, state.want, buf.Bytes())
	}
}

func TestReset(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Reset(); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x17, 0xff}
	assert(t, want, buf.Bytes())
}

func TestShutdown(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Shutdown(); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00, 0x00}
	assert(t, want, buf.Bytes())
}

func TestEnable(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Enable(); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00, 0x01}
	assert(t, want, buf.Bytes())
}

func TestSetLEDControlRegister(t *testing.T) {
	device, buf := openPiGlow(t)

	var states = []struct {
		register int
		enables  int
		want     []byte
		err      error
	}{
		{1, 0xFF, []byte{0x13, 0xFF, 0x16, 0xFF}, nil},
		{2, 0xFF, []byte{0x14, 0xFF, 0x16, 0xFF}, nil},
		{3, 0xFF, []byte{0x15, 0xFF, 0x16, 0xFF}, nil},
		{0, 0xFF, []byte{}, errors.New("0 is an unknown register")},
	}

	for _, state := range states {
		buf.Reset()

		err := device.SetLEDControlRegister(state.register, state.enables)
		assert(t, state.want, buf.Bytes())
		assert(t, state.err, err)
	}
}

func TestSetup(t *testing.T) {
	device, buf := openPiGlow(t)
	if err := device.Setup(); err != nil {
		t.Fatal(err)
	}
	got := []byte{
		0x17, 0xFF,
		0x00, 0x01,
		0x13, 0xFF,
		0x16, 0xFF,
		0x14, 0xFF,
		0x16, 0xFF,
		0x15, 0xFF,
		0x16, 0xFF,
	}
	assert(t, got, buf.Bytes())
}
