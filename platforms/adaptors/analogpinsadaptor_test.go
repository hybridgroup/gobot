//nolint:nonamedreturns // ok for tests
package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

const (
	analogReadPath            = "/sys/bus/iio/devices/iio:device0/in_voltage0_raw"
	analogWritePath           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/export"
	analogReadWritePath       = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/period"
	analogReadWriteStringPath = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/pwm44/polarity"
)

var analogMockPaths = []string{
	analogReadPath,
	analogWritePath,
	analogReadWritePath,
	analogReadWriteStringPath,
}

func initTestAnalogPinsAdaptorWithMockedFilesystem(mockPaths []string) (*AnalogPinsAdaptor, *system.MockFilesystem) {
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(mockPaths)
	a := NewAnalogPinsAdaptor(sys, testAnalogPinTranslator)
	fs.Files[analogReadPath].Contents = "54321"
	fs.Files[analogWritePath].Contents = "0"
	fs.Files[analogReadWritePath].Contents = "30000"
	fs.Files[analogReadWriteStringPath].Contents = "inverted"
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func testAnalogPinTranslator(id string) (string, bool, uint16, error) {
	switch id {
	case "read":
		return analogReadPath, false, 10, nil
	case "write":
		return analogWritePath, true, 0, nil
	case "read/write":
		return analogReadWritePath, true, 12, nil
	case "read/write_string":
		return analogReadWriteStringPath, true, 13, nil
	}

	return "", false, 0, fmt.Errorf("'%s' is not a valid id of an analog pin", id)
}

func TestAnalogPinsConnect(t *testing.T) {
	translate := func(id string) (path string, w bool, bufLen uint16, err error) { return }
	a := NewAnalogPinsAdaptor(system.NewAccesser(), translate)
	assert.Equal(t, (map[string]gobot.AnalogPinner)(nil), a.pins)

	err := a.AnalogWrite("write", 1)
	require.ErrorContains(t, err, "not connected")

	err = a.Connect()
	require.NoError(t, err)
	assert.NotEqual(t, (map[string]gobot.AnalogPinner)(nil), a.pins)
	assert.Empty(t, a.pins)
}

func TestAnalogPinsFinalize(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(analogMockPaths)
	a := NewAnalogPinsAdaptor(sys, testAnalogPinTranslator)
	fs.Files[analogReadPath].Contents = "0"
	// assert that finalize before connect is working
	require.NoError(t, a.Finalize())
	// arrange
	require.NoError(t, a.Connect())
	require.NoError(t, a.AnalogWrite("write", 1))
	assert.Len(t, a.pins, 1)
	// act
	err := a.Finalize()
	// assert
	require.NoError(t, err)
	assert.Empty(t, a.pins)
	// assert that finalize after finalize is working
	require.NoError(t, a.Finalize())
	// arrange missing file
	require.NoError(t, a.Connect())
	require.NoError(t, a.AnalogWrite("write", 2))
	delete(fs.Files, analogWritePath)
	err = a.Finalize()
	require.NoError(t, err) // because there is currently no access on finalize
	// arrange write error
	require.NoError(t, a.Connect())
	require.NoError(t, a.AnalogWrite("read/write_string", 5))
	fs.WithWriteError = true
	err = a.Finalize()
	require.NoError(t, err) // because there is currently no access on finalize
}

func TestAnalogPinsReConnect(t *testing.T) {
	// arrange
	a, _ := initTestAnalogPinsAdaptorWithMockedFilesystem(analogMockPaths)
	require.NoError(t, a.AnalogWrite("read/write_string", 1))
	assert.Len(t, a.pins, 1)
	require.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	require.NoError(t, err)
	assert.NotNil(t, a.pins)
	assert.Empty(t, a.pins)
}

func TestAnalogWrite(t *testing.T) {
	tests := map[string]struct {
		pin              string
		simulateWriteErr bool
		simulateReadErr  bool
		wantValW         string
		wantValRW        string
		wantValRWS       string
		wantErr          string
	}{
		"write_w_pin": {
			pin:        "write",
			wantValW:   "100",
			wantValRW:  "30000",
			wantValRWS: "inverted",
		},
		"write_rw_pin": {
			pin:        "read/write_string",
			wantValW:   "0",
			wantValRW:  "30000",
			wantValRWS: "100",
		},
		"ok_on_read_error": {
			pin:             "read/write_string",
			simulateReadErr: true,
			wantValW:        "0",
			wantValRW:       "30000",
			wantValRWS:      "100",
		},
		"error_write_error": {
			pin:              "read/write_string",
			simulateWriteErr: true,
			wantValW:         "0",
			wantValRW:        "30000",
			wantValRWS:       "inverted",
			wantErr:          "write error",
		},
		"error_notexist": {
			pin:        "notexist",
			wantValW:   "0",
			wantValRW:  "30000",
			wantValRWS: "inverted",
			wantErr:    "'notexist' is not a valid id of an analog pin",
		},
		"error_write_not_allowed": {
			pin:        "read",
			wantValW:   "0",
			wantValRW:  "30000",
			wantValRWS: "inverted",
			wantErr:    "the pin '/sys/bus/iio/devices/iio:device0/in_voltage0_raw' is not allowed to write (val: 100)",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a, fs := initTestAnalogPinsAdaptorWithMockedFilesystem(analogMockPaths)
			fs.WithWriteError = tc.simulateWriteErr
			fs.WithReadError = tc.simulateReadErr
			// act
			err := a.AnalogWrite(tc.pin, 100)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, "54321", fs.Files[analogReadPath].Contents)
			assert.Equal(t, tc.wantValW, fs.Files[analogWritePath].Contents)
			assert.Equal(t, tc.wantValRW, fs.Files[analogReadWritePath].Contents)
			assert.Equal(t, tc.wantValRWS, fs.Files[analogReadWriteStringPath].Contents)
		})
	}
}

func TestAnalogRead(t *testing.T) {
	tests := map[string]struct {
		pin              string
		simulateReadErr  bool
		simulateWriteErr bool
		wantVal          int
		wantErr          string
	}{
		"read_r_pin": {
			pin:     "read",
			wantVal: 54321,
		},
		"read_rw_pin": {
			pin:     "read/write",
			wantVal: 30000,
		},
		"ok_on_write_error": {
			pin:              "read",
			simulateWriteErr: true,
			wantVal:          54321,
		},
		"error_read_error": {
			pin:             "read",
			simulateReadErr: true,
			wantErr:         "read error",
		},
		"error_notexist": {
			pin:     "notexist",
			wantErr: "'notexist' is not a valid id of an analog pin",
		},
		"error_invalid_syntax": {
			pin:     "read/write_string",
			wantErr: "strconv.Atoi: parsing \"inverted\": invalid syntax",
		},
		"error_read_not_allowed": {
			pin:     "write",
			wantErr: "the pin '/sys/devices/platform/ff680020.pwm/pwm/pwmchip3/export' is not allowed to read",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a, fs := initTestAnalogPinsAdaptorWithMockedFilesystem(analogMockPaths)
			fs.WithReadError = tc.simulateReadErr
			fs.WithWriteError = tc.simulateWriteErr
			// act
			got, err := a.AnalogRead(tc.pin)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantVal, got)
		})
	}
}
