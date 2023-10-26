package system

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/v2"
)

var (
	_ gobot.DigitalPinner           = (*digitalPinSysfs)(nil)
	_ gobot.DigitalPinValuer        = (*digitalPinSysfs)(nil)
	_ gobot.DigitalPinOptioner      = (*digitalPinSysfs)(nil)
	_ gobot.DigitalPinOptionApplier = (*digitalPinSysfs)(nil)
)

func initTestDigitalPinSysfsWithMockedFilesystem(mockPaths []string) (*digitalPinSysfs, *MockFilesystem) {
	fs := newMockFilesystem(mockPaths)
	pin := newDigitalPinSysfs(fs, "10")
	return pin, fs
}

func Test_newDigitalPinSysfs(t *testing.T) {
	// arrange
	m := &MockFilesystem{}
	const pinID = "1"
	// act
	pin := newDigitalPinSysfs(m, pinID, WithPinOpenDrain())
	// assert
	assert.Equal(t, pinID, pin.pin)
	assert.Equal(t, m, pin.fs)
	assert.Equal(t, "gpio"+pinID, pin.label)
	assert.Equal(t, "in", pin.direction)
	assert.Equal(t, 1, pin.drive)
}

func TestApplyOptionsSysfs(t *testing.T) {
	tests := map[string]struct {
		changed    []bool
		simErr     bool
		wantExport string
		wantErr    string
	}{
		"both_changed": {
			changed:    []bool{true, true},
			wantExport: "10",
		},
		"first_changed": {
			changed:    []bool{true, false},
			wantExport: "10",
		},
		"second_changed": {
			changed:    []bool{false, true},
			wantExport: "10",
		},
		"none_changed": {
			changed:    []bool{false, false},
			wantExport: "",
		},
		"error_on_change": {
			changed:    []bool{false, true},
			simErr:     true,
			wantExport: "10",
			wantErr:    "gpio10/direction: no such file",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			mockPaths := []string{
				"/sys/class/gpio/export",
				"/sys/class/gpio/gpio10/value",
			}
			if !tc.simErr {
				mockPaths = append(mockPaths, "/sys/class/gpio/gpio10/direction")
			}
			pin, fs := initTestDigitalPinSysfsWithMockedFilesystem(mockPaths)

			optionFunction1 := func(gobot.DigitalPinOptioner) bool {
				pin.digitalPinConfig.direction = OUT
				return tc.changed[0]
			}
			optionFunction2 := func(gobot.DigitalPinOptioner) bool {
				pin.digitalPinConfig.drive = 15
				return tc.changed[1]
			}
			// act
			err := pin.ApplyOptions(optionFunction1, optionFunction2)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, OUT, pin.digitalPinConfig.direction)
			assert.Equal(t, 15, pin.digitalPinConfig.drive)
			// marker for call of reconfigure, correct reconfigure is tested independently
			assert.Equal(t, tc.wantExport, fs.Files["/sys/class/gpio/export"].Contents)
		})
	}
}

func TestDirectionBehaviorSysfs(t *testing.T) {
	// arrange
	pin := newDigitalPinSysfs(nil, "1")
	require.Equal(t, "in", pin.direction)
	pin.direction = "test"
	// act && assert
	assert.Equal(t, "test", pin.DirectionBehavior())
}

