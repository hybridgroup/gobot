package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*TM1638Driver)(nil)

// --------- HELPERS
func initTestTM1638Driver() (driver *TM1638Driver) {
	driver, _ = initTestTM1638DriverWithStubbedAdaptor()
	return
}

func initTestTM1638DriverWithStubbedAdaptor() (*TM1638Driver, *gpioTestAdaptor) {
	adaptor := newGpioTestAdaptor()
	return NewTM1638Driver(adaptor, "1", "2", "3"), adaptor
}

// --------- TESTS
func TestTM1638Driver(t *testing.T) {
	var a interface{} = initTestTM1638Driver()
	_, ok := a.(*TM1638Driver)
	if !ok {
		t.Errorf("NewTM1638Driver() should have returned a *TM1638Driver")
	}
}

func TestTM1638DriverStart(t *testing.T) {
	d := initTestTM1638Driver()
	assert.NoError(t, d.Start())
}

func TestTM1638DriverHalt(t *testing.T) {
	d := initTestTM1638Driver()
	assert.NoError(t, d.Halt())
}

func TestTM1638DriverDefaultName(t *testing.T) {
	d := initTestTM1638Driver()
	assert.True(t, strings.HasPrefix(d.Name(), "TM1638"))
}

func TestTM1638DriverSetName(t *testing.T) {
	d := initTestTM1638Driver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}

func TestTM1638DriverFromStringToByteArray(t *testing.T) {
	d := initTestTM1638Driver()
	data := d.fromStringToByteArray("Hello World")
	assert.Equal(t, data, []byte{0x76, 0x7B, 0x30, 0x30, 0x5C, 0x00, 0x1D, 0x5C, 0x50, 0x30, 0x5E})
}

func TestTM1638DriverAddFonts(t *testing.T) {
	d := initTestTM1638Driver()
	d.AddFonts(map[string]byte{"µ": 0x1C, "ß": 0x7F})
	data := d.fromStringToByteArray("µß")
	assert.Equal(t, data, []byte{0x1C, 0x7F})
}

func TestTM1638DriverClearFonts(t *testing.T) {
	d := initTestTM1638Driver()
	d.ClearFonts()
	data := d.fromStringToByteArray("Hello World")
	assert.Equal(t, data, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}
