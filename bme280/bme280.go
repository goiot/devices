// Package bme280 implements a driver for the bosch bme280 tempreture, humidity
// and barometric pressure sensor.
package bme280

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/davecgh/go-spew/spew"

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
)

// State state as read from the sensor during update
type State struct {
	Temperature float64
	Pressure    float64
	Humidity    float64
}

type calibrationData struct {
	digT1 uint16
	digT2 int16
	digT3 int16

	digP1 uint16
	digP2 int16
	digP3 int16
	digP4 int16
	digP5 int16
	digP6 int16
	digP7 int16
	digP8 int16
	digP9 int16

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
	cal   *calibrationData
	tfine int32
}

// Open connects to the passed driver and sets the device up.
func Open(o driver.Opener) (*Bme280, error) {
	device, err := i2c.Open(o, addr)
	if err != nil {
		return nil, err
	}
	sensor := &Bme280{Device: device, State: &State{}, cal: &calibrationData{}}

	if err = sensor.readChipID(); err != nil {
		return nil, err
	}

	if err = sensor.readCoefficients(); err != nil {
		return nil, err
	}

	if err = sensor.Device.WriteReg(controlRegHumid, []byte{0x05}); err != nil {
		return nil, err
	}

	if err = sensor.Device.WriteReg(controlReg, []byte{0xB7}); err != nil {
		return nil, err
	}

	return sensor, nil
}