func TestDigitalPinExportSysfs(t *testing.T) {
	// this tests mainly the function reconfigure()
	const (
		exportPath   = "/sys/class/gpio/export"
		dirPath      = "/sys/class/gpio/gpio10/direction"
		valuePath    = "/sys/class/gpio/gpio10/value"
		inversePath  = "/sys/class/gpio/gpio10/active_low"
		unexportPath = "/sys/class/gpio/unexport"
	)
	allMockPaths := []string{exportPath, dirPath, valuePath, inversePath, unexportPath}
	tests := map[string]struct {
		mockPaths             []string
		changeDirection       string
		changeOutInitialState int
		changeActiveLow       bool
		changeBias            int
		changeDrive           int
		changeDebouncePeriod  time.Duration
		changeEdge            int
		changePollInterval    time.Duration
		simEbusyOnWrite       int
		wantWrites            int
		wantExport            string
		wantUnexport          string
		wantDirection         string
		wantValue             string
		wantInverse           string
		wantErr               string
	}{
		"ok_without_option": {
			mockPaths:     allMockPaths,
			wantWrites:    2,
			wantExport:    "10",
			wantDirection: "in",
		},
		"ok_input_bias_dropped": {
			mockPaths:     allMockPaths,
			changeBias:    3,
			wantWrites:    2,
			wantExport:    "10",
			wantDirection: "in",
		},
		"ok_input_drive_dropped": {
			mockPaths:     allMockPaths,
			changeDrive:   2,
			wantWrites:    2,
			wantExport:    "10",
			wantDirection: "in",
		},
		"ok_input_debounce_dropped": {
			mockPaths:            allMockPaths,
			changeDebouncePeriod: 2 * time.Second,
			wantWrites:           2,
			wantExport:           "10",
			wantDirection:        "in",
		},
		"ok_input_inverse": {
			mockPaths:       allMockPaths,
			changeActiveLow: true,
			wantWrites:      3,
			wantExport:      "10",
			wantDirection:   "in",
			wantInverse:     "1",
		},
		"ok_output": {
			mockPaths:             allMockPaths,
			changeDirection:       "out",
			changeOutInitialState: 4,
			wantWrites:            3,
			wantExport:            "10",
			wantDirection:         "out",
			wantValue:             "4",
		},
		"ok_output_bias_dropped": {
			mockPaths:       allMockPaths,
			changeDirection: "out",
			changeBias:      3,
			wantWrites:      3,
			wantExport:      "10",
			wantDirection:   "out",
			wantValue:       "0",
		},
		"ok_output_drive_dropped": {
			mockPaths:       allMockPaths,
			changeDirection: "out",
			changeDrive:     2,
			wantWrites:      3,
			wantExport:      "10",
			wantDirection:   "out",
			wantValue:       "0",
		},
		"ok_output_debounce_dropped": {
			mockPaths:            allMockPaths,
			changeDirection:      "out",
			changeDebouncePeriod: 2 * time.Second,
			wantWrites:           3,
			wantExport:           "10",
			wantDirection:        "out",
			wantValue:            "0",
		},
		"ok_output_inverse": {
			mockPaths:       allMockPaths,
			changeDirection: "out",
			changeActiveLow: true,
			wantWrites:      4,
			wantExport:      "10",
			wantDirection:   "out",
			wantInverse:     "1",
			wantValue:       "0",
		},
		"ok_already_exported": {
			mockPaths:       allMockPaths,
			wantWrites:      2,
			wantExport:      "10",
			wantDirection:   "in",
			simEbusyOnWrite: 1, // just means "already exported"
		},
		"error_no_eventhandler_for_polling": { // this only tests the call of function, all other is tested separately
			mockPaths:          allMockPaths,
			changePollInterval: 3 * time.Second,
			wantWrites:         3,
			wantUnexport:       "10",
			wantDirection:      "in",
			wantErr:            "event handler is mandatory",
		},
		"error_no_export_file": {
			mockPaths: []string{unexportPath},
			wantErr:   "/export: no such file",
		},
		"error_no_direction_file": {
			mockPaths:    []string{exportPath, unexportPath},
			wantWrites:   2,
			wantUnexport: "10",
			wantErr:      "gpio10/direction: no such file",
		},
		"error_write_direction_file": {
			mockPaths:       allMockPaths,
			wantWrites:      3,
			wantUnexport:    "10",
			simEbusyOnWrite: 2,
			wantErr:         "device or resource busy",
		},
		"error_no_value_file": {
			mockPaths:    []string{exportPath, dirPath, unexportPath},
			wantWrites:   2,
			wantUnexport: "10",
			wantErr:      "gpio10/value: no such file",
		},
		"error_no_inverse_file": {
			mockPaths:       []string{exportPath, dirPath, valuePath, unexportPath},
			changeActiveLow: true,
			wantWrites:      3,
			wantUnexport:    "10",
			wantErr:         "gpio10/active_low: no such file",
		},
		"error_input_edge_without_poll": {
			mockPaths:    allMockPaths,
			changeEdge:   2,
			wantWrites:   3,
			wantUnexport: "10",
			wantErr:      "not implemented for sysfs without discrete polling",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			fs := newMockFilesystem(tc.mockPaths)
			pin := newDigitalPinSysfs(fs, "10")
			if tc.changeDirection != "" {
				pin.direction = tc.changeDirection
			}
			if tc.changeOutInitialState != 0 {
				pin.outInitialState = tc.changeOutInitialState
			}
			if tc.changeActiveLow {
				pin.activeLow = tc.changeActiveLow
			}
			if tc.changeBias != 0 {
				pin.bias = tc.changeBias
			}
			if tc.changeDrive != 0 {
				pin.drive = tc.changeDrive
			}
			if tc.changeDebouncePeriod != 0 {
				pin.debouncePeriod = tc.changeDebouncePeriod
			}
			if tc.changeEdge != 0 {
				pin.edge = tc.changeEdge
			}
			if tc.changePollInterval != 0 {
				pin.pollInterval = tc.changePollInterval
			}
			// arrange write function
			oldWriteFunc := writeFile
			numCallsWrite := 0
			writeFile = func(f File, data []byte) error {
				numCallsWrite++
				require.NoError(t, oldWriteFunc(f, data))
				if numCallsWrite == tc.simEbusyOnWrite {
					return &os.PathError{Err: Syscall_EBUSY}
				}
				return nil
			}
			defer func() { writeFile = oldWriteFunc }()
			// act
			err := pin.Export()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pin.valFile)
				assert.NotNil(t, pin.dirFile)
				assert.Equal(t, tc.wantDirection, fs.Files[dirPath].Contents)
				assert.Equal(t, tc.wantExport, fs.Files[exportPath].Contents)
				assert.Equal(t, tc.wantValue, fs.Files[valuePath].Contents)
				assert.Equal(t, tc.wantInverse, fs.Files[inversePath].Contents)
			}
			assert.Equal(t, tc.wantUnexport, fs.Files[unexportPath].Contents)
			assert.Equal(t, tc.wantWrites, numCallsWrite)
		})
	}
}

