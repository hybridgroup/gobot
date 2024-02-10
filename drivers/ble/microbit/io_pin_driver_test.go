package microbit

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// the MicrobitIOPinDriver is a Driver
var _ gobot.Driver = (*IOPinDriver)(nil)

// that supports the DigitalReader, DigitalWriter, & AnalogReader interfaces
var (
	_ gpio.DigitalReader = (*IOPinDriver)(nil)
	_ gpio.DigitalWriter = (*IOPinDriver)(nil)
	_ aio.AnalogReader   = (*IOPinDriver)(nil)
)

func TestNewIOPinDriver(t *testing.T) {
	d := NewIOPinDriver(testutil.NewBleTestAdaptor())
	assert.IsType(t, &IOPinDriver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit IO Pin"))
	assert.NotNil(t, d.Eventer)
}

func TestNewIOPinDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewIOPinDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestIOPinStartAndHalt(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return []byte{0, 1, 1, 0}, nil
	})
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestIOPinStartError(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return nil, errors.New("read error")
	})
	require.ErrorContains(t, d.Start(), "read error")
}

func TestIOPinDigitalRead(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return []byte{0, 1, 1, 0, 2, 1}, nil
	})

	val, err := d.DigitalRead("0")
	require.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = d.DigitalRead("1")
	require.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestIOPinDigitalReadInvalidPin(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	_, err := d.DigitalRead("A3")
	require.Error(t, err)

	_, err = d.DigitalRead("6")
	require.ErrorContains(t, err, "Invalid pin.")
}

func TestIOPinDigitalWrite(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	// TODO: a better test
	require.NoError(t, d.DigitalWrite("0", 1))
}

func TestIOPinDigitalWriteInvalidPin(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	require.Error(t, d.DigitalWrite("A3", 1))
	require.ErrorContains(t, d.DigitalWrite("6", 1), "Invalid pin.")
}

func TestIOPinAnalogRead(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, err := d.AnalogRead("0")
	require.NoError(t, err)
	assert.Equal(t, 0, val)

	val, err = d.AnalogRead("1")
	require.NoError(t, err)
	assert.Equal(t, 128, val)
}

func TestIOPinAnalogReadInvalidPin(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)

	_, err := d.AnalogRead("A3")
	require.Error(t, err)

	_, err = d.AnalogRead("6")
	require.ErrorContains(t, err, "Invalid pin.")
}

func TestIOPinDigitalAnalogRead(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, err := d.DigitalRead("0")
	require.NoError(t, err)
	assert.Equal(t, 0, val)

	val, err = d.AnalogRead("0")
	require.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestIOPinDigitalWriteAnalogRead(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	require.NoError(t, d.DigitalWrite("1", 0))

	val, err := d.AnalogRead("1")
	require.NoError(t, err)
	assert.Equal(t, 128, val)
}

func TestIOPinAnalogReadDigitalWrite(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewIOPinDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		return []byte{0, 0, 1, 128, 2, 1}, nil
	})

	val, err := d.AnalogRead("1")
	require.NoError(t, err)
	assert.Equal(t, 128, val)

	require.NoError(t, d.DigitalWrite("1", 0))
}
