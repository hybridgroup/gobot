package rockpi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2/system"
)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	_ = a.Connect()
	return a, fs
}

func TestDefaultI2cBus(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem([]string{})
	assert.Equal(t, 7, a.DefaultI2cBus())
}

func Test_getPinTranslatorFunction(t *testing.T) {
	cases := map[string]struct {
		pin          string
		model        string
		expectedLine int
		expectedErr  error
	}{
		"Rock Pi 4 specific pin": {
			pin:          "12",
			model:        "Radxa ROCK 4",
			expectedLine: 131,
			expectedErr:  nil,
		},
		"Rock Pi 4C+ specific pin": {
			pin:          "12",
			model:        "Radxa ROCK 4C+",
			expectedLine: 91,
			expectedErr:  nil,
		},
		"Generic pin": {
			pin:          "3",
			model:        "whatever",
			expectedLine: 71,
			expectedErr:  nil,
		},
		"Not a valid pin": {
			pin:          "666",
			model:        "whatever",
			expectedLine: 0,
			expectedErr:  fmt.Errorf("Not a valid pin"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			fn := a.getPinTranslatorFunction()
			fs := a.sys.UseMockFilesystem([]string{procDeviceTreeModel})
			fs.Files[procDeviceTreeModel].Contents = tc.model
			// act
			chip, line, err := fn(tc.pin)
			// assert
			assert.Equal(t, "", chip)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedLine, line)
		})
	}
}
