package gpio

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*PIRMotionDriver)(nil)

const motionTestDelay = 150

func initTestPIRMotionDriverWithStubbedAdaptor() (*PIRMotionDriver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	d := NewPIRMotionDriver(a, "1")
	return d, a
}

func TestNewPIRMotionDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewPIRMotionDriver(a, "1")
	// assert
	assert.IsType(t, &PIRMotionDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "PIRMotion"))
	assert.Equal(t, "1", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.False(t, d.active)
	assert.Equal(t, 10*time.Millisecond, d.pirMotionCfg.readInterval)
	assert.Nil(t, d.Eventer) // will be created on initialize
	assert.Nil(t, d.halt)    // will be created on initialize
}

func TestNewPIRMotionDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "voltage 1"
		cycReadDur = 30 * time.Millisecond
	)
	panicFunc := func() {
		NewPIRMotionDriver(newGpioTestAdaptor(), "1", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewPIRMotionDriver(newGpioTestAdaptor(), "1", WithName(myName), WithPIRMotionPollInterval(cycReadDur))
	// assert
	assert.Equal(t, cycReadDur, d.pirMotionCfg.readInterval)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestPIRMotionStart(t *testing.T) {
	// arrange
	sem := make(chan bool)
	nextVal := make(chan int, 1)
	a := newGpioTestAdaptor()
	d := NewPIRMotionDriver(a, "1")

	a.digitalReadFunc = func(string) (int, error) {
		val := 1
		var err error
		select {
		case val = <-nextVal:
			if val < 0 {
				err = fmt.Errorf("digital read error")
			}
			return val, err
		default:
			return val, err
		}
	}

	// act: start cyclic reading
	err := d.Start()

	_ = d.Once(MotionDetected, func(data interface{}) {
		assert.True(t, d.active)
		nextVal <- 0
		sem <- true
	})

	// assert & rearrange
	require.NoError(t, err)

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		require.Fail(t, "PIRMotionDriver Event \"MotionDetected\" was not published")
	}

	_ = d.Once(MotionStopped, func(data interface{}) {
		assert.False(t, d.active)
		nextVal <- -1
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		require.Fail(t, "PIRMotionDriver Event \"MotionStopped\" was not published")
	}

	_ = d.Once(Error, func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		require.Fail(t, "PIRMotionDriver Event \"Error\" was not published")
	}

	_ = d.Once(MotionDetected, func(data interface{}) {
		sem <- true
	})

	require.NoError(t, d.Halt())
	nextVal <- 1

	select {
	case <-sem:
		require.Fail(t, "PIRMotion Event \"MotionDetected\" should not published")
	case <-time.After(motionTestDelay * time.Millisecond):
	}
}

func TestPIRMotionHalt(t *testing.T) {
	// arrange
	d, _ := initTestPIRMotionDriverWithStubbedAdaptor()
	require.NoError(t, d.Start())
	timeout := 2 * d.pirMotionCfg.readInterval
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-d.halt: // wait until halt is broadcasted by close the channel
		case <-time.After(timeout): // otherwise run into the timeout
			assert.Fail(t, "halt was not received within %s", timeout)
		}
	}()
	// act & assert
	require.NoError(t, d.Halt())
	wg.Wait() // wait until the go function was really finished
}

func TestPIRMotionActive(t *testing.T) {
	tests := map[string]struct {
		want bool
	}{
		"active_true":  {want: true},
		"active_false": {want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := PIRMotionDriver{driver: newDriver(nil, "PIRMotion")} // just for mutex
			d.active = tc.want
			// act & assert
			assert.Equal(t, tc.want, d.Active())
		})
	}
}
