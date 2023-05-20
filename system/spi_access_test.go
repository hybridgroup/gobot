package system

import (
	"testing"

	"gobot.io/x/gobot/v2/gobottest"
)

func TestGpioSpi_isSupported(t *testing.T) {
	// arrange
	gsa := gpioSpiAccess{}
	// act
	got := gsa.isSupported()
	// assert
	gobottest.Assert(t, got, true)
}

func TestPeriphioSpi_isSupported(t *testing.T) {
	var tests = map[string]struct {
		mockPaths []string
		want      bool
	}{
		"supported": {
			mockPaths: []string{"/dev/spidev0.0", "/dev/spidev1.0"},
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
			psa := periphioSpiAccess{fs: fs}
			// act
			got := psa.isSupported()
			// assert
			gobottest.Assert(t, got, tc.want)
		})
	}
}
