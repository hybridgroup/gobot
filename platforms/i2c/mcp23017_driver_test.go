package i2c

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

type i2cMcpTestAdaptor struct {
	name            string
	i2cMcpReadImpl  func(int, int) ([]byte, error)
	i2cMcpWriteImpl func() error
	i2cMcpStartImpl func() error
}

func (t *i2cMcpTestAdaptor) I2cStart(int) (err error) {
	return t.i2cMcpStartImpl()
}
func (t *i2cMcpTestAdaptor) I2cRead(address int, numBytes int) (data []byte, err error) {
	return t.i2cMcpReadImpl(address, numBytes)
}
func (t *i2cMcpTestAdaptor) I2cWrite(int, []byte) (err error) {
	return t.i2cMcpWriteImpl()
}
func (t *i2cMcpTestAdaptor) Name() string             { return t.name }
func (t *i2cMcpTestAdaptor) Connect() (errs []error)  { return }
func (t *i2cMcpTestAdaptor) Finalize() (errs []error) { return }

func newMcpI2cTestAdaptor(name string) *i2cMcpTestAdaptor {
	return &i2cMcpTestAdaptor{
		name: name,
		i2cMcpReadImpl: func(address int, numBytes int) ([]byte, error) {
			return []byte{}, nil
		},
		i2cMcpWriteImpl: func() error {
			return nil
		},
		i2cMcpStartImpl: func() error {
			return nil
		},
	}
}

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

func initTestMCP23017DriverWithStubbedAdaptor(b uint8) (*MCP23017Driver, *i2cMcpTestAdaptor) {
	adaptor := newMcpI2cTestAdaptor("adaptor")
	return NewMCP23017Driver(adaptor, "bot", MCP23017Config{Bank: b}, 0x20), adaptor
}

func TestNewMCP23017Driver(t *testing.T) {
	var bm interface{} = NewMCP23017Driver(newMcpI2cTestAdaptor("adaptor"), "bot", MCP23017Config{}, 0x20)
	_, ok := bm.(*MCP23017Driver)
	if !ok {
		t.Errorf("NewMCP23017Driver() should have returned a *MCP23017Driver")
	}

	b := NewMCP23017Driver(newMcpI2cTestAdaptor("adaptor"), "bot", MCP23017Config{}, 0x20)
	gobottest.Assert(t, b.Name(), "bot")
	gobottest.Assert(t, b.Connection().Name(), "adaptor")
}

func TestMCP23017DriverStart(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)

	gobottest.Assert(t, len(mcp.Start()), 0)

	adaptor.i2cMcpWriteImpl = func() error {
		return errors.New("write error")
	}
	err := mcp.Start()
	gobottest.Assert(t, err[0], errors.New("write error"))

	adaptor.i2cMcpStartImpl = func() error {
		return errors.New("start error")
	}
	err = mcp.Start()
	gobottest.Assert(t, err[0], errors.New("start error"))
}

func TestMCP23017DriverHalt(t *testing.T) {
	mcp := initTestMCP23017Driver(0)
	gobottest.Assert(t, len(mcp.Halt()), 0)
}

func TestMCP23017DriverCommandsWriteGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	result := mcp.Command("WriteGPIO")(pinValPort)
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestMCP23017DriverCommandsReadGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	result := mcp.Command("ReadGPIO")(pinPort)
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestMCP23017DriverWriteGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobottest.Assert(t, err, nil)
}
func TestMCP23017DriverCommandsWriteGPIOErrIODIR(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return errors.New("write error")
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverCommandsWriteGPIOErrOLAT(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	numCalls := 1
	adaptor.i2cMcpWriteImpl = func() error {
		if numCalls == 2 {
			return errors.New("write error")
		}
		numCalls++
		return nil
	}
	err := mcp.WriteGPIO(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverReadGPIO(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	val, _ := mcp.ReadGPIO(7, "A")
	gobottest.Assert(t, val, uint8(0))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), errors.New("read error")
	}
	_, err := mcp.ReadGPIO(7, "A")
	gobottest.Assert(t, err, errors.New("read error"))

	// empty value from read
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), errors.New("Read came back with no data")
	}
	_, err = mcp.ReadGPIO(7, "A")
	gobottest.Assert(t, err, errors.New("Read came back with no data"))
}

func TestMCP23017DriverPinMode(t *testing.T) {
        mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
        adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
                return make([]byte, b), nil
        }
        adaptor.i2cMcpWriteImpl = func() error {
                return nil
        }
        err := mcp.PinMode(7, 0, "A")
        gobottest.Assert(t, err, nil)

        // write error
        mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
        adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
                return make([]byte, b), nil
        }
        adaptor.i2cMcpWriteImpl = func() error {
                return errors.New("write error")
        }
        err = mcp.PinMode(7, 0, "A")
        gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetPullUp(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	err := mcp.SetPullUp(7, 0, "A")
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.SetPullUp(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMCP23017DriverSetGPIOPolarity(t *testing.T) {
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	err := mcp.SetGPIOPolarity(7, 0, "A")
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.SetGPIOPolarity(7, 0, "A")
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestMCP23017DriverWrite(t *testing.T) {
	// clear bit
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	port := mcp.getPort("A")
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	err := mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, nil)

	// set bit
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("B")
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobottest.Assert(t, err, nil)

	// write error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return errors.New("write error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, errors.New("write error"))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), errors.New("read error")
	}
	err = mcp.write(port.IODIR, uint8(7), 0)
	gobottest.Assert(t, err, errors.New("read error"))

	//debug
	debug = true
	log.SetOutput(ioutil.Discard)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), nil
	}
	adaptor.i2cMcpWriteImpl = func() error {
		return nil
	}
	err = mcp.write(port.IODIR, uint8(7), 1)
	gobottest.Assert(t, err, nil)
	debug = false
	log.SetOutput(os.Stdout)
}

func TestMCP23017DriverReadPort(t *testing.T) {
	// read
	mcp, adaptor := initTestMCP23017DriverWithStubbedAdaptor(0)
	port := mcp.getPort("A")
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return []byte{255}, nil
	}
	val, _ := mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(255))

	// read error
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return make([]byte, b), errors.New("read error")
	}

	val, err := mcp.read(port.IODIR)
	gobottest.Assert(t, val, uint8(0))
	gobottest.Assert(t, err, errors.New("read error"))

	// read
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("A")
	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return []byte{}, nil
	}
	_, err = mcp.read(port.IODIR)
	gobottest.Assert(t, err, errors.New("Read was unable to get 1 bytes for register: 0x0\n"))

	// debug
	debug = true
	log.SetOutput(ioutil.Discard)
	mcp, adaptor = initTestMCP23017DriverWithStubbedAdaptor(0)
	port = mcp.getPort("A")

	adaptor.i2cMcpReadImpl = func(a int, b int) ([]byte, error) {
		return []byte{255}, nil
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
