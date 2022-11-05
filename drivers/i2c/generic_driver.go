package i2c

import (
	"fmt"
	"time"
)

// GenericDriver implements the interface gobot.Driver.
type GenericDriver struct {
	*Driver
}

// NewGenericDriver creates a new generic i2c gobot driver, which just forwards all connection functions.
func NewGenericDriver(c Connector, name string, address int, options ...func(Config)) *GenericDriver {
	return &GenericDriver{Driver: NewDriver(c, name, address, options...)}
}

// WriteByte writes one byte to the i2c device.
func (d *GenericDriver) WriteByte(val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.WriteByte(val)
}

// WriteByteData writes the given byte value to the given register of an i2c device.
func (d *GenericDriver) WriteByteData(reg uint8, val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.WriteByteData(reg, val)
}

// WriteWordData writes the given 16 bit value to the given register of an i2c device.
func (d *GenericDriver) WriteWordData(reg uint8, val uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.WriteWordData(reg, val)
}

// WriteBlockData writes the given buffer to the given register of an i2c device.
func (d *GenericDriver) WriteBlockData(reg uint8, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.WriteBlockData(reg, data)
}

// WriteData writes the given buffer to the given register of an i2c device.
// It uses plain write to prevent WriteBlockData(), which is sometimes not supported by adaptor.
func (d *GenericDriver) WriteData(reg uint8, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = reg

	return d.writeAndCheckCount(buf)
}

// Write writes the given buffer to the i2c device.
func (d *GenericDriver) Write(data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.writeAndCheckCount(data)
}

// ReadByte reads a byte from the current register of an i2c device.
func (d *GenericDriver) ReadByte() (byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.ReadByte()
}

// ReadByteData reads a byte from the given register of an i2c device.
func (d *GenericDriver) ReadByteData(reg uint8) (byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.ReadByteData(reg)
}

// ReadWordData reads a 16 bit value starting from the given register of an i2c device.
func (d *GenericDriver) ReadWordData(reg uint8) (uint16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.ReadWordData(reg)
}

// ReadBlockData fills the given buffer with reads starting from the given register of an i2c device.
func (d *GenericDriver) ReadBlockData(reg uint8, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.connection.ReadBlockData(reg, data)
}

// ReadData fills the given buffer with reads from the given register of an i2c device.
// It uses plain read to prevent ReadBlockData(), which is sometimes not supported by adaptor.
func (d *GenericDriver) ReadData(reg uint8, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.connection.WriteByte(reg); err != nil {
		return err
	}

	// write process needs some time, so wait at least 5ms before read a value
	// when decreasing to much, the check below will fail
	time.Sleep(10 * time.Millisecond)

	return d.readAndCheckCount(data)
}

// Read fills the given buffer with reads of an i2c device.
func (d *GenericDriver) Read(data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.readAndCheckCount(data)
}

func (d *GenericDriver) writeAndCheckCount(data []byte) error {
	n, err := d.connection.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("written count (%d) differ from expected (%d)", n, len(data))
	}
	return nil
}

func (d *GenericDriver) readAndCheckCount(data []byte) error {
	n, err := d.connection.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("read count (%d) differ from expected (%d)", n, len(data))
	}
	return nil
}
