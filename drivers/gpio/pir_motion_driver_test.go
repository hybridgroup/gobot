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
	assert.True(t, strings.HasPrefix(d.name, "PIRMotion"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	assert.NotNil(t, d.Eventer)
	assert.Equal(t, "1", d.pin)
	assert.False(t, d.active)
	assert.Equal(t, 10*time.Millisecond, d.interval)
	assert.NotNil(t, d.halt)
	// act & assert other interval
	d = NewPIRMotionDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	assert.Equal(t, 30*time.Second, d.interval)
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

	_ = d.Once(MotionDetected, func(data interface{}) {
		assert.True(t, d.active)
		nextVal <- 0
		sem <- true
	})

	// act
	err := d.Start()

	// assert & rearrange
	require.NoError(t, err)

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"MotionDetected\" was not published")
	}

	_ = d.Once(MotionStopped, func(data interface{}) {
		assert.False(t, d.active)
		nextVal <- -1
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"MotionStopped\" was not published")
	}

	_ = d.Once(Error, func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"Error\" was not published")
	}

	_ = d.Once(MotionDetected, func(data interface{}) {
		sem <- true
	})

	d.halt <- true
	nextVal <- 1

	select {
	case <-sem:
		t.Errorf("PIRMotion Event \"MotionDetected\" should not published")
	case <-time.After(motionTestDelay * time.Millisecond):
	}
}

func TestPIRMotionHalt(t *testing.T) {
	// arrange
	d, _ := initTestPIRMotionDriverWithStubbedAdaptor()
	const timeout = 10 * time.Microsecond
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-d.halt: // wait until halt was set to the channel
		case <-time.After(timeout): // otherwise run into the timeout
			t.Errorf("halt was not received within %s", timeout)
		}
	}()
	// act & assert
	require.NoError(t, d.Halt())
	wg.Wait() // wait until the go function was really finished
}

func TestPIRMotionPin(t *testing.T) {
	tests := map[string]struct {
		want string
	}{
		"10": {want: "10"},
		"36": {want: "36"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := PIRMotionDriver{pin: name}
			// act & assert
			assert.Equal(t, tc.want, d.Pin())
		})
	}
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
			d := PIRMotionDriver{Driver: NewDriver(nil, "PIRMotion")} // just for mutex
			d.active = tc.want
			// act & assert
			assert.Equal(t, tc.want, d.Active())
		})
	}
}
