package aio

import (
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

type DriverAndPinner interface {
	gobot.Driver
	gobot.Pinner
}

type DriverAndEventer interface {
	gobot.Driver
	gobot.Eventer
}

func TestDriverDefaults(t *testing.T) {
	testAdaptor := newAioTestAdaptor()
	pin := "456"

	drivers := []DriverAndPinner{
		NewGroveSoundSensorDriver(testAdaptor, pin),
		NewGroveLightSensorDriver(testAdaptor, pin),
		NewGrovePiezoVibrationSensorDriver(testAdaptor, pin),
		NewGroveRotaryDriver(testAdaptor, pin),
	}

	for _, driver := range drivers {
		gobottest.Assert(t, driver.Connection(), testAdaptor)
		gobottest.Assert(t, driver.Pin(), pin)
	}
}

func TestAnalogDriverHalt(t *testing.T) {
	testAdaptor := newAioTestAdaptor()
	pin := "456"

	drivers := []DriverAndEventer{
		NewGroveSoundSensorDriver(testAdaptor, pin),
		NewGroveLightSensorDriver(testAdaptor, pin),
		NewGrovePiezoVibrationSensorDriver(testAdaptor, pin),
		NewGroveRotaryDriver(testAdaptor, pin),
	}

	for _, driver := range drivers {
		var callCount int32
		testAdaptor.TestAdaptorAnalogRead(func() (int, error) {
			atomic.AddInt32(&callCount, 1)
			return 42, nil
		})

		// Start the driver and allow for multiple digital reads
		driver.Start()
		time.Sleep(20 * time.Millisecond)

		driver.Halt()
		lastCallCount := atomic.LoadInt32(&callCount)
		// If driver was not halted, digital reads would still continue
		time.Sleep(20 * time.Millisecond)
		if atomic.LoadInt32(&callCount) != lastCallCount {
			t.Errorf("AnalogRead was called after driver was halted")
		}
	}
}

func TestDriverPublishesError(t *testing.T) {
	testAdaptor := newAioTestAdaptor()
	pin := "456"

	drivers := []DriverAndEventer{
		NewGroveSoundSensorDriver(testAdaptor, pin),
		NewGroveLightSensorDriver(testAdaptor, pin),
		NewGrovePiezoVibrationSensorDriver(testAdaptor, pin),
		NewGroveRotaryDriver(testAdaptor, pin),
	}

	for _, driver := range drivers {
		sem := make(chan struct{}, 1)
		// send error
		returnErr := func() (val int, err error) {
			err = errors.New("read error")
			return
		}
		testAdaptor.TestAdaptorAnalogRead(returnErr)

		gobottest.Assert(t, driver.Start(), nil)

		// expect error
		driver.Once(driver.Event(Error), func(data interface{}) {
			gobottest.Assert(t, data.(error).Error(), "read error")
			close(sem)
		})

		select {
		case <-sem:
		case <-time.After(time.Second):
			t.Errorf("%s Event \"Error\" was not published", getType(driver))
		}

		// Cleanup
		driver.Halt()
	}
}

func getType(driver interface{}) string {
	d := reflect.TypeOf(driver)

	if d.Kind() == reflect.Ptr {
		return d.Elem().Name()
	}

	return d.Name()
}
