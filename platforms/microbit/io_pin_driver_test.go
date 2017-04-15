package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
)

// the IOPinDriver is a Driver
var _ gobot.Driver = (*IOPinDriver)(nil)

// and is also an Adaptor
var _ gobot.Adaptor = (*IOPinDriver)(nil)

// that supports the DigitalReader, DigitalWriter, & AnalogReader interfaces
var _ gpio.DigitalReader = (*IOPinDriver)(nil)
var _ gpio.DigitalWriter = (*IOPinDriver)(nil)
var _ aio.AnalogReader = (*IOPinDriver)(nil)

func initTestIOPinDriver() *IOPinDriver {
	d := NewIOPinDriver(NewBleTestAdaptor())
	return d
}

func TestIOPinDriver(t *testing.T) {
	d := initTestIOPinDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit IO Pin"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestIOPinDriverStartAndHalt(t *testing.T) {
	d := initTestIOPinDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestIOPinDriverDigitalRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 1, 1, 0, 2, 1}, nil
	})

	val, _ := d.DigitalRead("0")
	gobottest.Assert(t, val, 1)

	val, _ = d.DigitalRead("1")
	gobottest.Assert(t, val, 0)
}

func TestIOPinDriverDigitalWrite(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	// TODO: a better test
	gobottest.Assert(t, d.DigitalWrite("0", 1), nil)
}

func TestIOPinDriverAnalogRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, _ := d.AnalogRead("0")
	gobottest.Assert(t, val, 0)

	val, _ = d.AnalogRead("1")
	gobottest.Assert(t, val, 128)
}
