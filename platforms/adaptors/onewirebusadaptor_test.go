package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/system"
)

func initTestOneWireAdaptor() *OneWireBusAdaptor {
	a := NewOneWireBusAdaptor(system.NewAccesser())
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestNewOneWireBusAdaptor(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	// act
	a := NewOneWireBusAdaptor(sys)
	// assert
	assert.IsType(t, &OneWireBusAdaptor{}, a)
	assert.NotNil(t, a.mutex)
	assert.Nil(t, a.connections)
}

func TestOneWireGetOneWireConnection(t *testing.T) {
	// arrange
	const (
		familyCode   = 28
		serialNumber = 123456789
	)
	a := initTestOneWireAdaptor()
	// assert working connection
	c1, e1 := a.GetOneWireConnection(familyCode, serialNumber)
	require.NoError(t, e1)
	assert.NotNil(t, c1)
	assert.Len(t, a.connections, 1)
	// assert unconnected gets error
	require.NoError(t, a.Finalize())
	c2, e2 := a.GetOneWireConnection(familyCode, serialNumber+1)
	require.ErrorContains(t, e2, "not connected")
	assert.Nil(t, c2)
	assert.Empty(t, a.connections)
}

func TestOneWireFinalize(t *testing.T) {
	// arrange
	a := initTestOneWireAdaptor()
	// assert that finalize before connect is working
	require.NoError(t, a.Finalize())
	// arrange
	require.NoError(t, a.Connect())
	_, _ = a.GetOneWireConnection(28, 54321)
	assert.Len(t, a.connections, 1)
	// assert that Finalize after GetOneWireConnection is working and clean up
	require.NoError(t, a.Finalize())
	assert.Empty(t, a.connections)
	// assert that finalize after finalize is working
	require.NoError(t, a.Finalize())
}

func TestOneWireReConnect(t *testing.T) {
	// arrange
	a := initTestOneWireAdaptor()
	require.NoError(t, a.Finalize())
	// act
	require.NoError(t, a.Connect())
	// assert
	assert.NotNil(t, a.connections)
	assert.Empty(t, a.connections)
}
