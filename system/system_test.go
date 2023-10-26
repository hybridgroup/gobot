package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccesser(t *testing.T) {
	// act
	a := NewAccesser()
	// assert
	nativeSys := a.sys.(*nativeSyscall)
	nativeFsSys := a.fs.(*nativeFilesystem)
	perphioSpi := a.spiAccess.(*periphioSpiAccess)
	assert.NotNil(t, a)
	assert.NotNil(t, nativeSys)
	assert.NotNil(t, nativeFsSys)
	assert.NotNil(t, perphioSpi)
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
	assert.NoError(t, err)
	assert.NotNil(t, con)
	assert.Equal(t, busNum, spi.busNum)
	assert.Equal(t, chipNum, spi.chipNum)
	assert.Equal(t, mode, spi.mode)
	assert.Equal(t, bits, spi.bits)
	assert.Equal(t, maxSpeed, spi.maxSpeed)
}

func TestNewAccesser_IsSysfsDigitalPinAccess(t *testing.T) {
	tests := map[string]struct {
		gpiodAccesser bool
		wantSys       bool
	}{
		"default_accesser_sysfs": {
			wantSys: true,
		},
		"accesser_sysfs": {
			wantSys: true,
		},
		"accesser_gpiod": {
			gpiodAccesser: true,
			wantSys:       false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAccesser()
			if tc.gpiodAccesser {
				// there is no mock at this level, so if the system do not support
				// character device gpio, we skip the test
				dpa := &gpiodDigitalPinAccess{fs: &nativeFilesystem{}}
				if !dpa.isSupported() {
					t.Skip()
				}
				WithDigitalPinGpiodAccess()(a)
			}
			// act
			got := a.IsSysfsDigitalPinAccess()
			// assert
			assert.NotNil(t, a)
			if tc.wantSys {
				assert.True(t, got)
				dpaSys := a.digitalPinAccess.(*sysfsDigitalPinAccess)
				assert.NotNil(t, dpaSys)
				assert.Equal(t, a.fs.(*nativeFilesystem), dpaSys.fs)
			} else {
				assert.False(t, got)
				dpaGpiod := a.digitalPinAccess.(*gpiodDigitalPinAccess)
				assert.NotNil(t, dpaGpiod)
				assert.Equal(t, a.fs.(*nativeFilesystem), dpaGpiod.fs)
			}
		})
	}
}
