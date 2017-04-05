package microbit

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*TemperatureDriver)(nil)

func initTestTemperatureDriver() *TemperatureDriver {
	d := NewTemperatureDriver(NewBleTestAdaptor())
	return d
}

func TestTemperatureDriver(t *testing.T) {
	d := initTestTemperatureDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Temperature"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestTemperatureDriverStartAndHalt(t *testing.T) {
	d := initTestTemperatureDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestTemperatureDriverReadData(t *testing.T) {
	sem := make(chan bool, 0)
	a := NewBleTestAdaptor()
	d := NewTemperatureDriver(a)
	d.Start()
	d.On(Temperature, func(data interface{}) {
		gobottest.Assert(t, data, int8(0x22))
		sem <- true
	})

	a.TestReceiveNotification([]byte{0x22}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"Temperature\" was not published")
	}
}
