package adaptors

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/gobottest"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ i2c.Connector = (*I2cBusAdaptor)(nil)

const i2cBus1 = "/dev/i2c-1"

func initTestI2cAdaptorWithMockedFilesystem(mockPaths []string) (*I2cBusAdaptor, *system.MockFilesystem) {
	sys := system.NewAccesser()
	sys.UseMockSyscall()
	fs := sys.UseMockFilesystem(mockPaths)
	validator := func(busNr int) error {
		if busNr > 1 {
			return fmt.Errorf("%d not valid", busNr)
		}
		return nil
	}
	a := NewI2cBusAdaptor(sys, validator, 1)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestI2cWorkflow(t *testing.T) {
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	gobottest.Assert(t, len(a.buses), 0)

	con, err := a.GetI2cConnection(0xff, 1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.buses), 1)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
	gobottest.Assert(t, len(a.buses), 0)
}

func TestI2cGetI2cConnection(t *testing.T) {
	// arrange
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	// assert working connection
	c1, e1 := a.GetI2cConnection(0xff, 1)
	gobottest.Assert(t, e1, nil)
	gobottest.Refute(t, c1, nil)
	gobottest.Assert(t, len(a.buses), 1)
	// assert invalid bus gets error
	c2, e2 := a.GetI2cConnection(0x01, 99)
	gobottest.Assert(t, e2, fmt.Errorf("99 not valid"))
	gobottest.Assert(t, c2, nil)
	gobottest.Assert(t, len(a.buses), 1)
	// assert unconnected gets error
	gobottest.Assert(t, a.Finalize(), nil)
	c3, e3 := a.GetI2cConnection(0x01, 99)
	gobottest.Assert(t, e3, fmt.Errorf("not connected"))
	gobottest.Assert(t, c3, nil)
	gobottest.Assert(t, len(a.buses), 0)
}

func TestI2cFinalize(t *testing.T) {
	// arrange
	a, fs := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	// assert that finalize before connect is working
	gobottest.Assert(t, a.Finalize(), nil)
	// arrange
	gobottest.Assert(t, a.Connect(), nil)
	_, _ = a.GetI2cConnection(0xaf, 1)
	gobottest.Assert(t, len(a.buses), 1)
	// assert that Finalize after GetI2cConnection is working and clean up
	gobottest.Assert(t, a.Finalize(), nil)
	gobottest.Assert(t, len(a.buses), 0)
	// assert that finalize after finalize is working
	gobottest.Assert(t, a.Finalize(), nil)
	// assert that close error is recognized
	gobottest.Assert(t, a.Connect(), nil)
	con, _ := a.GetI2cConnection(0xbf, 1)
	gobottest.Assert(t, len(a.buses), 1)
	_, _ = con.Write([]byte{0xbf})
	fs.WithCloseError = true
	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "close error"), true)
}

func TestI2cReConnect(t *testing.T) {
	// arrange
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	gobottest.Assert(t, a.Finalize(), nil)
	// act
	gobottest.Assert(t, a.Connect(), nil)
	// assert
	gobottest.Refute(t, a.buses, nil)
	gobottest.Assert(t, len(a.buses), 0)
}

func TestI2cGetDefaultBus(t *testing.T) {
	a := NewI2cBusAdaptor(nil, nil, 2)
	gobottest.Assert(t, a.DefaultI2cBus(), 2)
}
