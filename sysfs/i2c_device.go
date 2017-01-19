package sysfs

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const (
	// ioctl signals
	I2C_SLAVE = 0x0703
	I2C_FUNCS = 0x0705
	I2C_SMBUS = 0x0720
	// Read/write markers
	I2C_SMBUS_READ  = 1
	I2C_SMBUS_WRITE = 0
	// Adapter functionality
	I2C_FUNC_SMBUS_READ_BLOCK_DATA  = 0x01000000
	I2C_FUNC_SMBUS_WRITE_BLOCK_DATA = 0x02000000
	// Transaction types
	I2C_SMBUS_BYTE                = 1
	I2C_SMBUS_BYTE_DATA           = 2
	I2C_SMBUS_WORD_DATA           = 3
	I2C_SMBUS_PROC_CALL           = 4
	I2C_SMBUS_BLOCK_DATA          = 5
	I2C_SMBUS_I2C_BLOCK_DATA      = 6
	I2C_SMBUS_BLOCK_PROC_CALL     = 7  /* SMBus 2.0 */
	I2C_SMBUS_BLOCK_DATA_PEC      = 8  /* SMBus 2.0 */
	I2C_SMBUS_PROC_CALL_PEC       = 9  /* SMBus 2.0 */
	I2C_SMBUS_BLOCK_PROC_CALL_PEC = 10 /* SMBus 2.0 */
	I2C_SMBUS_WORD_DATA_PEC       = 11 /* SMBus 2.0 */
)

type i2cSmbusIoctlData struct {
	readWrite byte
	command   byte
	size      uint32
	data      uintptr
}

type SMBusOperations interface {
	ReadByte() (val uint8, err error)
	ReadByteData(reg uint8) (val uint8, err error)
	ReadWordData(reg uint8) (val uint16, err error)
	ReadBlockData(b []byte) (n int, err error)
	WriteByte(val uint8) (err error)
	WriteByteData(reg uint8, val uint8) (err error)
	WriteBlockData(b []byte) (err error)
}

// I2cDevice is the interface to a specific i2c bus
type I2cDevice interface {
	io.ReadWriteCloser
	SMBusOperations
	SetAddress(int) error
}

type i2cDevice struct {
	file  File
	funcs uint64 // adapter functionality mask
}

// NewI2cDevice returns an io.ReadWriteCloser with the proper ioctrl given
// an i2c bus location.
// Device address parameter is optional and kept for compatibility reasons.
func NewI2cDevice(location string, address ...int) (d *i2cDevice, err error) {
	d = &i2cDevice{}

	if d.file, err = OpenFile(location, os.O_RDWR, os.ModeExclusive); err != nil {
		return
	}
	if err = d.queryFunctionality(); err != nil {
		return
	}

	if len(address) > 0 {
		err = d.SetAddress(address[0])
	}

	return
}

func (d *i2cDevice) queryFunctionality() (err error) {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_FUNCS,
		uintptr(unsafe.Pointer(&d.funcs)),
	)

	if errno != 0 {
		err = fmt.Errorf("Querying functionality failed with syscall.Errno %v", errno)
	}
	return
}

func (d *i2cDevice) SetAddress(address int) (err error) {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SLAVE,
		uintptr(byte(address)),
	)

	if errno != 0 {
		err = fmt.Errorf("Setting address failed with syscall.Errno %v", errno)
	}

	return
}

func (d *i2cDevice) Close() (err error) {
	return d.file.Close()
}

func (d *i2cDevice) ReadByte() (val uint8, err error) {
	var data uint8
	err = d.smbusAccess(I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, uintptr(unsafe.Pointer(&data)))
	return data, err
}

func (d *i2cDevice) ReadByteData(reg uint8) (val uint8, err error) {
	var data uint8
	err = d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data)))
	return data, err
}

func (d *i2cDevice) ReadWordData(reg uint8) (val uint16, err error) {
	var data uint16
	err = d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&data)))
	return data, err
}

func (d *i2cDevice) ReadBlockData(b []byte) (n int, err error) {
	// Command byte - a data byte which often selects a register on the device:
	// 	https://www.kernel.org/doc/Documentation/i2c/smbus-protocol
	command := byte(b[0])
	buf := b[1:]

	data := make([]byte, len(buf)+1)
	data[0] = byte(len(buf))

	copy(data[1:], buf)

	err = d.smbusAccess(I2C_SMBUS_READ, command, I2C_SMBUS_I2C_BLOCK_DATA, uintptr(unsafe.Pointer(&data[0])))

	copy(b, data[1:])
	return int(data[0]), err
}

func (d *i2cDevice) WriteByte(val uint8) (err error) {
	err = d.smbusAccess(I2C_SMBUS_WRITE, val, I2C_SMBUS_BYTE, uintptr(0))
	return err
}

func (d *i2cDevice) WriteByteData(reg uint8, val uint8) (err error) {
	var data uint8 = val
	err = d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data)))
	return err
}

func (d *i2cDevice) WriteBlockData(b []byte) (err error) {
	// Command byte - a data byte which often selects a register on the device:
	// 	https://www.kernel.org/doc/Documentation/i2c/smbus-protocol
	command := byte(b[0])
	buf := b[1:]

	data := make([]byte, len(buf)+1)
	data[0] = byte(len(buf))

	copy(data[1:], buf)

	err = d.smbusAccess(I2C_SMBUS_WRITE, command, I2C_SMBUS_I2C_BLOCK_DATA, uintptr(unsafe.Pointer(&data[0])))

	return err
}

// Read implements the io.ReadWriteCloser method by SMBus block data reads.
// If SMBus block read is not supported, direct read (i2c mode) is performed
// instead.
func (d *i2cDevice) Read(b []byte) (n int, err error) {
	if d.funcs&I2C_FUNC_SMBUS_READ_BLOCK_DATA == 0 {
		// Adapter doesn't support SMBus block read
		return d.file.Read(b)
	}

	return d.ReadBlockData(b)
}

// Write implements the io.ReadWriteCloser method by SMBus block data writes.
// If SMBus block write is not supported, direct write (i2c mode) is performed
// instead.
func (d *i2cDevice) Write(b []byte) (n int, err error) {
	if d.funcs&I2C_FUNC_SMBUS_WRITE_BLOCK_DATA == 0 {
		// Adapter doesn't support SMBus block write
		return d.file.Write(b)
	}

	return len(b), d.WriteBlockData(b)
}

func (d *i2cDevice) smbusAccess(readWrite byte, command byte, size uint32, data uintptr) error {
	smbus := &i2cSmbusIoctlData{
		readWrite: readWrite,
		command:   command,
		size:      size,
		data:      data,
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		return fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return nil
}
