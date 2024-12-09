package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2/system"
)

func TestNewDigitalPinTranslator(t *testing.T) {
	// arrange
	sys := &system.Accesser{}
	pinDef := DigitalPinDefinitions{}
	// act
	pt := NewDigitalPinTranslator(sys, pinDef)
	// assert
	assert.IsType(t, &DigitalPinTranslator{}, pt)
	assert.Equal(t, sys, pt.sys)
	assert.Equal(t, pinDef, pt.pinDefinitions)
}

func TestDigitalPinTranslatorTranslate(t *testing.T) {
	pinDefinitions := DigitalPinDefinitions{
		"7":  {Sysfs: 17, Cdev: CdevPin{Chip: 0, Line: 17}},
		"22": {Sysfs: 171, Cdev: CdevPin{Chip: 5, Line: 19}},
		"5":  {Sysfs: 253, Cdev: CdevPin{Chip: 8, Line: 5}},
	}
	tests := map[string]struct {
		access   string
		pin      string
		wantChip string
		wantLine int
		wantErr  error
	}{
		"cdev_ok_7": {
			access:   "cdev",
			pin:      "7",
			wantChip: "gpiochip0",
			wantLine: 17,
		},
		"cdev_ok_22": {
			access:   "cdev",
			pin:      "22",
			wantChip: "gpiochip5",
			wantLine: 19,
		},
		"cdev_ok_5": {
			access:   "cdev",
			pin:      "5",
			wantChip: "gpiochip8",
			wantLine: 5,
		},
		"sysfs_ok_7": {
			access:   "sysfs",
			pin:      "7",
			wantChip: "",
			wantLine: 17,
		},
		"sysfs_ok_22": {
			access:   "sysfs",
			pin:      "22",
			wantChip: "",
			wantLine: 171,
		},
		"sysfs_ok_5": {
			access:   "sysfs",
			pin:      "5",
			wantChip: "",
			wantLine: 253,
		},
		"unknown_pin": {
			pin:      "99",
			wantChip: "",
			wantLine: -1,
			wantErr:  fmt.Errorf("'99' is not a valid id for a digital pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			sys := system.NewAccesser(system.WithDigitalPinGpiodAccess())
			sys.UseDigitalPinAccessWithMockFs(tc.access, []string{})
			pt := NewDigitalPinTranslator(sys, pinDefinitions)
			// act
			chip, line, err := pt.Translate(tc.pin)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantChip, chip)
			assert.Equal(t, tc.wantLine, line)
		})
	}
}
