package i2c

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	"errors"
	"bytes"
)

var _ gobot.Driver = (*INA3221Driver)(nil)

func initTestINA3221Driver() (*INA3221Driver) {
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

//func TestINA3221Driver_Measurements(t *testing.T) {
//	d, a := initTestINA3221DriverWithStubbedAdaptor()
//	a.i2cReadImpl = func(b []byte) (int, error) {
//		buf := new(bytes.Buffer)
//		if a.written[len(a.written)-1] == INA3221_REG_BUSVOLTAGE_1 {
//			buf.Write([]byte{0x09, 0x33})
//		}
//		copy(b, buf.Bytes())
//		return buf.Len(), nil
//	}
//
//	d.Start()
//
//	bv, err := d.GetBusVoltage(INA3221Channel1)
//	gobottest.Assert(t, err, nil)
//	t.Logf("bv1: %f", bv)
//}
