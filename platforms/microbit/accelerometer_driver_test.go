package microbit

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AccelerometerDriver)(nil)

func initTestAccelerometerDriver() *AccelerometerDriver {
	d := NewAccelerometerDriver(NewBleTestAdaptor())
	return d
}

func TestAccelerometerDriver(t *testing.T) {
	d := initTestAccelerometerDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Accelerometer"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestAccelerometerDriverStartAndHalt(t *testing.T) {
	d := initTestAccelerometerDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestAccelerometerDriverReadData(t *testing.T) {
	sem := make(chan bool, 0)
	a := NewBleTestAdaptor()
	d := NewAccelerometerDriver(a)
	d.Start()
	d.On(Accelerometer, func(data interface{}) {
		gobottest.Assert(t, data.(*AccelerometerData).X, float32(8.738))
		gobottest.Assert(t, data.(*AccelerometerData).Y, float32(8.995))
		gobottest.Assert(t, data.(*AccelerometerData).Z, float32(9.252))
		sem <- true
	})

	a.TestReceiveNotification([]byte{0x22, 0x22, 0x23, 0x23, 0x24, 0x24}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"Accelerometer\" was not published")
	}
}
