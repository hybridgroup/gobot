package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

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

func initTestMCP23017Driver(b uint8) (driver *MCP23017Driver) {
	// create the driver without starting it
	adaptor := newI2cTestAdaptor()
	mcp := NewMCP23017Driver(adaptor, WithMCP23017Bank(b))
	return mcp
}

func initTestMCP23017DriverWithStubbedAdaptor(b uint8) (*MCP23017Driver, *i2cTestAdaptor) {
	// create the driver, ready to use for tests
	adaptor := newI2cTestAdaptor()
	mcp := NewMCP23017Driver(adaptor, WithMCP23017Bank(b))
	mcp.Start()
	return mcp, adaptor
}

func TestNewMCP23017Driver(t *testing.T) {
	var bm interface{} = NewMCP23017Driver(newI2cTestAdaptor())
	_, ok := bm.(*MCP23017Driver)
	if !ok {
		t.Errorf("NewMCP23017Driver() should have returned a *MCP23017Driver")
	}

	b := NewMCP23017Driver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

func TestNewMCP23017DriverBank(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Bank(1))
	gobottest.Assert(t, b.mcpConf.bank, uint8(1))
}

func TestNewMCP23017DriverMirror(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Mirror(1))
	gobottest.Assert(t, b.mcpConf.mirror, uint8(1))
}

func TestNewMCP23017DriverSeqop(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Seqop(1))
	gobottest.Assert(t, b.mcpConf.seqop, uint8(1))
}

func TestNewMCP23017DriverDisslw(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Disslw(1))
	gobottest.Assert(t, b.mcpConf.disslw, uint8(1))
}

func TestNewMCP23017DriverHaen(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Haen(1))
	gobottest.Assert(t, b.mcpConf.haen, uint8(1))
}

func TestNewMCP23017DriverOdr(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Odr(1))
	gobottest.Assert(t, b.mcpConf.odr, uint8(1))
}

func TestNewMCP23017DriverIntpol(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Intpol(1))
	gobottest.Assert(t, b.mcpConf.intpol, uint8(1))
}

func TestNewMCP23017DriverForceRefresh(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017ForceRefresh(1))
	gobottest.Assert(t, b.mvpBehav.forceRefresh, true)
}

func TestNewMCP23017DriverAutoIODirOff(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017AutoIODirOff(1))
	gobottest.Assert(t, b.mvpBehav.autoIODirOff, true)
}

func TestMCP23017DriverStart(t *testing.T) {
	// arrange
	mcp := initTestMCP23017Driver(0)
	// act & assert
	gobottest.Assert(t, mcp.Start(), nil)
}

func TestMCP23017DriverStartErr(t *testing.T) {
	// arrange
	adaptor := newI2cTestAdaptor()
	mcp := NewMCP23017Driver(adaptor, WithMCP23017Bank(0))
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := mcp.Start()
	// assert
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverHalt(t *testing.T) {
	mcp := initTestMCP23017Driver(0)
	gobottest.Assert(t, mcp.Halt(), nil)
}

func TestMCP23017DriverCommandsWriteGPIO(t *testing.T) {
	// arrange
	mcp, _ := initTestMCP23017DriverWithStubbedAdaptor(0)
	// act
	result := mcp.Command("WriteGPIO")(pinValPort)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestMCP23017DriverCommandsReadGPIO(t *testing.T) {
	// arrange
	mcp, _ := initTestMCP23017DriverWithStubbedAdaptor(0)
	// act
	result := mcp.Command("ReadGPIO")(pinPort)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestMCP23017DriverWriteGPIO(t *testing.T) {
	// sequence to write (we force the refresh by preset with inverse bit state):
	// * read current state of IODIR (write reg, read val) => see also PinMode()
	// * set IODIR of pin to input (manipulate val, write reg, write val) => see also PinMode()
	// * read current state of OLAT (write reg, read val)
	// * write OLAT (manipulate val, write reg, write val)
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		err := mcp.WriteGPIO(testPin, uint8(bitState), testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(adaptor.written), 6)
		gobottest.Assert(t, adaptor.written[0], wantReg1)
		gobottest.Assert(t, adaptor.written[1], wantReg1)
		gobottest.Assert(t, adaptor.written[2], wantReg1Val)
		gobottest.Assert(t, adaptor.written[3], wantReg2)
		gobottest.Assert(t, adaptor.written[4], wantReg2)
		gobottest.Assert(t, adaptor.written[5], wantReg2Val)
		gobottest.Assert(t, numCallsRead, 2)
	}
}

func TestMCP23017DriverWriteGPIONoRefresh(t *testing.T) {
	// sequence to write with take advantage of refresh optimization (see forceRefresh):
	// * read current state of IODIR (write reg, read val) => by PinMode()
	// * read current state of OLAT (write reg, read val)
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		err := mcp.WriteGPIO(testPin, uint8(bitState), testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(adaptor.written), 2)
		gobottest.Assert(t, adaptor.written[0], wantReg1)
		gobottest.Assert(t, adaptor.written[1], wantReg2)
		gobottest.Assert(t, numCallsRead, 2)
	}
}

func TestMCP23017DriverWriteGPIONoAutoDir(t *testing.T) {
	// sequence to write with suppressed automatic setting of IODIR:
	// * read current state of OLAT (write reg, read val)
	// * write OLAT (manipulate val, write reg, write val)
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	mcp.mvpBehav.autoIODirOff = true
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := mcp.WriteGPIO(testPin, uint8(bitState), testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(adaptor.written), 3)
		gobottest.Assert(t, adaptor.written[0], wantReg)
		gobottest.Assert(t, adaptor.written[1], wantReg)
		gobottest.Assert(t, adaptor.written[2], wantRegVal)
		gobottest.Assert(t, numCallsRead, 1)
	}
}

