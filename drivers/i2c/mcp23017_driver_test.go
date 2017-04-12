package i2c

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
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
	driver, _ = initTestMCP23017DriverWithStubbedAdaptor(b)
	return
}

func initTestMCP23017DriverWithStubbedAdaptor(b uint8) (*MCP23017Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewMCP23017Driver(adaptor, WithMCP23017Bank(b)), adaptor
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
	gobottest.Assert(t, b.MCPConf.Bank, uint8(1))
}

func TestNewMCP23017DriverMirror(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Mirror(1))
	gobottest.Assert(t, b.MCPConf.Mirror, uint8(1))
}

func TestNewMCP23017DriverSeqop(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Seqop(1))
	gobottest.Assert(t, b.MCPConf.Seqop, uint8(1))
}

func TestNewMCP23017DriverDisslw(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Disslw(1))
	gobottest.Assert(t, b.MCPConf.Disslw, uint8(1))
}

func TestNewMCP23017DriverHaen(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Haen(1))
	gobottest.Assert(t, b.MCPConf.Haen, uint8(1))
}

func TestNewMCP23017DriverOdr(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Odr(1))
	gobottest.Assert(t, b.MCPConf.Odr, uint8(1))
}

func TestNewMCP23017DriverIntpol(t *testing.T) {
	b := NewMCP23017Driver(newI2cTestAdaptor(), WithMCP23017Intpol(1))
	gobottest.Assert(t, b.MCPConf.Intpol, uint8(1))
}

func TestMCP23017DriverStart(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)

	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err := mcp.Start()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017StartConnectError(t *testing.T) {
	d, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestMCP23017DriverHalt(t *testing.T) {
	mcp := initTestMCP23017Driver(0)
	gobottest.Assert(t, mcp.Halt(), nil)
}

func TestMCP23017DriverCommandsWriteGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	result := mcp.Command("WriteGPIO")(pinValPort)
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestMCP23017DriverCommandsReadGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	result := mcp.Command("ReadGPIO")(pinPort)
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestMCP23017DriverWriteGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobottest.Assert(t, err, nil)
}

func TestMCP23017DriverCommandsWriteGPIOErrIODIR(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverCommandsWriteGPIOErrOLAT(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	numCalls := 1
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		if numCalls == 2 {
			return 0, errors.New("write error")
		}
		numCalls++
		return 0, nil
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverReadGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	val, _ := mcp.ReadGPIO(7, "A")
	gobottest.Assert(t, val, uint8(0))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}
	_, err := mcp.ReadGPIO(7, "A")
	gobottest.Assert(t, err, errors.New("read error"))

	// empty value from read
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("Read came back with no data")
	}
	_, err = mcp.ReadGPIO(7, "A")
	gobottest.Assert(t, err, errors.New("Read came back with no data"))
}

func TestMCP23017DriverPinMode(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err := mcp.PinMode(7, 0, "A")
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err = mcp.PinMode(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetPullUp(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err := mcp.SetPullUp(7, 0, "A")
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err = mcp.SetPullUp(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetGPIOPolarity(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err := mcp.SetGPIOPolarity(7, 0, "A")
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err = mcp.SetGPIOPolarity(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestMCP23017DriverWrite(t *testing.T) {
	// clear bit
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	port := mcp.getPort("A")
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err := mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, nil)

	// set bit
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	port = mcp.getPort("B")
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, errors.New("write error"))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, errors.New("read error"))

	//debug
	debug = true
	log.SetOutput(ioutil.Discard)
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, nil
	}
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobottest.Assert(t, err, nil)
	debug = false
	log.SetOutput(os.Stdout)
}

func TestMCP23017DriverReadPort(t *testing.T) {
	// read
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	port := mcp.getPort("A")
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{255})
		return 1, nil
	}
	val, _ := mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(255))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return len(b), errors.New("read error")
	}

	val, err := mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(0))
	gobottest.Assert(t, err, errors.New("read error"))

	// read
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	port = mcp.getPort("A")
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, nil
	}
	_, err = mcp.read(port.IODIR)
	gobottest.Assert(t, err, errors.New("Read was unable to get 1 bytes for register: 0x0\n"))

	// debug
	debug = true
	log.SetOutput(ioutil.Discard)
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	gobottest.Assert(t, mcp.Start(), nil)

	port = mcp.getPort("A")

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{255})
		return 1, nil
	}

	val, _ = mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(255))
	debug = false
	log.SetOutput(os.Stdout)
}

func TestMCP23017DriverGetPort(t *testing.T) {
	// port A
	mcp := initTestMCP23017Driver(0)
	expectedPort := getBank(0).PortA
	actualPort := mcp.getPort("A")
	gobottest.Assert(t, expectedPort, actualPort)

	// port B
	mcp = initTestMCP23017Driver(0)
	expectedPort = getBank(0).PortB
	actualPort = mcp.getPort("B")
	gobottest.Assert(t, expectedPort, actualPort)

	// default
	mcp = initTestMCP23017Driver(0)
	expectedPort = getBank(0).PortA
	actualPort = mcp.getPort("")
	gobottest.Assert(t, expectedPort, actualPort)

	// port A bank 1
	mcp = initTestMCP23017Driver(1)
	expectedPort = getBank(1).PortA
	actualPort = mcp.getPort("")
	gobottest.Assert(t, expectedPort, actualPort)
}

func TestSetBit(t *testing.T) {
	var expectedVal uint8 = 129
	actualVal := setBit(1, 7)
	gobottest.Assert(t, expectedVal, actualVal)
}

func TestClearBit(t *testing.T) {
	var expectedVal uint8
	actualVal := clearBit(128, 7)
	gobottest.Assert(t, expectedVal, actualVal)
}

func TestMCP23017DriverSetName(t *testing.T) {
	d := initTestMCP23017Driver(0)
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}
