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
	assert.True(t, strings.HasPrefix(d.name, "Button"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	assert.NotNil(t, d.Eventer)
	assert.Equal(t, "1", d.pin)
	assert.False(t, d.active)
	assert.Equal(t, 0, d.defaultState)
	assert.Equal(t, 10*time.Millisecond, d.interval)
	assert.NotNil(t, d.halt)
	// act & assert other interval
	d = NewButtonDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	assert.Equal(t, 30*time.Second, d.interval)
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

	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active())
		nextVal <- 0
		sem <- true
	})

	// act
	err := d.Start()

	// assert & rearrange
	require.NoError(t, err)

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	_ = d.Once(ButtonRelease, func(data interface{}) {
		assert.False(t, d.Active())
		nextVal <- -1
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}

	_ = d.Once(Error, func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Error\" was not published")
	}

	_ = d.Once(ButtonPush, func(data interface{}) {
		sem <- true
	})

	d.halt <- true
	nextVal <- 1

	select {
	case <-sem:
		t.Errorf("Button Event \"Press\" should not published")
	case <-time.After(buttonTestDelay * time.Millisecond):
	}
}

func TestButtonSetDefaultState(t *testing.T) {
	// arrange
	sem := make(chan bool)
	nextVal := make(chan int, 1)
	d, a := initTestButtonDriverWithStubbedAdaptor()

	a.digitalReadFunc = func(string) (int, error) {
		val := 0
		select {
		case val = <-nextVal:
			return val, nil
		default:
			return val, nil
		}
	}
	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active())
		nextVal <- 1
		sem <- true
	})

	// act
	d.SetDefaultState(1)

	// assert & rearrange
	require.Equal(t, 1, d.defaultState)
	require.NoError(t, d.Start())

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	_ = d.Once(ButtonRelease, func(data interface{}) {
		assert.False(t, d.Active())

		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}
}

func TestButtonHalt(t *testing.T) {
	// arrange
	d, _ := initTestButtonDriverWithStubbedAdaptor()
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

func TestButtonPin(t *testing.T) {
	tests := map[string]struct {
		want string
	}{
		"10": {want: "10"},
		"36": {want: "36"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := ButtonDriver{pin: name}
			// act & assert
			assert.Equal(t, tc.want, d.Pin())
		})
	}
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
			d := ButtonDriver{Driver: NewDriver(nil, "Button")} // just for mutex
			d.active = tc.want
			// act & assert
			assert.Equal(t, tc.want, d.Active())
		})
	}
}
