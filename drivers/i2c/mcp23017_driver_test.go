//nolint:forcetypeassert,gosec // ok here
package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MCP23017Driver)(nil)

var pinValPort = map[string]interface{}{
	"pin":  uint8(7),
	"val":  uint8(0),
	"port": "A",
}

var pinPort = map[string]interface{}{
	"pin":  uint8(7),
	"port": "A",
}

func initTestMCP23017(b uint8) *MCP23017Driver {
	// create the driver without starting it
	a := newI2cTestAdaptor()
	d := NewMCP23017Driver(a, WithMCP23017Bank(b))
	return d
}

func initTestMCP23017WithStubbedAdaptor(b uint8) (*MCP23017Driver, *i2cTestAdaptor) { //nolint:unparam // keep for tests
	// create the driver, ready to use for tests
	a := newI2cTestAdaptor()
	d := NewMCP23017Driver(a, WithMCP23017Bank(b))
	_ = d.Start()
	return d, a
}

func TestNewMCP23017Driver(t *testing.T) {
	var di interface{} = NewMCP23017Driver(newI2cTestAdaptor())
	d, ok := di.(*MCP23017Driver)
	if !ok {
		require.Fail(t, "NewMCP23017Driver() should have returned a *MCP23017Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "MCP23017"))
	assert.Equal(t, 0x20, d.defaultAddress)
	assert.NotNil(t, d.mcpConf)
	assert.NotNil(t, d.mcpBehav)
}

func TestWithMCP23017Bank(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Bank(1))
	assert.Equal(t, uint8(1), b.mcpConf.bank)
}

func TestWithMCP23017Mirror(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Mirror(1))
	assert.Equal(t, uint8(1), b.mcpConf.mirror)
}

func TestWithMCP23017Seqop(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Seqop(1))
	assert.Equal(t, uint8(1), b.mcpConf.seqop)
}

func TestWithMCP23017Disslw(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Disslw(1))
	assert.Equal(t, uint8(1), b.mcpConf.disslw)
}

func TestWithMCP23017Haen(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Haen(1))
	assert.Equal(t, uint8(1), b.mcpConf.haen)
}

func TestWithMCP23017Odr(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Odr(1))
	assert.Equal(t, uint8(1), b.mcpConf.odr)
}

func TestWithMCP23017Intpol(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Intpol(1))
	assert.Equal(t, uint8(1), b.mcpConf.intpol)
}

func TestWithMCP23017ForceRefresh(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017ForceRefresh(1))
	assert.True(t, b.mcpBehav.forceRefresh)
}

func TestWithMCP23017AutoIODirOff(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017AutoIODirOff(1))
	assert.True(t, b.mcpBehav.autoIODirOff)
}

