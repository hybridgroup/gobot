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

type I2cOperations interface {
	io.ReadWriteCloser
	ReadByte() (val uint8, err error)
	ReadByteData(reg uint8) (val uint8, err error)
	ReadWordData(reg uint8) (val uint16, err error)
	ReadBlockData(reg uint8, b []byte) (n int, err error)
	WriteByte(val uint8) (err error)
	WriteByteData(reg uint8, val uint8) (err error)
	WriteWordData(reg uint8, val uint16) (err error)
	WriteBlockData(reg uint8, b []byte) (err error)
}

// I2cDevice is the interface to a specific i2c bus
type I2cDevice interface {
	I2cOperations
	SetAddress(int) error
}

type i2cDevice struct {
	file  File
	funcs uint64 // adapter functionality mask
}

// NewI2cDevice returns an io.ReadWriteCloser with the proper ioctrl given
// an i2c bus location.
func NewI2cDevice(location string) (d *i2cDevice, err error) {
	d = &i2cDevice{}

	if d.file, err = OpenFile(location, os.O_RDWR, os.ModeExclusive); err != nil {
		return
	}
	if err = d.queryFunctionality(); err != nil {
		return
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

func (d *i2cDevice) ReadBlockData(reg uint8, buf []byte) (n int, err error) {
	if d.funcs&I2C_FUNC_SMBUS_READ_BLOCK_DATA == 0 {
		return 0, fmt.Errorf("SMBus block data reading not supported")
	}

	data := make([]byte, 32+1) // Max message + size as defined by SMBus standard

	err = d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_I2C_BLOCK_DATA, uintptr(unsafe.Pointer(&data[0])))

	copy(buf, data[1:])
	return int(data[0]), err
}

func (d *i2cDevice) WriteByte(val uint8) (err error) {
	err = d.smbusAccess(I2C_SMBUS_WRITE, val, I2C_SMBUS_BYTE, uintptr(0))
	return err
}

func (d *i2cDevice) WriteByteData(reg uint8, val uint8) (err error) {
	var data = val
	err = d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data)))
	return err
}

func (d *i2cDevice) WriteWordData(reg uint8, val uint16) (err error) {
	var data = val
	err = d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&data)))
	return err
}

func (d *i2cDevice) WriteBlockData(reg uint8, data []byte) (err error) {
	if d.funcs&I2C_FUNC_SMBUS_WRITE_BLOCK_DATA == 0 {
		return fmt.Errorf("SMBus block data writing not supported")
	}

	if len(data) > 32 {
		data = data[:32]
	}

	buf := make([]byte, len(data)+1)
	copy(buf[:1], data)
	buf[0] = uint8(len(data))

	err = d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_I2C_BLOCK_DATA, uintptr(unsafe.Pointer(&buf[0])))

	return err
}

// Read implements the io.ReadWriteCloser method by direct I2C read operations.
func (d *i2cDevice) Read(b []byte) (n int, err error) {
	return d.file.Read(b)
}

// Write implements the io.ReadWriteCloser method by direct I2C write operations.
func (d *i2cDevice) Write(b []byte) (n int, err error) {
	return d.file.Write(b)
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
