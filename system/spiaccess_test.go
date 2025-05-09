package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGpioSpi_isType(t *testing.T) {
	// arrange
	gsa := gpioSpiAccess{}
	// act & assert GPIO
	assert.True(t, gsa.isType(spiBusAccesserTypeGPIO))
	// act & assert Periphio
	assert.False(t, gsa.isType(spiBusAccesserTypePeriphio))
}

func TestPeriphioSpi_isType(t *testing.T) {
	// arrange
	psa := periphioSpiAccess{}
	// act & assert Periphio
	assert.True(t, psa.isType(spiBusAccesserTypePeriphio))
	// act & assert GPIO
	assert.False(t, psa.isType(spiBusAccesserTypeGPIO))
}

func TestGpioSpi_isSupported(t *testing.T) {
	// arrange
	gsa := gpioSpiAccess{}
	// act
	got := gsa.isSupported()
	// assert
	assert.True(t, got)
}

func TestPeriphioSpi_isSupported(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.want, got)
		})
	}
}
