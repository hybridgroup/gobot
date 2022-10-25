package sysfs

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const i2cDeviceDebug = false

const (
	// From  /usr/include/linux/i2c-dev.h:
	// ioctl signals
	I2C_SLAVE = 0x0703
	I2C_FUNCS = 0x0705
	I2C_SMBUS = 0x0720
	// Read/write markers
	I2C_SMBUS_READ  = 1
	I2C_SMBUS_WRITE = 0

	// From  /usr/include/linux/i2c.h:
	// Adapter functionality
	I2C_FUNC_SMBUS_READ_BYTE        = 0x00020000
	I2C_FUNC_SMBUS_WRITE_BYTE       = 0x00040000
	I2C_FUNC_SMBUS_READ_BYTE_DATA   = 0x00080000
	I2C_FUNC_SMBUS_WRITE_BYTE_DATA  = 0x00100000
	I2C_FUNC_SMBUS_READ_WORD_DATA   = 0x00200000
	I2C_FUNC_SMBUS_WRITE_WORD_DATA  = 0x00400000
	I2C_FUNC_SMBUS_READ_BLOCK_DATA  = 0x01000000
	I2C_FUNC_SMBUS_WRITE_BLOCK_DATA = 0x02000000
	// Transaction types
	I2C_SMBUS_BYTE             = 1
	I2C_SMBUS_BYTE_DATA        = 2
	I2C_SMBUS_WORD_DATA        = 3
	I2C_SMBUS_PROC_CALL        = 4
	I2C_SMBUS_BLOCK_DATA       = 5
	I2C_SMBUS_I2C_BLOCK_BROKEN = 6
	I2C_SMBUS_BLOCK_PROC_CALL  = 7 /* SMBus 2.0 */
	I2C_SMBUS_I2C_BLOCK_DATA   = 8 /* SMBus 2.0 */
)

type i2cSmbusIoctlData struct {
	readWrite byte
	command   byte
	size      uint32
	data      uintptr
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

func (d *i2cDevice) queryFunctionality() error {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_FUNCS,
		uintptr(unsafe.Pointer(&d.funcs)),
	)

	if errno != 0 {
		return fmt.Errorf("Querying functionality failed with syscall.Errno %v", errno)
	}
	return nil
}

func (d *i2cDevice) SetAddress(address int) error {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SLAVE,
		uintptr(byte(address)),
	)

	if errno != 0 {
		return fmt.Errorf("Setting address failed with syscall.Errno %v", errno)
	}
	return nil
}

func (d *i2cDevice) Close() error {
	return d.file.Close()
}

func (d *i2cDevice) ReadByte() (val byte, err error) {
	if d.funcs&I2C_FUNC_SMBUS_READ_BYTE == 0 {
		return 0, fmt.Errorf("SMBus read byte not supported")
	}

	var data uint8
	err = d.smbusAccess(I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, uintptr(unsafe.Pointer(&data)))
	return data, err
}

func (d *i2cDevice) ReadByteData(reg uint8) (val uint8, err error) {
	if d.funcs&I2C_FUNC_SMBUS_READ_BYTE_DATA == 0 {
		return 0, fmt.Errorf("SMBus read byte data not supported")
	}

	var data uint8
	err = d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data)))
	return data, err
}

func (d *i2cDevice) ReadWordData(reg uint8) (val uint16, err error) {
	if d.funcs&I2C_FUNC_SMBUS_READ_WORD_DATA == 0 {
		return 0, fmt.Errorf("SMBus read word data not supported")
	}

	var data uint16
	err = d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&data)))
	return data, err
}

func (d *i2cDevice) ReadBlockData(reg uint8, data []byte) error {
	if len(data) > 32 {
		return fmt.Errorf("Reading blocks larger than 32 bytes (%v) not supported", len(data))
	}

	if d.funcs&I2C_FUNC_SMBUS_READ_BLOCK_DATA == 0 {
		if i2cDeviceDebug {
			log.Printf("SMBus read block data not supported, use fallback\n")
		}
		return d.readBlockDataFallback(reg, data)
	}

	return d.smbusAccess(I2C_SMBUS_READ, reg, I2C_SMBUS_BLOCK_DATA, uintptr(unsafe.Pointer(&data)))
}

func (d *i2cDevice) WriteByte(val byte) error {
	if d.funcs&I2C_FUNC_SMBUS_WRITE_BYTE == 0 {
		return fmt.Errorf("SMBus write byte not supported")
	}

	return d.smbusAccess(I2C_SMBUS_WRITE, val, I2C_SMBUS_BYTE, uintptr(0))
}

func (d *i2cDevice) WriteByteData(reg uint8, val uint8) error {
	if d.funcs&I2C_FUNC_SMBUS_WRITE_BYTE_DATA == 0 {
		return fmt.Errorf("SMBus write byte data not supported")
	}

	var data = val
	return d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data)))
}

func (d *i2cDevice) WriteWordData(reg uint8, val uint16) error {
	if d.funcs&I2C_FUNC_SMBUS_WRITE_WORD_DATA == 0 {
		return fmt.Errorf("SMBus write word data not supported")
	}

	var data = val
	return d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&data)))
}

func (d *i2cDevice) WriteBlockData(reg uint8, data []byte) error {
	if len(data) > 32 {
		return fmt.Errorf("Writing blocks larger than 32 bytes (%v) not supported", len(data))
	}

	if d.funcs&I2C_FUNC_SMBUS_WRITE_BLOCK_DATA == 0 {
		if i2cDeviceDebug {
			log.Printf("SMBus write block data not supported, use fallback\n")
		}
		return d.writeBlockDataFallback(reg, data)
	}

	return d.smbusAccess(I2C_SMBUS_WRITE, reg, I2C_SMBUS_BLOCK_DATA, uintptr(unsafe.Pointer(&data)))
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

func (d *i2cDevice) readBlockDataFallback(reg uint8, data []byte) error {
	if err := d.writeAndCheckCount([]byte{reg}); err != nil {
		return err
	}
	if err := d.readAndCheckCount(data); err != nil {
		return err
	}
	return nil
}

func (d *i2cDevice) writeBlockDataFallback(reg uint8, data []byte) error {
	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = reg

	if err := d.writeAndCheckCount(buf); err != nil {
		return err
	}
	return nil
}

func (d *i2cDevice) readAndCheckCount(data []byte) error {
	n, err := d.file.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("Read %v bytes from device by sysfs, expected %v", n, len(data))
	}
	return nil
}

func (d *i2cDevice) writeAndCheckCount(data []byte) error {
	n, err := d.file.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("Write %v bytes to device by sysfs, expected %v", n, len(data))
	}
	return nil
}
