package adaptors

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/gobottest"
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
	gobottest.Assert(t, a.SpiDefaultBusNumber(), 1)
	gobottest.Assert(t, a.SpiDefaultChipNumber(), 2)
	gobottest.Assert(t, a.SpiDefaultMode(), 3)
	gobottest.Assert(t, a.SpiDefaultBitCount(), 4)
	gobottest.Assert(t, a.SpiDefaultMaxSpeed(), int64(5))
	_, err := a.GetSpiConnection(10, 0, 0, 8, 10000000)
	gobottest.Assert(t, err.Error(), "not connected")
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
	gobottest.Assert(t, len(a.connections), 0)
	// act
	con1, err1 := a.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	// assert
	gobottest.Assert(t, err1, nil)
	gobottest.Refute(t, con1, nil)
	gobottest.Assert(t, len(a.connections), 1)
	// assert cached connection
	con1a, err2 := a.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	gobottest.Assert(t, err2, nil)
	gobottest.Assert(t, con1a, con1)
	gobottest.Assert(t, len(a.connections), 1)
	// assert second connection
	con2, err3 := a.GetSpiConnection(busNum, chipNum+1, mode, bits, maxSpeed)
	gobottest.Assert(t, err3, nil)
	gobottest.Refute(t, con2, nil)
	gobottest.Refute(t, con2, con1)
	gobottest.Assert(t, len(a.connections), 2)
	// assert bus validation error
	con, err := a.GetSpiConnection(busNum+1, chipNum, mode, bits, maxSpeed)
	gobottest.Assert(t, err.Error(), "16 not valid")
	gobottest.Assert(t, con, nil)
	// assert create error
	spi.CreateError = true
	con, err = a.GetSpiConnection(busNum, chipNum+2, mode, bits, maxSpeed)
	gobottest.Assert(t, err.Error(), "error while create SPI connection in mock")
	gobottest.Assert(t, con, nil)
}

func TestSpiFinalize(t *testing.T) {
	// arrange
	a, _ := initTestSpiBusAdaptorWithMockedSpi()
	_, e := a.GetSpiConnection(spiTestAllowedBus, 2, 3, 4, 5)
	gobottest.Assert(t, e, nil)
	gobottest.Assert(t, len(a.connections), 1)
	// act
	err := a.Finalize()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.connections), 0)
}

func TestSpiFinalizeWithError(t *testing.T) {
	// arrange
	a, spi := initTestSpiBusAdaptorWithMockedSpi()
	_, e := a.GetSpiConnection(spiTestAllowedBus, 2, 3, 4, 5)
	gobottest.Assert(t, e, nil)
	spi.SetCloseError(true)
	// act
	err := a.Finalize()
	// assert
	gobottest.Assert(t, strings.Contains(err.Error(), "error while SPI close"), true)
}

func TestSpiReConnect(t *testing.T) {
	// arrange
	a, _ := initTestSpiBusAdaptorWithMockedSpi()
	gobottest.Assert(t, a.Finalize(), nil)
	// act
	gobottest.Assert(t, a.Connect(), nil)
	// assert
	gobottest.Refute(t, a.connections, nil)
	gobottest.Assert(t, len(a.connections), 0)
}
