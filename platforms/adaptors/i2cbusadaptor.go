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

// i2cBusConfiguration contains all changeable attributes of the adaptor.
type i2cBusConfiguration struct {
	debug bool
}

// I2cBusAdaptor is a adaptor for i2c bus, normally used for composition in platforms.
type I2cBusAdaptor struct {
	sys              *system.Accesser
	i2cBusCfg        *i2cBusConfiguration
	validateNumber   i2cBusNumberValidator
	defaultBusNumber int
	mutex            sync.Mutex
	buses            map[int]gobot.I2cSystemDevicer
}

// NewI2cBusAdaptor provides the access to i2c buses of the board. The validator is used to check the bus number,
// which is given by user, to the abilities of the board.
//
// Options:
//
//	"WithI2cDebug"
func NewI2cBusAdaptor(
	sys *system.Accesser,
	v i2cBusNumberValidator,
	defaultBusNr int,
	opts ...I2CBusOptionApplier,
) *I2cBusAdaptor {
	a := I2cBusAdaptor{
		sys:              sys,
		i2cBusCfg:        &i2cBusConfiguration{},
		validateNumber:   v,
		defaultBusNumber: defaultBusNr,
	}

	for _, o := range opts {
		o.apply(a.i2cBusCfg)
	}

	sys.AddI2CSupport()

	return &a
}

// WithI2cDebug can be used to switch on debugging for I2C implementation.
func WithI2cDebug() i2cBusDebugOption {
	return i2cBusDebugOption(true)
}

// Connect prepares the connection to i2c buses.
func (a *I2cBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.buses != nil {
		return fmt.Errorf("I2C bus adaptor already connected, please call Finalize() for re-connect")
	}

	a.buses = make(map[int]gobot.I2cSystemDevicer)
	a.debuglnf("connect the I2C bus adaptor done")

	return nil
}

// Finalize closes all i2c buses.
func (a *I2cBusAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.debuglnf("finalize the I2C bus adaptor for %d buses...", len(a.buses))

	var err error
	for busNum, bus := range a.buses {
		if bus != nil {
			e := bus.Close()
			if e != nil {
				err = multierror.Append(err, e)
			}
			a.debuglnf("I2C bus %d closed with error: %v", busNum, e)
		}
	}
	a.buses = nil
	a.debuglnf("finalize the I2C bus adaptor done with error: %v", err)

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

func (a *I2cBusAdaptor) debuglnf(format string, p ...interface{}) {
	gobot.Debuglnf(a.i2cBusCfg.debug, format, p...)
}
