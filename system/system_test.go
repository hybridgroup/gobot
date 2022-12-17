package system

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestNewAccesser(t *testing.T) {
	// act
	a := NewAccesser()
	// assert
	nativeSys := a.sys.(*nativeSyscall)
	nativeFsSys := a.fs.(*nativeFilesystem)
	perphioSpi := a.spiAccess.(*periphioSpiAccess)
	gobottest.Refute(t, a, nil)
	gobottest.Refute(t, nativeSys, nil)
	gobottest.Refute(t, nativeFsSys, nil)
	gobottest.Refute(t, perphioSpi, nil)
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
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, con, nil)
	gobottest.Assert(t, spi.busNum, busNum)
	gobottest.Assert(t, spi.chipNum, chipNum)
	gobottest.Assert(t, spi.mode, mode)
	gobottest.Assert(t, spi.bits, bits)
	gobottest.Assert(t, spi.maxSpeed, maxSpeed)
}

func TestNewAccesser_IsSysfsDigitalPinAccess(t *testing.T) {
	const gpiodTestCase = "accesser_gpiod"
	var tests = map[string]struct {
		accesser string
		wantSys  bool
	}{
		"default_accesser_sysfs": {
			wantSys: true,
		},
		"accesser_sysfs": {
			accesser: "sysfs",
			wantSys:  true,
		},
		gpiodTestCase: {
			accesser: "cdev",
			wantSys:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			if name == gpiodTestCase {
				// there is no mock at this level, so if the system do not support
				// character device gpio, we skip the test
				dpa := &gpiodDigitalPinAccess{fs: &nativeFilesystem{}}
				if !dpa.isSupported() {
					t.Skip()
				}
			}
			// act
			a := NewAccesser(tc.accesser)
			got := a.IsSysfsDigitalPinAccess()
			// assert
			gobottest.Refute(t, a, nil)
			if tc.wantSys {
				gobottest.Assert(t, got, true)
				dpaSys := a.digitalPinAccess.(*sysfsDigitalPinAccess)
				gobottest.Refute(t, dpaSys, nil)
				gobottest.Assert(t, dpaSys.fs, a.fs.(*nativeFilesystem))
			} else {
				gobottest.Assert(t, got, false)
				dpaGpiod := a.digitalPinAccess.(*gpiodDigitalPinAccess)
				gobottest.Refute(t, dpaGpiod, nil)
				gobottest.Assert(t, dpaGpiod.fs, a.fs.(*nativeFilesystem))
			}
		})
	}
}
