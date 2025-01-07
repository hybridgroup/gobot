//nolint:forcetypeassert // ok here
package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithDigitalPinOptionAccess_HasDigitalPinCdevAccess_HasDigitalPinSysfsAccess(t *testing.T) {
	// note: default accesser already tested on tests with AddDigitalPinSupport()
	tests := map[string]struct {
		initialAccesser digitalPinAccesser
		option          AccesserOptionApplier
		wantCdev        bool
	}{
		"with_cdev": {
			initialAccesser: &sysfsDigitalPinAccess{},
			option:          WithDigitalPinCdevAccess(),
			wantCdev:        true,
		},
		"with_sysfs": {
			initialAccesser: &cdevDigitalPinAccess{},
			option:          WithDigitalPinSysfsAccess(),
			wantCdev:        false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAccesser()
			a.digitalPinAccess = tc.initialAccesser
			a.AddDigitalPinSupport(tc.option)
			// act
			gotCdev := a.HasDigitalPinCdevAccess()
			gotSysfs := a.HasDigitalPinSysfsAccess()
			// assert
			if tc.wantCdev {
				assert.True(t, gotCdev)
				assert.False(t, gotSysfs)
				assert.IsType(t, &cdevDigitalPinAccess{}, a.digitalPinAccess)
				assert.IsType(t, &nativeFilesystem{}, a.digitalPinAccess.(*cdevDigitalPinAccess).fs)
			} else {
				assert.False(t, gotCdev)
				assert.True(t, gotSysfs)
				assert.IsType(t, &sysfsDigitalPinAccess{}, a.digitalPinAccess)
				assert.IsType(t, &nativeFilesystem{}, a.digitalPinAccess.(*sysfsDigitalPinAccess).sfa.fs)
			}
		})
	}
}

func TestWithSpiOptionAccess_HasSpiPeriphioAccess_HasSpiGpioAccess(t *testing.T) {
	// note: default accesser already tested on tests with AddSPISupport()
	tests := map[string]struct {
		initialAccesser spiAccesser
		option          AccesserOptionApplier
		wantPeriphio    bool
	}{
		"with_periphio": {
			initialAccesser: &gpioSpiAccess{},
			wantPeriphio:    true,
		},
		"withr_sysfs": {
			initialAccesser: &periphioSpiAccess{},
			option:          WithSpiGpioAccess(nil, "2", "3", "4", "5"),
			wantPeriphio:    false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAccesser()
			a.spiAccess = tc.initialAccesser
			// arrange for periphio needed
			a.UseMockFilesystem([]string{"/dev/spidev"})
			a.AddDigitalPinSupport()
			a.AddSPISupport(tc.option)
			// act
			gotPeriphio := a.HasSpiPeriphioAccess()
			gotGpio := a.HasSpiGpioAccess()
			// assert
			require.NotNil(t, a.spiAccess)
			if tc.wantPeriphio {
				assert.True(t, gotPeriphio)
				assert.False(t, gotGpio)
				assert.IsType(t, &periphioSpiAccess{}, a.spiAccess)
			} else {
				assert.False(t, gotPeriphio)
				assert.True(t, gotGpio)
				assert.IsType(t, &gpioSpiAccess{}, a.spiAccess)
			}
		})
	}
}
