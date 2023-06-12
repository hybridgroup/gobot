package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initDriverWithStubbedAdaptor() (*Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewDriver(a, "I2C_BASIC", 0x15)
	return d, a
}

func initTestDriver() *Driver {
	d, _ := initDriverWithStubbedAdaptor()
	return d
}

func TestNewDriver(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	// act
	var di interface{} = NewDriver(a, "I2C_BASIC", 0x15)
	// assert
	d, ok := di.(*Driver)
	if !ok {
		t.Errorf("NewDriver() should have returned a *Driver")
	}
	gobottest.Assert(t, strings.Contains(d.name, "I2C_BASIC"), true)
	gobottest.Assert(t, d.defaultAddress, 0x15)
	gobottest.Assert(t, d.connector, a)
	gobottest.Assert(t, d.connection, nil)
	gobottest.Assert(t, d.afterStart(), nil)
	gobottest.Assert(t, d.beforeHalt(), nil)
	gobottest.Refute(t, d.Config, nil)
	gobottest.Refute(t, d.Commander, nil)
	gobottest.Refute(t, d.mutex, nil)
}

func TestSetName(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act
	d.SetName("TESTME")
	// assert
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestConnection(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	gobottest.Refute(t, d.Connection(), nil)
}

func TestStart(t *testing.T) {
	// arrange
	d, a := initDriverWithStubbedAdaptor()
	// act, assert
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, 0x15, a.address)
}

func TestStartConnectError(t *testing.T) {
	// arrange
	d, a := initDriverWithStubbedAdaptor()
	a.Testi2cConnectErr(true)
	// act, assert
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestHalt(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	gobottest.Assert(t, d.Halt(), nil)
}

func TestWrite(t *testing.T) {
	// arrange
	const (
		address     = "82"
		wantAddress = uint8(0x52)
		value       = 0x25
	)
	d, a := initDriverWithStubbedAdaptor()
	_ = d.Start()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// act
	err := d.Write(address, value)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, numCallsWrite, 1)
	gobottest.Assert(t, a.written[0], wantAddress)
	gobottest.Assert(t, a.written[1], uint8(value))
}

func TestRead(t *testing.T) {
	// arrange
	const (
		address     = "83"
		wantAddress = uint8(0x53)
		want        = uint8(0x44)
	)
	d, a := initDriverWithStubbedAdaptor()
	_ = d.Start()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func(b []byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[0] = want
		return len(b), nil
	}
	// act
	val, err := d.Read(address)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, int(want))
	gobottest.Assert(t, numCallsWrite, 1)
	gobottest.Assert(t, a.written[0], wantAddress)
	gobottest.Assert(t, numCallsRead, 1)
}
