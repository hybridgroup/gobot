//nolint:forcetypeassert // ok here
package i2c

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCF8583Driver)(nil)

func initTestPCF8583WithStubbedAdaptor() (*PCF8583Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewPCF8583Driver(a)
	_ = d.Start()
	return d, a
}

func TestNewPCF8583Driver(t *testing.T) {
	var di interface{} = NewPCF8583Driver(newI2cTestAdaptor())
	d, ok := di.(*PCF8583Driver)
	if !ok {
		require.Fail(t, "NewPCF8583Driver() should have returned a *PCF8583Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.name, "PCF8583"))
	assert.Equal(t, 0x50, d.defaultAddress)
	assert.Equal(t, PCF8583Control(0x00), d.mode)
	assert.Equal(t, 0, d.yearOffset)
	assert.Equal(t, uint8(0x10), d.ramOffset)
}

func TestPCF8583Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewPCF8583Driver(newI2cTestAdaptor(), WithBus(2), WithPCF8583Mode(PCF8583CtrlModeClock50))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, PCF8583CtrlModeClock50, d.mode)
}

func TestPCF8583CommandsWriteTime(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	readCtrlState := uint8(0x10) // clock 50Hz
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	result := d.Command("WriteTime")(map[string]interface{}{"val": time.Now()})
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCF8583CommandsReadTime(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	d.yearOffset = 2019
	milliSec := 550 * time.Millisecond // 0.55 sec = 550 ms
	want := time.Date(2021, time.December, 24, 18, 0, 0, int(milliSec), time.UTC)
	reg0Val := uint8(0x00) // clock mode 32.768 kHz
	reg1Val := uint8(0x55) // BCD: 1/10 and 1/100 sec (55)
	reg2Val := uint8(0x00) // BCD: 10 and 1 sec (00)
	reg3Val := uint8(0x00) // BCD: 10 and 1 min (00)
	reg4Val := uint8(0x18) // BCD: 10 and 1 hour (18)
	reg5Val := uint8(0xA4) // year (2) and BCD: date (24)
	reg6Val := uint8(0xB2) // weekday 5, bit 5 and bit 7 (0xA0) and BCD: month (0x12)
	returnRead := [2][]uint8{
		{reg0Val},
		{reg1Val, reg2Val, reg3Val, reg4Val, reg5Val, reg6Val},
	}
	// arrange reads
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		rr := returnRead[numCallsRead-1]
		for i := 0; i < len(b); i++ {
			b[i] = rr[i]
		}
		return len(b), nil
	}
	// act
	result := d.Command("ReadTime")(map[string]interface{}{})
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
	assert.Equal(t, want, result.(map[string]interface{})["val"])
}

func TestPCF8583CommandsWriteCounter(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	readCtrlState := uint8(0x20) // counter
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	result := d.Command("WriteCounter")(map[string]interface{}{"val": int32(123456)})
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCF8583CommandsReadCounter(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	want := int32(123456)
	reg0Val := uint8(0x20) // counter mode
	reg1Val := uint8(0x56) // BCD: 56
	reg2Val := uint8(0x34) // BCD: 34
	reg3Val := uint8(0x12) // BCD: 12
	returnRead := [2][]uint8{
		{reg0Val},
		{reg1Val, reg2Val, reg3Val},
	}
	// arrange reads
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		rr := returnRead[numCallsRead-1]
		for i := 0; i < len(b); i++ {
			b[i] = rr[i]
		}
		return len(b), nil
	}
	// act
	result := d.Command("ReadCounter")(map[string]interface{}{})
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
	assert.Equal(t, want, result.(map[string]interface{})["val"])
}

