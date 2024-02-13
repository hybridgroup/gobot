package i2c

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*GrovePiDriver)(nil)

// must implement the DigitalReader interface
var _ gpio.DigitalReader = (*GrovePiDriver)(nil)

// must implement the DigitalWriter interface
var _ gpio.DigitalWriter = (*GrovePiDriver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*GrovePiDriver)(nil)

// must implement the AnalogWriter interface
var _ aio.AnalogWriter = (*GrovePiDriver)(nil)

// must implement the Adaptor interface
var _ gobot.Adaptor = (*GrovePiDriver)(nil)

func initGrovePiDriverWithStubbedAdaptor() (*GrovePiDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGrovePiDriver(adaptor), adaptor
}

func TestNewGrovePiDriver(t *testing.T) {
	var di interface{} = NewGrovePiDriver(newI2cTestAdaptor())
	d, ok := di.(*GrovePiDriver)
	if !ok {
		require.Fail(t, "NewGrovePiDriver() should have returned a *GrovePiDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "GrovePi"))
	assert.Equal(t, 0x04, d.defaultAddress)
	assert.NotNil(t, d.pins)
}

func TestGrovePiOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewGrovePiDriver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestGrovePiSomeRead(t *testing.T) {
	// arrange
	tests := map[string]struct {
		usedPin          int
		wantWritten      []uint8
		simResponse      [][]uint8
		wantErr          error
		wantCallsRead    int
		wantResult       int
		wantResultF1     float32
		wantResultF2     float32
		wantResultString string
	}{
		"DigitalRead": {
			usedPin:       2,
			wantWritten:   []uint8{commandSetPinMode, 2, 0, 0, commandReadDigital, 2, 0, 0},
			simResponse:   [][]uint8{{0}, {commandReadDigital, 3}},
			wantCallsRead: 2,
			wantResult:    3,
		},
		"AnalogRead": {
			usedPin:     3,
			wantWritten: []uint8{commandSetPinMode, 3, 0, 0, commandReadAnalog, 3, 0, 0},
			simResponse: [][]uint8{{0}, {commandReadAnalog, 4, 5}},
			wantResult:  1029,
		},
		"UltrasonicRead": {
			usedPin:     4,
			wantWritten: []uint8{commandSetPinMode, 4, 0, 0, commandReadUltrasonic, 4, 0, 0},
			simResponse: [][]uint8{{0}, {commandReadUltrasonic, 5, 6}},
			wantResult:  1281,
		},
		"FirmwareVersionRead": {
			wantWritten:      []uint8{commandReadFirmwareVersion, 0, 0, 0},
			simResponse:      [][]uint8{{commandReadFirmwareVersion, 7, 8, 9}},
			wantResultString: "7.8.9",
		},
		"DHTRead": {
			usedPin:      5,
			wantWritten:  []uint8{commandSetPinMode, 5, 0, 0, commandReadDHT, 5, 1, 0},
			simResponse:  [][]uint8{{0}, {commandReadDHT, 164, 112, 69, 193, 20, 174, 54, 66}},
			wantResultF1: -12.34,
			wantResultF2: 45.67,
		},
		"DigitalRead_error_wrong_return_cmd": {
			usedPin:     15,
			wantWritten: []uint8{commandSetPinMode, 15, 0, 0, commandReadDigital, 15, 0, 0},
			simResponse: [][]uint8{{0}, {0, 2}},
			wantErr:     fmt.Errorf("answer (0) was not for command (1)"),
		},
		"AnalogRead_error_wrong_return_cmd": {
			usedPin:     16,
			wantWritten: []uint8{commandSetPinMode, 16, 0, 0, commandReadAnalog, 16, 0, 0},
			simResponse: [][]uint8{{0}, {0, 3, 4}},
			wantErr:     fmt.Errorf("answer (0) was not for command (3)"),
		},
		"UltrasonicRead_error_wrong_return_cmd": {
			usedPin:     17,
			wantWritten: []uint8{commandSetPinMode, 17, 0, 0, commandReadUltrasonic, 17, 0, 0},
			simResponse: [][]uint8{{0}, {0, 5, 6}},
			wantErr:     fmt.Errorf("answer (0) was not for command (7)"),
		},
		"FirmwareVersionRead_error_wrong_return_cmd": {
			wantWritten: []uint8{commandReadFirmwareVersion, 0, 0, 0},
			simResponse: [][]uint8{{0, 7, 8, 9}},
			wantErr:     fmt.Errorf("answer (0) was not for command (8)"),
		},
		"DHTRead_error_wrong_return_cmd": {
			usedPin:     18,
			wantWritten: []uint8{commandSetPinMode, 18, 0, 0, commandReadDHT, 18, 1, 0},
			simResponse: [][]uint8{{0}, {0, 164, 112, 69, 193, 20, 174, 54, 66}},
			wantErr:     fmt.Errorf("answer (0) was not for command (40)"),
		},
		"DigitalRead_error_wrong_data_count": {
			usedPin:     28,
			wantWritten: []uint8{commandSetPinMode, 28, 0, 0, commandReadDigital, 28, 0, 0},
			simResponse: [][]uint8{{0}, {commandReadDigital, 2, 3}},
			wantErr:     fmt.Errorf("read count mismatch (3 should be 2)"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g, a := initGrovePiDriverWithStubbedAdaptor()
			_ = g.Start()
			a.written = []byte{} // reset writes of former test and start
			numCallsRead := 0
			a.i2cReadImpl = func(bytes []byte) (int, error) {
				numCallsRead++
				copy(bytes, tc.simResponse[numCallsRead-1])
				return len(tc.simResponse[numCallsRead-1]), nil
			}
			var got int
			var gotF1, gotF2 float32
			var gotString string
			var err error
			// act
			switch {
			case strings.Contains(name, "DigitalRead"):
				got, err = g.DigitalRead(strconv.Itoa(tc.usedPin))
			case strings.Contains(name, "AnalogRead"):
				got, err = g.AnalogRead(strconv.Itoa(tc.usedPin))
			case strings.Contains(name, "UltrasonicRead"):
				got, err = g.UltrasonicRead(strconv.Itoa(tc.usedPin), 2)
			case strings.Contains(name, "FirmwareVersionRead"):
				gotString, err = g.FirmwareVersionRead()
			case strings.Contains(name, "DHTRead"):
				gotF1, gotF2, err = g.DHTRead(strconv.Itoa(tc.usedPin), 1, 2)
			default:
				require.Fail(t, "unknown command %s", name)
				return
			}
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantWritten, a.written)
			assert.Len(t, tc.simResponse, numCallsRead)
			assert.Equal(t, tc.wantResult, got)
			assert.InDelta(t, tc.wantResultF1, gotF1, 0.0)
			assert.InDelta(t, tc.wantResultF2, gotF2, 0.0)
			assert.Equal(t, tc.wantResultString, gotString)
		})
	}
}

func TestGrovePiSomeWrite(t *testing.T) {
	// arrange
	tests := map[string]struct {
		usedPin     int
		usedValue   int
		wantWritten []uint8
		simResponse []uint8
	}{
		"DigitalWrite": {
			usedPin:     2,
			usedValue:   3,
			wantWritten: []uint8{commandSetPinMode, 2, 1, 0, commandWriteDigital, 2, 3, 0},
			simResponse: []uint8{4},
		},
		"AnalogWrite": {
			usedPin:     5,
			usedValue:   6,
			wantWritten: []uint8{commandSetPinMode, 5, 1, 0, commandWriteAnalog, 5, 6, 0},
			simResponse: []uint8{7},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g, a := initGrovePiDriverWithStubbedAdaptor()
			_ = g.Start()
			a.written = []byte{} // reset writes of former test and start
			a.i2cReadImpl = func(bytes []byte) (int, error) {
				copy(bytes, tc.simResponse)
				return len(bytes), nil
			}
			var err error
			// act
			switch name {
			case "DigitalWrite":
				err = g.DigitalWrite(strconv.Itoa(tc.usedPin), byte(tc.usedValue))
			case "AnalogWrite":
				err = g.AnalogWrite(strconv.Itoa(tc.usedPin), tc.usedValue)
			default:
				require.Fail(t, "unknown command %s", name)
				return
			}
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestGrovePi_getPin(t *testing.T) {
	assert.Equal(t, "1", getPin("a1"))
	assert.Equal(t, "16", getPin("A16"))
	assert.Equal(t, "3", getPin("D3"))
	assert.Equal(t, "22", getPin("d22"))
	assert.Equal(t, "22", getPin("22"))
}
