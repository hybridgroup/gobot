package i2c

import (
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCF8591Driver)(nil)

func initTestPCF8591DriverWithStubbedAdaptor() (*PCF8591Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewPCF8591Driver(a, WithPCF8591With400kbitStabilization(0, 2))
	d.lastCtrlByte = 0xFF // prevent skipping of write
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewPCF8591Driver(t *testing.T) {
	var di interface{} = NewPCF8591Driver(newI2cTestAdaptor())
	d, ok := di.(*PCF8591Driver)
	if !ok {
		t.Errorf("NewPCF8591Driver() should have returned a *PCF8591Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "PCF8591"), true)
	gobottest.Assert(t, d.defaultAddress, 0x48)
}

func TestPCF8591Start(t *testing.T) {
	d := NewPCF8591Driver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestPCF8591Halt(t *testing.T) {
	d := NewPCF8591Driver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}

func TestPCF8591WithPCF8591With400kbitStabilization(t *testing.T) {
	d := NewPCF8591Driver(newI2cTestAdaptor(), WithPCF8591With400kbitStabilization(5, 6))
	gobottest.Assert(t, d.additionalReadWrite, uint8(5))
	gobottest.Assert(t, d.additionalRead, uint8(6))
}

func TestPCF8591AnalogReadSingle(t *testing.T) {
	// sequence to read the input channel:
	// * prepare value (with channel and mode) and write control register
	// * read 3 values to drop (see description in implementation)
	// * read the analog value
	//
	// arrange
	d, a := initTestPCF8591DriverWithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	description := "s.1"
	d.lastCtrlByte = 0x00
	ctrlByteOn := uint8(pcf8591_ALLSINGLE) | uint8(pcf8591_CHAN1)
	returnRead := []uint8{0x01, 0x02, 0x03, 0xFF}
	want := int(returnRead[3])
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b = returnRead[0:len(b)]
		}
		if numCallsRead == 2 {
			b[0] = returnRead[len(returnRead)-1]
		}
		return len(b), nil
	}
	// act
	got, err := d.AnalogRead(description)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 1)
	gobottest.Assert(t, a.written[0], ctrlByteOn)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, got, want)
}

func TestPCF8591AnalogReadDiff(t *testing.T) {
	// sequence to read the input channel:
	// * prepare value (with channel and mode) and write control register
	// * read 3 values to drop (see description in implementation)
	// * read the analog value
	// * convert to 8-bit two's complement (-127...128)
	//
	// arrange
	d, a := initTestPCF8591DriverWithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	description := "m.2-3"
	d.lastCtrlByte = 0x00
	ctrlByteOn := uint8(pcf8591_MIXED) | uint8(pcf8591_CHAN2)
	// some two' complements
	// 0x80 => -128
	// 0xFF => -1
	// 0x00 => 0
	// 0x7F => 127
	returnRead := []uint8{0x01, 0x02, 0x03, 0xFF}
	want := -1
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b = returnRead[0:len(b)]
		}
		if numCallsRead == 2 {
			b[0] = returnRead[len(returnRead)-1]
		}
		return len(b), nil
	}
	// act
	got, err := d.AnalogRead(description)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 1)
	gobottest.Assert(t, a.written[0], ctrlByteOn)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, got, want)
}

func TestPCF8591AnalogWrite(t *testing.T) {
	// sequence to write the output:
	// * create new value for the control register (ANAON)
	// * write the control register and value
	//
	// arrange
	d, a := initTestPCF8591DriverWithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	d.lastCtrlByte = 0x00
	d.lastAnaOut = 0x00
	ctrlByteOn := uint8(pcf8591_ANAON)
	want := uint8(0x15)
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// act
	err := d.AnalogWrite("", int(want))
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 2)
	gobottest.Assert(t, a.written[0], ctrlByteOn)
	gobottest.Assert(t, a.written[1], want)
}

func TestPCF8591AnalogOutputState(t *testing.T) {
	// sequence to set the state:
	// * create the new value (ctrlByte) for the control register (ANAON)
	// * write the register value
	//
	// arrange
	d, a := initTestPCF8591DriverWithStubbedAdaptor()
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		d.lastCtrlByte = uint8(0x00)
		wantCtrlByteVal := uint8(pcf8591_ANAON)
		if bitState == 0 {
			d.lastCtrlByte = uint8(0xFF)
			wantCtrlByteVal = uint8(0xFF & ^pcf8591_ANAON)
		}
		// act
		err := d.AnalogOutputState(bitState == 1)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(a.written), 1)
		gobottest.Assert(t, a.written[0], wantCtrlByteVal)
	}
}
