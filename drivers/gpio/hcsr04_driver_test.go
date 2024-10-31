package gpio

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/system"
)

func initTestHCSR04DriverWithStubbedAdaptor(triggerPinID string, echoPinID string) (*HCSR04Driver, *digitalPinMock) {
	a := newGpioTestAdaptor()
	tpin := a.addDigitalPin(triggerPinID)
	_ = a.addDigitalPin(echoPinID)
	d := NewHCSR04Driver(a, triggerPinID, echoPinID)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, tpin
}

func TestNewHCSR04Driver(t *testing.T) {
	// arrange
	const (
		triggerPinID = "3"
		echoPinID    = "4"
	)
	a := newGpioTestAdaptor()
	tpin := a.addDigitalPin(triggerPinID)
	epin := a.addDigitalPin(echoPinID)
	// act
	d := NewHCSR04Driver(a, triggerPinID, echoPinID)
	// assert
	assert.IsType(t, &HCSR04Driver{}, d)
	// assert: gpio.driver attributes
	assert.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "HCSR04"))
	assert.Equal(t, a, d.connection)
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.False(t, d.hcsr04Cfg.useEdgePolling)
	assert.Equal(t, triggerPinID, d.triggerPinID)
	assert.Equal(t, echoPinID, d.echoPinID)
	assert.NotNil(t, d.measureMutex)
	assert.Equal(t, tpin, d.triggerPin)
	assert.Equal(t, epin, d.echoPin)
}

func TestNewHCSR04Driver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "count up"
		cycReadDur = 30 * time.Millisecond
	)
	panicFunc := func() {
		NewHCSR04Driver(newGpioTestAdaptor(), "1", "2", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewHCSR04Driver(newGpioTestAdaptor(), "1", "2", WithName(myName), WithHCSR04UseEdgePolling())
	// assert
	assert.True(t, d.hcsr04Cfg.useEdgePolling)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestHCSR04MeasureDistance(t *testing.T) {
	tests := map[string]struct {
		measureMicroSec  int64
		simulateWriteErr string
		wantCallsWrite   int
		wantVal          float64
		wantErr          string
	}{
		"measure_ok": {
			measureMicroSec: 5831,
			wantCallsWrite:  2,
			wantVal:         1.0,
		},
		"error_timeout": {
			measureMicroSec: 170000, // > 160 ms
			wantCallsWrite:  2,
			wantErr:         "timeout 160ms reached",
		},
		"error_write": {
			measureMicroSec:  5831,
			simulateWriteErr: "write error",
			wantCallsWrite:   1,
			wantErr:          "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, tpin := initTestHCSR04DriverWithStubbedAdaptor("3", "4")
			// arrange sensor and event handler simulation
			waitForTriggerChan := make(chan struct{})
			loopWg := sync.WaitGroup{}
			defer func() {
				close(waitForTriggerChan)
				loopWg.Wait()
			}()
			loopWg.Add(1)
			go func() {
				<-waitForTriggerChan
				m := tc.measureMicroSec // to prevent data race together with wait group
				loopWg.Done()
				time.Sleep(time.Duration(m) * time.Microsecond)
				d.delayMicroSecChan <- m
			}()
			// arrange writes
			numCallsWrite := 0
			var oldVal int
			tpin.writeFunc = func(val int) error {
				numCallsWrite++
				if val == 0 && oldVal == 1 {
					// falling edge detected
					waitForTriggerChan <- struct{}{}
				}
				oldVal = val
				var err error
				if tc.simulateWriteErr != "" {
					err = errors.New(tc.simulateWriteErr)
				}
				return err
			}
			// act
			got, err := d.MeasureDistance()
			// assert
			assert.Equal(t, tc.wantCallsWrite, numCallsWrite)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.InDelta(t, tc.wantVal, got, 0.0)
		})
	}
}

func TestHCSR04Distance(t *testing.T) {
	tests := map[string]struct {
		measureMicroSec  int64
		simulateWriteErr string
		wantVal          float64
		wantErr          string
	}{
		"distance_0mm": {
			measureMicroSec: 0, // no validity test yet
			wantVal:         0.0,
		},
		"distance_2cm": {
			measureMicroSec: 117, // 117us ~ 0.12ms => ~2cm
			wantVal:         0.02,
		},
		"distance_4m": {
			measureMicroSec: 23324, // 23324us ~ 24ms => ~4m
			wantVal:         4.0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := HCSR04Driver{lastMeasureMicroSec: tc.measureMicroSec}
			// act
			got := d.Distance()
			// assert
			assert.InDelta(t, tc.wantVal, got, 0.0)
		})
	}
}

