package bme280

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPressureOversample(t *testing.T) {
	assert.Equal(t, oversamplePressure(Oversample1), uint8(32))
	assert.Equal(t, oversamplePressure(Oversample2), uint8(64))
	assert.Equal(t, oversamplePressure(Oversample4), uint8(96))
	assert.Equal(t, oversamplePressure(Oversample8), uint8(128))
	assert.Equal(t, oversamplePressure(Oversample16), uint8(160))

}

func TestTempOversample(t *testing.T) {
	oversampleTemp := shiftFunc(tempShift)
	assert.Equal(t, oversampleTemp(Oversample1), uint8(4))
	assert.Equal(t, oversampleTemp(Oversample2), uint8(8))
	assert.Equal(t, oversampleTemp(Oversample4), uint8(12))
	assert.Equal(t, oversampleTemp(Oversample8), uint8(16))
	assert.Equal(t, oversampleTemp(Oversample16), uint8(20))
}

func TestMode(t *testing.T) {
	assert.Equal(t, mode(modeSleep), uint8(0))
	assert.Equal(t, mode(modeForced), uint8(2))
	assert.Equal(t, mode(modeNormal), uint8(3))
}

func TestConfig(t *testing.T) {
	assert.Equal(t, mode(modeNormal)|oversampleTemp(Oversample16)|oversamplePressure(Oversample16), uint8(0xB7))
}

func TestMeasureSleeptime(t *testing.T) {
	assert.Equal(t, measureSleeptime(Oversample2), 17*time.Millisecond)
}
