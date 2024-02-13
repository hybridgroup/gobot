package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*TSL2561Driver)(nil)

func testIDReader(b []byte) (int, error) {
	buf := new(bytes.Buffer)
	// Mock device responding 0xA
	_ = binary.Write(buf, binary.LittleEndian, uint8(0x0A))
	copy(b, buf.Bytes())
	return buf.Len(), nil
}

func initTestTSL2561Driver() (*TSL2561Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a)
	a.i2cReadImpl = testIDReader
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewTSL2561Driver(t *testing.T) {
	var di interface{} = NewTSL2561Driver(newI2cTestAdaptor())
	d, ok := di.(*TSL2561Driver)
	if !ok {
		require.Fail(t, "NewTSL2561Driver() should have returned a *TSL2561Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "TSL2561"))
	assert.Equal(t, 0x39, d.defaultAddress)
	assert.False(t, d.autoGain)
	assert.Equal(t, TSL2561Gain(0), d.gain)
	assert.Equal(t, TSL2561IntegrationTime(2), d.integrationTime)
}

func TestTSL2561DriverOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewTSL2561Driver(newI2cTestAdaptor(), WithBus(2), WithTSL2561AutoGain)
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.True(t, d.autoGain)
}

func TestTSL2561DriverStart(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a)
	a.i2cReadImpl = testIDReader

	require.NoError(t, d.Start())
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
	require.ErrorContains(t, d.Start(), "TSL2561 device not found (0x1)")
}

func TestTSL2561DriverHalt(t *testing.T) {
	d, _ := initTestTSL2561Driver()
	require.NoError(t, d.Halt())
}

func TestTSL2561DriverRead16(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cReadImpl = testIDReader
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// send low
		_ = binary.Write(buf, binary.LittleEndian, uint8(0xEA))
		// send high
		_ = binary.Write(buf, binary.LittleEndian, uint8(0xAE))
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	val, err := d.connection.ReadWordData(1)
	require.NoError(t, err)
	assert.Equal(t, uint16(0xAEEA), val)
}

func TestTSL2561DriverValidOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)

	assert.NotNil(t, d)
	assert.True(t, d.autoGain)
	assert.Equal(t, TSL2561IntegrationTime101MS, d.integrationTime)
}

func TestTSL2561DriverMoreOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain16X)

	assert.NotNil(t, d)
	assert.False(t, d.autoGain)
	assert.Equal(t, TSL2561Gain(TSL2561Gain16X), d.gain)
}

func TestTSL2561DriverEvenMoreOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561Gain1X)

	assert.NotNil(t, d)
	assert.False(t, d.autoGain)
	assert.Equal(t, TSL2561Gain1X, d.gain)
	assert.Equal(t, TSL2561IntegrationTime13MS, d.integrationTime)
}

func TestTSL2561DriverYetEvenMoreOptions(t *testing.T) {
	a := newI2cTestAdaptor()

	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithAddress(TSL2561AddressLow),
		WithTSL2561AutoGain)

	assert.NotNil(t, d)
	assert.True(t, d.autoGain)
	assert.Equal(t, TSL2561IntegrationTime402MS, d.integrationTime)
}

func TestTSL2561DriverGetDataWriteError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, _, err := d.getData()
	require.ErrorContains(t, err, "write error")
}

func TestTSL2561DriverGetDataReadError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, err := d.getData()
	require.ErrorContains(t, err, "read error")
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
	require.NoError(t, err)
	assert.Equal(t, uint16(12365), bb)
	assert.Equal(t, uint16(12365), ir)
	assert.Equal(t, uint32(72), d.CalculateLux(bb, ir))
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

	_ = d.Start()
	bb, ir, err := d.GetLuminocity()
	require.NoError(t, err)
	assert.Equal(t, uint16(12365), bb)
	assert.Equal(t, uint16(12365), ir)
	assert.Equal(t, uint32(72), d.CalculateLux(bb, ir))
}

func TestTSL2561SetIntegrationTimeError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.SetIntegrationTime(TSL2561IntegrationTime101MS), "write error")
}

func TestTSL2561SetGainError(t *testing.T) {
	d, a := initTestTSL2561Driver()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.SetGain(TSL2561Gain16X), "write error")
}

func TestTSL2561getHiLo13MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	assert.Equal(t, uint16(tsl2561AgcTHi13MS), hi)
	assert.Equal(t, uint16(tsl2561AgcTLo13MS), lo)
}

func TestTSL2561getHiLo101MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	assert.Equal(t, uint16(tsl2561AgcTHi101MS), hi)
	assert.Equal(t, uint16(tsl2561AgcTLo101MS), lo)
}

func TestTSL2561getHiLo402MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithTSL2561AutoGain)

	hi, lo := d.getHiLo()
	assert.Equal(t, uint16(tsl2561AgcTHi402MS), hi)
	assert.Equal(t, uint16(tsl2561AgcTLo402MS), lo)
}

func TestTSL2561getClipScaling13MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	assert.Equal(t, uint16(tsl2561Clipping13MS), c)
	assert.Equal(t, uint32(tsl2561LuxCHScaleTInt0), s)
}

func TestTSL2561getClipScaling101MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime101MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	assert.Equal(t, uint16(tsl2561Clipping101MS), c)
	assert.Equal(t, uint32(tsl2561LuxChScaleTInt1), s)
}

func TestTSL2561getClipScaling402MS(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime402MS,
		WithTSL2561AutoGain)

	c, s := d.getClipScaling()
	d.waitForADC()
	assert.Equal(t, uint16(tsl2561Clipping402MS), c)
	assert.Equal(t, uint32(1<<tsl2561LuxChScale), s)
}

func TestTSL2561getBM(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewTSL2561Driver(a,
		WithTSL2561IntegrationTime13MS,
		WithTSL2561AutoGain)

	b, m := d.getBM(tsl2561LuxK1T)
	assert.Equal(t, uint32(tsl2561LuxB1T), b)
	assert.Equal(t, uint32(tsl2561LuxM1T), m)

	b, m = d.getBM(tsl2561LuxK2T)
	assert.Equal(t, uint32(tsl2561LuxB2T), b)
	assert.Equal(t, uint32(tsl2561LuxM2T), m)

	b, m = d.getBM(tsl2561LuxK3T)
	assert.Equal(t, uint32(tsl2561LuxB3T), b)
	assert.Equal(t, uint32(tsl2561LuxM3T), m)

	b, m = d.getBM(tsl2561LuxK4T)
	assert.Equal(t, uint32(tsl2561LuxB4T), b)
	assert.Equal(t, uint32(tsl2561LuxM4T), m)

	b, m = d.getBM(tsl2561LuxK5T)
	assert.Equal(t, uint32(tsl2561LuxB5T), b)
	assert.Equal(t, uint32(tsl2561LuxM5T), m)

	b, m = d.getBM(tsl2561LuxK6T)
	assert.Equal(t, uint32(tsl2561LuxB6T), b)
	assert.Equal(t, uint32(tsl2561LuxM6T), m)

	b, m = d.getBM(tsl2561LuxK7T)
	assert.Equal(t, uint32(tsl2561LuxB7T), b)
	assert.Equal(t, uint32(tsl2561LuxM7T), m)

	b, m = d.getBM(tsl2561LuxK8T + 1)
	assert.Equal(t, uint32(tsl2561LuxB8T), b)
	assert.Equal(t, uint32(tsl2561LuxM8T), m)
}
