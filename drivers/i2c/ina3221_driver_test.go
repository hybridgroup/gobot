package i2c

import (
	"testing"

	"errors"

	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*INA3221Driver)(nil)

func initTestINA3221Driver() *INA3221Driver {
	d, _ := initTestINA3221DriverWithStubbedAdaptor()
	return d
}

func initTestINA3221DriverWithStubbedAdaptor() (*INA3221Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	return NewINA3221Driver(a), a
}

func TestNewINA3221Driver(t *testing.T) {
	var d interface{} = NewINA3221Driver(newI2cTestAdaptor())
	if _, ok := d.(*INA3221Driver); !ok {
		t.Error("NewINA3221Driver() should return a *INA3221Driver")
	}
}

func TestINA3221Driver_Connection(t *testing.T) {
	d := initTestINA3221Driver()
	gobottest.Refute(t, d.Connection(), nil)
}

func TestINA3221Driver_Start(t *testing.T) {
	d := initTestINA3221Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestINA3221Driver_ConnectError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestINA3221Driver_StartWriteError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestINA3221Driver_Halt(t *testing.T) {
	d := initTestINA3221Driver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestINA3221DriverGetBusVoltage(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		// bus voltage sensor values from 12V battery
		copy(b, []byte{0x36, 0x68})
		return 2, nil
	}

	v, err := d.GetBusVoltage(INA3221Channel1)
	gobottest.Assert(t, v, float64(13.928))
	gobottest.Assert(t, err, nil)
}

func TestINA3221DriverGetBusVoltageReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetBusVoltage(INA3221Channel1)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestINA3221DriverGetShuntVoltage(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		// shunt voltage sensor values from 12V battery
		copy(b, []byte{0x05, 0xD8})
		return 2, nil
	}

	v, err := d.GetShuntVoltage(INA3221Channel1)
	gobottest.Assert(t, v, float64(7.48))
	gobottest.Assert(t, err, nil)
}

func TestINA3221DriverGetShuntVoltageReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetShuntVoltage(INA3221Channel1)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestINA3221DriverGetCurrent(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		// shunt voltage sensor values from 12V battery
		copy(b, []byte{0x05, 0x0D8})
		return 2, nil
	}

	v, err := d.GetCurrent(INA3221Channel1)
	gobottest.Assert(t, v, float64(74.8))
	gobottest.Assert(t, err, nil)
}

func TestINA3221DriverCurrentReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetCurrent(INA3221Channel1)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestINA3221DriverGetLoadVoltage(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	i := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		// TODO: return test data as read from actual sensor
		copy(b, []byte{0x36, 0x68, 0x05, 0xd8}[i:i+2])
		i += 2
		return 2, nil
	}

	v, err := d.GetLoadVoltage(INA3221Channel2)
	gobottest.Assert(t, v, float64(13.935480))
	gobottest.Assert(t, err, nil)
}

func TestINA3221DriverGetLoadVoltageReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetLoadVoltage(INA3221Channel2)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestINA3221DriverName(t *testing.T) {
	d := initTestINA3221Driver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "INA3221"), true)
}

func TestINA3221DriverSetName(t *testing.T) {
	d := initTestINA3221Driver()
	d.SetName("foobot")
	gobottest.Assert(t, d.Name(), "foobot")
}
