package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isSupportedSysfs(t *testing.T) {
	// arrange
	dpa := sysfsDigitalPinAccess{}
	// act
	got := dpa.isSupported()
	// assert
	assert.True(t, got)
}

func Test_isSupportedGpiod(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_createAsSysfs(t *testing.T) {
	// arrange
	dpa := sysfsDigitalPinAccess{}
	// act
	dp := dpa.createPin("chip", 8)
	// assert
	assert.NotNil(t, dp)
	dps := dp.(*digitalPinSysfs)
	// chip is dropped
	assert.Equal(t, "gpio8", dps.label)
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
	assert.NotNil(t, dp)
	dpg := dp.(*digitalPinGpiod)
	assert.Equal(t, label, dpg.label)
	assert.Equal(t, chip, dpg.chipName)
}

func Test_createPinWithOptionsSysfs(t *testing.T) {
	// This is a general test, that options are applied by using "create" with the WithPinLabel() option.
	// All other configuration options will be tested in tests for "digitalPinConfig".
	//
	// arrange
	const label = "my sysfs label"
	dpa := sysfsDigitalPinAccess{}
	// act
	dp := dpa.createPin("", 9, WithPinLabel(label))
	dps := dp.(*digitalPinSysfs)
	// assert
	assert.Equal(t, label, dps.label)
}

func Test_createPinWithOptionsGpiod(t *testing.T) {
	// This is a general test, that options are applied by using "create" with the WithPinLabel() option.
	// All other configuration options will be tested in tests for "digitalPinConfig".
	//
	// arrange
	const label = "my gpiod label"
	dpa := gpiodDigitalPinAccess{}
	// act
	dp := dpa.createPin("", 19, WithPinLabel(label))
	dpg := dp.(*digitalPinGpiod)
	// assert
	assert.Equal(t, label, dpg.label)
	// test fallback for empty chip
	assert.Equal(t, "gpiochip0", dpg.chipName)
}
