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
	connections      map[string]i2c.Connection
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

	a.connections = make(map[string]i2c.Connection)
	return nil
}

// Finalize closes all i2c connections.
func (a *I2cBusAdaptor) Finalize() error {
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
	a.connections = nil
	return err
}

// GetI2cConnection returns a connection to a device on a specified i2c bus
func (a *I2cBusAdaptor) GetI2cConnection(address int, busNum int) (connection i2c.Connection, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.connections == nil {
		return nil, fmt.Errorf("not connected")
	}

	id := fmt.Sprintf("%d_%d", busNum, address)

	con := a.connections[id]
	if con == nil {
		if err := a.validateNumber(busNum); err != nil {
			return nil, err
		}
		bus, err := a.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", busNum))
		if err != nil {
			return nil, err
		}
		con = i2c.NewConnection(bus, address)
		a.connections[id] = con
	}
	return con, err
}

// DefaultI2cBus returns the default i2c bus number for this platform.
func (a *I2cBusAdaptor) DefaultI2cBus() int {
	return a.defaultBusNumber
}
