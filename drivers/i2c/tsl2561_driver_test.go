package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
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

func TestTSL2561DriverStart(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.i2cReadImpl = idReader
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "TSL2561"), true)
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Refute(t, d.Connection(), nil)
}

func TestTSL2561StartConnectError(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestTSL2561DriverStartError(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestTSL2561DriverStartNotFound(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{1})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	gobottest.Assert(t, d.Start(), errors.New("TSL2561 device not found (0x1)"))
}

func TestTSL2561DriverHalt(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.i2cReadImpl = idReader

	d.Start()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "TSL2561"), true)
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
		WithTSL2561AutoGain)

	gobottest.Refute(t, device, nil)
	gobottest.Assert(t, device.autoGain, true)
	gobottest.Assert(t, device.integrationTime, TSL2561IntegrationTime101MS)
}

func TestTSL2561DriverMoreOptions(t *testing.T) {
	adaptor := newI2cTestAdaptor()

	device := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime101MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain16X)

	gobottest.Refute(t, device, nil)
	gobottest.Assert(t, device.autoGain, false)
	gobottest.Assert(t, device.gain, TSL2561Gain(TSL2561Gain16X))
}

func TestTSL2561DriverEvenMoreOptions(t *testing.T) {
	adaptor := newI2cTestAdaptor()

	device := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime13MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain1X)

	gobottest.Refute(t, device, nil)
	gobottest.Assert(t, device.autoGain, false)
	gobottest.Assert(t, device.gain, TSL2561Gain(TSL2561Gain1X))
	gobottest.Assert(t, device.integrationTime, TSL2561IntegrationTime13MS)
}

func TestTSL2561DriverYetEvenMoreOptions(t *testing.T) {
	adaptor := newI2cTestAdaptor()

	device := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime402MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)

	gobottest.Refute(t, device, nil)
	gobottest.Assert(t, device.autoGain, true)
	gobottest.Assert(t, device.integrationTime, TSL2561IntegrationTime402MS)
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

func TestTSL2561DriverGetDataWriteError(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.i2cReadImpl = idReader
	gobottest.Assert(t, d.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, _, err := d.getData()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestTSL2561DriverGetDataReadError(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	adaptor.i2cReadImpl = idReader
	gobottest.Assert(t, d.Start(), nil)

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, err := d.getData()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestTSL2561DriverGetLuminocity(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()

	// TODO: obtain real sensor data here for testing
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{77, 48})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	d.Start()
	bb, ir, err := d.GetLuminocity()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, bb, uint16(12365))
	gobottest.Assert(t, ir, uint16(12365))
	gobottest.Assert(t, d.CalculateLux(bb, ir), uint32(72))
}

func TestTSL2561DriverGetLuminocityAutoGain(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime402MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)

	// TODO: obtain real sensor data here for testing
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{77, 48})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	d.Start()
	bb, ir, err := d.GetLuminocity()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, bb, uint16(12365))
	gobottest.Assert(t, ir, uint16(12365))
	gobottest.Assert(t, d.CalculateLux(bb, ir), uint32(72))
}

func TestTSL2561SetIntegrationTimeError(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	d.Start()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetIntegrationTime(TSL2561IntegrationTime101MS), errors.New("write error"))
}

func TestTSL2561SetGainError(t *testing.T) {
	d, adaptor := initTestTSL2561Driver()
	d.Start()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetGain(TSL2561Gain16X), errors.New("write error"))
}

func TestTSL2561getHiLo13MS(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	gobottest.Assert(t, hi, uint16(tsl2561AgcTHi13MS))
	gobottest.Assert(t, lo, uint16(tsl2561AgcTLo13MS))
}

func TestTSL2561getHiLo101MS(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime101MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	gobottest.Assert(t, hi, uint16(tsl2561AgcTHi101MS))
	gobottest.Assert(t, lo, uint16(tsl2561AgcTLo101MS))
}

func TestTSL2561getHiLo402MS(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime402MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	gobottest.Assert(t, hi, uint16(tsl2561AgcTHi402MS))
	gobottest.Assert(t, lo, uint16(tsl2561AgcTLo402MS))
}

func TestTSL2561getClipScaling13MS(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	gobottest.Assert(t, c, uint16(tsl2561Clipping13MS))
	gobottest.Assert(t, s, uint32(tsl2561LuxCHScaleTInt0))
}

func TestTSL2561getClipScaling101MS(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime101MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	gobottest.Assert(t, c, uint16(tsl2561Clipping101MS))
	gobottest.Assert(t, s, uint32(tsl2561LuxChScaleTInt1))
}

func TestTSL2561getClipScaling402MS(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime402MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	gobottest.Assert(t, c, uint16(tsl2561Clipping402MS))
	gobottest.Assert(t, s, uint32(1<<tsl2561LuxChScale))
}

func TestTSL2561getBM(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	d := NewTSL2561Driver(adaptor,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	b, m := d.getBM(tsl2561LuxK1T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB1T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM1T))

	b, m = d.getBM(tsl2561LuxK2T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB2T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM2T))

	b, m = d.getBM(tsl2561LuxK3T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB3T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM3T))

	b, m = d.getBM(tsl2561LuxK4T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB4T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM4T))

	b, m = d.getBM(tsl2561LuxK5T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB5T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM5T))

	b, m = d.getBM(tsl2561LuxK6T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB6T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM6T))

	b, m = d.getBM(tsl2561LuxK7T)
	gobottest.Assert(t, b, uint32(tsl2561LuxB7T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM7T))

	b, m = d.getBM(tsl2561LuxK8T + 1)
	gobottest.Assert(t, b, uint32(tsl2561LuxB8T))
	gobottest.Assert(t, m, uint32(tsl2561LuxM8T))
}
