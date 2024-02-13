//nolint:forcetypeassert // ok here
package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCA9501Driver)(nil)

var (
	pinVal = map[string]interface{}{
		"pin": uint8(7),
		"val": uint8(0),
	}
	pin = map[string]interface{}{
		"pin": uint8(7),
	}
	addressVal = map[string]interface{}{
		"address": uint8(15),
		"val":     uint8(7),
	}
	address = map[string]interface{}{
		"address": uint8(15),
	}
)

func initPCA9501WithStubbedAdaptor() (*PCA9501Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewPCA9501Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewPCA9501Driver(t *testing.T) {
	// arrange, act
	var di interface{} = NewPCA9501Driver(newI2cTestAdaptor())
	// assert
	d, ok := di.(*PCA9501Driver)
	if !ok {
		require.Fail(t, "NewPCA9501Driver() should have returned a *PCA9501Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "PCA9501"))
	assert.Equal(t, 0x3f, d.defaultAddress)
}

func TestPCA9501Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewPCA9501Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestPCA9501CommandsWriteGPIO(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	// act
	result := d.Command("WriteGPIO")(pinVal)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCA9501CommandsReadGPIO(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// act
	result := d.Command("ReadGPIO")(pin)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCA9501CommandsWriteEEPROM(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	// act
	result := d.Command("WriteEEPROM")(addressVal)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCA9501CommandsReadEEPROM(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// act
	result := d.Command("ReadEEPROM")(address)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCA9501WriteGPIO(t *testing.T) {
	tests := map[string]struct {
		setVal          uint8
		ioDirAllInput   uint8
		ioStateAllInput uint8
		pin             uint8
		wantPin         uint8
		wantState       uint8
	}{
		"clear_bit": {
			setVal:          0,
			ioDirAllInput:   0xF1,
			ioStateAllInput: 0xF2,
			pin:             6,
			wantPin:         0xB1,
			wantState:       0xB2,
		},
		"set_bit": {
			setVal:          2,
			ioDirAllInput:   0x1F,
			ioStateAllInput: 0x20,
			pin:             3,
			wantPin:         0x17,
			wantState:       0x28,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA9501WithStubbedAdaptor()
			// prepare all reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				if numCallsRead == 1 {
					// first call read current io direction of all pins
					b[0] = tc.ioDirAllInput
				}
				if numCallsRead == 2 {
					// second call read current state of all pins
					b[0] = tc.ioStateAllInput
				}
				return len(b), nil
			}
			// act
			err := d.WriteGPIO(tc.pin, tc.setVal)
			// assert
			require.NoError(t, err)
			assert.Equal(t, 2, numCallsRead)
			assert.Len(t, a.written, 2)
			assert.Equal(t, tc.wantPin, a.written[0])
			assert.Equal(t, tc.wantState, a.written[1])
		})
	}
}

func TestPCA9501WriteGPIOErrorAtWriteDirection(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	wantErr := errors.New("write error")
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		return len(b), nil
	}
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		if numCallsWrite == 1 {
			// first call writes the CTRL register for port direction
			return 0, wantErr
		}
		return 0, nil
	}
	// act
	err := d.WriteGPIO(7, 0)
	// assert
	assert.Equal(t, wantErr, err)
	assert.Less(t, numCallsRead, 2)
	assert.Equal(t, 1, numCallsWrite)
}

func TestPCA9501WriteGPIOErrorAtWriteValue(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	wantErr := errors.New("write error")
	// prepare all reads
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		if numCallsWrite == 2 {
			// second call writes the value to IO port
			return 0, wantErr
		}
		return 0, nil
	}
	// act
	err := d.WriteGPIO(7, 0)
	// assert
	assert.Equal(t, wantErr, err)
	assert.Equal(t, 2, numCallsWrite)
}

func TestPCA9501ReadGPIO(t *testing.T) {
	tests := map[string]struct {
		ctrlState uint8
		want      uint8
	}{
		"pin_is_set":     {ctrlState: 0x80, want: 1},
		"pin_is_not_set": {ctrlState: 0x7F, want: 0},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const (
				pin           = uint8(7)
				wantCtrlState = uint8(0x80)
			)
			d, a := initPCA9501WithStubbedAdaptor()
			// prepare all reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				if numCallsRead == 1 {
					// current state of io
					b[0] = 0x00
				}
				if numCallsRead == 2 {
					b[0] = tc.ctrlState
				}
				return len(b), nil
			}
			// act
			got, err := d.ReadGPIO(pin)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, 2, numCallsRead)
			assert.Len(t, a.written, 1)
			assert.Equal(t, wantCtrlState, a.written[0])
		})
	}
}

func TestPCA9501ReadGPIOErrorAtReadDirection(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	wantErr := errors.New("read error")
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			// first read gets the CTRL register for pin direction
			return 0, wantErr
		}
		return len(b), nil
	}
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// act
	_, err := d.ReadGPIO(1)
	// assert
	assert.Equal(t, wantErr, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Equal(t, 0, numCallsWrite)
}

func TestPCA9501ReadGPIOErrorAtReadValue(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	wantErr := errors.New("read error")
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 2 {
			// second read gets the value from IO port
			return 0, wantErr
		}
		return len(b), nil
	}
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// act
	_, err := d.ReadGPIO(2)
	// assert
	assert.Equal(t, wantErr, err)
	assert.Equal(t, 1, numCallsWrite)
}

func TestPCA9501WriteEEPROM(t *testing.T) {
	// arrange
	const (
		addressEEPROM = uint8(0x52)
		want          = uint8(0x25)
	)
	d, a := initPCA9501WithStubbedAdaptor()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// act
	err := d.WriteEEPROM(addressEEPROM, want)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, addressEEPROM, a.written[0])
	assert.Equal(t, want, a.written[1])
}

func TestPCA9501ReadEEPROM(t *testing.T) {
	// arrange
	const (
		addressEEPROM = uint8(51)
		want          = uint8(0x44)
	)
	d, a := initPCA9501WithStubbedAdaptor()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func(b []byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[0] = want
		return len(b), nil
	}
	// act
	val, err := d.ReadEEPROM(addressEEPROM)
	// assert
	require.NoError(t, err)
	assert.Equal(t, want, val)
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, addressEEPROM, a.written[0])
	assert.Equal(t, 1, numCallsRead)
}

func TestPCA9501ReadEEPROMErrorWhileWriteAddress(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	wantErr := errors.New("error while write")
	// prepare all writes
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, wantErr
	}
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		return len(b), nil
	}
	// act
	_, err := d.ReadEEPROM(15)
	// assert
	assert.Equal(t, wantErr, err)
	assert.Equal(t, 0, numCallsRead)
}

func TestPCA9501ReadEEPROMErrorWhileReadValue(t *testing.T) {
	// arrange
	d, a := initPCA9501WithStubbedAdaptor()
	wantErr := errors.New("error while read")
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// prepare all reads
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), wantErr
	}
	// act
	_, err := d.ReadEEPROM(15)
	// assert
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, wantErr, err)
}

func TestPCA9501_initialize(t *testing.T) {
	// arrange
	const want = 0x7f
	d, a := initPCA9501WithStubbedAdaptor()
	// act
	err := d.initialize()
	// assert
	require.NoError(t, err)
	assert.Equal(t, want, a.address)
}
