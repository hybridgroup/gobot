package microbit

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// the IOPinDriver is a Driver
var _ gobot.Driver = (*IOPinDriver)(nil)

// that supports the DigitalReader, DigitalWriter, & AnalogReader interfaces
var (
	_ gpio.DigitalReader = (*IOPinDriver)(nil)
	_ gpio.DigitalWriter = (*IOPinDriver)(nil)
	_ aio.AnalogReader   = (*IOPinDriver)(nil)
)

func initTestIOPinDriver() *IOPinDriver {
	d := NewIOPinDriver(NewBleTestAdaptor())
	return d
}

func TestIOPinDriver(t *testing.T) {
	d := initTestIOPinDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit IO Pin"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestIOPinDriverStartAndHalt(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 1, 1, 0}, nil
	})
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestIOPinDriverStartError(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return nil, errors.New("read error")
	})
	assert.ErrorContains(t, d.Start(), "read error")
}

func TestIOPinDriverDigitalRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 1, 1, 0, 2, 1}, nil
	})

	val, _ := d.DigitalRead("0")
	assert.Equal(t, 1, val)

	val, _ = d.DigitalRead("1")
	assert.Equal(t, 0, val)
}

func TestIOPinDriverDigitalReadInvalidPin(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	_, err := d.DigitalRead("A3")
	assert.NotNil(t, err)

	_, err = d.DigitalRead("6")
	assert.ErrorContains(t, err, "Invalid pin.")
}

func TestIOPinDriverDigitalWrite(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	// TODO: a better test
	assert.NoError(t, d.DigitalWrite("0", 1))
}

func TestIOPinDriverDigitalWriteInvalidPin(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	assert.NotNil(t, d.DigitalWrite("A3", 1))
	assert.ErrorContains(t, d.DigitalWrite("6", 1), "Invalid pin.")
}

func TestIOPinDriverAnalogRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, _ := d.AnalogRead("0")
	assert.Equal(t, 0, val)

	val, _ = d.AnalogRead("1")
	assert.Equal(t, 128, val)
}

func TestIOPinDriverAnalogReadInvalidPin(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	_, err := d.AnalogRead("A3")
	assert.NotNil(t, err)

	_, err = d.AnalogRead("6")
	assert.ErrorContains(t, err, "Invalid pin.")
}

func TestIOPinDriverDigitalAnalogRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, _ := d.DigitalRead("0")
	assert.Equal(t, 0, val)

	val, _ = d.AnalogRead("0")
	assert.Equal(t, 0, val)
}

func TestIOPinDriverDigitalWriteAnalogRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	assert.NoError(t, d.DigitalWrite("1", 0))

	val, _ := d.AnalogRead("1")
	assert.Equal(t, 128, val)
}

func TestIOPinDriverAnalogReadDigitalWrite(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, _ := d.AnalogRead("1")
	assert.Equal(t, 128, val)

	assert.NoError(t, d.DigitalWrite("1", 0))
}
