package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
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
	assert.NoError(t, d.Start())
}

func TestAIP1640DriverHalt(t *testing.T) {
	d := initTestAIP1640Driver()
	assert.NoError(t, d.Halt())
}

func TestAIP1640DriverDefaultName(t *testing.T) {
	d := initTestAIP1640Driver()
	assert.True(t, strings.HasPrefix(d.Name(), "AIP1640Driver"))
}

func TestAIP1640DriverSetName(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}

func TestAIP1640DriveDrawPixel(t *testing.T) {
	d := initTestAIP1640Driver()
	d.DrawPixel(2, 3, true)
	d.DrawPixel(0, 3, true)
	assert.Equal(t, d.buffer[7-3], uint8(5))
}

func TestAIP1640DriverDrawRow(t *testing.T) {
	d := initTestAIP1640Driver()
	d.DrawRow(4, 0x3C)
	assert.Equal(t, d.buffer[7-4], uint8(0x3C))
}

func TestAIP1640DriverDrawMatrix(t *testing.T) {
	d := initTestAIP1640Driver()
	drawing := [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	d.DrawMatrix(drawing)
	assert.Equal(t, d.buffer, [8]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01})
}

func TestAIP1640DriverClear(t *testing.T) {
	d := initTestAIP1640Driver()
	drawing := [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	d.DrawMatrix(drawing)
	assert.Equal(t, d.buffer, [8]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01})
	d.Clear()
	assert.Equal(t, d.buffer, [8]byte{})
}

func TestAIP1640DriverSetIntensity(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetIntensity(3)
	assert.Equal(t, d.intensity, uint8(3))
}

func TestAIP1640DriverSetIntensityHigherThan7(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetIntensity(19)
	assert.Equal(t, d.intensity, uint8(7))
}
