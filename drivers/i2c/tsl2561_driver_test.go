package i2c

import (
	"bytes"
	"encoding/binary"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*TSL2561Driver)(nil)

func initTestTSL2561Driver() (*TSL2561Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewTSL2561Driver(adaptor), adaptor
}

func idReader(b []byte) (int, error) {
	buf := new(bytes.Buffer)
	// Mock device responding 0xA
	binary.Write(buf, binary.LittleEndian, uint8(0x0A))
	copy(b, buf.Bytes())
	return buf.Len(), nil
}

func TestTSL2561Driver(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()

	gobottest.Assert(t, d.Name(), "TSL2561")

	adaptor.i2cReadImpl = idReader

	gobottest.Assert(t, d.Start(), nil)

	gobottest.Assert(t, d.Halt(), nil)
}

func TestRead16(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()

	adaptor.i2cReadImpl = idReader

	gobottest.Assert(t, d.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// send low
		binary.Write(buf, binary.LittleEndian, uint8(0xEA))
		// send high
		binary.Write(buf, binary.LittleEndian, uint8(0xAE))
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	val, err := d.connection.ReadWordData(1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, uint16(0xAEEA))
}

func TestBadOption(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	options := map[string]int{
		"hej": 12,
	}

	defer func() {
		x := recover()
		gobottest.Refute(t, x, nil)
	}()

	device := NewTSL2561Driver(adaptor, options)

	gobottest.Refute(t, device, nil)
}

func TestBadOptionValue(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	options := map[string]int{
		"integrationTime": 47,
	}

	defer func() {
		x := recover()
		gobottest.Refute(t, x, nil)
	}()

	device := NewTSL2561Driver(adaptor, options)

	gobottest.Refute(t, device, nil)
}

func TestValidOptions(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	options := map[string]int{
		"integrationTime": int(TSL2561IntegrationTime101MS),
		"address":         TSL2561AddressLow,
		"gain":            TSL2561Gain16X,
		"autoGain":        1,
	}

	device := NewTSL2561Driver(adaptor, options)

	gobottest.Refute(t, device, nil)
}