func TestDigitalPinSysfs(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	}
	pin, fs := initTestDigitalPinSysfsWithMockedFilesystem(mockPaths)

	assert.Equal(t, "10", pin.pin)
	assert.Equal(t, "gpio10", pin.label)
	assert.Nil(t, pin.valFile)

	err := pin.Unexport()
	assert.NoError(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/gpio/unexport"].Contents)

	require.NoError(t, pin.Export())

	err = pin.Write(1)
	assert.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio10/value"].Contents)

	err = pin.ApplyOptions(WithPinDirectionInput())
	assert.NoError(t, err)
	assert.Equal(t, "in", fs.Files["/sys/class/gpio/gpio10/direction"].Contents)

	data, _ := pin.Read()
	assert.Equal(t, data, 1)

	pin2 := newDigitalPinSysfs(fs, "30")
	err = pin2.Write(1)
	assert.ErrorContains(t, err, "pin has not been exported")

	data, err = pin2.Read()
	assert.ErrorContains(t, err, "pin has not been exported")
	assert.Equal(t, 0, data)

	writeFile = func(File, []byte) error {
		return &os.PathError{Err: Syscall_EINVAL}
	}

	err = pin.Unexport()
	assert.NoError(t, err)

	writeFile = func(File, []byte) error {
		return &os.PathError{Err: errors.New("write error")}
	}

	err = pin.Unexport()
	assert.ErrorContains(t, err.(*os.PathError).Err, "write error")
}

func TestDigitalPinUnexportErrorSysfs(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/unexport",
	}
	pin, _ := initTestDigitalPinSysfsWithMockedFilesystem(mockPaths)

	writeFile = func(File, []byte) error {
		return &os.PathError{Err: Syscall_EBUSY}
	}

	err := pin.Unexport()
	assert.ErrorContains(t, err, " : device or resource busy")
}
