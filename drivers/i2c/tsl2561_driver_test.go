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

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*TSL2561Driver)(nil)

func testIdReader(b []byte) (int, error) {
	buf := new(bytes.Buffer)
	// Mock device responding 0xA
	binary.Write(buf, binary.LittleEndian, uint8(0x0A))
	copy(b, buf.Bytes())
	return buf.Len(), nil
}

func initTestTSL2561Driver() (*TSL2561Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a)
	a.i2cReadImpl = testIdReader
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewTSL2561Driver(t *testing.T) {
	var di interface{} = NewTSL2561Driver(newI2cTestAdaptor())
	d, ok := di.(*TSL2561Driver)
	if !ok {
		t.Errorf("NewTSL2561Driver() should have returned a *TSL2561Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "TSL2561"), true)
	gobottest.Assert(t, d.defaultAddress, 0x39)
	gobottest.Assert(t, d.autoGain, false)
	gobottest.Assert(t, d.gain, TSL2561Gain(0))
	gobottest.Assert(t, d.integrationTime, TSL2561IntegrationTime(2))
}

func TestTSL2561DriverOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewTSL2561Driver(newI2cTestAdaptor(), WithBus(2), WithTSL2561AutoGain)
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.autoGain, true)
}

func TestTSL2561DriverStart(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a)
	a.i2cReadImpl = testIdReader

	gobottest.Assert(t, d.Start(), nil)
}

func TestTSL2561DriverStartNotFound(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a)
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{1})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	gobottest.Assert(t, d.Start(), errors.New("TSL2561 device not found (0x1)"))
}

func TestTSL2561DriverHalt(t *testing.T) {
	d, _ := initTestTSL2561Driver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestTSL2561DriverRead16(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cReadImpl = testIdReader
	a.i2cReadImpl = func(b []byte) (int, error) {
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
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)

	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.autoGain, true)
	gobottest.Assert(t, d.integrationTime, TSL2561IntegrationTime101MS)
}

func TestTSL2561DriverMoreOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain16X)

	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.autoGain, false)
	gobottest.Assert(t, d.gain, TSL2561Gain(TSL2561Gain16X))
}

func TestTSL2561DriverEvenMoreOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain1X)

	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.autoGain, false)
	gobottest.Assert(t, d.gain, TSL2561Gain(TSL2561Gain1X))
	gobottest.Assert(t, d.integrationTime, TSL2561IntegrationTime13MS)
}

func TestTSL2561DriverYetEvenMoreOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)

	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.autoGain, true)
	gobottest.Assert(t, d.integrationTime, TSL2561IntegrationTime402MS)
}

func TestTSL2561DriverGetDataWriteError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, _, err := d.getData()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestTSL2561DriverGetDataReadError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, err := d.getData()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestTSL2561DriverGetLuminocity(t *testing.T) {
	d, a := initTestTSL2561Driver()
	// TODO: obtain real sensor data here for testing
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{77, 48})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bb, ir, err := d.GetLuminocity()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, bb, uint16(12365))
	gobottest.Assert(t, ir, uint16(12365))
	gobottest.Assert(t, d.CalculateLux(bb, ir), uint32(72))
}

func TestTSL2561DriverGetLuminocityAutoGain(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)
	// TODO: obtain real sensor data here for testing
	a.i2cReadImpl = func(b []byte) (int, error) {
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
	d, a := initTestTSL2561Driver()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetIntegrationTime(TSL2561IntegrationTime101MS), errors.New("write error"))
}

func TestTSL2561SetGainError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetGain(TSL2561Gain16X), errors.New("write error"))
}

func TestTSL2561getHiLo13MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	gobottest.Assert(t, hi, uint16(tsl2561AgcTHi13MS))
	gobottest.Assert(t, lo, uint16(tsl2561AgcTLo13MS))
}

func TestTSL2561getHiLo101MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	gobottest.Assert(t, hi, uint16(tsl2561AgcTHi101MS))
	gobottest.Assert(t, lo, uint16(tsl2561AgcTLo101MS))
}

func TestTSL2561getHiLo402MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	gobottest.Assert(t, hi, uint16(tsl2561AgcTHi402MS))
	gobottest.Assert(t, lo, uint16(tsl2561AgcTLo402MS))
}

func TestTSL2561getClipScaling13MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	gobottest.Assert(t, c, uint16(tsl2561Clipping13MS))
	gobottest.Assert(t, s, uint32(tsl2561LuxCHScaleTInt0))
}

func TestTSL2561getClipScaling101MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	gobottest.Assert(t, c, uint16(tsl2561Clipping101MS))
	gobottest.Assert(t, s, uint32(tsl2561LuxChScaleTInt1))
}

func TestTSL2561getClipScaling402MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	gobottest.Assert(t, c, uint16(tsl2561Clipping402MS))
	gobottest.Assert(t, s, uint32(1<<tsl2561LuxChScale))
}

func TestTSL2561getBM(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
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
