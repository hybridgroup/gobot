package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AIP1640Driver)(nil)

// --------- HELPERS
func initTestAIP1640Driver() (driver *AIP1640Driver) {
	driver, _ = initTestAIP1640DriverWithStubbedAdaptor()
	return
}

func initTestAIP1640DriverWithStubbedAdaptor() (*AIP1640Driver, *gpioTestAdaptor) {
	adaptor := newGpioTestAdaptor()
	return NewAIP1640Driver(adaptor, "1", "2"), adaptor
}

// --------- TESTS
func TestAIP1640Driver(t *testing.T) {
	var a interface{} = initTestAIP1640Driver()
	_, ok := a.(*AIP1640Driver)
	if !ok {
		t.Errorf("NewAIP1640Driver() should have returned a *AIP1640Driver")
	}
}

func TestAIP1640DriverStart(t *testing.T) {
	d := initTestAIP1640Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestAIP1640DriverHalt(t *testing.T) {
	d := initTestAIP1640Driver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestAIP1640DriverDefaultName(t *testing.T) {
	d := initTestAIP1640Driver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "AIP1640Driver"), true)
}

func TestAIP1640DriverSetName(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}

func TestAIP1640DriveDrawPixel(t *testing.T) {
	d := initTestAIP1640Driver()
	d.DrawPixel(2, 3, true)
	d.DrawPixel(0, 3, true)
	gobottest.Assert(t, uint8(5), d.buffer[7-3])
}

func TestAIP1640DriverDrawRow(t *testing.T) {
	d := initTestAIP1640Driver()
	d.DrawRow(4, 0x3C)
	gobottest.Assert(t, uint8(0x3C), d.buffer[7-4])
}

func TestAIP1640DriverDrawMatrix(t *testing.T) {
	d := initTestAIP1640Driver()
	drawing := [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	d.DrawMatrix(drawing)
	gobottest.Assert(t, [8]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}, d.buffer)
}

func TestAIP1640DriverClear(t *testing.T) {
	d := initTestAIP1640Driver()
	drawing := [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	d.DrawMatrix(drawing)
	gobottest.Assert(t, [8]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}, d.buffer)
	d.Clear()
	gobottest.Assert(t, [8]byte{}, d.buffer)
}

func TestAIP1640DriverSetIntensity(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetIntensity(3)
	gobottest.Assert(t, uint8(3), d.intensity)
}

func TestAIP1640DriverSetIntensityHigherThan7(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetIntensity(19)
	gobottest.Assert(t, uint8(7), d.intensity)
}