// Update read from the sensor and update the state
func (bme *Bme280) Update() error {

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

/*
// Returns pressure in Pa as unsigned 32 bit integer in Q24.8 format (24 integer bits and 8 fractional bits).
// Output value of “24674867” represents 24674867/256 = 96386.2 Pa = 963.862 hPa
BME280_U32_t BME280_compensate_P_int64(BME280_S32_t adc_P)
{
BME280_S64_t var1, var2, p;
var1 = ((BME280_S64_t)t_fine) – 128000;
var2 = var1 * var1 * (BME280_S64_t)dig_P6;
var2 = var2 + ((var1*(BME280_S64_t)dig_P5)<<17);
var2 = var2 + (((BME280_S64_t)dig_P4)<<35);
var1 = ((var1 * var1 * (BME280_S64_t)dig_P3)>>8) + ((var1 * (BME280_S64_t)dig_P2)<<12);
var1 = (((((BME280_S64_t)1)<<47)+var1))*((BME280_S64_t)dig_P1)>>33;
if (var1 == 0)
{
return 0; // avoid exception caused by division by zero
}
p = 1048576-adc_P;
p = (((p<<31)-var2)*3125)/var1;
var1 = (((BME280_S64_t)dig_P9) * (p>>13) * (p>>13)) >> 25;
var2 = (((BME280_S64_t)dig_P8) * p) >> 19;
p = ((p + var1 + var2) >> 8) + (((BME280_S64_t)dig_P7)<<4);
return (BME280_U32_t)p;
}
*/
func (bme *Bme280) readPressure() error {
	buf := make([]byte, 3)
	err := bme.Device.ReadReg(pressureReg, buf)
	if err != nil {
		return err
	}

	adcP := int32(buf[0])<<16 | int32(buf[1])<<8 | int32(buf[2])
	adcP = adcP >> 4 // only want 0xF9 (bit 7, 6, 5, 4)

	var1 := int64(bme.tfine) - 128000
	var2 := var1 * var1 * int64(bme.cal.digP6)
	var2 = var2 + ((var1 * int64(bme.cal.digP5)) << 17)
	var2 = var2 + ((int64(bme.cal.digP4)) << 35)
	var1 = (var1 * var1 * (int64(bme.cal.digP3) >> 8)) + ((var1 * int64(bme.cal.digP2)) << 12)
	var1 = ((int64(1) << 47) + var1) * (int64(bme.cal.digP1)) >> 33
	if var1 == 0 {
		fmt.Printf("press = %d\n", 0)
		return nil
	}
	p := int64(1048576) - int64(adcP)
	p = ((int64(p<<31) - var2) * 3125) / var1
	var1 = (int64(bme.cal.digP8) * (p >> 13) * (p >> 13)) >> 25
	var2 = (int64(bme.cal.digP8) * p) >> 19
	p = ((p + var1 + var2) >> 8) + (int64(bme.cal.digP7) << 4)

	bme.State.Pressure = float64(p) / 256 / 100

	return nil
}

/*
// Returns temperature in DegC, resolution is 0.01 DegC. Output value of “5123” equals 51.23 DegC.
// t_fine carries fine temperature as global value
BME280_S32_t t_fine;
BME280_S32_t BME280_compensate_T_int32(BME280_S32_t adc_T)
{
BME280_S32_t var1, var2, T;
var1 = ((((adc_T>>3) – ((BME280_S32_t)dig_T1<<1))) * ((BME280_S32_t)dig_T2)) >> 11;
var2 = (((((adc_T>>4) – ((BME280_S32_t)dig_T1)) * ((adc_T>>4) – ((BME280_S32_t)dig_T1))) >> 12) *
((BME280_S32_t)dig_T3)) >> 14;
t_fine = var1 + var2;
T = (t_fine * 5 + 128) >> 8;
return T;
}
*/
func (bme *Bme280) readTemperature() error {
	buf := make([]byte, 3)
	err := bme.Device.ReadReg(tempReg, buf)
	if err != nil {
		return err
	}

	adcT := int32(buf[0])<<16 | int32(buf[1])<<8 | int32(buf[2])
	adcT = adcT >> 4 // only want 0xFC (bit 7, 6, 5, 4)

	var1 := ((((adcT) >> 3) - ((int32)(bme.cal.digT1) << 1)) * ((int32)(bme.cal.digT2))) >> 11
	var2 := (((((adcT) >> 4) - ((int32)(bme.cal.digT1))) * ((adcT >> 4) - ((int32)(bme.cal.digT1))) >> 12) * ((int32)(bme.cal.digT3))) >> 14
	bme.tfine = var1 + var2
	t := (bme.tfine*5 + 128) >> 8

	bme.State.Temperature = float64(t) * 0.01

	return nil
}

/*
// Returns humidity in %RH as unsigned 32 bit integer in Q22.10 format (22 integer and 10 fractional bits).
// Output value of “47445” represents 47445/1024 = 46.333 %RH
BME280_U32_t bme280_compensate_H_int32(BME280_S32_t adc_H)
{
BME280_S32_t v_x1_u32r;
v_x1_u32r = (t_fine – ((BME280_S32_t)76800));
v_x1_u32r = (((((adc_H << 14) – (((BME280_S32_t)dig_H4) << 20) – (((BME280_S32_t)dig_H5) * v_x1_u32r)) +
((BME280_S32_t)16384)) >> 15) * (((((((v_x1_u32r * ((BME280_S32_t)dig_H6)) >> 10) * (((v_x1_u32r *
((BME280_S32_t)dig_H3)) >> 11) + ((BME280_S32_t)32768))) >> 10) + ((BME280_S32_t)2097152)) *
((BME280_S32_t)dig_H2) + 8192) >> 14));
v_x1_u32r = (v_x1_u32r – (((((v_x1_u32r >> 15) * (v_x1_u32r >> 15)) >> 7) * ((BME280_S32_t)dig_H1)) >> 4));
v_x1_u32r = (v_x1_u32r < 0 ? 0 : v_x1_u32r);
v_x1_u32r = (v_x1_u32r > 419430400 ? 419430400 : v_x1_u32r);
return (BME280_U32_t)(v_x1_u32r>>12);
}
*/
func (bme *Bme280) readHumidity() error {
	buf := make([]byte, 2)
	err := bme.Device.ReadReg(humidReg, buf)
	if err != nil {
		return err
	}
	fmt.Printf("humid = % x\n", buf)

	adcH := int32(buf[0])<<8 | int32(buf[1])
	fmt.Printf("humid = %d\n", adcH)

	vx1u32r := (bme.tfine - (int32(76800)))
	vx1u32r = ((((adcH << 14) - (int32(bme.cal.digH4) << 20) - (int32(bme.cal.digH5) * vx1u32r)) + int32(16384)) >> 15) *
		(((((((vx1u32r * int32(bme.cal.digH6)) >> 10) * ((vx1u32r * (int32(bme.cal.digH3) >> 11)) + (int32(32768)))) >> 10) + (int32(2097152))) * (int32(bme.cal.digH2) + 8192)) >> 14)
	vx1u32r = (vx1u32r - (((((vx1u32r >> 15) * (vx1u32r >> 15)) >> 7) * (int32(bme.cal.digH1))) >> 4))

	if vx1u32r < 0 {
		vx1u32r = 0
	}

	if vx1u32r > 419430400 {
		vx1u32r = 419430400
	}

	bme.State.Humidity = float64(vx1u32r) / 1024

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

func (bme *Bme280) readCoefficients() error {
	// 0x88…0xA1
	buf := make([]byte, 26)

	err := bme.Device.ReadReg(0x88, buf)

	if err != nil {
		return err
	}

	p := bytes.NewBuffer(buf)

	binary.Read(p, binary.LittleEndian, &bme.cal.digT1)
	binary.Read(p, binary.LittleEndian, &bme.cal.digT2)
	binary.Read(p, binary.LittleEndian, &bme.cal.digT3)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP1)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP2)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP3)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP4)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP5)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP6)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP7)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP8)
	binary.Read(p, binary.LittleEndian, &bme.cal.digP9)

	buf = make([]byte, 1)
	err = bme.Device.ReadReg(digH1Reg, buf)
	if err != nil {
		return err
	}

	bme.cal.digH1 = buf[0]

	buf = make([]byte, 8)
	err = bme.Device.ReadReg(0xE1, buf)
	if err != nil {
		return err
	}

	bme.cal.digH2 = int16(buf[0])<<8 | int16(buf[1])
	bme.cal.digH3 = buf[2]
	bme.cal.digH4 = int16(buf[3])<<4 | (int16(buf[4]) & 0xF)
	bme.cal.digH5 = int16(buf[6])<<4 | int16(buf[5])>>4
	bme.cal.digH6 = int8(buf[7])

	spew.Dump(bme.cal)

	return nil
}
