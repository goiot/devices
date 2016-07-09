// Package bme280 implements a driver for the bosch bme280 temperature, humidity
// and barometric pressure sensor.
//
// Datasheet https://cdn-shop.adafruit.com/datasheets/BST-BME280_DS001-10.pdf
package bme280

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

const (
	// addr is the i2c address of this sensor
	addr = 0x76

	chipIDReg       = 0xD0
	versionReg      = 0xD1
	softResetReg    = 0xE0
	controlReg      = 0xF4
	controlRegHumid = 0xF2

	pressureReg = 0xF7
	tempReg     = 0xFA
	humidReg    = 0xFD

	digT1Reg = 0x88
	digT2Reg = 0x8A
	digT3Reg = 0x8C

	digP1Reg = 0x8E
	digP2Reg = 0x90
	digP3Reg = 0x92
	digP4Reg = 0x94
	digP5Reg = 0x96
	digP6Reg = 0x98
	digP7Reg = 0x9A
	digP8Reg = 0x9C
	digP9Reg = 0x9E

	digH1Reg = 0xA1
	digH2Reg = 0xE1
	digH3Reg = 0xE3
	digH4Reg = 0xE4
	digH5Reg = 0xE5
	digH6Reg = 0xE7

	chipID = 0x60

	seaLevelhPa = 1013.25

	// Controls oversampling of data
	pressureShift = 5
	tempShift     = 2

	// Oversample0 oversample 0 times
	Oversample0 = 0
	// Oversample1 oversample 1 times
	Oversample1 = 1
	// Oversample2 oversample 2 times
	Oversample2 = 2
	// Oversample4 oversample 4 times
	Oversample4 = 3
	// Oversample8 oversample 8 times
	Oversample8 = 4
	// Oversample16 oversample 16 times
	Oversample16 = 5

	// Controls the sensor mode of the device
	modeSleep  = 0
	modeForced = 2
	modeNormal = 3

	setupMaxPressure  = 10
	setupMaxHumidity  = 10
	initMax           = 20
	measurePerOsrsMax = 37
)

var (
	oversamplePressure = shiftFunc(pressureShift)
	oversampleTemp     = shiftFunc(tempShift)
	oversampleHumidity = shiftFunc(0)
	mode               = shiftFunc(0)
)

// State state as read from the sensor during update
type State struct {
	Temperature float64
	Pressure    float64
	Humidity    float64
}

type tempCalData struct {
	DigT1 uint16
	DigT2 int16
	DigT3 int16
}
type pressureCalData struct {
	DigP1 uint16
	DigP2 int16
	DigP3 int16
	DigP4 int16
	DigP5 int16
	DigP6 int16
	DigP7 int16
	DigP8 int16
	DigP9 int16
}

type humidCalData struct {
	digH1 uint8
	digH2 int16
	digH3 uint8
	digH4 int16
	digH5 int16
	digH6 int8
}

// Bme280 reresents the bosch bme280 tempreture, humidity and
// barometric pressure sensor
type Bme280 struct {
	Device *i2c.Device

	State *State
	tcal  *tempCalData
	pcal  *pressureCalData
	hcal  *humidCalData
	tfine int32
}

// Open connects to the passed driver and sets the device up.
func Open(o driver.Opener) (*Bme280, error) {
	device, err := i2c.Open(o, addr)
	if err != nil {
		return nil, err
	}
	sensor := &Bme280{
		Device: device,
		State:  &State{},
		tcal:   &tempCalData{},
		pcal:   &pressureCalData{},
		hcal:   &humidCalData{},
	}

	if err = sensor.readChipID(); err != nil {
		return nil, err
	}

	if err = sensor.readCoefficients(); err != nil {
		return nil, err
	}

	return sensor, nil
}

// Update read from the sensor and update the state
func (bme *Bme280) Update(oversample uint8) error {

	if err := bme.Device.WriteReg(controlRegHumid,
		[]byte{oversampleHumidity(oversample)}); err != nil {
		return err
	}

	controlSettings := mode(modeNormal) | oversampleTemp(oversample) | oversamplePressure(oversample)

	if err := bme.Device.WriteReg(controlReg,
		[]byte{controlSettings}); err != nil {
		return err
	}

	sleeptime := measureSleeptime(oversample)

	time.Sleep(sleeptime)

	err := bme.readTemperature()
	if err != nil {
		return err
	}

	err = bme.readPressure()
	if err != nil {
		return err
	}

	err = bme.readHumidity()
	if err != nil {
		return err
	}

	return nil

}

// Returns pressure in Pa as unsigned 32 bit integer in Q24.8 format (24 integer bits and 8 fractional bits).
// Output value of “24674867” represents 24674867/256 = 96386.2 Pa = 963.862 hPa
func (bme *Bme280) readPressure() error {
	adcP, err := bme.read24(pressureReg)
	if err != nil {
		return err
	}

	adcP = adcP >> 4 // only want 0xF9 (bit 7, 6, 5, 4)

	cal := bme.pcal

	// ported from 4.2.3 Compensation formulas section of datasheet

	var1 := int64(bme.tfine) - 128000
	var2 := var1 * var1 * int64(cal.DigP6)
	var2 = var1 * var1 * int64(cal.DigP6)
	var2 = var2 + ((var1 * int64(cal.DigP5)) << 17)
	var2 = var2 + ((int64(cal.DigP4)) << 35)
	var1 = ((var1 * var1 * int64(cal.DigP3)) >> 8) + ((var1 * int64(cal.DigP2)) << 12)
	var1 = (((int64(1)) << 47) + var1) * (int64(cal.DigP1)) >> 33

	if var1 == 0 {
		bme.State.Pressure = 0
		return nil
	}

	p := int64(1048576) - int64(adcP)
	p = (((p << 31) - var2) * 3125) / var1
	var1 = ((int64(cal.DigP9)) * (p >> 13) * (p >> 13)) >> 25
	var2 = ((int64(cal.DigP8)) * p) >> 19
	p = ((p + var1 + var2) >> 8) + ((int64(cal.DigP7)) << 4)

	bme.State.Pressure = float64(p) / 256 / 100

	return nil
}