func TestMCP23017DriverCommandsWriteGPIOErrIODIR(t *testing.T) {
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := mcp.WriteGPIO(7, 0, "A")
	// assert
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverCommandsWriteGPIOErrOLAT(t *testing.T) {
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	numCalls := 1
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		if numCalls == 2 {
			return 0, errors.New("write error")
		}
		numCalls++
		return 0, nil
	}
	// act
	err := mcp.WriteGPIO(7, 0, "A")
	// assert
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverReadGPIO(t *testing.T) {
	// sequence to read:
	// * read current state of IODIR (write reg, read val) => see also PinMode()
	// * set IODIR of pin to input (manipulate val, write reg, write val) => see also PinMode()
	// * read GPIO (write reg, read val)
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		val, err := mcp.ReadGPIO(testPin, testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, numCallsRead, 2)
		gobottest.Assert(t, len(adaptor.written), 4)
		gobottest.Assert(t, adaptor.written[0], wantReg1)
		gobottest.Assert(t, adaptor.written[1], wantReg1)
		gobottest.Assert(t, adaptor.written[2], wantReg1Val)
		gobottest.Assert(t, adaptor.written[3], wantReg2)
		gobottest.Assert(t, val, uint8(bitState))
	}
}

func TestMCP23017DriverReadGPIONoRefresh(t *testing.T) {
	// sequence to read with take advantage of refresh optimization (see forceRefresh):
	// * read current state of IODIR (write reg, read val) => by PinMode()
	// * read GPIO (write reg, read val)
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead[numCallsRead-1]
			return len(b), nil
		}
		// act
		val, err := mcp.ReadGPIO(testPin, testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, numCallsRead, 2)
		gobottest.Assert(t, len(adaptor.written), 2)
		gobottest.Assert(t, adaptor.written[0], wantReg1)
		gobottest.Assert(t, adaptor.written[1], wantReg2)
		gobottest.Assert(t, val, uint8(bitState))
	}
}

func TestMCP23017DriverReadGPIONoAutoDir(t *testing.T) {
	// sequence to read with suppressed automatic setting of IODIR:
	// * read GPIO (write reg, read val)
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	mcp.mvpBehav.autoIODirOff = true
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		val, err := mcp.ReadGPIO(testPin, testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, numCallsRead, 1)
		gobottest.Assert(t, len(adaptor.written), 1)
		gobottest.Assert(t, adaptor.written[0], wantReg2)
		gobottest.Assert(t, val, uint8(bitState))
	}
}