func TestMCP23017CommandsWriteGPIO(t *testing.T) {
	// arrange
	d, _ := initTestMCP23017WithStubbedAdaptor(0)
	// act
	result := d.Command("WriteGPIO")(pinValPort)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestMCP23017CommandsReadGPIO(t *testing.T) {
	// arrange
	d, _ := initTestMCP23017WithStubbedAdaptor(0)
	// act
	result := d.Command("ReadGPIO")(pinPort)
	// assert
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestMCP23017WriteGPIO(t *testing.T) {
	// sequence to write (we force the refresh by preset with inverse bit state):
	// * read current state of IODIR (write reg, read val) => see also SetPinMode()
	// * set IODIR of pin to input (manipulate val, write reg, write val) => see also SetPinMode()
	// * read current state of OLAT (write reg, read val)
	// * write OLAT (manipulate val, write reg, write val)
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "A"
		testPin := uint8(7)
		wantReg1 := uint8(0x00)             // IODIRA
		wantReg2 := uint8(0x14)             // OLATA
		returnRead := []uint8{0xFF, 0xFF}   // emulate all IO's are inputs, emulate bit is on
		wantReg1Val := returnRead[0] & 0x7F // IODIRA: bit 7 reset, all other untouched
		wantReg2Val := returnRead[1] & 0x7F // OLATA: bit 7 reset, all other untouched
		if bitState == 1 {
			returnRead[1] = 0x7F               // emulate bit is off
			wantReg2Val = returnRead[1] | 0x80 // OLATA: bit 7 set, all other untouched
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		err := d.WriteGPIO(testPin, testPort, uint8(bitState))
		// assert
		require.NoError(t, err)
		assert.Len(t, a.written, 6)
		assert.Equal(t, wantReg1, a.written[0])
		assert.Equal(t, wantReg1, a.written[1])
		assert.Equal(t, wantReg1Val, a.written[2])
		assert.Equal(t, wantReg2, a.written[3])
		assert.Equal(t, wantReg2, a.written[4])
		assert.Equal(t, wantReg2Val, a.written[5])
		assert.Equal(t, 2, numCallsRead)
	}
}

func TestMCP23017WriteGPIONoRefresh(t *testing.T) {
	// sequence to write with take advantage of refresh optimization (see forceRefresh):
	// * read current state of IODIR (write reg, read val) => by SetPinMode()
	// * read current state of OLAT (write reg, read val)
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "B"
		testPin := uint8(3)
		wantReg1 := uint8(0x01)           // IODIRB
		wantReg2 := uint8(0x15)           // OLATB
		returnRead := []uint8{0xF7, 0xF7} // emulate all IO's are inputs except pin 3, emulate bit is already off
		if bitState == 1 {
			returnRead[1] = 0x08 // emulate bit is already on
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		err := d.WriteGPIO(testPin, testPort, uint8(bitState))
		// assert
		require.NoError(t, err)
		assert.Len(t, a.written, 2)
		assert.Equal(t, wantReg1, a.written[0])
		assert.Equal(t, wantReg2, a.written[1])
		assert.Equal(t, 2, numCallsRead)
	}
}

func TestMCP23017WriteGPIONoAutoDir(t *testing.T) {
	// sequence to write with suppressed automatic setting of IODIR:
	// * read current state of OLAT (write reg, read val)
	// * write OLAT (manipulate val, write reg, write val)
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	d.mcpBehav.autoIODirOff = true
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "A"
		testPin := uint8(7)
		wantReg := uint8(0x14)          // OLATA
		returnRead := uint8(0xFF)       // emulate bit is on
		wantRegVal := returnRead & 0x7F // OLATA: bit 7 reset, all other untouched
		if bitState == 1 {
			returnRead = 0x7F              // emulate bit is off
			wantRegVal = returnRead | 0x80 // OLATA: bit 7 set, all other untouched
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := d.WriteGPIO(testPin, testPort, uint8(bitState))
		// assert
		require.NoError(t, err)
		assert.Len(t, a.written, 3)
		assert.Equal(t, wantReg, a.written[0])
		assert.Equal(t, wantReg, a.written[1])
		assert.Equal(t, wantRegVal, a.written[2])
		assert.Equal(t, 1, numCallsRead)
	}
}

func TestMCP23017CommandsWriteGPIOErrIODIR(t *testing.T) {
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := d.WriteGPIO(7, "A", 0)
	// assert
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=0): write error")
}

func TestMCP23017CommandsWriteGPIOErrOLAT(t *testing.T) {
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	numCalls := 1
	a.i2cWriteImpl = func([]byte) (int, error) {
		if numCalls == 2 {
			return 0, errors.New("write error")
		}
		numCalls++
		return 0, nil
	}
	// act
	err := d.WriteGPIO(7, "A", 0)
	// assert
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=20): write error")
}

func TestMCP23017ReadGPIO(t *testing.T) {
	// sequence to read:
	// * read current state of IODIR (write reg, read val) => see also SetPinMode()
	// * set IODIR of pin to input (manipulate val, write reg, write val) => see also SetPinMode()
	// * read GPIO (write reg, read val)
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "A"
		testPin := uint8(7)
		wantReg1 := uint8(0x00)             // IODIRA
		wantReg2 := uint8(0x12)             // GPIOA
		returnRead := []uint8{0x00, 0x7F}   // emulate all IO's are outputs, emulate bit is off
		wantReg1Val := returnRead[0] | 0x80 // IODIRA: bit 7 set, all other untouched
		if bitState == 1 {
			returnRead[1] = 0xFF // emulate bit is set
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		val, err := d.ReadGPIO(testPin, testPort)
		// assert
		require.NoError(t, err)
		assert.Equal(t, 2, numCallsRead)
		assert.Len(t, a.written, 4)
		assert.Equal(t, wantReg1, a.written[0])
		assert.Equal(t, wantReg1, a.written[1])
		assert.Equal(t, wantReg1Val, a.written[2])
		assert.Equal(t, wantReg2, a.written[3])
		assert.Equal(t, uint8(bitState), val)
	}
}

