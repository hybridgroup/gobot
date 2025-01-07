package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ i2c.Connector = (*I2cBusAdaptor)(nil)

const i2cBus1 = "/dev/i2c-1"

func initTestI2cAdaptorWithMockedFilesystem(mockPaths []string) (*I2cBusAdaptor, *system.MockFilesystem) {
	validator := func(busNr int) error {
		if busNr > 1 {
			return fmt.Errorf("%d not valid", busNr)
		}
		return nil
	}
	a := NewI2cBusAdaptor(system.NewAccesser(), validator, 1)
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestI2cWorkflow(t *testing.T) {
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	assert.Empty(t, a.buses)

	con, err := a.GetI2cConnection(0xff, 1)
	require.NoError(t, err)
	assert.Len(t, a.buses, 1)

	_, err = con.Write([]byte{0x00, 0x01})
	require.NoError(t, err)

	data := []byte{42, 42}
	_, err = con.Read(data)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x01}, data)

	require.NoError(t, a.Finalize())
	assert.Empty(t, a.buses)
}

func TestI2cGetI2cConnection(t *testing.T) {
	// arrange
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	// assert working connection
	c1, e1 := a.GetI2cConnection(0xff, 1)
	require.NoError(t, e1)
	assert.NotNil(t, c1)
	assert.Len(t, a.buses, 1)
	// assert invalid bus gets error
	c2, e2 := a.GetI2cConnection(0x01, 99)
	require.ErrorContains(t, e2, "99 not valid")
	assert.Nil(t, c2)
	assert.Len(t, a.buses, 1)
	// assert unconnected gets error
	require.NoError(t, a.Finalize())
	c3, e3 := a.GetI2cConnection(0x01, 99)
	require.ErrorContains(t, e3, "not connected")
	assert.Nil(t, c3)
	assert.Empty(t, a.buses)
}

func TestI2cFinalize(t *testing.T) {
	// arrange
	a, fs := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	// assert that finalize before connect is working
	require.NoError(t, a.Finalize())
	// arrange
	require.NoError(t, a.Connect())
	_, _ = a.GetI2cConnection(0xaf, 1)
	assert.Len(t, a.buses, 1)
	// assert that Finalize after GetI2cConnection is working and clean up
	require.NoError(t, a.Finalize())
	assert.Empty(t, a.buses)
	// assert that finalize after finalize is working
	require.NoError(t, a.Finalize())
	// assert that close error is recognized
	require.NoError(t, a.Connect())
	con, _ := a.GetI2cConnection(0xbf, 1)
	assert.Len(t, a.buses, 1)
	_, _ = con.Write([]byte{0xbf})
	fs.WithCloseError = true
	err := a.Finalize()
	require.ErrorContains(t, err, "close error")
}

func TestI2cReConnect(t *testing.T) {
	// arrange
	a, _ := initTestI2cAdaptorWithMockedFilesystem([]string{i2cBus1})
	require.NoError(t, a.Finalize())
	// act
	require.NoError(t, a.Connect())
	// assert
	assert.NotNil(t, a.buses)
	assert.Empty(t, a.buses)
}

func TestI2cGetDefaultBus(t *testing.T) {
	a := NewI2cBusAdaptor(system.NewAccesser(), nil, 2)
	assert.Equal(t, 2, a.DefaultI2cBus())
}
