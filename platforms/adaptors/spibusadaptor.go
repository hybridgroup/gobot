package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/system"
)

type spiBusNumberValidator func(busNumber int) error

type SpiBusAdaptor struct {
	sys               *system.Accesser
	validateBusNumber spiBusNumberValidator
	defaultBusNumber  int
	defaultChipNumber int
	defaultMode       int
	defaultBitsNumber int
	defaultMaxSpeed   int64
	mutex             sync.Mutex
	connections       map[string]spi.Connection
}

func NewSpiBusAdaptor(sys *system.Accesser, v spiBusNumberValidator, busNr, chipNum, mode, bits int, maxSpeed int64) *SpiBusAdaptor {
	a := &SpiBusAdaptor{
		sys:               sys,
		validateBusNumber: v,
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

	a.connections = make(map[string]spi.Connection)
	return nil
}

func (a *SpiBusAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var err error
	for _, con := range a.connections {
		if con != nil {
			if e := con.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// TODO: all client get the default values and call GetConnection
//      --> introduce with Functions to change and remove from Interface

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus numbers range between 0 and 65536, valid chip numbers are 0 ... 255.
func (a *SpiBusAdaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (spi.Connection, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.connections == nil {
		return nil, fmt.Errorf("not connected")
	}

	id := fmt.Sprintf("%d_%d", busNum, chipNum)

	con := a.connections[id]
	if con == nil {
		if err := a.validateBusNumber(busNum); err != nil {
			return nil, err
		}
		var err error
		if con, err = a.sys.NewSpiConnection(busNum, chipNum, mode, bits, maxSpeed); err != nil {
			return nil, err
		}
		a.connections[id] = con
	}

	return con, nil
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
