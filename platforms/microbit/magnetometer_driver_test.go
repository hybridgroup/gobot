package microbit

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MagnetometerDriver)(nil)

func initTestMagnetometerDriver() *MagnetometerDriver {
	d := NewMagnetometerDriver(NewBleTestAdaptor())
	return d
}

func TestMagnetometerDriver(t *testing.T) {
	d := initTestMagnetometerDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Magnetometer"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestMagnetometerDriverStartAndHalt(t *testing.T) {
	d := initTestMagnetometerDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMagnetometerDriverReadData(t *testing.T) {
	sem := make(chan bool, 0)
	a := NewBleTestAdaptor()
	d := NewMagnetometerDriver(a)
	d.Start()
	d.On(Magnetometer, func(data interface{}) {
		gobottest.Assert(t, data.(*MagnetometerData).X, float32(8.738))
		gobottest.Assert(t, data.(*MagnetometerData).Y, float32(8.995))
		gobottest.Assert(t, data.(*MagnetometerData).Z, float32(9.252))
		sem <- true
	})

	a.TestReceiveNotification([]byte{0x22, 0x22, 0x23, 0x23, 0x24, 0x24}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"Magnetometer\" was not published")
	}
}