func TestPCF8583CommandsWriteRAM(t *testing.T) {
	// arrange
	d, _ := initTestPCF8583WithStubbedAdaptor()
	addressValue := map[string]interface{}{
		"address": uint8(0x12),
		"val":     uint8(0x45),
	}
	// act
	result := d.Command("WriteRAM")(addressValue)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestPCF8583CommandsReadRAM(t *testing.T) {
	// arrange
	d, _ := initTestPCF8583WithStubbedAdaptor()
	address := map[string]interface{}{
		"address": uint8(0x34),
	}
	// act
	result := d.Command("ReadRAM")(address)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
	assert.Equal(t, uint8(0), result.(map[string]interface{})["val"])
}

func TestPCF8583WriteTime(t *testing.T) {
	// sequence to write the time:
	// * read control register for get current state and ensure an clock mode is set
	// * write the control register (stop counting)
	// * create the values for date registers (default is 24h mode)
	// * write the clock and calendar registers with auto increment
	// * write the control register (start counting)
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{}               // reset writes of Start() and former test
	readCtrlState := uint8(0x07)       // 32.768kHz clock mode
	milliSec := 210 * time.Millisecond // 0.21 sec = 210 ms
	initDate := time.Date(2022, time.December, 16, 15, 14, 13, int(milliSec), time.UTC)
	wantCtrlStop := uint8(0x87)  // stop counting bit is set
	wantReg1Val := uint8(0x21)   // BCD: 1/10 and 1/100 sec (21)
	wantReg2Val := uint8(0x13)   // BCD: 10 and 1 sec (13)
	wantReg3Val := uint8(0x14)   // BCD: 10 and 1 min (14)
	wantReg4Val := uint8(0x15)   // BCD: 10 and 1 hour (15)
	wantReg5Val := uint8(0x16)   // year (0) and BCD: date (16)
	wantReg6Val := uint8(0xB2)   // weekday 5, bit 5 and bit 7 (0xA0) and BCD: month (0x12)
	wantCrtlStart := uint8(0x07) // stop counting bit is reset
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	err := d.WriteTime(initDate)
	// assert
	require.NoError(t, err)
	assert.Equal(t, initDate.Year(), d.yearOffset)
	assert.Equal(t, 1, numCallsRead)
	assert.Len(t, a.written, 11)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[1])
	assert.Equal(t, wantCtrlStop, a.written[2])
	assert.Equal(t, wantReg1Val, a.written[3])
	assert.Equal(t, wantReg2Val, a.written[4])
	assert.Equal(t, wantReg3Val, a.written[5])
	assert.Equal(t, wantReg4Val, a.written[6])
	assert.Equal(t, wantReg5Val, a.written[7])
	assert.Equal(t, wantReg6Val, a.written[8])
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[9])
	assert.Equal(t, wantCrtlStart, a.written[10])
}

func TestPCF8583WriteTimeNoTimeModeFails(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{}         // reset writes of Start() and former test
	readCtrlState := uint8(0x30) // test mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	err := d.WriteTime(time.Now())
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "wrong mode 0x30")
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, 1, numCallsRead)
}

func TestPCF8583ReadTime(t *testing.T) {
	// sequence to read the time:
	// * read the control register to determine mask flag and ensure an clock mode is set
	// * read the clock and calendar registers with auto increment
	// * create the value out of registers content
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	d.yearOffset = 2020
	milliSec := 210 * time.Millisecond // 0.21 sec = 210 ms
	want := time.Date(2022, time.December, 16, 15, 14, 13, int(milliSec), time.UTC)
	reg0Val := uint8(0x10) // clock mode 50Hz
	reg1Val := uint8(0x21) // BCD: 1/10 and 1/100 sec (21)
	reg2Val := uint8(0x13) // BCD: 10 and 1 sec (13)
	reg3Val := uint8(0x14) // BCD: 10 and 1 min (14)
	reg4Val := uint8(0x15) // BCD: 10 and 1 hour (15)
	reg5Val := uint8(0x96) // year (2) and BCD: date (16)
	reg6Val := uint8(0xB2) // weekday 5, bit 5 and bit 7 (0xA0) and BCD: month (0x12)
	returnRead := [2][]uint8{
		{reg0Val},
		{reg1Val, reg2Val, reg3Val, reg4Val, reg5Val, reg6Val},
	}
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		rr := returnRead[numCallsRead-1]
		for i := 0; i < len(b); i++ {
			b[i] = rr[i]
		}
		return len(b), nil
	}
	// act
	got, err := d.ReadTime()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, 2, numCallsRead)
	assert.Equal(t, want, got)
}

func TestPCF8583ReadTimeNoTimeModeFails(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{}         // reset writes of Start() and former test
	readCtrlState := uint8(0x20) // counter mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	got, err := d.ReadTime()
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "wrong mode 0x20")
	assert.Equal(t, time.Time{}, got)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, 1, numCallsRead)
}

func TestPCF8583WriteCounter(t *testing.T) {
	// sequence to write the counter:
	// * read control register for get current state and ensure the event counter mode is set
	// * write the control register (stop counting)
	// * create the values for counter registers
	// * write the counter registers
	// * write the control register (start counting)
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{}         // reset writes of Start() and former test
	readCtrlState := uint8(0x27) // counter mode
	initCount := int32(654321)   // 6 digits used of 10 possible with int32
	wantCtrlStop := uint8(0xA7)  // stop counting bit is set
	wantReg1Val := uint8(0x21)   // BCD: 21
	wantReg2Val := uint8(0x43)   // BCD: 43
	wantReg3Val := uint8(0x65)   // BCD: 65
	wantCtrlStart := uint8(0x27) // counter mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	err := d.WriteCounter(initCount)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Len(t, a.written, 8)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[1])
	assert.Equal(t, wantCtrlStop, a.written[2])
	assert.Equal(t, wantReg1Val, a.written[3])
	assert.Equal(t, wantReg2Val, a.written[4])
	assert.Equal(t, wantReg3Val, a.written[5])
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[6])
	assert.Equal(t, wantCtrlStart, a.written[7])
}