// Returns temperature in DegC, resolution is 0.01 DegC. Output value of “5123” equals 51.23 DegC.
// t_fine carries fine temperature as global value
func (bme *Bme280) readTemperature() error {
	adcT, err := bme.read24(tempReg)
	if err != nil {
		return err
	}

	adcT = adcT >> 4 // only want 0xFC (bit 7, 6, 5, 4)
	cal := bme.tcal

	// ported from 4.2.3 Compensation formulas section of datasheet

	var1 := (((adcT >> 3) - (int32(cal.DigT1) << 1)) * (int32(cal.DigT2))) >> 11
	var2 := (((((adcT >> 4) - (int32(cal.DigT1))) * ((adcT >> 4) - (int32(cal.DigT1)))) >> 12) * (int32(cal.DigT3))) >> 14

	bme.tfine = var1 + var2
	t := (bme.tfine*5 + 128) >> 8

	bme.State.Temperature = float64(t) / 100

	return nil
}

// Returns humidity in %RH as unsigned 32 bit integer in Q22.10 format (22 integer and 10 fractional bits).
// Output value of “47445” represents 47445/1024 = 46.333 %RH
func (bme *Bme280) readHumidity() error {

	adcH, err := bme.read16(humidReg)
	if err != nil {
		return err
	}

	cal := bme.hcal

	// ported from 4.2.3 Compensation formulas section of datasheet

	vx1u32r := (bme.tfine - (int32(76800)))

	vx1u32r = (((((adcH << 14) - ((int32(cal.digH4)) << 20) - ((int32(cal.digH5)) * vx1u32r)) +
		(int32(16384))) >> 15) * (((((((vx1u32r*(int32(cal.digH6)))>>10)*(((vx1u32r*
		(int32(cal.digH3)))>>11)+(int32(32768))))>>10)+(int32(2097152)))*
		(int32(cal.digH2)) + 8192) >> 14))

	vx1u32r = (vx1u32r - (((((vx1u32r >> 15) * (vx1u32r >> 15)) >> 7) * (int32(cal.digH1))) >> 4))

	if vx1u32r < 0 {
		vx1u32r = 0
	}

	if vx1u32r > 419430400 {
		vx1u32r = 419430400
	}

	bme.State.Humidity = float64(vx1u32r>>12) / 1024

	return nil

}

func (bme *Bme280) readChipID() error {
	chipIDBuff := make([]byte, 1)

	err := bme.Device.ReadReg(chipIDReg, chipIDBuff)
	if err != nil {
		return err
	}

	if chipIDBuff[0] != chipID {
		return fmt.Errorf("ChipId mismatch expected %x got %x", chipID, chipIDBuff[0])

	}

	return nil
}

func (bme *Bme280) read24(reg byte) (int32, error) {
	var v int32
	buf := make([]byte, 3)
	err := bme.Device.ReadReg(reg, buf)
	if err != nil {
		return v, err
	}

	v = int32(buf[0])<<16 | int32(buf[1])<<8 | int32(buf[2])

	return v, nil
}

func (bme *Bme280) read16(reg byte) (int32, error) {
	var v int32
	buf := make([]byte, 2)
	err := bme.Device.ReadReg(reg, buf)
	if err != nil {
		return v, err
	}

	v = int32(buf[0])<<8 | int32(buf[1])

	return v, nil

}

func (bme *Bme280) readCoefficients() error {
	// 0x88…0xA1
	buf := make([]byte, 26)

	err := bme.Device.ReadReg(0x88, buf)

	if err != nil {
		return err
	}

	p := bytes.NewBuffer(buf)

	if err := binary.Read(p, binary.LittleEndian, bme.tcal); err != nil {
		return err
	}

	if err := binary.Read(p, binary.LittleEndian, bme.pcal); err != nil {
		return err
	}

	buf = make([]byte, 1)
	err = bme.Device.ReadReg(digH1Reg, buf)
	if err != nil {
		return err
	}

	bme.hcal.digH1 = buf[0]

	buf = make([]byte, 7)
	err = bme.Device.ReadReg(digH2Reg, buf)
	if err != nil {
		return err
	}

	bme.hcal.digH2 = int16(buf[0]) | int16(buf[1])<<8
	bme.hcal.digH3 = buf[2]
	bme.hcal.digH4 = int16(buf[3])<<4 | (int16(buf[4]<<4))>>4
	bme.hcal.digH5 = int16(buf[4])>>4 | int16(buf[5])<<4
	bme.hcal.digH6 = int8(buf[6])

	return nil
}

func measureSleeptime(oversample uint8) time.Duration {
	maxPressure := uint16(0)
	maxHumidity := uint16(0)

	if oversample > 0 {
		maxPressure = setupMaxPressure
		maxHumidity = setupMaxHumidity
	}

	v := uint16(((1 << oversample) >> 1) + ((1 << oversample) >> 1) + ((1 << oversample) >> 1))

	sleepms := (initMax + measurePerOsrsMax*v + (maxPressure) + (maxHumidity) + 15) / 16
	fmt.Printf("sleepms = %d\n", sleepms)
	return time.Duration(sleepms) * time.Millisecond
}

func shiftFunc(shift uint8) func(v uint8) uint8 {
	return func(v uint8) uint8 {
		// fmt.Printf("%d %b\n", v, v)
		if v < 6 {
			return v << shift
		}
		return 1
	}
}
