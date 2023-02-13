package rockpi

import (
	"fmt"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
	"testing"
)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	a.Connect()
	return a, fs
}

func TestDefaultI2cBus(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem([]string{})
	gobottest.Assert(t, a.DefaultI2cBus(), 7)
}

func Test_getPinTranslatorFunction(t *testing.T) {
	var cases = map[string]struct {
		pin          string
		model        string
		expectedLine int
		expectedErr  error
	}{
		"Rock Pi 4 specific pin": {
			pin:          "12",
			model:        "Radxa ROCK 4",
			expectedLine: 131,
			expectedErr:  nil,
		},
		"Rock Pi 4C+ specific pin": {
			pin:          "12",
			model:        "Radxa ROCK 4C+",
			expectedLine: 91,
			expectedErr:  nil,
		},
		"Generic pin": {
			pin:          "3",
			model:        "whatever",
			expectedLine: 71,
			expectedErr:  nil,
		},
		"Not a valid pin": {
			pin:          "666",
			model:        "whatever",
			expectedLine: 0,
			expectedErr:  fmt.Errorf("Not a valid pin"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			fn := a.getPinTranslatorFunction()
			fs := a.sys.UseMockFilesystem([]string{procDeviceTreeModel})
			fs.Files[procDeviceTreeModel].Contents = tc.model
			// act
			_, line, err := fn(tc.pin)
			// assert
			gobottest.Assert(t, err, tc.expectedErr)
			gobottest.Assert(t, line, tc.expectedLine)
		})
	}
}

func Test_validateSpiBusNumber(t *testing.T) {
	var cases = map[string]struct {
		busNr       int
		expectedErr error
	}{
		"number_1_ok": {
			busNr: 2,
		},
		"number_2_ok": {
			busNr: 2,
		},
		"number_0_not_ok": {
			busNr:       0,
			expectedErr: fmt.Errorf("SPI Bus number 0 invalid: only 1, 2 supported by current Rockchip."),
		},
		"number_6_not_ok": {
			busNr:       6,
			expectedErr: fmt.Errorf("SPI Bus number 6 invalid: only 1, 2 supported by current Rockchip."),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateSpiBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.expectedErr)
		})
	}
}

func Test_validateI2cBusNumber(t *testing.T) {
	var cases = map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("I2C Bus number -1 invalid: only 2, 6, 7 supported by current Rockchip."),
		},
		"number_2_ok": {
			busNr: 2,
		},
		"number_6_ok": {
			busNr: 6,
		},
		"number_7_ok": {
			busNr: 7,
		},
		"number_1_not_ok": {
			busNr:   1,
			wantErr: fmt.Errorf("I2C Bus number 1 invalid: only 2, 6, 7 supported by current Rockchip."),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateI2cBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}