func TestPCF8583WriteCounterNoCounterModeFails(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{}         // reset writes of Start() and former test
	readCtrlState := uint8(0x10) // 50Hz mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	err := d.WriteCounter(123)
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "wrong mode 0x10")
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, 1, numCallsRead)
}

func TestPCF8583ReadCounter(t *testing.T) {
	// sequence to read the counter:
	// * read the control register to ensure the event counter mode is set
	// * read the counter registers
	// * create the value out of registers content
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	want := int32(654321)
	reg0Val := uint8(0x20) // counter mode
	reg1Val := uint8(0x21) // BCD: 21
	reg2Val := uint8(0x43) // BCD: 43
	reg3Val := uint8(0x65) // BCD: 65
	returnRead := [2][]uint8{
		{reg0Val},
		{reg1Val, reg2Val, reg3Val},
	}
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		rr := returnRead[numCallsRead-1]
		for i := 0; i < len(b); i++ {
			b[i] = rr[i]
		}
		return len(b), nil
	}
	// act
	got, err := d.ReadCounter()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, 2, numCallsRead)
	assert.Equal(t, want, got)
}

func TestPCF8583ReadCounterNoCounterModeFails(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{}         // reset writes of Start() and former test
	readCtrlState := uint8(0x30) // test mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act
	got, err := d.ReadCounter()
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "wrong mode 0x30")
	assert.Equal(t, int32(0), got)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, 1, numCallsRead)
}

func TestPCF8583WriteRam(t *testing.T) {
	// sequence to write the RAM:
	// * calculate the RAM address and check for valid range
	// * write the given value to the given RAM address
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	wantRAMAddress := uint8(0xFF)
	wantRAMValue := uint8(0xEF)
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// act
	err := d.WriteRAM(wantRAMAddress-pcf8583RamOffset, wantRAMValue)
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 2)
	assert.Equal(t, wantRAMAddress, a.written[0])
	assert.Equal(t, wantRAMValue, a.written[1])
}

func TestPCF8583WriteRamAddressOverflowFails(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	// act
	err := d.WriteRAM(uint8(0xF0), 15)
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "overflow 256")
	assert.Empty(t, a.written)
}

func TestPCF8583ReadRam(t *testing.T) {
	// sequence to read the RAM:
	// * calculate the RAM address and check for valid range
	// * read the value from the given RAM address
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	wantRAMAddress := uint8(pcf8583RamOffset)
	want := uint8(0xAB)
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = want
		return len(b), nil
	}
	// act
	got, err := d.ReadRAM(wantRAMAddress - pcf8583RamOffset)
	// assert
	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Len(t, a.written, 1)
	assert.Equal(t, wantRAMAddress, a.written[0])
	assert.Equal(t, 1, numCallsRead)
}

func TestPCF8583ReadRamAddressOverflowFails(t *testing.T) {
	// arrange
	d, a := initTestPCF8583WithStubbedAdaptor()
	a.written = []byte{} // reset writes of Start() and former test
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		return len(b), nil
	}
	// act
	got, err := d.ReadRAM(uint8(0xF0))
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "overflow 256")
	assert.Equal(t, uint8(0), got)
	assert.Empty(t, a.written)
	assert.Equal(t, 0, numCallsRead)
}

func TestPCF8583_initializeNoModeSwitch(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	d := NewPCF8583Driver(a)
	a.written = []byte{}         // reset writes of former tests
	readCtrlState := uint8(0x01) // 32.768kHz clock mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act, assert - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
}

func TestPCF8583_initializeWithModeSwitch(t *testing.T) {
	// sequence to change mode:
	// * read control register for get current state
	// * reset old mode bits and set new mode bit
	// * write the control register
	// arrange
	a := newI2cTestAdaptor()
	d := NewPCF8583Driver(a)
	d.mode = PCF8583CtrlModeCounter
	a.written = []byte{}         // reset writes of former tests
	readCtrlState := uint8(0x02) // 32.768kHz clock mode
	wantReg0Val := uint8(0x22)   // event counter mode
	// arrange writes
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[len(b)-1] = readCtrlState
		return len(b), nil
	}
	// act, assert - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Len(t, a.written, 3)
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[0])
	assert.Equal(t, uint8(pcf8583Reg_CTRL), a.written[1])
	assert.Equal(t, wantReg0Val, a.written[2])
}
