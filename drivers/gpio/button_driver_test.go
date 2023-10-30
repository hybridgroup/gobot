package gpio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

const buttonTestDelay = 250

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(newGpioTestAdaptor(), "1")
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	go func() {
		<-d.halt
	}()
	assert.NoError(t, d.Halt())
}

func TestButtonDriver(t *testing.T) {
	d := NewButtonDriver(newGpioTestAdaptor(), "1")
	assert.NotNil(t, d.Connection())

	d = NewButtonDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	assert.Equal(t, 30*time.Second, d.interval)
}

func TestButtonDriverStart(t *testing.T) {
	sem := make(chan bool)
	a := newGpioTestAdaptor()
	d := NewButtonDriver(a, "1")

	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 1
		return
	}

	assert.NoError(t, d.Start())

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	_ = d.Once(ButtonRelease, func(data interface{}) {
		assert.False(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 0
		return
	}

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}

	_ = d.Once(Error, func(data interface{}) {
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		err = errors.New("digital read error")
		return
	}

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Error\" was not published")
	}

	_ = d.Once(ButtonPush, func(data interface{}) {
		sem <- true
	})

	d.halt <- true

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 1
		return
	}

	select {
	case <-sem:
		t.Errorf("Button Event \"Press\" should not published")
	case <-time.After(buttonTestDelay * time.Millisecond):
	}
}

func TestButtonDriverDefaultState(t *testing.T) {
	sem := make(chan bool)
	a := newGpioTestAdaptor()
	d := NewButtonDriver(a, "1")
	d.DefaultState = 1

	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 0
		return
	}

	assert.NoError(t, d.Start())

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Push\" was not published")
	}

	_ = d.Once(ButtonRelease, func(data interface{}) {
		assert.False(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 1
		return
	}

	select {
	case <-sem:
	case <-time.After(buttonTestDelay * time.Millisecond):
		t.Errorf("Button Event \"Release\" was not published")
	}
}

func TestButtonDriverDefaultName(t *testing.T) {
	g := initTestButtonDriver()
	assert.True(t, strings.HasPrefix(g.Name(), "Button"))
}

func TestButtonDriverSetName(t *testing.T) {
	g := initTestButtonDriver()
	g.SetName("mybot")
	assert.Equal(t, "mybot", g.Name())
}
