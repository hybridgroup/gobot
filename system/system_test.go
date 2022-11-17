package system

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

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
			nativeSys := a.sys.(*nativeSyscall)
			nativeFsSys := a.fs.(*nativeFilesystem)
			gobottest.Refute(t, a, nil)
			gobottest.Refute(t, nativeSys, nil)
			gobottest.Refute(t, nativeFsSys, nil)
			if tc.wantSys {
				gobottest.Assert(t, got, true)
				dpaSys := a.digitalPinAccess.(*sysfsDigitalPinAccess)
				gobottest.Refute(t, dpaSys, nil)
				gobottest.Assert(t, dpaSys.fs, nativeFsSys)
			} else {
				gobottest.Assert(t, got, false)
				dpaGpiod := a.digitalPinAccess.(*gpiodDigitalPinAccess)
				gobottest.Refute(t, dpaGpiod, nil)
				gobottest.Assert(t, dpaGpiod.fs, nativeFsSys)
			}
		})
	}
}
