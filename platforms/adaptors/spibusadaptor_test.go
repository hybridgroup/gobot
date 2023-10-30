package adaptors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	spi := sys.UseMockSpi()
	a := NewSpiBusAdaptor(sys, validator, 1, 2, 3, 4, 5)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, spi
}

func TestNewSpiAdaptor(t *testing.T) {
	// arrange
	a := NewSpiBusAdaptor(nil, nil, 1, 2, 3, 4, 5)
	// act & assert
	assert.Equal(t, 1, a.SpiDefaultBusNumber())
	assert.Equal(t, 2, a.SpiDefaultChipNumber())
	assert.Equal(t, 3, a.SpiDefaultMode())
	assert.Equal(t, 4, a.SpiDefaultBitCount())
	assert.Equal(t, int64(5), a.SpiDefaultMaxSpeed())
	_, err := a.GetSpiConnection(10, 0, 0, 8, 10000000)
	assert.ErrorContains(t, err, "not connected")
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
	assert.Equal(t, 0, len(a.connections))
	// act
	con1, err1 := a.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	// assert
	assert.NoError(t, err1)
	assert.NotNil(t, con1)
	assert.Equal(t, 1, len(a.connections))
	// assert cached connection
	con1a, err2 := a.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	assert.NoError(t, err2)
	assert.Equal(t, con1, con1a)
	assert.Equal(t, 1, len(a.connections))
	// assert second connection
	con2, err3 := a.GetSpiConnection(busNum, chipNum+1, mode, bits, maxSpeed)
	assert.NoError(t, err3)
	assert.NotNil(t, con2)
	assert.NotEqual(t, con1, con2)
	assert.Equal(t, 2, len(a.connections))
	// assert bus validation error
	con, err := a.GetSpiConnection(busNum+1, chipNum, mode, bits, maxSpeed)
	assert.ErrorContains(t, err, "16 not valid")
	assert.Nil(t, con)
	// assert create error
	spi.CreateError = true
	con, err = a.GetSpiConnection(busNum, chipNum+2, mode, bits, maxSpeed)
	assert.ErrorContains(t, err, "error while create SPI connection in mock")
	assert.Nil(t, con)
}

func TestSpiFinalize(t *testing.T) {
	// arrange
	a, _ := initTestSpiBusAdaptorWithMockedSpi()
	_, e := a.GetSpiConnection(spiTestAllowedBus, 2, 3, 4, 5)
	assert.NoError(t, e)
	assert.Equal(t, 1, len(a.connections))
	// act
	err := a.Finalize()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(a.connections))
}

func TestSpiFinalizeWithError(t *testing.T) {
	// arrange
	a, spi := initTestSpiBusAdaptorWithMockedSpi()
	_, e := a.GetSpiConnection(spiTestAllowedBus, 2, 3, 4, 5)
	assert.NoError(t, e)
	spi.SetCloseError(true)
	// act
	err := a.Finalize()
	// assert
	assert.Contains(t, err.Error(), "error while SPI close")
}

func TestSpiReConnect(t *testing.T) {
	// arrange
	a, _ := initTestSpiBusAdaptorWithMockedSpi()
	assert.NoError(t, a.Finalize())
	// act
	assert.NoError(t, a.Connect())
	// assert
	assert.NotNil(t, a.connections)
	assert.Equal(t, 0, len(a.connections))
}
