package gpio

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	stepAngle   = 0.5 // use non int step angle to check int math
	stepsPerRev = 720
)

func initTestEasyDriverWithStubbedAdaptor() (*EasyDriver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	d := NewEasyDriver(a, stepAngle, "1", "2", "3", "4")
	return d, a
}

func TestNewEasyDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewEasyDriver(a, stepAngle, "1", "2", "3", "4")
	// assert
	assert.IsType(t, &EasyDriver{}, d)
	assert.True(t, strings.HasPrefix(d.name, "EasyDriver"))
	assert.Equal(t, a, d.connection)
	assert.NoError(t, d.afterStart())
	assert.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	assert.Equal(t, "1", d.stepPin)
	assert.Equal(t, "2", d.dirPin)
	assert.Equal(t, "3", d.enPin)
	assert.Equal(t, "4", d.sleepPin)
	assert.Equal(t, float32(stepAngle), d.angle)
	assert.Equal(t, uint(180), d.rpm)
	assert.Equal(t, int8(1), d.dir)
	assert.Equal(t, 0, d.stepNum)
	assert.Equal(t, true, d.enabled)
	assert.Equal(t, false, d.sleeping)
	assert.Nil(t, d.runStopChan)
}

func TestEasyDriverHalt(t *testing.T) {
	// arrange
	d, _ := initTestEasyDriverWithStubbedAdaptor()
	require.NoError(t, d.Run())
	require.True(t, d.IsMoving())
	// act
	err := d.Halt()
	// assert
	assert.NoError(t, err)
	assert.False(t, d.IsMoving())
}

func TestEasyDriverMove(t *testing.T) {
	tests := map[string]struct {
		inputSteps             int
		simulateDisabled       bool
		simulateAlreadyRunning bool
		simulateWriteErr       bool
		wantWrites             int
		wantSteps              int
		wantMoving             bool
		wantErr                string
	}{
		"move_one": {
			inputSteps: 1,
			wantWrites: 4,
			wantSteps:  2,
			wantMoving: false,
		},
		"move_more": {
			inputSteps: 20,
			wantWrites: 80,
			wantSteps:  40,
			wantMoving: false,
		},
		"error_disabled": {
			simulateDisabled: true,
			wantMoving:       false,
			wantErr:          "is disabled",
		},
		"error_already_running": {
			simulateAlreadyRunning: true,
			wantMoving:             true,
			wantErr:                "already running or moving",
		},
		"error_write": {
			inputSteps:       1,
			simulateWriteErr: true,
			wantWrites:       1,
			wantMoving:       false,
			wantErr:          "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestEasyDriverWithStubbedAdaptor()
			d.enabled = !tc.simulateDisabled
			if tc.simulateAlreadyRunning {
				d.runStopChan = make(chan struct{})
				defer func() { close(d.runStopChan); d.runStopChan = nil }()
			}
			var numCallsWrite int
			a.digitalWriteFunc = func(string, byte) error {
				numCallsWrite++
				if tc.simulateWriteErr {
					return fmt.Errorf("write error")
				}
				return nil
			}
			// act
			err := d.Move(tc.inputSteps)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantSteps, d.stepNum)
			assert.Equal(t, tc.wantWrites, numCallsWrite)
			assert.Equal(t, tc.wantMoving, d.IsMoving())
		})
	}
}

func TestEasyDriverRun_IsMoving(t *testing.T) {
	tests := map[string]struct {
		simulateDisabled       bool
		simulateAlreadyRunning bool
		simulateWriteErr       bool
		wantMoving             bool
		wantErr                string
	}{
		"run": {
			wantMoving: true,
		},
		"error_disabled": {
			simulateDisabled: true,
			wantMoving:       false,
			wantErr:          "is disabled",
		},
		"write_error_skipped": {
			simulateWriteErr: true,
			wantMoving:       true,
		},
		"error_already_running": {
			simulateAlreadyRunning: true,
			wantMoving:             true,
			wantErr:                "already running or moving",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestEasyDriverWithStubbedAdaptor()
			d.enabled = !tc.simulateDisabled
			if tc.simulateAlreadyRunning {
				d.runStopChan = make(chan struct{})
				defer func() { close(d.runStopChan); d.runStopChan = nil }()
			}
			simWriteErr := tc.simulateWriteErr // to prevent data race in write function (go-called)
			a.digitalWriteFunc = func(string, byte) error {
				if simWriteErr {
					simWriteErr = false // to prevent to much output
					return fmt.Errorf("write error")
				}
				return nil
			}
			// act
			err := d.Run()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantMoving, d.IsMoving())
		})
	}
}

func TestEasyDriverStop_IsMoving(t *testing.T) {
	// arrange
	d, _ := initTestEasyDriverWithStubbedAdaptor()
	require.NoError(t, d.Run())
	require.True(t, d.IsMoving())
	// act
	err := d.Stop()
	// assert
	assert.NoError(t, err)
	assert.False(t, d.IsMoving())
}

