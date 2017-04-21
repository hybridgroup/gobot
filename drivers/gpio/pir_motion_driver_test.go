package gpio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
	gobottest.Assert(t, d.Halt(), nil)
}

func TestPIRMotionDriver(t *testing.T) {
	d := NewPIRMotionDriver(newGpioTestAdaptor(), "1")
	gobottest.Refute(t, d.Connection(), nil)

	d = NewPIRMotionDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	gobottest.Assert(t, d.interval, 30*time.Second)
}

func TestPIRMotionDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	a := newGpioTestAdaptor()
	d := NewPIRMotionDriver(a, "1")

	gobottest.Assert(t, d.Start(), nil)

	d.Once(MotionDetected, func(data interface{}) {
		gobottest.Assert(t, d.Active, true)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 1
		return
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"MotionDetected\" was not published")
	}

	d.Once(MotionStopped, func(data interface{}) {
		gobottest.Assert(t, d.Active, false)
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		val = 0
		return
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"MotionStopped\" was not published")
	}

	d.Once(Error, func(data interface{}) {
		sem <- true
	})

	a.TestAdaptorDigitalRead(func() (val int, err error) {
		err = errors.New("digital read error")
		return
	})

	select {
	case <-sem:
	case <-time.After(motionTestDelay * time.Millisecond):
		t.Errorf("PIRMotionDriver Event \"Error\" was not published")
	}
}

func TestPIRDriverDefaultName(t *testing.T) {
	d := initTestPIRMotionDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "PIR"), true)
}

func TestPIRDriverSetName(t *testing.T) {
	d := initTestPIRMotionDriver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