func TestMCP23017ReadGPIONoRefresh(t *testing.T) {
	// sequence to read with take advantage of refresh optimization (see forceRefresh):
	// * read current state of IODIR (write reg, read val) => by SetPinMode()
	// * read GPIO (write reg, read val)
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "A"
		testPin := uint8(7)
		wantReg1 := uint8(0x00)           // IODIRA
		wantReg2 := uint8(0x12)           // GPIOA
		returnRead := []uint8{0x80, 0x7F} // emulate all IO's are outputs except pin 7, emulate bit is off
		if bitState == 1 {
			returnRead[1] = 0xFF // emulate bit is set
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		val, err := d.ReadGPIO(testPin, testPort)
		// assert
		require.NoError(t, err)
		assert.Equal(t, 2, numCallsRead)
		assert.Len(t, a.written, 2)
		assert.Equal(t, wantReg1, a.written[0])
		assert.Equal(t, wantReg2, a.written[1])
		assert.Equal(t, uint8(bitState), val)
	}
}

func TestMCP23017ReadGPIONoAutoDir(t *testing.T) {
	// sequence to read with suppressed automatic setting of IODIR:
	// * read GPIO (write reg, read val)
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	d.mcpBehav.autoIODirOff = true
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "A"
		testPin := uint8(7)
		wantReg2 := uint8(0x12)   // GPIOA
		returnRead := uint8(0x7F) // emulate bit is off
		if bitState == 1 {
			returnRead = 0xFF // emulate bit is set
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		val, err := d.ReadGPIO(testPin, testPort)
		// assert
		require.NoError(t, err)
		assert.Equal(t, 1, numCallsRead)
		assert.Len(t, a.written, 1)
		assert.Equal(t, wantReg2, a.written[0])
		assert.Equal(t, uint8(bitState), val)
	}
}

func TestMCP23017ReadGPIOErr(t *testing.T) {
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	// arrange reads
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}
	// act
	_, err := d.ReadGPIO(7, "A")
	// assert
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=0): read error")
}

func TestMCP23017SetPinMode(t *testing.T) {
	// sequence for setting pin direction:
	// * read current state of IODIR (write reg, read val)
	// * set IODIR of pin to input or output (manipulate val, write reg, write val)
	// TODO: can be optimized by not writing, when value is already fine
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start() and former test
		// arrange some values
		testPort := "A"
		testPin := uint8(7)
		wantReg := uint8(0x00)          // IODIRA
		returnRead := uint8(0xFF)       // emulate all ports are inputs
		wantRegVal := returnRead & 0x7F // bit 7 reset, all other untouched
		if bitState == 1 {
			returnRead = 0x00              // emulate all ports are outputs
			wantRegVal = returnRead | 0x80 // bit 7 set, all other untouched
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := d.SetPinMode(testPin, testPort, uint8(bitState))
		// assert
		require.NoError(t, err)
		assert.Len(t, a.written, 3)
		assert.Equal(t, wantReg, a.written[0])
		assert.Equal(t, wantReg, a.written[1])
		assert.Equal(t, wantRegVal, a.written[2])
		assert.Equal(t, 1, numCallsRead)
	}
}

func TestMCP23017SetPinModeErr(t *testing.T) {
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := d.SetPinMode(7, "A", 0)
	// assert
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=0): write error")
}

func TestMCP23017SetPullUp(t *testing.T) {
	// sequence for setting input pin pull up:
	// * read current state of GPPU (write reg, read val)
	// * set GPPU of pin to target state (manipulate val, write reg, write val)
	// TODO: can be optimized by not writing, when value is already fine
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start()
		// arrange some values
		testPort := "A"
		wantReg := uint8(0x0C) // GPPUA
		testPin := uint8(5)
		returnRead := uint8(0xFF)       // emulate all I's with pull up
		wantSetVal := returnRead & 0xDF // bit 5 cleared, all other unchanged
		if bitState == 1 {
			returnRead = uint8(0x00)       // emulate all I's without pull up
			wantSetVal = returnRead | 0x20 // bit 5 set, all other unchanged
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := d.SetPullUp(testPin, testPort, uint8(bitState))
		// assert
		require.NoError(t, err)
		assert.Len(t, a.written, 3)
		assert.Equal(t, wantReg, a.written[0])
		assert.Equal(t, wantReg, a.written[1])
		assert.Equal(t, wantSetVal, a.written[2])
		assert.Equal(t, 1, numCallsRead)
	}
}

