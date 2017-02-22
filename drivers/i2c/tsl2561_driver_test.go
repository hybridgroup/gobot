package i2c

import (
	"bytes"
	"encoding/binary"
	"strings"
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

	gobottest.Assert(t, strings.HasPrefix(d.Name(), "TSL2561"), true)

	adaptor.i2cReadImpl = idReader

	gobottest.Assert(t, d.Start(), nil)

	gobottest.Assert(t, d.Halt(), nil)
}

func TestTSL2561DriverRead16(t *testing.T) {
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

func TestTSL2561DriverValidOptions(t *testing.T) {
	adaptor := newI2cTestAdaptor()

	device := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime101MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain16X,
		WithTSL2561AutoGain)

	gobottest.Refute(t, device, nil)
}

func TestTSL2561DriverSetName(t *testing.T) {
	d, _ := initTestTSL2561Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestTSL2561DriverOptions(t *testing.T) {
	d := NewTSL2561Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
