package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCA9685Driver)(nil)

// and also the PwmWriter and ServoWriter interfaces
var (
	_ gpio.PwmWriter   = (*PCA9685Driver)(nil)
	_ gpio.ServoWriter = (*PCA9685Driver)(nil)
)

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
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "PCA9685"))
	assert.Equal(t, 0x40, d.defaultAddress)
}

func TestPCA9685Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewPCA9685Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestPCA9685Start(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewPCA9685Driver(a)
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())
}

func TestPCA9685Halt(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())
	assert.Nil(t, d.Halt())
}

func TestPCA9685SetPWM(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())
	assert.Nil(t, d.SetPWM(0, 0, 256))
}

func TestPCA9685SetPWMError(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	assert.Error(t, d.SetPWM(0, 0, 256), "write error")
}

func TestPCA9685SetPWMFreq(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())

	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.SetPWMFreq(60))
}

func TestPCA9685SetPWMFreqReadError(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	assert.Error(t, d.SetPWMFreq(60), "read error")
}

func TestPCA9685SetPWMFreqWriteError(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	assert.Nil(t, d.Start())

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	assert.Error(t, d.SetPWMFreq(60), "write error")
}

func TestPCA9685Commands(t *testing.T) {
	d, a := initTestPCA9685DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	_ = d.Start()

	err := d.Command("PwmWrite")(map[string]interface{}{"pin": "1", "val": "1"})
	assert.Nil(t, err)

	err = d.Command("ServoWrite")(map[string]interface{}{"pin": "1", "val": "1"})
	assert.Nil(t, err)

	err = d.Command("SetPWM")(map[string]interface{}{"channel": "1", "on": "0", "off": "1024"})
	assert.Nil(t, err)

	err = d.Command("SetPWMFreq")(map[string]interface{}{"freq": "60"})
	assert.Nil(t, err)
}