func TestMCP23017DriverReadGPIOErr(t *testing.T) {
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	// arrange reads
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}
	// act
	_, err := mcp.ReadGPIO(7, "A")
	// assert
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestMCP23017DriverPinMode(t *testing.T) {
	// sequence for setting pin direction:
	// * read current state of IODIR (write reg, read val)
	// * set IODIR of pin to input or output (manipulate val, write reg, write val)
	// TODO: can be optimized by not writing, when value is already fine
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start() and former test
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := mcp.PinMode(testPin, uint8(bitState), testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(adaptor.written), 3)
		gobottest.Assert(t, adaptor.written[0], wantReg)
		gobottest.Assert(t, adaptor.written[1], wantReg)
		gobottest.Assert(t, adaptor.written[2], wantRegVal)
		gobottest.Assert(t, numCallsRead, 1)
	}
}

func TestMCP23017DriverPinModeErr(t *testing.T) {
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := mcp.PinMode(7, 0, "A")
	// assert
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetPullUp(t *testing.T) {
	// sequence for setting input pin pull up:
	// * read current state of GPPU (write reg, read val)
	// * set GPPU of pin to target state (manipulate val, write reg, write val)
	// TODO: can be optimized by not writing, when value is already fine
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start()
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := mcp.SetPullUp(testPin, uint8(bitState), testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(adaptor.written), 3)
		gobottest.Assert(t, adaptor.written[0], wantReg)
		gobottest.Assert(t, adaptor.written[1], wantReg)
		gobottest.Assert(t, adaptor.written[2], wantSetVal)
		gobottest.Assert(t, numCallsRead, 1)
	}
}

func TestMCP23017DriverSetPullUpErr(t *testing.T) {
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := mcp.SetPullUp(7, 0, "A")
	// assert
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetGPIOPolarity(t *testing.T) {
	// sequence for setting input pin polarity:
	// * read current state of IPOL (write reg, read val)
	// * set IPOL of pin to target state (manipulate val, write reg, write val)
	// TODO: can be optimized by not writing, when value is already fine
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	for bitState := 0; bitState <= 1; bitState++ {
		adaptor.written = []byte{} // reset writes of Start()
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
		adaptor.i2cReadImpl = func(b []byte) (int, error) {
			numCallsRead++
			b[len(b)-1] = returnRead
			return len(b), nil
		}
		// act
		err := mcp.SetGPIOPolarity(testPin, uint8(bitState), testPort)
		// assert
		gobottest.Assert(t, err, nil)
		gobottest.Assert(t, len(adaptor.written), 3)
		gobottest.Assert(t, adaptor.written[0], wantReg)
		gobottest.Assert(t, adaptor.written[1], wantReg)
		gobottest.Assert(t, adaptor.written[2], wantSetVal)
		gobottest.Assert(t, numCallsRead, 1)
	}
}

func TestMCP23017DriverSetGPIOPolarityErr(t *testing.T) {
	// arrange
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act
	err := mcp.SetGPIOPolarity(7, 0, "A")
	// assert
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetName(t *testing.T) {
	d := initTestMCP23017Driver(0)
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestMCP23017Driver_write(t *testing.T) {
	// clear bit
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	port := mcp.getPort("A")
	err := mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, nil)

	// set bit
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("B")
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, errors.New("write error"))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, errors.New("read error"))
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobottest.Assert(t, err, nil)
}

func TestMCP23017Driver_read(t *testing.T) {
	// read
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	port := mcp.getPort("A")
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{255})
		return 1, nil
	}
	val, _ := mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(255))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}

	val, err := mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(0))
	gobottest.Assert(t, err, errors.New("read error"))

	// read
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("A")
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{255})
		return 1, nil
	}
	val, _ = mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(255))
}

func TestMCP23017Driver_getPort(t *testing.T) {
	// port A
	mcp := initTestMCP23017Driver(0)
	expectedPort := getBank(0).portA
	actualPort := mcp.getPort("A")
	gobottest.Assert(t, expectedPort, actualPort)

	// port B
	mcp = initTestMCP23017Driver(0)
	expectedPort = getBank(0).portB
	actualPort = mcp.getPort("B")
	gobottest.Assert(t, expectedPort, actualPort)

	// default
	mcp = initTestMCP23017Driver(0)
	expectedPort = getBank(0).portA
	actualPort = mcp.getPort("")
	gobottest.Assert(t, expectedPort, actualPort)

	// port A bank 1
	mcp = initTestMCP23017Driver(1)
	expectedPort = getBank(1).portA
	actualPort = mcp.getPort("")
	gobottest.Assert(t, expectedPort, actualPort)
}
