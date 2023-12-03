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

var _ gobot.Driver = (*ButtonDriver)(nil)

const buttonTestDelay = 250

func initTestButtonDriverWithStubbedAdaptor() (*ButtonDriver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	d := NewButtonDriver(a, "1")
	return d, a
}

func TestNewButtonDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewButtonDriver(a, "1")
	// assert
	assert.IsType(t, &ButtonDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "Button"))
	assert.Equal(t, "1", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.False(t, d.active)
	assert.Nil(t, d.Eventer) // will be created on initialize
	assert.Nil(t, d.halt)    // will be created on initialize
	require.NotNil(t, d.buttonCfg)
	assert.Equal(t, 0, d.buttonCfg.defaultState)
	assert.Equal(t, 10*time.Millisecond, d.buttonCfg.readInterval)
}

func TestNewButtonDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "count up"
		cycReadDur = 30 * time.Millisecond
	)
	panicFunc := func() {
		NewButtonDriver(newGpioTestAdaptor(), "1", WithName("crazy"), aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewButtonDriver(newGpioTestAdaptor(), "1", WithName(myName), WithButtonPollInterval(cycReadDur))
	// assert
	assert.Equal(t, cycReadDur, d.buttonCfg.readInterval)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestButton_WithButtonDefaultState(t *testing.T) {
	// arrange
	const myDefaultState = 5 // only for test, usually it would be 0 or 1
	cfg := buttonConfiguration{}
	// act
	WithButtonDefaultState(myDefaultState).apply(&cfg)
	// assert
	assert.Equal(t, myDefaultState, cfg.defaultState)
}

func TestButtonStart(t *testing.T) {
	// arrange
	sem := make(chan bool)
	nextVal := make(chan int, 1)
	d, a := initTestButtonDriverWithStubbedAdaptor()

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

	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active())
		nextVal <- 0
		sem <- true
	})

	// assert & rearrange
	require.NoError(t, err)

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		assert.Fail(t, "Button Event \"Push\" was not published")
	}

	_ = d.Once(ButtonRelease, func(data interface{}) {
		assert.False(t, d.Active())
		nextVal <- -1
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		assert.Fail(t, "Button Event \"Release\" was not published")
	}

	_ = d.Once(Error, func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		assert.Fail(t, "Button Event \"Error\" was not published")
	}

	_ = d.Once(ButtonPush, func(data interface{}) {
		sem <- true
	})

	require.NoError(t, d.Halt())
	nextVal <- 1

	select {
	case <-sem:
		assert.Fail(t, "Button Event \"Press\" should not published")
	case <-time.After(buttonTestDelay * time.Millisecond):
	}
}

func TestButtonStart_WithDefaultState(t *testing.T) {
	// arrange
	sem := make(chan bool)
	nextVal := make(chan int, 1)
	a := newGpioTestAdaptor()
	d := NewButtonDriver(a, "1", WithButtonDefaultState(1))

	a.digitalReadFunc = func(string) (int, error) {
		val := 0
		select {
		case val = <-nextVal:
			return val, nil
		default:
			return val, nil
		}
	}

	// act: start cyclic reading
	require.NoError(t, d.Start())
	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active())
		nextVal <- 1
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		assert.Fail(t, "Button Event \"Push\" was not published")
	}

	_ = d.Once(ButtonRelease, func(data interface{}) {
		assert.False(t, d.Active())

		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		assert.Fail(t, "Button Event \"Release\" was not published")
	}
}

func TestButtonHalt(t *testing.T) {
	// arrange
	d, _ := initTestButtonDriverWithStubbedAdaptor()
	require.NoError(t, d.Start())
	timeout := 2 * d.buttonCfg.readInterval
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

func TestButtonActive(t *testing.T) {
	tests := map[string]struct {
		want bool
	}{
		"active_true":  {want: true},
		"active_false": {want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := ButtonDriver{driver: newDriver(nil, "Button")} // just for mutex
			d.active = tc.want
			// act & assert
			assert.Equal(t, tc.want, d.Active())
		})
	}
}
