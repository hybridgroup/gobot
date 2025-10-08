package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/onewire"
	"gobot.io/x/gobot/v2/system"
)

// oneWireBusConfiguration contains all changeable attributes of the adaptor.
type oneWireBusConfiguration struct {
	debug bool
}

// OneWireBusAdaptor is a adaptor for the 1-wire bus, normally used for composition in platforms.
// note: currently only one controller is supported by most platforms, but it would be possible to activate more,
// see https://forums.raspberrypi.com/viewtopic.php?t=65137
//
// Options:
//
//	"WithOneWireDebug"
type OneWireBusAdaptor struct {
	sys           *system.Accesser
	oneWireBusCfg *oneWireBusConfiguration
	mutex         *sync.Mutex
	connections   map[string]onewire.Connection
}

// NewOneWireBusAdaptor provides the access to 1-wire devices of the board.
func NewOneWireBusAdaptor(sys *system.Accesser, opts ...OneWireBusOptionApplier) *OneWireBusAdaptor {
	a := OneWireBusAdaptor{
		sys:           sys,
		oneWireBusCfg: &oneWireBusConfiguration{},
		mutex:         &sync.Mutex{},
	}

	for _, o := range opts {
		o.apply(a.oneWireBusCfg)
	}

	sys.AddOneWireSupport()

	return &a
}

// WithOneWireDebug can be used to switch on debugging for 1-wire implementation.
func WithOneWireDebug() oneWireBusDebugOption {
	return oneWireBusDebugOption(true)
}

// Connect prepares the connection to 1-wire devices.
func (a *OneWireBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.connections != nil {
		return fmt.Errorf("1-wire bus adaptor already connected, please call Finalize() for re-connect")
	}

	a.connections = make(map[string]onewire.Connection)
	a.debuglnf("connect the 1-wire bus adaptor done")

	return nil
}

// Finalize closes all 1-wire connections.
func (a *OneWireBusAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.debuglnf("finalize the 1-wire bus adaptor for %d buses...", len(a.connections))

	var err error
	for id, con := range a.connections {
		if con != nil {
			e := con.Close()
			if e != nil {
				err = multierror.Append(err, e)
			}
			a.debuglnf("1-wire bus %s closed with error: %v", id, e)
		}
	}
	a.connections = nil
	a.debuglnf("finalize the 1-wire bus adaptor done with error: %v", err)

	return err
}

// GetOneWireConnection returns a 1-wire connection to a device with the given family code and serial number.
func (a *OneWireBusAdaptor) GetOneWireConnection(familyCode byte, serialNumber uint64) (onewire.Connection, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.connections == nil {
		return nil, fmt.Errorf("not connected")
	}

	id := fmt.Sprintf("%d_%d", familyCode, serialNumber)

	con := a.connections[id]
	if con == nil {
		var err error
		dev, err := a.sys.NewOneWireDevice(familyCode, serialNumber)
		if err != nil {
			return nil, err
		}
		con = onewire.NewConnection(dev)
		a.connections[id] = con
	}

	return con, nil
}

func (a *OneWireBusAdaptor) debuglnf(format string, p ...interface{}) {
	gobot.Debuglnf(a.oneWireBusCfg.debug, format, p...)
}
