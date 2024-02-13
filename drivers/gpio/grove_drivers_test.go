//nolint:forcetypeassert // ok here
package gpio

import (
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	a := newGpioTestAdaptor()
	pin := "456"

	drivers := []DriverAndPinner{
		NewGroveTouchDriver(a, pin),
		NewGroveButtonDriver(a, pin),
		NewGroveBuzzerDriver(a, pin),
		NewGroveLedDriver(a, pin),
		NewGroveRelayDriver(a, pin),
		NewGroveMagneticSwitchDriver(a, pin),
	}

	for _, driver := range drivers {
		assert.Equal(t, a, driver.Connection())
		assert.Equal(t, pin, driver.Pin())
	}
}

func TestDigitalDriverHalt(t *testing.T) {
	a := newGpioTestAdaptor()
	pin := "456"

	drivers := []DriverAndEventer{
		NewGroveTouchDriver(a, pin),
		NewGroveButtonDriver(a, pin),
		NewGroveMagneticSwitchDriver(a, pin),
	}

	for _, driver := range drivers {

		var callCount int32
		a.digitalReadFunc = func(string) (int, error) {
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
		// note: if a reading is already in progress, it will be finished before halt have an impact
		if atomic.LoadInt32(&callCount) > lastCallCount+1 {
			require.Fail(t, "DigitalRead was called more than once after driver was halted")
		}
	}
}

func TestDriverPublishesError(t *testing.T) {
	a := newGpioTestAdaptor()
	pin := "456"

	drivers := []DriverAndEventer{
		NewGroveTouchDriver(a, pin),
		NewGroveButtonDriver(a, pin),
		NewGroveMagneticSwitchDriver(a, pin),
	}

	for _, driver := range drivers {
		sem := make(chan struct{}, 1)
		// send error
		returnErr := func(string) (int, error) {
			return 0, errors.New("read error")
		}
		a.digitalReadFunc = returnErr

		require.NoError(t, driver.Start())

		// expect error
		_ = driver.Once(driver.Event(Error), func(data interface{}) {
			assert.Equal(t, "read error", data.(error).Error())
			close(sem)
		})

		select {
		case <-sem:
		case <-time.After(time.Second):
			require.Fail(t, "%s Event \"Error\" was not published", getType(driver))
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
