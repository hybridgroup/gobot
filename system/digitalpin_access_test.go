package system

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func Test_isSupportedSysfs(t *testing.T) {
	// arrange
	dpa := sysfsDigitalPinAccess{}
	// act
	got := dpa.isSupported()
	// assert
	gobottest.Assert(t, got, true)
}

func Test_isSupportedGpiod(t *testing.T) {
	var tests = map[string]struct {
		mockPaths []string
		want      bool
	}{
		"supported": {
			mockPaths: []string{"/sys/class/gpio/", "/dev/gpiochip3"},
			want:      true,
		},
		"not_supported": {
			mockPaths: []string{"/sys/class/gpio/"},
			want:      false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			fs := newMockFilesystem(tc.mockPaths)
			dpa := gpiodDigitalPinAccess{fs: fs}
			// act
			got := dpa.isSupported()
			// assert
			gobottest.Assert(t, got, tc.want)
		})
	}
}

func Test_createAsSysfs(t *testing.T) {
	// arrange
	dpa := sysfsDigitalPinAccess{}
	// act
	dp := dpa.createPin("chip", 8)
	// assert
	gobottest.Refute(t, dp, nil)
	dps := dp.(*digitalPinSysfs)
	// chip is dropped
	gobottest.Assert(t, dps.label, "gpio8")
}

func Test_createAsGpiod(t *testing.T) {
	// arrange
	const (
		pin   = 18
		label = "gobotio18"
		chip  = "gpiochip1"
	)
	dpa := gpiodDigitalPinAccess{}
	// act
	dp := dpa.createPin(chip, 18)
	// assert
	gobottest.Refute(t, dp, nil)
	dpg := dp.(*digitalPinGpiod)
	gobottest.Assert(t, dpg.label, label)
	gobottest.Assert(t, dpg.chipName, chip)
}

func Test_createPinWithOptionsSysfs(t *testing.T) {
	// This is a general test, that options are applied by using "create" with the WithLabel() option.
	// All other configuration options will be tested in tests for "digitalPinConfig".
	//
	// arrange
	const label = "my sysfs label"
	dpa := sysfsDigitalPinAccess{}
	// act
	dp := dpa.createPin("", 9, WithLabel(label))
	dps := dp.(*digitalPinSysfs)
	// assert
	gobottest.Assert(t, dps.label, label)
}

func Test_createPinWithOptionsGpiod(t *testing.T) {
	// This is a general test, that options are applied by using "create" with the WithLabel() option.
	// All other configuration options will be tested in tests for "digitalPinConfig".
	//
	// arrange
	const label = "my gpiod label"
	dpa := gpiodDigitalPinAccess{}
	// act
	dp := dpa.createPin("", 19, WithLabel(label))
	dpg := dp.(*digitalPinGpiod)
	// assert
	gobottest.Assert(t, dpg.label, label)
	// test fallback for empty chip
	gobottest.Assert(t, dpg.chipName, "gpiochip0")
}
