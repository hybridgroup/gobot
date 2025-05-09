package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this SpiBusAdaptor fulfills all the required interfaces
var _ spi.Connector = (*SpiBusAdaptor)(nil)

const spiTestAllowedBus = 15

func initTestSpiBusAdaptorWithMockedSpi() (*SpiBusAdaptor, *system.MockSpiAccess) {
	validator := func(busNr int) error {
		if busNr != spiTestAllowedBus {
			return fmt.Errorf("%d not valid", busNr)
		}
		return nil
	}
	sys := system.NewAccesser()
	// arrange for periphio needed
	sys.UseMockFilesystem([]string{"/dev/spidev"})
	dpa := sys.UseMockDigitalPinAccess()
	a := NewSpiBusAdaptor(sys, validator, 1, 2, 3, 4, 5, dpa)
	spi := a.sys.UseMockSpi()
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, spi
}

func TestNewSpiAdaptor(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	// arrange for periphio needed
	sys.UseMockFilesystem([]string{"/dev/spidev"})
	dpa := sys.UseMockDigitalPinAccess()
	a := NewSpiBusAdaptor(sys, nil, 1, 2, 3, 4, 5, dpa)
	// act & assert
	assert.Equal(t, 1, a.SpiDefaultBusNumber())
	assert.Equal(t, 2, a.SpiDefaultChipNumber())
	assert.Equal(t, 3, a.SpiDefaultMode())
	assert.Equal(t, 4, a.SpiDefaultBitCount())
	assert.Equal(t, int64(5), a.SpiDefaultMaxSpeed())
	_, err := a.GetSpiConnection(10, 0, 0, 8, 10000000)
	require.ErrorContains(t, err, "not connected")
}

func TestGetSpiConnection(t *testing.T) {
	// arrange
	const (
		busNum   = spiTestAllowedBus
		chipNum  = 14
		mode     = 13
		bits     = 12
		maxSpeed = int64(11)
	)
	a, spi := initTestSpiBusAdaptorWithMockedSpi()
	assert.Empty(t, a.connections)
	// act
	con1, err1 := a.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	// assert
	require.NoError(t, err1)
	assert.NotNil(t, con1)
	assert.Len(t, a.connections, 1)
	// assert cached connection
	con1a, err2 := a.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	require.NoError(t, err2)
	assert.Equal(t, con1, con1a)
	assert.Len(t, a.connections, 1)
	// assert second connection
	con2, err3 := a.GetSpiConnection(busNum, chipNum+1, mode, bits, maxSpeed)
	require.NoError(t, err3)
	assert.NotNil(t, con2)
	assert.NotEqual(t, con1, con2)
	assert.Len(t, a.connections, 2)
	// assert bus validation error
	con, err := a.GetSpiConnection(busNum+1, chipNum, mode, bits, maxSpeed)
	require.ErrorContains(t, err, "16 not valid")
	assert.Nil(t, con)
	// assert create error
	spi.CreateError = true
	con, err = a.GetSpiConnection(busNum, chipNum+2, mode, bits, maxSpeed)
	require.ErrorContains(t, err, "error while create SPI connection in mock")
	assert.Nil(t, con)
}

func TestSpiFinalize(t *testing.T) {
	// arrange
	a, _ := initTestSpiBusAdaptorWithMockedSpi()
	_, e := a.GetSpiConnection(spiTestAllowedBus, 2, 3, 4, 5)
	require.NoError(t, e)
	assert.Len(t, a.connections, 1)
	// act
	err := a.Finalize()
	// assert
	require.NoError(t, err)
	assert.Empty(t, a.connections)
}

func TestSpiFinalizeWithError(t *testing.T) {
	// arrange
	a, spi := initTestSpiBusAdaptorWithMockedSpi()
	_, e := a.GetSpiConnection(spiTestAllowedBus, 2, 3, 4, 5)
	require.NoError(t, e)
	spi.SetCloseError(true)
	// act
	err := a.Finalize()
	// assert
	require.ErrorContains(t, err, "error while SPI close")
}

func TestSpiReConnect(t *testing.T) {
	// arrange
	a, _ := initTestSpiBusAdaptorWithMockedSpi()
	require.NoError(t, a.Finalize())
	// act
	require.NoError(t, a.Connect())
	// assert
	assert.NotNil(t, a.connections)
	assert.Empty(t, a.connections)
}