func TestHCSR04StartDistanceMonitor(t *testing.T) {
	tests := map[string]struct {
		simulateIsStarted bool
		simulateWriteErr  bool
		wantErr           string
	}{
		"start_ok": {},
		"start_ok_measure_error": {
			simulateWriteErr: true,
		},
		"error_already_started": {
			simulateIsStarted: true,
			wantErr:           "already started for 'HCSR04-",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, tpin := initTestHCSR04DriverWithStubbedAdaptor("3", "4")
			defer func() {
				if d.distanceMonitorStopChan != nil {
					close(d.distanceMonitorStopChan)
				}
				if d.distanceMonitorStopWaitGroup != nil {
					d.distanceMonitorStopWaitGroup.Wait()
				}
			}()
			if tc.simulateIsStarted {
				d.distanceMonitorStopChan = make(chan struct{})
			}
			tpin.writeFunc = func(val int) error {
				if tc.simulateWriteErr {
					return fmt.Errorf("write error")
				}
				return nil
			}
			// act
			err := d.StartDistanceMonitor()
			time.Sleep(1 * time.Millisecond) // < 160 ms
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, d.distanceMonitorStopChan)
				assert.NotNil(t, d.distanceMonitorStopWaitGroup)
			}
		})
	}
}

func TestHCSR04StopDistanceMonitor(t *testing.T) {
	tests := map[string]struct {
		start   bool
		wantErr string
	}{
		"stop_ok": {
			start: true,
		},
		"error_not_started": {
			wantErr: "not yet started for 'HCSR04-",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, _ := initTestHCSR04DriverWithStubbedAdaptor("3", "4")
			defer func() {
				if d.distanceMonitorStopChan != nil {
					close(d.distanceMonitorStopChan)
				}
				if d.distanceMonitorStopWaitGroup != nil {
					d.distanceMonitorStopWaitGroup.Wait()
				}
			}()
			if tc.start {
				err := d.StartDistanceMonitor()
				require.NoError(t, err)
			}
			// act
			err := d.StopDistanceMonitor()
			time.Sleep(1 * time.Millisecond) // < 160 ms
			// assert
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Nil(t, d.distanceMonitorStopChan)
			}
		})
	}
}

func TestHCSR04_createEventHandler(t *testing.T) {
	type eventCall struct {
		timeStamp time.Duration
		eventType string
	}
	tests := map[string]struct {
		calls []eventCall
		wants []int64
	}{
		"only_rising": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 2 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
			},
		},
		"only_falling": {
			calls: []eventCall{
				{timeStamp: 2 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
				{timeStamp: 3 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
			},
		},
		"event_normal": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 10 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
			},
			wants: []int64{9},
		},
		"event_falling_before": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
				{timeStamp: 2 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 10 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
			},
			wants: []int64{8},
		},
		"event_falling_after": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 10 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
				{timeStamp: 12 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
			},
			wants: []int64{9},
		},
		"event_rising_before": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 5 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 10 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
			},
			wants: []int64{5},
		},
		"event_rising_after": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 10 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
				{timeStamp: 12 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
			},
			wants: []int64{9},
		},
		"event_multiple": {
			calls: []eventCall{
				{timeStamp: 1 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 10 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
				{timeStamp: 11 * time.Microsecond, eventType: system.DigitalPinEventRisingEdge},
				{timeStamp: 13 * time.Microsecond, eventType: system.DigitalPinEventFallingEdge},
			},
			wants: []int64{9, 2},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := HCSR04Driver{delayMicroSecChan: make(chan int64, len(tc.wants))}
			// act
			eh := d.createEventHandler()
			for _, call := range tc.calls {
				eh(0, call.timeStamp, call.eventType, 0, 0)
			}
			// assert
			for _, want := range tc.wants {
				got := <-d.delayMicroSecChan
				assert.Equal(t, want, got)
			}
		})
	}
}