func TestEasyDriverStep(t *testing.T) {
	tests := map[string]struct {
		countCallsForth        int
		countCallsBack         int
		simulateAlreadyRunning bool
		simulateWriteErr       bool
		wantSteps              int
		wantWritten            []byte
		wantErr                string
	}{
		"single": {
			countCallsForth: 1,
			wantSteps:       1,
			wantWritten:     []byte{0x00, 0x01},
		},
		"many": {
			countCallsForth: 4,
			wantSteps:       4,
			wantWritten:     []byte{0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1},
		},
		"forth_and_back": {
			countCallsForth: 5,
			countCallsBack:  3,
			wantSteps:       2,
			wantWritten:     []byte{0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1},
		},
		"reverse": {
			countCallsBack: 3,
			wantSteps:      -3,
			wantWritten:    []byte{0x0, 0x1, 0x0, 0x1, 0x0, 0x1},
		},
		"error_already_running": {
			countCallsForth:        1,
			simulateAlreadyRunning: true,
			wantErr:                "already running or moving",
		},
		"error_write": {
			simulateWriteErr: true,
			wantWritten:      []byte{0x00, 0x00},
			countCallsBack:   2,
			wantErr:          "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestEasyDriverWithStubbedAdaptor()
			if tc.simulateAlreadyRunning {
				d.runStopChan = make(chan struct{})
				defer func() { close(d.runStopChan); d.runStopChan = nil }()
			}
			var writtenValues []byte
			a.digitalWriteFunc = func(pin string, val byte) error {
				assert.Equal(t, d.stepPin, pin)
				writtenValues = append(writtenValues, val)
				if tc.simulateWriteErr {
					return fmt.Errorf("write error")
				}
				return nil
			}
			var errs []string
			// act
			for i := 0; i < tc.countCallsForth; i++ {
				if err := d.Step(); err != nil {
					errs = append(errs, err.Error())
				}
			}
			d.dir = -1
			for i := 0; i < tc.countCallsBack; i++ {
				if err := d.Step(); err != nil {
					errs = append(errs, err.Error())
				}
			}
			// assert
			if tc.wantErr != "" {
				assert.Contains(t, strings.Join(errs, ","), tc.wantErr)
			} else {
				assert.Nil(t, errs)
			}
			assert.Equal(t, tc.wantSteps, d.stepNum)
			assert.Equal(t, tc.wantSteps, d.CurrentStep())
			assert.Equal(t, tc.wantWritten, writtenValues)
		})
	}
}

func TestEasyDriverSetDirection(t *testing.T) {
	tests := map[string]struct {
		dirPin  string
		input   string
		wantVal int8
		wantErr string
	}{
		"cw": {
			input:   "cw",
			dirPin:  "10",
			wantVal: 1,
		},
		"ccw": {
			input:   "ccw",
			dirPin:  "11",
			wantVal: -1,
		},
		"unknown": {
			input:   "unknown",
			dirPin:  "12",
			wantVal: 1,
		},
		"error_no_pin": {
			dirPin:  "",
			wantVal: 1,
			wantErr: "dirPin is not set",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newGpioTestAdaptor()
			d := NewEasyDriver(a, stepAngle, "1", tc.dirPin, "3", "4")
			require.Equal(t, int8(1), d.dir)
			// act
			err := d.SetDirection(tc.input)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantVal, d.dir)
		})
	}
}

func TestEasyDriverSetSpeed(t *testing.T) {
	const (
		angle = 10
		max   = 36 // 360/angle
	)

	tests := map[string]struct {
		input uint
		want  uint
	}{
		"below_minimum": {
			input: 0,
			want:  1,
		},
		"minimum": {
			input: 1,
			want:  1,
		},
		"maximum": {
			input: max,
			want:  max,
		},
		"above_maximum": {
			input: max + 1,
			want:  max,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := EasyDriver{angle: angle}
			// act
			err := d.SetSpeed(tc.input)
			// assert
			assert.NoError(t, err)
			assert.Equal(t, tc.want, d.rpm)
		})
	}
}

func TestEasyDriverMaxSpeed(t *testing.T) {
	tests := map[string]struct {
		angle float32
		want  uint
	}{
		"180": {
			angle: 2.0,
			want:  180,
		},
		"360": {
			angle: 1.0,
			want:  360,
		},
		"720": {
			angle: 0.5,
			want:  720,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := EasyDriver{angle: tc.angle}
			// act & assert
			assert.Equal(t, tc.want, d.MaxSpeed())
		})
	}
}

