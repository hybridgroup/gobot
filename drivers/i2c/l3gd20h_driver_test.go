package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*HMC6352Driver)(nil)

// --------- HELPERS
func initTestL3GD20HDriver() (driver *L3GD20HDriver) {
	driver, _ = initTestL3GD20HDriverWithStubbedAdaptor()
	return
}

func initTestL3GD20HDriverWithStubbedAdaptor() (*L3GD20HDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewL3GD20HDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewL3GD20HDriver(t *testing.T) {
	// Does it return a pointer to an instance of HMC6352Driver?
	var d interface{} = NewL3GD20HDriver(newI2cTestAdaptor())
	_, ok := d.(*L3GD20HDriver)
	if !ok {
		t.Errorf("NewL3GD20HDriver() should have returned a *L3GD20HDriver")
	}
}

func TestL3GD20HDriver(t *testing.T) {
	d := initTestL3GD20HDriver()
	gobottest.Refute(t, d.Connection(), nil)
}

// Methods
func TestL3GD20HDriverStart(t *testing.T) {
	d, _ := initTestL3GD20HDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.Start(), nil)
}

func TestL3GD20HStartConnectError(t *testing.T) {
	d, adaptor := initTestL3GD20HDriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestL3GD20HDriverStartWriteError(t *testing.T) {
	d, adaptor := initTestL3GD20HDriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestL3GD20HDriverHalt(t *testing.T) {
	d := initTestL3GD20HDriver()

	gobottest.Assert(t, d.Halt(), nil)
}

func TestL3GD20HDriverScale(t *testing.T) {
	d := initTestL3GD20HDriver()
	gobottest.Assert(t, d.Scale(), L3GD20HScale250dps)
	gobottest.Assert(t, d.getSensitivity(), float32(0.00875))

	d.SetScale(L3GD20HScale500dps)
	gobottest.Assert(t, d.Scale(), L3GD20HScale500dps)
	gobottest.Assert(t, d.getSensitivity(), float32(0.0175))

	d.SetScale(L3GD20HScale2000dps)
	gobottest.Assert(t, d.Scale(), L3GD20HScale2000dps)
	gobottest.Assert(t, d.getSensitivity(), float32(0.07))
}

func TestL3GD20HDriverMeasurement(t *testing.T) {
	d, adaptor := initTestL3GD20HDriverWithStubbedAdaptor()
	rawX := 5
	rawY := 8
	rawZ := -3
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, int16(rawX))
		binary.Write(buf, binary.LittleEndian, int16(rawY))
		binary.Write(buf, binary.LittleEndian, int16(rawZ))
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	d.Start()
	x, y, z, err := d.XYZ()
	gobottest.Assert(t, err, nil)
	var sensitivity float32 = 0.00875
	gobottest.Assert(t, x, float32(rawX)*sensitivity)
	gobottest.Assert(t, y, float32(rawY)*sensitivity)
	gobottest.Assert(t, z, float32(rawZ)*sensitivity)
}

func TestL3GD20HDriverMeasurementError(t *testing.T) {
	d, adaptor := initTestL3GD20HDriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	d.Start()
	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestL3GD20HDriverMeasurementWriteError(t *testing.T) {
	d, adaptor := initTestL3GD20HDriverWithStubbedAdaptor()
	d.Start()
	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestL3GD20HDriverSetName(t *testing.T) {
	d := initTestL3GD20HDriver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestL3GD20HDriverOptions(t *testing.T) {
	d := NewL3GD20HDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
