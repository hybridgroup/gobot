package adaptors

import (
	"fmt"
	"testing"

	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

// make sure that this SpiBusAdaptor fulfills all the required interfaces
var _ spi.Connector = (*SpiBusAdaptor)(nil)

func initTestSpiBusAdaptorWithMockedFilesystem() *SpiBusAdaptor {
	validator := func(busNr int) error {
		if busNr > 1 {
			return fmt.Errorf("%d not valid", busNr)
		}
		return nil
	}
	sys := system.NewAccesser()
	a := NewSpiBusAdaptor(sys, validator, 1, 2, 3, 4, 5)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestNewSpiAdaptor(t *testing.T) {
	// arrange
	sys := system.NewAccesser()
	validator := func(busNr int) error {
		if busNr > 1 {
			return fmt.Errorf("%d not valid", busNr)
		}
		return nil
	}
	a := NewSpiBusAdaptor(sys, validator, 1, 2, 3, 4, 5)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	// act & assert
	gobottest.Assert(t, a.SpiDefaultBusNumber(), 1)
	gobottest.Assert(t, a.SpiDefaultChipNumber(), 2)
	gobottest.Assert(t, a.SpiDefaultMode(), 3)
	gobottest.Assert(t, a.SpiDefaultBitCount(), 4)
	gobottest.Assert(t, a.SpiDefaultMaxSpeed(), int64(5))

	_, err := a.GetSpiConnection(10, 0, 0, 8, 10000000)
	gobottest.Assert(t, err.Error(), "10 not valid")

	// TODO: tests for real connection currently not possible, because not using system.Accessor
	// TODO: test tx/rx here...
}

func TestSpiFinalizeWithErrors(t *testing.T) {
	// arrange
	a := initTestSpiBusAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.Connect(), nil)
	a.GetSpiConnection(1, 2, 3, 4, 5)
	//gobottest.Assert(t, err, nil)
	//err = con.Tx([]byte{}, []byte{})
	//gobottest.Assert(t, err, nil)
	// act
	a.Finalize()
	// assert
	//gobottest.Assert(t, strings.Contains(err.Error(), "close error"), true)
}
