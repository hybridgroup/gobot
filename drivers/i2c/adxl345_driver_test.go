package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*ADXL345Driver)(nil)

// --------- HELPERS
func initTestADXL345Driver() (driver *ADXL345Driver) {
	driver, _ = initTestADXL345DriverWithStubbedAdaptor()
	return
}

func initTestADXL345DriverWithStubbedAdaptor() (*ADXL345Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewADXL345Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewADXL345Driver(t *testing.T) {
	// Does it return a pointer to an instance of ADXL345Driver?
	var mma interface{} = NewADXL345Driver(newI2cTestAdaptor())
	_, ok := mma.(*ADXL345Driver)
	if !ok {
		t.Errorf("NewADXL345Driver() should have returned a *ADXL345Driver")
	}
}

// Methods
func TestADXL345Driver(t *testing.T) {
	mma := initTestADXL345Driver()

	gobottest.Refute(t, mma.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(mma.Name(), "ADXL345"), true)
}

func TestADXL345DriverSetName(t *testing.T) {
	d := initTestADXL345Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestADXL345DriverOptions(t *testing.T) {
	d := NewADXL345Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestADXL345DriverStart(t *testing.T) {
	d := initTestADXL345Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestADXL345StartConnectError(t *testing.T) {
	d, adaptor := initTestADXL345DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestADXL345DriverStartWriteError(t *testing.T) {
	mma, adaptor := initTestADXL345DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, mma.Start(), errors.New("write error"))
}

func TestADXL345DriverHalt(t *testing.T) {
	d := initTestADXL345Driver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestADXL345DriverNullXYZ(t *testing.T) {
	d, _ := initTestADXL345DriverWithStubbedAdaptor()
	d.Start()
	x, y, z, _ := d.XYZ()
	gobottest.Assert(t, x, 0.0)
	gobottest.Assert(t, y, 0.0)
	gobottest.Assert(t, z, 0.0)
}

func TestADXL345DriverXYZ(t *testing.T) {
	d, adaptor := initTestADXL345DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{218, 0, 251, 255, 100, 0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	x, y, z, _ := d.XYZ()
	gobottest.Assert(t, x, 0.8515625)
	gobottest.Assert(t, y, -0.01953125)
	gobottest.Assert(t, z, 0.390625)
}

func TestADXL345DriverXYZError(t *testing.T) {
	d, adaptor := initTestADXL345DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("read error"))
}


func TestADXL345DriverRawXYZ(t *testing.T) {
	d, adaptor := initTestADXL345DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{218, 0, 251, 255, 100, 0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	x, y, z, _ := d.RawXYZ()
	gobottest.Assert(t, int(x), 218)
	gobottest.Assert(t, int(y), -5)
	gobottest.Assert(t, int(z), 100)
}

func TestADXL345DriverRawXYZError(t *testing.T) {
	d, adaptor := initTestADXL345DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, _, err := d.RawXYZ()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestADXL345DriverSetRange(t *testing.T) {
	d := initTestADXL345Driver()
	d.Start()
	gobottest.Assert(t, d.SetRange(ADXL345_RANGE_16G), nil)
}
