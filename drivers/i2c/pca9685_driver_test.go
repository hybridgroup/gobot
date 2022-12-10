package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCA9685Driver)(nil)

// and also the PwmWriter and ServoWriter interfaces
var _ gpio.PwmWriter = (*PCA9685Driver)(nil)
var _ gpio.ServoWriter = (*PCA9685Driver)(nil)

func initTestPCA9685DriverWithStubbedAdaptor() (*PCA9685Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewPCA9685Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewPCA9685Driver(t *testing.T) {
	var di interface{} = NewPCA9685Driver(newI2cTestAdaptor())
	d, ok := di.(*PCA9685Driver)
	if !ok {
		t.Errorf("NewPCA9685Driver() should have returned a *PCA9685Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "PCA9685"), true)
}

func TestPCA9685Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewPCA9685Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestPCA9685Start(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewPCA9685Driver(a)
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)
}

func TestPCA9685Halt(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestPCA9685SetPWM(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.SetPWM(0, 0, 256), nil)
}

func TestPCA9685SetPWMError(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetPWM(0, 0, 256), errors.New("write error"))
}

func TestPCA9685SetPWMFreq(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.SetPWMFreq(60), nil)
}

func TestPCA9685SetPWMFreqReadError(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, d.SetPWMFreq(60), errors.New("read error"))
}

func TestPCA9685SetPWMFreqWriteError(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, d.Start(), nil)

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetPWMFreq(60), errors.New("write error"))
}

func TestPCA9685Commands(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	d.Start()

	err := d.Command("PwmWrite")(map[string]interface{}{"pin": "1", "val": "1"})
	gobottest.Assert(t, err, nil)

	err = d.Command("ServoWrite")(map[string]interface{}{"pin": "1", "val": "1"})
	gobottest.Assert(t, err, nil)

	err = d.Command("SetPWM")(map[string]interface{}{"channel": "1", "on": "0", "off": "1024"})
	gobottest.Assert(t, err, nil)

	err = d.Command("SetPWMFreq")(map[string]interface{}{"freq": "60"})
	gobottest.Assert(t, err, nil)
}
