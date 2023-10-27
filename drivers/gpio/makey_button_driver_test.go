package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*MakeyButtonDriver)(nil)

const makeyTestDelay = 250

func initTestMakeyButtonDriver() *MakeyButtonDriver {
	return NewMakeyButtonDriver(newGpioTestAdaptor(), "1")
}

func TestMakeyButtonDriverHalt(t *testing.T) {
	d := initTestMakeyButtonDriver()
	done := make(chan struct{})
	go func() {
		<-d.halt
		close(done)
	}()
	assert.NoError(t, d.Halt())
	select {
	case <-done:
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton was not halted")
	}
}

func TestMakeyButtonDriver(t *testing.T) {
	d := initTestMakeyButtonDriver()
	assert.Equal(t, "1", d.Pin())
	assert.NotNil(t, d.Connection())
	assert.Equal(t, 10*time.Millisecond, d.interval)

	d = NewMakeyButtonDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	assert.Equal(t, 30*time.Second, d.interval)
}

func TestMakeyButtonDriverStart(t *testing.T) {
	sem := make(chan bool)
	a := newGpioTestAdaptor()
	d := NewMakeyButtonDriver(a, "1")

	assert.NoError(t, d.Start())

	_ = d.Once(ButtonPush, func(data interface{}) {
		assert.True(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 0
		return
	}

	select {
	case <-sem:
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton Event \"Push\" was not published")
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
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton Event \"Release\" was not published")
	}

	_ = d.Once(Error, func(data interface{}) {
		assert.Equal(t, "digital read error", data.(error).Error())
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		err = errors.New("digital read error")
		return
	}

	select {
	case <-sem:
	case <-time.After(makeyTestDelay * time.Millisecond):
		t.Errorf("MakeyButton Event \"Error\" was not published")
	}

	// send a halt message
	_ = d.Once(ButtonRelease, func(data interface{}) {
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 1
		return
	}

	d.halt <- true

	select {
	case <-sem:
		t.Errorf("MakeyButton Event should not have been published")
	case <-time.After(makeyTestDelay * time.Millisecond):
	}
}