func TestEasyDriverEnable_IsEnabled(t *testing.T) {
	tests := map[string]struct {
		enPin            string
		simulateWriteErr bool
		wantWrites       int
		wantEnabled      bool
		wantErr          string
	}{
		"basic": {
			enPin:       "10",
			wantWrites:  1,
			wantEnabled: true,
		},
		"with_run": {
			enPin:       "11",
			wantWrites:  1,
			wantEnabled: true,
		},
		"error_no_pin": {
			enPin:       "",
			wantWrites:  0,
			wantEnabled: true, // is enabled by default
			wantErr:     "enPin is not set",
		},
		"error_write": {
			enPin:            "12",
			simulateWriteErr: true,
			wantWrites:       1,
			wantEnabled:      false,
			wantErr:          "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newGpioTestAdaptor()
			d := NewEasyDriver(a, stepAngle, "1", "2", tc.enPin, "4")
			var numCallsWrite int
			var writtenPin string
			writtenValue := byte(0xFF)
			a.digitalWriteFunc = func(pin string, val byte) error {
				numCallsWrite++
				writtenPin = pin
				writtenValue = val
				if tc.simulateWriteErr {
					return fmt.Errorf("write error")
				}
				return nil
			}
			d.enabled = false
			require.False(t, d.IsEnabled())
			// act
			err := d.Enable()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, byte(0), writtenValue) // enable pin is active low
			}
			assert.Equal(t, tc.wantEnabled, d.IsEnabled())
			assert.Equal(t, tc.wantWrites, numCallsWrite)
			assert.Equal(t, tc.enPin, writtenPin)
		})
	}
}

func TestEasyDriverDisable_IsEnabled(t *testing.T) {
	tests := map[string]struct {
		enPin            string
		runBefore        bool
		simulateWriteErr string
		wantWrites       int
		wantEnabled      bool
		wantErr          string
	}{
		"basic": {
			enPin:       "10",
			wantWrites:  1,
			wantEnabled: false,
		},
		"with_run": {
			enPin:       "10",
			runBefore:   true,
			wantWrites:  1,
			wantEnabled: false,
		},
		"error_no_pin": {
			enPin:       "",
			wantWrites:  0,
			wantEnabled: true, // is enabled by default
			wantErr:     "enPin is not set",
		},
		"error_write": {
			enPin:            "12",
			simulateWriteErr: "write error",
			wantWrites:       1,
			wantEnabled:      true,
			wantErr:          "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newGpioTestAdaptor()
			d := NewEasyDriver(a, stepAngle, "1", "2", tc.enPin, "4")
			writeMutex := sync.Mutex{}
			var numCallsWrite int
			var writtenPin string
			writtenValue := byte(0xFF)
			a.digitalWriteFunc = func(pin string, val byte) error {
				writeMutex.Lock()
				defer writeMutex.Unlock()
				if pin == d.stepPin {
					// we do not consider call of step()
					return nil
				}
				numCallsWrite++
				writtenPin = pin
				writtenValue = val
				if tc.simulateWriteErr != "" {
					return fmt.Errorf(tc.simulateWriteErr)
				}
				return nil
			}
			if tc.runBefore {
				require.NoError(t, d.Run())
				require.True(t, d.IsMoving())
				time.Sleep(time.Millisecond)
			}
			d.enabled = true
			require.True(t, d.IsEnabled())
			// act
			err := d.Disable()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, byte(1), writtenValue) // enable pin is active low
			}
			assert.Equal(t, tc.wantEnabled, d.IsEnabled())
			assert.False(t, d.IsMoving())
			assert.Equal(t, tc.wantWrites, numCallsWrite)
			assert.Equal(t, tc.enPin, writtenPin)
		})
	}
}

func TestEasyDriverSleep_IsSleeping(t *testing.T) {
	tests := map[string]struct {
		sleepPin  string
		runBefore bool
		wantSleep bool
		wantErr   string
	}{
		"basic": {
			sleepPin:  "10",
			wantSleep: true,
		},
		"with_run": {
			sleepPin:  "11",
			runBefore: true,
			wantSleep: true,
		},
		"error_no_pin": {
			sleepPin:  "",
			wantSleep: false,
			wantErr:   "sleepPin is not set",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newGpioTestAdaptor()
			d := NewEasyDriver(a, stepAngle, "1", "2", "3", tc.sleepPin)
			if tc.runBefore {
				require.NoError(t, d.Run())
			}
			d.sleeping = false
			require.False(t, d.IsSleeping())
			// act
			err := d.Sleep()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantSleep, d.IsSleeping())
		})
	}
}

func TestEasyDriverWake_IsSleeping(t *testing.T) {
	tests := map[string]struct {
		sleepPin  string
		wantSleep bool
		wantErr   string
	}{
		"basic": {
			sleepPin:  "10",
			wantSleep: false,
		},
		"error_no_pin": {
			sleepPin:  "",
			wantSleep: true,
			wantErr:   "sleepPin is not set",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := newGpioTestAdaptor()
			d := NewEasyDriver(a, stepAngle, "1", "2", "3", tc.sleepPin)
			d.sleeping = true
			require.True(t, d.IsSleeping())
			// act
			err := d.Wake()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantSleep, d.IsSleeping())
		})
	}
}
