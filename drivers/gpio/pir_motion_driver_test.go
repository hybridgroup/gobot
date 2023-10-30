package gpio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*PIRMotionDriver)(nil)

const motionTestDelay = 150

func initTestPIRMotionDriver() *PIRMotionDriver {
	return NewPIRMotionDriver(newGpioTestAdaptor(), "1")
}

func TestPIRMotionDriverHalt(t *testing.T) {
	d := initTestPIRMotionDriver()
	go func() {
		<-d.halt
	}()
	assert.NoError(t, d.Halt())
}

func TestPIRMotionDriver(t *testing.T) {
	d := NewPIRMotionDriver(newGpioTestAdaptor(), "1")
	assert.NotNil(t, d.Connection())

	d = NewPIRMotionDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	assert.Equal(t, 30*time.Second, d.interval)
}

func TestPIRMotionDriverStart(t *testing.T) {
	sem := make(chan bool)
	a := newGpioTestAdaptor()
	d := NewPIRMotionDriver(a, "1")

	assert.NoError(t, d.Start())

	_ = d.Once(MotionDetected, func(data interface{}) {
		assert.True(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 1
		return
	}

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"MotionDetected\" was not published")
	}

	_ = d.Once(MotionStopped, func(data interface{}) {
		assert.False(t, d.Active)
		sem <- true
	})

	a.digitalReadFunc = func(string) (val int, err error) {
		val = 0
		return
	}

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"MotionStopped\" was not published")
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
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"Error\" was not published")
	}
}

func TestPIRDriverDefaultName(t *testing.T) {
	d := initTestPIRMotionDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "PIR"))
}

func TestPIRDriverSetName(t *testing.T) {
	d := initTestPIRMotionDriver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
