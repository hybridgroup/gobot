package i2c

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
	"time"

	"gobot.io/x/gobot"
)

// Driver implements the interface gobot.Driver.
type Driver struct {
	name           string
	defaultAddress int
	connector      Connector
	connection     Connection
	afterStart     func() error
	beforeHalt     func() error
	Config
	gobot.Commander
	mutex *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
}

// NewDriver creates a new generic and basic i2c gobot driver.
func NewDriver(c Connector, name string, address int, options ...func(Config)) *Driver {
	d := &Driver{
		name:           gobot.DefaultName(name),
		defaultAddress: address,
		connector:      c,
		afterStart:     func() error { return nil },
		beforeHalt:     func() error { return nil },
		Config:         NewConfig(),
		Commander:      gobot.NewCommander(),
		mutex:          &sync.Mutex{},
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// Name returns the name of the i2c device.
func (d *Driver) Name() string {
	return d.name
}

// SetName sets the name of the i2c device.
func (d *Driver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the i2c device.
func (d *Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the i2c device.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var err error
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(int(d.defaultAddress))

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	return d.afterStart()
}

// Halt halts the i2c device.
func (d *Driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.beforeHalt(); err != nil {
		return err
	}

	// currently there is nothing to do here for the driver
	return nil
}

// Write writes one byte to the i2c device.
func (d *Driver) WriteByte(val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.WriteByte(val)
}

// Write implements a simple write mechanism, starting from the given register of an i2c device.
func (d *Driver) Write(pin string, val int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	if val > 0xFFFF {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(val))
		return d.connection.WriteBlockData(register, buf)
	}
	if val > 0xFF {
		return d.connection.WriteWordData(register, uint16(val))
	}
	return d.connection.WriteByteData(register, uint8(val))
}

// WriteByteData writes the given byte value to the given register of an i2c device.
func (d *Driver) WriteByteData(pin string, val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	return d.connection.WriteByteData(register, val)
}

// WriteWordData writes the given 16 bit value to the given register of an i2c device.
func (d *Driver) WriteWordData(pin string, val uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	return d.connection.WriteWordData(register, val)
}

// WriteBlockData writes the given buffer to the given register of an i2c device.
func (d *Driver) WriteBlockData(pin string, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	return d.connection.WriteBlockData(uint8(register), data)
}

// WriteData writes the given buffer to the given register of an i2c device.
// It uses plain write to prevent WriteBlockData(), which is sometimes not supported by adaptor.
func (d *Driver) WriteData(pin string, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = register

	cnt, err := d.connection.Write(buf)
	if cnt != len(buf) {
		return fmt.Errorf("written count (%d) differ from expected (%d)", cnt, len(buf))
	}

	return err
}

// Read implements a simple read mechanism from the given register of an i2c device.
func (d *Driver) Read(pin string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return 0, err
	}

	val, err := d.connection.ReadByteData(register)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

// ReadByte reads a byte from the current register of an i2c device.
func (d *Driver) ReadByte() (byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.ReadByte()
}

// ReadByteData reads a byte from the given register of an i2c device.
func (d *Driver) ReadByteData(pin string) (byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return 0, err
	}

	return d.connection.ReadByteData(register)
}

// ReadWordData reads a 16 bit value starting from the given register of an i2c device.
func (d *Driver) ReadWordData(pin string) (uint16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return 0, err
	}

	return d.connection.ReadWordData(register)
}

// ReadBlockData fills the given buffer with reads starting from the given register of an i2c device.
func (d *Driver) ReadBlockData(pin string, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	return d.connection.ReadBlockData(register, data)
}

// ReadBlockData fills the given buffer with reads from the given register of an i2c device.
// It uses plain read to prevent ReadBlockData(), which is sometimes not supported by adaptor.
func (d *Driver) ReadData(pin string, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	if err := d.connection.WriteByte(register); err != nil {
		return err
	}

	// write process needs some time, so wait at least 5ms before read a value
	// when decreasing to much, the check below will fail
	time.Sleep(10 * time.Millisecond)

	n, err := d.connection.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("Read %v bytes from device by sysfs, expected %v", n, len(data))
	}
	return nil
}

func driverParseRegister(pin string) (uint8, error) {
	register, err := strconv.ParseUint(pin, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("Could not parse the register from given pin '%s'", pin)
	}
	return uint8(register), nil
}
