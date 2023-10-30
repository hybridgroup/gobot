package aio

import (
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
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
		assert.Equal(t, testAdaptor, driver.Connection())
		assert.Equal(t, pin, driver.Pin())
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
		testAdaptor.analogReadFunc = func() (int, error) {
			atomic.AddInt32(&callCount, 1)
			return 42, nil
		}

		// Start the driver and allow for multiple digital reads
		_ = driver.Start()
		time.Sleep(20 * time.Millisecond)

		_ = driver.Halt()
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
		testAdaptor.analogReadFunc = func() (val int, err error) {
			err = errors.New("read error")
			return
		}

		assert.NoError(t, driver.Start())

		// expect error
		_ = driver.Once(driver.Event(Error), func(data interface{}) {
			assert.Equal(t, "read error", data.(error).Error())
			close(sem)
		})

		select {
		case <-sem:
		case <-time.After(time.Second):
			t.Errorf("%s Event \"Error\" was not published", getType(driver))
		}

		// Cleanup
		_ = driver.Halt()
	}
}

func getType(driver interface{}) string {
	d := reflect.TypeOf(driver)

	if d.Kind() == reflect.Ptr {
		return d.Elem().Name()
	}

	return d.Name()
}
