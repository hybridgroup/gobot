package onewire

import (
	"sync"

	"gobot.io/x/gobot/v2"
)

// onewireConnection is the common implementation of the 1-wire bus interface.
type onewireConnection struct {
	onewireSystem gobot.OneWireSystemDevicer
	mutex         *sync.Mutex
}

// NewConnection uses the given 1-wire system device and provides it as gobot.OneWireOperations.
func NewConnection(onewireSystem gobot.OneWireSystemDevicer) *onewireConnection {
	return &onewireConnection{onewireSystem: onewireSystem, mutex: &sync.Mutex{}}
}

// ID returns the device id in the form "family code"-"serial number". Implements gobot.OneWireOperations.
func (d *onewireConnection) ID() string {
	return d.onewireSystem.ID()
}

// ReadData reads the data according the command, e.g. from the specified file on sysfs bus.
// Implements gobot.OneWireOperations.
func (c *onewireConnection) ReadData(command string, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.onewireSystem.ReadData(command, data)
}

// WriteData writes the data according the command, e.g. to the specified file on sysfs bus.
// Implements gobot.OneWireOperations.
func (c *onewireConnection) WriteData(command string, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.onewireSystem.WriteData(command, data)
}

// ReadInteger reads the value according the command, e.g. to the specified file on sysfs bus.
// Implements gobot.OneWireOperations.
func (c *onewireConnection) ReadInteger(command string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.onewireSystem.ReadInteger(command)
}

// WriteInteger writes the value according the command, e.g. to the specified file on sysfs bus.
// Implements gobot.OneWireOperations.
func (c *onewireConnection) WriteInteger(command string, val int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.onewireSystem.WriteInteger(command, val)
}

// Close connection to underlying 1-wire device. Implements functions of onewire.Connection respectively
// gobot.OneWireOperations.
func (c *onewireConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.onewireSystem.Close()
}
