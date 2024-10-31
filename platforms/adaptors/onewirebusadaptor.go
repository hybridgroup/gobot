package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2/drivers/onewire"
	"gobot.io/x/gobot/v2/system"
)

// OneWireBusAdaptor is a adaptor for the 1-wire bus, normally used for composition in platforms.
// note: currently only one controller is supported by most platforms, but it would be possible to activate more,
// see https://forums.raspberrypi.com/viewtopic.php?t=65137
type OneWireBusAdaptor struct {
	sys         *system.Accesser
	mutex       *sync.Mutex
	connections map[string]onewire.Connection
}

// NewOneWireBusAdaptor provides the access to 1-wire devices of the board.
func NewOneWireBusAdaptor(sys *system.Accesser) *OneWireBusAdaptor {
	a := &OneWireBusAdaptor{sys: sys, mutex: &sync.Mutex{}}
	return a
}

// Connect prepares the connection to 1-wire devices.
func (a *OneWireBusAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.connections = make(map[string]onewire.Connection)
	return nil
}

// Finalize closes all 1-wire connections.
func (a *OneWireBusAdaptor) Finalize() error {
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
