package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2/drivers/i2c"
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
	assert.Equal(t, 0, len(a.buses))

	con, err := a.GetI2cConnection(0xff, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(a.buses))

	_, err = con.Write([]byte{0x00, 0x01})
	assert.NoError(t, err)

	data := []byte{42, 42}
	_, err = con.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x01}, data)

	assert.NoError(t, a.Finalize())
	assert.Equal(t, 0, len(a.buses))
}

func TestI2cGetI2cConnection(t *testing.T) {
	// arrange
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	// assert working connection
	c1, e1 := a.GetI2cConnection(0xff, 1)
	assert.NoError(t, e1)
	assert.NotNil(t, c1)
	assert.Equal(t, 1, len(a.buses))
	// assert invalid bus gets error
	c2, e2 := a.GetI2cConnection(0x01, 99)
	assert.ErrorContains(t, e2, "99 not valid")
	assert.Nil(t, c2)
	assert.Equal(t, 1, len(a.buses))
	// assert unconnected gets error
	assert.NoError(t, a.Finalize())
	c3, e3 := a.GetI2cConnection(0x01, 99)
	assert.ErrorContains(t, e3, "not connected")
	assert.Nil(t, c3)
	assert.Equal(t, 0, len(a.buses))
}

func TestI2cFinalize(t *testing.T) {
	// arrange
	a, fs := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	// assert that finalize before connect is working
	assert.NoError(t, a.Finalize())
	// arrange
	assert.NoError(t, a.Connect())
	_, _ = a.GetI2cConnection(0xaf, 1)
	assert.Equal(t, 1, len(a.buses))
	// assert that Finalize after GetI2cConnection is working and clean up
	assert.NoError(t, a.Finalize())
	assert.Equal(t, 0, len(a.buses))
	// assert that finalize after finalize is working
	assert.NoError(t, a.Finalize())
	// assert that close error is recognized
	assert.NoError(t, a.Connect())
	con, _ := a.GetI2cConnection(0xbf, 1)
	assert.Equal(t, 1, len(a.buses))
	_, _ = con.Write([]byte{0xbf})
	fs.WithCloseError = true
	err := a.Finalize()
	assert.Contains(t, err.Error(), "close error")
}

func TestI2cReConnect(t *testing.T) {
	// arrange
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	assert.NoError(t, a.Finalize())
	// act
	assert.NoError(t, a.Connect())
	// assert
	assert.NotNil(t, a.buses)
	assert.Equal(t, 0, len(a.buses))
}

func TestI2cGetDefaultBus(t *testing.T) {
	a := NewI2cBusAdaptor(nil, nil, 2)
	assert.Equal(t, 2, a.DefaultI2cBus())
}
