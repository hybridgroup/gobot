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
	gpiodDigitalPin := a.digitalPinAccess.(*cdevDigitalPinAccess)
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
		wantCdev      bool
	}{
		"default_accesser_gpiod": {
			wantCdev: true,
		},
		"accesser_sysfs": {
			sysfsAccesser: true,
			wantCdev:      false,
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
			gotCdev := a.IsCdevDigitalPinAccess()
			gotSysfs := a.IsSysfsDigitalPinAccess()
			// assert
			assert.NotNil(t, a)
			if tc.wantCdev {
				assert.True(t, gotCdev)
				assert.False(t, gotSysfs)
				dpaGpioCdev := a.digitalPinAccess.(*cdevDigitalPinAccess)
				assert.NotNil(t, dpaGpioCdev)
				assert.Equal(t, a.fs.(*nativeFilesystem), dpaGpioCdev.fs)
			} else {
				assert.False(t, gotCdev)
				assert.True(t, gotSysfs)
				dpaSys := a.digitalPinAccess.(*sysfsDigitalPinAccess)
				assert.NotNil(t, dpaSys)
				assert.Equal(t, a.fs.(*nativeFilesystem), dpaSys.sfa.fs)
			}
		})
	}
}
