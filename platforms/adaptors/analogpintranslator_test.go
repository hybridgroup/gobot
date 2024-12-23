package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/system"
)

func TestNewAnalogPinTranslator(t *testing.T) {
	// arrange
	sys := &system.Accesser{}
	pinDef := AnalogPinDefinitions{}
	// act
	pt := NewAnalogPinTranslator(sys, pinDef)
	// assert
	assert.IsType(t, &AnalogPinTranslator{}, pt)
	assert.Equal(t, sys, pt.sys)
	assert.Equal(t, pinDef, pt.pinDefinitions)
}

func TestAnalogPinTranslatorTranslate(t *testing.T) {
	pinDefinitions := AnalogPinDefinitions{
		"thermal_zone0": {Path: "/sys/class/thermal/thermal_zone0/temp", W: false, ReadBufLen: 7},
		"thermal_zone1": {Path: "/sys/class/thermal/thermal_zone1/temp", W: false, ReadBufLen: 7},
	}
	mockedPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
	}
	tests := map[string]struct {
		id         string
		wantPath   string
		wantBufLen uint16
		wantErr    string
	}{
		"translate_thermal_zone0": {
			id:         "thermal_zone0",
			wantPath:   "/sys/class/thermal/thermal_zone0/temp",
			wantBufLen: 7,
		},
		"translate_thermal_zone1": {
			id:         "thermal_zone1",
			wantPath:   "/sys/class/thermal/thermal_zone1/temp",
			wantBufLen: 7,
		},
		"unknown_id": {
			id:      "99",
			wantErr: "'99' is not a valid id for an analog pin",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			sys := system.NewAccesser()
			_ = sys.UseMockFilesystem(mockedPaths)
			pt := NewAnalogPinTranslator(sys, pinDefinitions)
			// act
			path, w, buf, err := pt.Translate(tc.id)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantPath, path)
			assert.False(t, w)
			assert.Equal(t, tc.wantBufLen, buf)
		})
	}
}
