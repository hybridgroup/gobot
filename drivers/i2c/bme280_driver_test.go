package i2c

import (
	"bytes"
	"encoding/binary"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BME280Driver)(nil)

// --------- HELPERS
func initTestBME280Driver() (driver *BME280Driver) {
	driver, _ = initTestBME280DriverWithStubbedAdaptor()
	return
}

func initTestBME280DriverWithStubbedAdaptor() (*BME280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBME280Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBME280Driver(t *testing.T) {
	// Does it return a pointer to an instance of BME280Driver?
	var bme280 interface{} = NewBME280Driver(newI2cTestAdaptor())
	_, ok := bme280.(*BME280Driver)
	if !ok {
		t.Errorf("NewBME280Driver() should have returned a *BME280Driver")
	}
}

func TestBME280Driver(t *testing.T) {
	bme280 := initTestBME280Driver()
	gobottest.Refute(t, bme280.Connection(), nil)
}

func TestBME280DriverStart(t *testing.T) {
	bme280, _ := initTestBME280DriverWithStubbedAdaptor()
	gobottest.Assert(t, bme280.Start(), nil)
}

func TestBME280DriverHalt(t *testing.T) {
	bme280 := initTestBME280Driver()

	gobottest.Assert(t, bme280.Halt(), nil)
}

func TestBME280DriverMeasurements(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		if adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB {
			binary.Write(buf, binary.BigEndian, int16(408))
			binary.Write(buf, binary.BigEndian, int16(-72))
			binary.Write(buf, binary.BigEndian, int16(-14383))
			binary.Write(buf, binary.BigEndian, uint16(32741))
			binary.Write(buf, binary.BigEndian, uint16(32757))
			binary.Write(buf, binary.BigEndian, uint16(23153))
			binary.Write(buf, binary.BigEndian, int16(6190))
			binary.Write(buf, binary.BigEndian, int16(4))
			binary.Write(buf, binary.BigEndian, int16(-32768))
			binary.Write(buf, binary.BigEndian, int16(-8711))
			binary.Write(buf, binary.BigEndian, int16(2868))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterTempMSB {
			binary.Write(buf, binary.BigEndian, int16(27898))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterPressureMSB {
			binary.Write(buf, binary.BigEndian, int16(23843))
			// XLSB, not used in this test.
			buf.WriteByte(0)
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bme280.Start()
	temp, err := bme280.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(15.0))
	pressure, err := bme280.Pressure(BMP180UltraLowPower)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(69964))
}

func TestBME280DriverSetName(t *testing.T) {
	b := initTestBME280Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBME280DriverOptions(t *testing.T) {
	b := NewBME280Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
