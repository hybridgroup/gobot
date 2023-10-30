package gpio

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
	testAdaptor := newGpioTestAdaptor()
	pin := "456"

	drivers := []DriverAndPinner{
		NewGroveTouchDriver(testAdaptor, pin),
		NewGroveButtonDriver(testAdaptor, pin),
		NewGroveBuzzerDriver(testAdaptor, pin),
		NewGroveLedDriver(testAdaptor, pin),
		NewGroveRelayDriver(testAdaptor, pin),
		NewGroveMagneticSwitchDriver(testAdaptor, pin),
	}

	for _, driver := range drivers {
		assert.Equal(t, testAdaptor, driver.Connection())
		assert.Equal(t, pin, driver.Pin())
	}
}

func TestDigitalDriverHalt(t *testing.T) {
	testAdaptor := newGpioTestAdaptor()
	pin := "456"

	drivers := []DriverAndEventer{
		NewGroveTouchDriver(testAdaptor, pin),
		NewGroveButtonDriver(testAdaptor, pin),
		NewGroveMagneticSwitchDriver(testAdaptor, pin),
	}

	for _, driver := range drivers {

		var callCount int32
		testAdaptor.digitalReadFunc = func(string) (int, error) {
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
			t.Errorf("DigitalRead was called after driver was halted")
		}
	}
}

func TestDriverPublishesError(t *testing.T) {
	testAdaptor := newGpioTestAdaptor()
	pin := "456"

	drivers := []DriverAndEventer{
		NewGroveTouchDriver(testAdaptor, pin),
		NewGroveButtonDriver(testAdaptor, pin),
		NewGroveMagneticSwitchDriver(testAdaptor, pin),
	}

	for _, driver := range drivers {
		sem := make(chan struct{}, 1)
		// send error
		returnErr := func(string) (val int, err error) {
			err = errors.New("read error")
			return
		}
		testAdaptor.digitalReadFunc = returnErr

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
