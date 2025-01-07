package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/system"
)

type i2cBusNumberValidator func(busNumber int) error

// I2cBusAdaptor is a adaptor for i2c bus, normally used for composition in platforms.
type I2cBusAdaptor struct {
	sys              *system.Accesser
	validateNumber   i2cBusNumberValidator
	defaultBusNumber int
	mutex            sync.Mutex
	buses            map[int]gobot.I2cSystemDevicer
}

// NewI2cBusAdaptor provides the access to i2c buses of the board. The validator is used to check the bus number,
// which is given by user, to the abilities of the board.
func NewI2cBusAdaptor(sys *system.Accesser, v i2cBusNumberValidator, defaultBusNr int) *I2cBusAdaptor {
	a := I2cBusAdaptor{
		sys:              sys,
		validateNumber:   v,
		defaultBusNumber: defaultBusNr,
	}

	sys.AddI2CSupport()

	return &a
}

// Connect prepares the connection to i2c buses.
func (a *I2cBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.buses = make(map[int]gobot.I2cSystemDevicer)
	return nil
}

// Finalize closes all i2c buses.
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

// GetI2cConnection returns a connection to a device on a specified i2c bus
func (a *I2cBusAdaptor) GetI2cConnection(address int, busNum int) (i2c.Connection, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.buses == nil {
		return nil, fmt.Errorf("not connected")
	}

	bus := a.buses[busNum]
	if bus == nil {
		err := a.validateNumber(busNum)
		if err != nil {
			return nil, err
		}
		bus, err = a.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", busNum))
		if err != nil {
			return nil, err
		}
		a.buses[busNum] = bus
	}
	return i2c.NewConnection(bus, address), nil
}

// DefaultI2cBus returns the default i2c bus number for this platform.
func (a *I2cBusAdaptor) DefaultI2cBus() int {
	return a.defaultBusNumber
}
