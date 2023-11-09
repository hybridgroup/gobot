package gpio

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTestStepperDriverWithStubbedAdaptor() (*StepperDriver, *gpioTestAdaptor) {
	const stepsPerRev = 32

	a := newGpioTestAdaptor()
	d := NewStepperDriver(a, [4]string{"7", "11", "13", "15"}, StepperModes.DualPhaseStepping, stepsPerRev)
	return d, a
}

func TestNewStepperDriver(t *testing.T) {
	// arrange
	const stepsPerRev = 32

	a := newGpioTestAdaptor()
	// act
	d := NewStepperDriver(a, [4]string{"7", "11", "13", "15"}, StepperModes.DualPhaseStepping, stepsPerRev)
	// assert
	assert.IsType(t, &StepperDriver{}, d)
	assert.True(t, strings.HasPrefix(d.name, "Stepper"))
	assert.Equal(t, a, d.connection)
	assert.NoError(t, d.afterStart())
	assert.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	assert.Equal(t, "forward", d.direction)
	assert.Equal(t, StepperModes.DualPhaseStepping, d.phase)
	assert.Equal(t, float32(stepsPerRev), d.stepsPerRev)
	assert.Equal(t, 0, d.stepNum)
	assert.Nil(t, d.stopAsynchRunFunc)
}

func TestStepperMove_IsMoving(t *testing.T) {
	const stepsPerRev = 32

	tests := map[string]struct {
		inputSteps             int
		noAutoStopIfRunning    bool
		simulateAlreadyRunning bool
		simulateWriteErr       bool
		wantWrites             int
		wantSteps              int
		wantMoving             bool
		wantErr                string
	}{
		"move_forward": {
			inputSteps: 2,
			wantWrites: 8,
			wantSteps:  2,
			wantMoving: false,
		},
		"move_more_forward": {
			inputSteps: 10,
			wantWrites: 40,
			wantSteps:  10,
			wantMoving: false,
		},
		"move_forward_full_revolution": {
			inputSteps: stepsPerRev,
			wantWrites: 128,
			wantSteps:  0, // will be reset after each revision
			wantMoving: false,
		},
		"move_backward": {
			inputSteps: -2,
			wantWrites: 8,
			wantSteps:  stepsPerRev - 2,
			wantMoving: false,
		},
		"move_more_backward": {
			inputSteps: -10,
			wantWrites: 40,
			wantSteps:  stepsPerRev - 10,
			wantMoving: false,
		},
		"move_backward_full_revolution": {
			inputSteps: -stepsPerRev,
			wantWrites: 128,
			wantSteps:  0, // will be reset after each revision
			wantMoving: false,
		},
		"already_running_autostop": {
			inputSteps:             3,
			simulateAlreadyRunning: true,
			wantWrites:             12,
			wantSteps:              3,
			wantMoving:             false,
		},
		"error_already_running": {
			noAutoStopIfRunning:    true,
			simulateAlreadyRunning: true,
			wantMoving:             true,
			wantErr:                "already running or moving",
		},
		"error_no_steps": {
			inputSteps: 0,
			wantWrites: 0,
			wantSteps:  0,
			wantMoving: false,
			wantErr:    "no steps to do",
		},
		"error_write": {
			inputSteps:       1,
			simulateWriteErr: true,
			wantWrites:       0,
			wantMoving:       false,
			wantErr:          "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestStepperDriverWithStubbedAdaptor()
			defer func() {
				// for cleanup dangling channels
				if d.stopAsynchRunFunc != nil {
					err := d.stopAsynchRunFunc(true)
					assert.NoError(t, err)
				}
			}()
			// arrange: different behavior
			d.haltIfRunning = !tc.noAutoStopIfRunning
			if tc.simulateAlreadyRunning {
				d.stopAsynchRunFunc = func(bool) error { log.Println("former run stopped"); return nil }
			}
			// arrange: writes
			a.written = nil // reset writes of Start()
			a.simulateWriteError = tc.simulateWriteErr
			// act
			err := d.Move(tc.inputSteps)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantSteps, d.stepNum)
			assert.Equal(t, tc.wantWrites, len(a.written))
			assert.Equal(t, tc.wantMoving, d.IsMoving())
		})
	}
}

