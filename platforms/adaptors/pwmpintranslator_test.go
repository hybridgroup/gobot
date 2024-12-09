package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2/system"
)

func TestNewPWMPinTranslator(t *testing.T) {
	// arrange
	sys := &system.Accesser{}
	pinDef := PWMPinDefinitions{}
	// act
	pt := NewPWMPinTranslator(sys, pinDef)
	// assert
	assert.IsType(t, &PWMPinTranslator{}, pt)
	assert.Equal(t, sys, pt.sys)
	assert.Equal(t, pinDef, pt.pinDefinitions)
}

func TestPWMPinTranslatorTranslate(t *testing.T) {
	pinDefinitions := PWMPinDefinitions{
		"33": {Dir: "/sys/devices/platform/ff680020.pwm/pwm/", DirRegexp: "pwmchip[0|1|2]$", Channel: 0},
		"32": {Dir: "/sys/devices/platform/ff680030.pwm/pwm/", DirRegexp: "pwmchip[0|1|2|3]$", Channel: 0},
	}
	basePaths := []string{
		"/sys/devices/platform/ff680020.pwm/pwm/",
		"/sys/devices/platform/ff680030.pwm/pwm/",
	}
	tests := map[string]struct {
		pin         string
		chip        string
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"32_chip0": {
			pin:         "32",
			chip:        "pwmchip0",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip0",
			wantChannel: 0,
		},
		"32_chip1": {
			pin:         "32",
			chip:        "pwmchip1",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip1",
			wantChannel: 0,
		},
		"32_chip2": {
			pin:         "32",
			chip:        "pwmchip2",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip2",
			wantChannel: 0,
		},
		"32_chip3": {
			pin:         "32",
			chip:        "pwmchip3",
			wantDir:     "/sys/devices/platform/ff680030.pwm/pwm/pwmchip3",
			wantChannel: 0,
		},
		"33_chip0": {
			pin:         "33",
			chip:        "pwmchip0",
			wantDir:     "/sys/devices/platform/ff680020.pwm/pwm/pwmchip0",
			wantChannel: 0,
		},
		"33_chip1": {
			pin:         "33",
			chip:        "pwmchip1",
			wantDir:     "/sys/devices/platform/ff680020.pwm/pwm/pwmchip1",
			wantChannel: 0,
		},
		"33_chip2": {
			pin:         "33",
			chip:        "pwmchip2",
			wantDir:     "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2",
			wantChannel: 0,
		},
		"invalid_pin": {
			pin:         "7",
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'7' is not a valid id for a PWM pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			mockedPaths := []string{}
			for _, base := range basePaths {
				mockedPaths = append(mockedPaths, base+tc.chip+"/")
			}
			sys := system.NewAccesser()
			_ = sys.UseMockFilesystem(mockedPaths)
			pt := NewPWMPinTranslator(sys, pinDefinitions)
			// act
			dir, channel, err := pt.Translate(tc.pin)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, dir)
			assert.Equal(t, tc.wantChannel, channel)
		})
	}
}
