package i2c

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestMCP23017Driver(b uint8) (driver *MCP23017Driver) {
	driver, _ = initTestMCP23017DriverWithStubbedAdaptor(b)
	return
}

func initTestMCP23017DriverWithStubbedAdaptor(b uint8) (*MCP23017Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewMCP23017Driver(adaptor, "bot", MCP23017Config{Bank: b}, 0x20), adaptor
}

func TestNewMCP23017Driver(t *testing.T) {
	var bm interface{} = NewMCP23017Driver(newI2cTestAdaptor("adaptor"), "bot", MCP23017Config{}, 0x20)
	_, ok := bm.(*MCP23017Driver)
	if !ok {
		t.Errorf("NewMCP23017Driver() should have returned a *MCP23017Driver")
	}

	b := NewMCP23017Driver(newI2cTestAdaptor("adaptor"), "bot", MCP23017Config{}, 0x20)
	gobot.Assert(t, b.Name(), "bot")
	gobot.Assert(t, b.Connection().Name(), "adaptor")
}

func TestMCP23017DriverStart(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)

	gobot.Assert(t, len(mcp.Start()), 0)

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	err := mcp.Start()
	gobot.Assert(t, err[0], errors.New("write error"))

	adaptor.i2cStartImpl = func() error {
		return errors.New("start error")
	}
	err = mcp.Start()
	gobot.Assert(t, err[0], errors.New("start error"))
}

func TestMCP23017DriverHalt(t *testing.T) {
	mcp := initTestMCP23017Driver(0)

	gobot.Assert(t, len(mcp.Halt()), 0)
}

func TestMCP23017DriverWriteGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return nil
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobot.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.WriteGPIO(7, 0, "A")
	gobot.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverReadGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	val, _ := mcp.ReadGPIO(7, "A")
	gobot.Assert(t, val, true)

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return nil, errors.New("read error")
	}
	_, err := mcp.ReadGPIO(7, "A")
	gobot.Assert(t, err, errors.New("read error"))
}

func TestMCP23017DriverSetPullUp(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return nil
	}
	err := mcp.SetPullUp(7, 0, "A")
	gobot.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.SetPullUp(7, 0, "A")
	gobot.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetGPIOPolarity(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return nil
	}
	err := mcp.SetGPIOPolarity(7, 0, "A")
	gobot.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.SetGPIOPolarity(7, 0, "A")
	gobot.Assert(t, err, errors.New("write error"))

}

func TestMCP23017DriverWrite(t *testing.T) {
	// clear bit
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	port := mcp.getPort("A")
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return nil
	}
	err := mcp.write(port.IODIR, uint8(7), 0)
	gobot.Assert(t, err, nil)

	// set bit
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("B")
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return nil
	}
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobot.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, nil
	}
	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobot.Assert(t, err, errors.New("write error"))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{}, errors.New("read error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobot.Assert(t, err, errors.New("read error"))
}

func TestMCP23017DriverReadPort(t *testing.T) {
	// read
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	port := mcp.getPort("A")

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil
	}
	val, _ := mcp.read(port.IODIR)
	gobot.Assert(t, val, uint8(255))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{}, errors.New("read error")
	}

	val, err := mcp.read(port.IODIR)
	gobot.Assert(t, val, uint8(0))
	gobot.Assert(t, err, errors.New("read error"))

	// debug
	Debug = true
	log.SetOutput(ioutil.Discard)
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("A")

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil
	}

	val, _ = mcp.read(port.IODIR)
	gobot.Assert(t, val, uint8(255))
	Debug = false
	log.SetOutput(os.Stdout)
}

func TestMCP23017DriverGetPort(t *testing.T) {
	// port a
	mcp := initTestMCP23017Driver(0)
	expectedPort := getBank(0).PortA
	actualPort := mcp.getPort("A")
	gobot.Assert(t, expectedPort, actualPort)

	// port b
	mcp = initTestMCP23017Driver(0)
	expectedPort = getBank(0).PortB
	actualPort = mcp.getPort("B")
	gobot.Assert(t, expectedPort, actualPort)

	// default
	mcp = initTestMCP23017Driver(0)
	expectedPort = getBank(0).PortA
	actualPort = mcp.getPort("")
	gobot.Assert(t, expectedPort, actualPort)

	// port a bank 1
	mcp = initTestMCP23017Driver(1)
	expectedPort = getBank(1).PortA
	actualPort = mcp.getPort("")
	gobot.Assert(t, expectedPort, actualPort)
}

func TestSetBit(t *testing.T) {
	var expectedVal uint8 = 129
	actualVal := setBit(1, 7)
	gobot.Assert(t, expectedVal, actualVal)
}

func TestClearBit(t *testing.T) {
	var expectedVal uint8 = 0
	actualVal := clearBit(128, 7)
	gobot.Assert(t, expectedVal, actualVal)
}
