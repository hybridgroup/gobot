package microbit

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

func initTestButtonDriver() *ButtonDriver {
	d := NewButtonDriver(NewBleTestAdaptor())
	return d
}

func TestButtonDriver(t *testing.T) {
	d := initTestButtonDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Button"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestButtonDriverStartAndHalt(t *testing.T) {
	d := initTestButtonDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestButtonDriverReadData(t *testing.T) {
	sem := make(chan bool, 0)
	a := NewBleTestAdaptor()
	d := NewButtonDriver(a)
	d.Start()
	d.On(ButtonB, func(data interface{}) {
		sem <- true
	})

	a.TestReceiveNotification([]byte{1}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"ButtonB\" was not published")
	}
}
