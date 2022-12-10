package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
		t.Errorf("NewPCA9501Driver() should have returned a *PCA9501Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "PCA9501"), true)
}

func TestPCA9501Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewPCA9501Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
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
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
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
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
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
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
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
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestPCA9501WriteGPIO(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 2)
			gobottest.Assert(t, len(a.written), 2)
			gobottest.Assert(t, a.written[0], tc.wantPin)
			gobottest.Assert(t, a.written[1], tc.wantState)
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
	gobottest.Assert(t, err, wantErr)
	gobottest.Assert(t, numCallsRead < 2, true)
	gobottest.Assert(t, numCallsWrite, 1)
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
	gobottest.Assert(t, err, wantErr)
	gobottest.Assert(t, numCallsWrite, 2)
}

func TestPCA9501ReadGPIO(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, numCallsRead, 2)
			gobottest.Assert(t, len(a.written), 1)
			gobottest.Assert(t, a.written[0], wantCtrlState)
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
	gobottest.Assert(t, err, wantErr)
	gobottest.Assert(t, numCallsRead, 1)
	gobottest.Assert(t, numCallsWrite, 0)
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
	gobottest.Assert(t, err, wantErr)
	gobottest.Assert(t, numCallsWrite, 1)
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, numCallsWrite, 1)
	gobottest.Assert(t, a.written[0], addressEEPROM)
	gobottest.Assert(t, a.written[1], want)
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, want)
	gobottest.Assert(t, numCallsWrite, 1)
	gobottest.Assert(t, a.written[0], addressEEPROM)
	gobottest.Assert(t, numCallsRead, 1)
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
	gobottest.Assert(t, err, wantErr)
	gobottest.Assert(t, numCallsRead, 0)
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
	gobottest.Assert(t, numCallsWrite, 1)
	gobottest.Assert(t, err, wantErr)
}

func TestPCA9501_initialize(t *testing.T) {
	// arrange
	const want = 0x7f
	d, a := initPCA9501WithStubbedAdaptor()
	// act
	err := d.initialize()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.address, want)
}