func TestStepperRun_IsMoving(t *testing.T) {
	tests := map[string]struct {
		noAutoStopIfRunning    bool
		simulateAlreadyRunning bool
		simulateWriteErr       bool
		wantMoving             bool
		wantErr                string
	}{
		"run": {
			wantMoving: true,
		},
		"error_write": {
			simulateWriteErr: true,
			wantMoving:       true,
		},
		"error_already_running": {
			noAutoStopIfRunning:    true,
			simulateAlreadyRunning: true,
			wantMoving:             true,
			wantErr:                "already running or moving",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestStepperDriverWithStubbedAdaptor()
			defer func() {
				// for cleanup dangling channels
				if d.stopAsynchRunFunc != nil {
					err := d.stopAsynchRunFunc(true)
					assert.NoError(t, err)
				}
			}()
			// arrange: different behavior
			writeChan := make(chan struct{})
			if tc.noAutoStopIfRunning {
				// in this case no write should be called
				close(writeChan)
				writeChan = nil
				d.haltIfRunning = false
			} else {
				d.haltIfRunning = true
			}
			if tc.simulateAlreadyRunning {
				d.stopAsynchRunFunc = func(bool) error { return nil }
			}
			// arrange: writes
			simWriteErr := tc.simulateWriteErr // to prevent data race in write function (go-called)
			var firstWriteDone bool
			a.digitalWriteFunc = func(string, byte) error {
				if firstWriteDone {
					return nil // to prevent to much output and write to channel
				}
				writeChan <- struct{}{}
				firstWriteDone = true
				if simWriteErr {
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
			if writeChan != nil {
				// wait until the first write was called and a little bit longer
				<-writeChan
				time.Sleep(10 * time.Millisecond)
				var asynchErr error
				if d.stopAsynchRunFunc != nil {
					asynchErr = d.stopAsynchRunFunc(false)
					d.stopAsynchRunFunc = nil
				}
				if tc.simulateWriteErr {
					assert.Error(t, asynchErr)
				} else {
					assert.NoError(t, asynchErr)
				}
			}
		})
	}
}

func TestStepperStop_IsMoving(t *testing.T) {
	tests := map[string]struct {
		running bool
		wantErr string
	}{
		"stop_running": {
			running: true,
		},
		"errro_not_started": {
			running: false,
			wantErr: "is not yet started",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, _ := initTestStepperDriverWithStubbedAdaptor()
			if tc.running {
				require.NoError(t, d.Run())
				require.True(t, d.IsMoving())
			}
			// act
			err := d.Stop()
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.False(t, d.IsMoving())
		})
	}
}

func TestStepperHalt_IsMoving(t *testing.T) {
	tests := map[string]struct {
		running bool
	}{
		"halt_running": {
			running: true,
		},
		"halt_not_started": {
			running: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, _ := initTestStepperDriverWithStubbedAdaptor()
			if tc.running {
				require.NoError(t, d.Run())
				require.True(t, d.IsMoving())
			}
			// act
			err := d.Halt()
			// assert
			assert.NoError(t, err)
			assert.False(t, d.IsMoving())
		})
	}
}

func TestStepperSetDirection(t *testing.T) {
	tests := map[string]struct {
		input   string
		wantVal string
		wantErr string
	}{
		"direction_forward": {
			input:   "forward",
			wantVal: "forward",
		},
		"direction_backward": {
			input:   "backward",
			wantVal: "backward",
		},
		"error_invalid_direction": {
			input:   "reverse",
			wantVal: "forward",
			wantErr: "Invalid direction 'reverse'",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, _ := initTestStepperDriverWithStubbedAdaptor()
			require.Equal(t, "forward", d.direction)
			// act
			err := d.SetDirection(tc.input)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantVal, d.direction)
		})
	}
}

func TestStepperMaxSpeed(t *testing.T) {
	const delayForMaxSpeed = 1428 * time.Microsecond // 1/700Hz

	tests := map[string]struct {
		stepsPerRev float32
		want        uint
	}{
		"maxspeed_for_20spr": {
			stepsPerRev: 20,
			want:        2100,
		},
		"maxspeed_for_50spr": {
			stepsPerRev: 50,
			want:        840,
		},
		"maxspeed_for_100spr": {
			stepsPerRev: 100,
			want:        420,
		},
		"maxspeed_for_400spr": {
			stepsPerRev: 400,
			want:        105,
		},
		"maxspeed_for_1000spr": {
			stepsPerRev: 1000,
			want:        42,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := StepperDriver{stepsPerRev: tc.stepsPerRev}
			// act
			got := d.MaxSpeed()
			d.speedRpm = got
			got2 := d.getDelayPerStep()
			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, delayForMaxSpeed, got2)
		})
	}
}

func TestStepperSetSpeed(t *testing.T) {
	const maxRpm = 1166

	tests := map[string]struct {
		input   uint
		want    uint
		wantErr string
	}{
		"below_minimum": {
			input:   0,
			want:    0,
			wantErr: "RPM (0) cannot be a zero or negative value",
		},
		"minimum": {
			input: 1,
			want:  1,
		},
		"maximum": {
			input: maxRpm,
			want:  maxRpm,
		},
		"above_maximum": {
			input:   maxRpm + 1,
			want:    maxRpm,
			wantErr: "cannot be greater then maximal value 1166",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, _ := initTestStepperDriverWithStubbedAdaptor()
			d.stepsPerRev = 36
			// act
			err := d.SetSpeed(tc.input)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, d.speedRpm)
		})
	}
}
