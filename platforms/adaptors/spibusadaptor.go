package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/drivers/spi"
)

type spiBusNumberValidator func(busNumber int) error

type SpiBusAdaptor struct {
	validateNumber    spiBusNumberValidator
	defaultBusNumber  int
	defaultChipNumber int
	defaultMode       int
	defaultBitsNumber int
	defaultMaxSpeed   int64
	mutex             sync.Mutex
	buses             map[int]spi.Connection
}

func NewSpiBusAdaptor(v spiBusNumberValidator, busNr, chipNum, mode, bits int, maxSpeed int64) *SpiBusAdaptor {
	a := &SpiBusAdaptor{
		validateNumber:    v,
		defaultBusNumber:  busNr,
		defaultChipNumber: chipNum,
		defaultMode:       mode,
		defaultBitsNumber: bits,
		defaultMaxSpeed:   maxSpeed,
	}
	return a
}

// Connect prepares the connection to SPI buses.
func (a *SpiBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.buses = make(map[int]spi.Connection)
	return nil
}

func (a *SpiBusAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var err error
	for _, bus := range a.buses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// TODO: all client get the default values and call GetConnection
//      --> introduce with Functions to change and remove from Interface

// GetSpiConnection returns an spi connection to a device on a specified bus.
func (a *SpiBusAdaptor) GetSpiConnection(busNr, chipNum, mode, bits int, maxSpeed int64) (spi.Connection, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.buses == nil {
		return nil, fmt.Errorf("not connected")
	}

	bus := a.buses[busNr]
	if bus == nil {
		if err := a.validateNumber(busNr); err != nil {
			return nil, err
		}
		var err error
		if bus, err = spi.GetSpiConnection(busNr, chipNum, mode, bits, maxSpeed); err != nil {
			return nil, err
		}
		a.buses[busNr] = bus
	}

	return bus, nil
}

// TODO: remove Get

// GetSpiDefaultBus returns the default spi bus for this platform.
func (a *SpiBusAdaptor) GetSpiDefaultBus() int {
	return a.defaultBusNumber
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (a *SpiBusAdaptor) GetSpiDefaultChip() int {
	return a.defaultChipNumber
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (a *SpiBusAdaptor) GetSpiDefaultMode() int {
	return a.defaultMode
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (a *SpiBusAdaptor) GetSpiDefaultBits() int {
	return a.defaultBitsNumber
}

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (a *SpiBusAdaptor) GetSpiDefaultMaxSpeed() int64 {
	return a.defaultMaxSpeed
}
