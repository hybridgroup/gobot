package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/system"
)

type i2cBusNumberValidator func(busNumber int) error

// I2cBusAdaptor is a adaptor for i2c bus, normally used for composition in platforms.
type I2cBusAdaptor struct {
	sys              *system.Accesser
	validateNumber   i2cBusNumberValidator
	defaultBusNumber int
	mutex            sync.Mutex
	buses            map[int]i2c.I2cDevice
}

// NewI2cBusAdaptor provides the access to i2c buses of the board. The validator is used to check the bus number,
// which is given by user, to the abilities of the board.
func NewI2cBusAdaptor(sys *system.Accesser, v i2cBusNumberValidator, defaultBusNr int) *I2cBusAdaptor {
	a := &I2cBusAdaptor{
		sys:              sys,
		validateNumber:   v,
		defaultBusNumber: defaultBusNr,
	}
	return a
}

// Connect prepares the connection to i2c buses.
func (a *I2cBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.buses = make(map[int]i2c.I2cDevice)
	return nil
}

// Finalize closes all i2c connections.
func (a *I2cBusAdaptor) Finalize() error {
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
	a.buses = nil
	return err
}

// GetConnection returns a connection to a device on a specified i2c bus
func (a *I2cBusAdaptor) GetConnection(address int, busNr int) (connection i2c.Connection, err error) {
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
		bus, err = a.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", busNr))
		if err != nil {
			return nil, err
		}
		a.buses[busNr] = bus
	}
	return i2c.NewConnection(bus, address), err
}

// GetDefaultBus returns the default i2c bus number for this platform.
func (a *I2cBusAdaptor) GetDefaultBus() int {
	return a.defaultBusNumber
}