func TestMCP23017SetPullUpErr(t *testing.T) {
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := d.SetPullUp(7, "A", 0)
	// assert
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=12): write error")
}

func TestMCP23017SetGPIOPolarity(t *testing.T) {
	// sequence for setting input pin polarity:
	// * read current state of IPOL (write reg, read val)
	// * set IPOL of pin to target state (manipulate val, write reg, write val)
	// TODO: can be optimized by not writing, when value is already fine
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		a.written = []byte{} // reset writes of Start()
		// arrange some values
		testPort := "B"
		wantReg := uint8(0x03) // IPOLB
		testPin := uint8(6)
		returnRead := uint8(0xFF)       // emulate all I's negotiated
		wantSetVal := returnRead & 0xBF // bit 6 cleared, all other unchanged
		if bitState == 1 {
			returnRead = uint8(0x00)       // emulate all I's not negotiated
			wantSetVal = returnRead | 0x40 // bit 6 set, all other unchanged
		}
		// arrange reads
		numCallsRead := 0
		a.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := d.SetGPIOPolarity(testPin, testPort, uint8(bitState))
		// assert
		require.NoError(t, err)
		assert.Len(t, a.written, 3)
		assert.Equal(t, wantReg, a.written[0])
		assert.Equal(t, wantReg, a.written[1])
		assert.Equal(t, wantSetVal, a.written[2])
		assert.Equal(t, 1, numCallsRead)
	}
}

func TestMCP23017SetGPIOPolarityErr(t *testing.T) {
	// arrange
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := d.SetGPIOPolarity(7, "A", 0)
	// assert
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=2): write error")
}

func TestMCP23017_write(t *testing.T) {
	// clear bit
	d, _ := initTestMCP23017WithStubbedAdaptor(0)
	port := d.getPort("A")
	err := d.write(port.IODIR, uint8(7), 0)
	require.NoError(t, err)

	// set bit
	d, _ = initTestMCP23017WithStubbedAdaptor(0)
	port = d.getPort("B")
	err = d.write(port.IODIR, uint8(7), 1)
	require.NoError(t, err)

	// write error
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err = d.write(port.IODIR, uint8(7), 0)
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=1): write error")

	// read error
	d, a = initTestMCP23017WithStubbedAdaptor(0)
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}
	err = d.write(port.IODIR, uint8(7), 0)
	require.ErrorContains(t, err, "MCP write-read: MCP write-ReadByteData(reg=1): read error")
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	err = d.write(port.IODIR, uint8(7), 1)
	require.NoError(t, err)
}

func TestMCP23017_read(t *testing.T) {
	// read
	d, a := initTestMCP23017WithStubbedAdaptor(0)
	port := d.getPort("A")
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{255})
		return 1, nil
	}
	val, _ := d.read(port.IODIR)
	assert.Equal(t, uint8(255), val)

	// read error
	d, a = initTestMCP23017WithStubbedAdaptor(0)
	a.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}

	val, err := d.read(port.IODIR)
	assert.Equal(t, uint8(0), val)
	require.ErrorContains(t, err, "MCP write-ReadByteData(reg=0): read error")

	// read
	d, a = initTestMCP23017WithStubbedAdaptor(0)
	port = d.getPort("A")
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{255})
		return 1, nil
	}
	val, _ = d.read(port.IODIR)
	assert.Equal(t, uint8(255), val)
}

func TestMCP23017_getPort(t *testing.T) {
	// port A
	d := initTestMCP23017(0)
	wantPort := mcp23017GetBank(0).portA
	gotPort := d.getPort("A")
	assert.Equal(t, wantPort, gotPort)

	// port B
	d = initTestMCP23017(0)
	wantPort = mcp23017GetBank(0).portB
	gotPort = d.getPort("B")
	assert.Equal(t, wantPort, gotPort)

	// default
	d = initTestMCP23017(0)
	wantPort = mcp23017GetBank(0).portA
	gotPort = d.getPort("")
	assert.Equal(t, wantPort, gotPort)

	// port A bank 1
	d = initTestMCP23017(1)
	wantPort = mcp23017GetBank(1).portA
	gotPort = d.getPort("")
	assert.Equal(t, wantPort, gotPort)
}
