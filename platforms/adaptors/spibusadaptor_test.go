package adaptors

import (
	"fmt"
	"testing"

	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/gobottest"
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

	a := NewSpiBusAdaptor(validator, 1, 2, 3, 4, 5)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestNewSpiAdaptor(t *testing.T) {
	// arrange
	validator := func(busNr int) error {
		if busNr > 1 {
			return fmt.Errorf("%d not valid", busNr)
		}
		return nil
	}
	a := NewSpiBusAdaptor(validator, 1, 2, 3, 4, 5)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	// act & assert
	gobottest.Assert(t, a.GetSpiDefaultBus(), 1)
	gobottest.Assert(t, a.GetSpiDefaultChip(), 2)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 3)
	gobottest.Assert(t, a.GetSpiDefaultBits(), 4)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(5))

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
