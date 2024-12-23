//nolint:forcetypeassert // ok here
package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccesser(t *testing.T) {
	// act
	a := NewAccesser()
	// assert
	nativeSys := a.sys.(*nativeSyscall)
	nativeFsSys := a.fs.(*nativeFilesystem)
	perphioSpi := a.spiAccess.(*periphioSpiAccess)
	gpiodDigitalPin := a.digitalPinAccess.(*gpiodDigitalPinAccess)
	assert.NotNil(t, a)
	assert.NotNil(t, nativeSys)
	assert.NotNil(t, nativeFsSys)
	assert.NotNil(t, perphioSpi)
	assert.NotNil(t, gpiodDigitalPin)
}

func TestNewAccesser_NewSpiDevice(t *testing.T) {
	// arrange

	const (
		busNum   = 15
		chipNum  = 14
		mode     = 13
		bits     = 12
		maxSpeed = int64(11)
	)
	a := NewAccesser()
	spi := a.UseMockSpi()
	// act
	con, err := a.NewSpiDevice(busNum, chipNum, mode, bits, maxSpeed)
	// assert
	require.NoError(t, err)
	assert.NotNil(t, con)
	assert.Equal(t, busNum, spi.busNum)
	assert.Equal(t, chipNum, spi.chipNum)
	assert.Equal(t, mode, spi.mode)
	assert.Equal(t, bits, spi.bits)
	assert.Equal(t, maxSpeed, spi.maxSpeed)
}

func TestNewAccesser_IsSysfsDigitalPinAccess(t *testing.T) {
	tests := map[string]struct {
		sysfsAccesser bool
		wantGpiod     bool
	}{
		"default_accesser_gpiod": {
			wantGpiod: true,
		},
		"accesser_sysfs": {
			sysfsAccesser: true,
			wantGpiod:     false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAccesser()
			if tc.sysfsAccesser {
				WithDigitalPinSysfsAccess()(a)
			}
			// act
			gotGpiod := a.IsGpiodDigitalPinAccess()
			gotSysfs := a.IsSysfsDigitalPinAccess()
			// assert
			assert.NotNil(t, a)
			if tc.wantGpiod {
				assert.True(t, gotGpiod)
				assert.False(t, gotSysfs)
				dpaGpiod := a.digitalPinAccess.(*gpiodDigitalPinAccess)
				assert.NotNil(t, dpaGpiod)
				assert.Equal(t, a.fs.(*nativeFilesystem), dpaGpiod.fs)
			} else {
				assert.False(t, gotGpiod)
				assert.True(t, gotSysfs)
				dpaSys := a.digitalPinAccess.(*sysfsDigitalPinAccess)
				assert.NotNil(t, dpaSys)
				assert.Equal(t, a.fs.(*nativeFilesystem), dpaSys.sfa.fs)
			}
		})
	}
}
