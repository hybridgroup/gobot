package system

import (
	"fmt"
	"log"
	"os"
	"sync"
	"unsafe"
)

const (
	i2cDeviceDebug  = false
	forceSetAddress = false // normally address will be written only when changed, this behavior can be overridden
)

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
	I2C_FUNC_SMBUS_READ_I2C_BLOCK   = 0x04000000 // I2C-like block transfer with 1-byte reg. addr.
	I2C_FUNC_SMBUS_WRITE_I2C_BLOCK  = 0x08000000 // I2C-like block transfer with 1-byte reg. addr.
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
	protocol  uint32
	data      unsafe.Pointer
}

type i2cDevice struct {
	location    string
	sys         systemCaller
	fs          filesystem
	file        File
	funcs       uint64 // adapter functionality mask
	lastAddress int
	mutex       sync.Mutex
}

// NewI2cDevice returns a Linux Kernel access by ioctrl to the given i2c bus location (character device).
// Important note for "os.ModeExclusive": this is undefined without create the file for character devices, this means
// a second open will not return an error e.g. due to a busy resource. If this is not wanted, e.g. to minimize count of
// open fd's this needs to be prevented at caller side by implementing a caching mechanism. Furthermore this behavior
// can lead to problems with multiple devices on the same bus because the cycle SetAddress()...Read()/Write() etc. can
// be interrupted when using multiple instances for the same location.
func (a *Accesser) NewI2cDevice(location string) (*i2cDevice, error) {
	if location == "" {
		return nil, fmt.Errorf("the given character device location is empty")
	}

	d := &i2cDevice{
		location:    location,
		sys:         a.sys,
		fs:          a.fs,
		lastAddress: -1,
	}
	return d, nil
}

// Close closes the character device file and resets the lazy variables.
func (d *i2cDevice) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.funcs = 0
	d.lastAddress = -1
	if d.file != nil {
		return d.file.Close()
	}
	return nil
}

// ReadByte reads a byte from the current register of an i2c device.
func (d *i2cDevice) ReadByte(address int) (byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_READ_BYTE, "read byte"); err != nil {
		return 0, err
	}

	var data uint8 = 0xFC // set value for debugging purposes
	err := d.smbusAccess(address, I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, unsafe.Pointer(&data))
	return data, err
}

// ReadByteData reads a byte from the given register of an i2c device.
func (d *i2cDevice) ReadByteData(address int, reg uint8) (val uint8, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_READ_BYTE_DATA, "read byte data"); err != nil {
		return 0, err
	}

	var data uint8 = 0xFD // set value for debugging purposes
	err = d.smbusAccess(address, I2C_SMBUS_READ, reg, I2C_SMBUS_BYTE_DATA, unsafe.Pointer(&data))
	return data, err
}

// ReadWordData reads a 16 bit value starting from the given register of an i2c device.
func (d *i2cDevice) ReadWordData(address int, reg uint8) (val uint16, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_READ_WORD_DATA, "read word data"); err != nil {
		return 0, err
	}

	var data uint16 = 0xFFFE // set value for debugging purposes
	err = d.smbusAccess(address, I2C_SMBUS_READ, reg, I2C_SMBUS_WORD_DATA, unsafe.Pointer(&data))
	return data, err
}

// ReadBlockData fills the given buffer with reads starting from the given register of an i2c device.
func (d *i2cDevice) ReadBlockData(address int, reg uint8, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	dataLen := len(data)
	if dataLen > 32 {
		return fmt.Errorf("Reading blocks larger than 32 bytes (%v) not supported", len(data))
	}

	data[0] = 0xFF // set value for debugging purposes
	if err := d.queryFunctionality(I2C_FUNC_SMBUS_READ_I2C_BLOCK, "read block data"); err != nil {
		if i2cDeviceDebug {
			log.Printf("%s, use fallback\n", err.Error())
		}
		return d.readBlockDataFallback(address, reg, data)
	}

	// set the first element with the data size
	buf := make([]byte, dataLen+1)
	buf[0] = byte(dataLen)
	copy(buf[1:], data)
	if err := d.smbusAccess(address, I2C_SMBUS_READ, reg, I2C_SMBUS_I2C_BLOCK_DATA, unsafe.Pointer(&buf[0])); err != nil {
		return err
	}
	// get data from buffer without first size element
	copy(data, buf[1:])
	return nil
}

// WriteByte writes the given byte value to the current register of an i2c device.
func (d *i2cDevice) WriteByte(address int, val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_WRITE_BYTE, "write byte"); err != nil {
		return err
	}

	return d.smbusAccess(address, I2C_SMBUS_WRITE, val, I2C_SMBUS_BYTE, nil)
}

// WriteByteData writes the given byte value to the given register of an i2c device.
func (d *i2cDevice) WriteByteData(address int, reg uint8, val uint8) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_WRITE_BYTE_DATA, "write byte data"); err != nil {
		return err
	}

	data := val
	return d.smbusAccess(address, I2C_SMBUS_WRITE, reg, I2C_SMBUS_BYTE_DATA, unsafe.Pointer(&data))
}

// WriteWordData writes the given 16 bit value starting from the given register of an i2c device.
func (d *i2cDevice) WriteWordData(address int, reg uint8, val uint16) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_WRITE_WORD_DATA, "write word data"); err != nil {
		return err
	}

	data := val
	return d.smbusAccess(address, I2C_SMBUS_WRITE, reg, I2C_SMBUS_WORD_DATA, unsafe.Pointer(&data))
}

// WriteBlockData writes the given buffer starting from the given register of an i2c device.
func (d *i2cDevice) WriteBlockData(address int, reg uint8, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	dataLen := len(data)
	if dataLen > 32 {
		return fmt.Errorf("Writing blocks larger than 32 bytes (%v) not supported", len(data))
	}

	if err := d.queryFunctionality(I2C_FUNC_SMBUS_WRITE_I2C_BLOCK, "write i2c block"); err != nil {
		if i2cDeviceDebug {
			log.Printf("%s, use fallback\n", err.Error())
		}
		return d.writeBlockDataFallback(address, reg, data)
	}

	// set the first element with the data size
	buf := make([]byte, dataLen+1)
	buf[0] = byte(dataLen)
	copy(buf[1:], data)

	return d.smbusAccess(address, I2C_SMBUS_WRITE, reg, I2C_SMBUS_I2C_BLOCK_DATA, unsafe.Pointer(&buf[0]))
}

// WriteBytes writes the given buffer starting from the current register of an i2c device.
func (d *i2cDevice) WriteBytes(address int, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.writeBytes(address, data)
}

// Read implements direct I2C read operations.
func (d *i2cDevice) Read(address int, b []byte) (n int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.read(address, b)
}

// Write implements the io.ReadWriteCloser method by direct I2C write operations.
func (d *i2cDevice) Write(address int, b []byte) (n int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.write(address, b)
}

func (d *i2cDevice) readBlockDataFallback(address int, reg uint8, data []byte) error {
	if err := d.writeBytes(address, []byte{reg}); err != nil {
		return err
	}
	if err := d.readAndCheckCount(address, data); err != nil {
		return err
	}
	return nil
}

func (d *i2cDevice) writeBlockDataFallback(address int, reg uint8, data []byte) error {
	buf := make([]byte, len(data)+1)
	copy(buf[1:], data)
	buf[0] = reg

	if err := d.writeBytes(address, buf); err != nil {
		return err
	}
	return nil
}

func (d *i2cDevice) writeBytes(address int, data []byte) error {
	n, err := d.write(address, data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("Write %v bytes to device by sysfs, expected %v", n, len(data))
	}
	return nil
}

func (d *i2cDevice) write(address int, b []byte) (n int, err error) {
	if err = d.setAddress(address); err != nil {
		return 0, err
	}
	if err := d.openFileLazy("Write"); err != nil {
		return 0, err
	}
	return d.file.Write(b)
}

func (d *i2cDevice) readAndCheckCount(address int, data []byte) error {
	n, err := d.read(address, data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("Read %v bytes from device by sysfs, expected %v", n, len(data))
	}
	return nil
}

func (d *i2cDevice) read(address int, b []byte) (n int, err error) {
	if err = d.setAddress(address); err != nil {
		return 0, err
	}
	if err := d.openFileLazy("Read"); err != nil {
		return 0, err
	}

	return d.file.Read(b)
}

func (d *i2cDevice) queryFunctionality(requested uint64, sender string) error {
	// lazy initialization
	if d.funcs == 0 {
		if err := d.syscallIoctl(I2C_FUNCS, unsafe.Pointer(&d.funcs), "Querying functionality"); err != nil {
			return err
		}
	}

	if d.funcs&requested == 0 {
		return fmt.Errorf("SMBus %s not supported", sender)
	}

	return nil
}

func (d *i2cDevice) smbusAccess(address int, readWrite byte, command byte, protocol uint32, dataStart unsafe.Pointer) error {
	if err := d.setAddress(address); err != nil {
		return err
	}

	smbus := i2cSmbusIoctlData{
		readWrite: readWrite,
		command:   command,
		protocol:  protocol,
		data:      dataStart, // the reflected value of unsafePointer equals uintptr(dataStart),
	}

	sender := fmt.Sprintf("SMBus access r/w: %d, command: %d, protocol: %d, address: %d", readWrite, command, protocol, d.lastAddress)
	if err := d.syscallIoctl(I2C_SMBUS, unsafe.Pointer(&smbus), sender); err != nil {
		return err
	}

	return nil
}

// setAddress sets the address of the i2c device to use.
func (d *i2cDevice) setAddress(address int) error {
	if d.lastAddress == address && !forceSetAddress {
		if i2cDeviceDebug {
			log.Printf("I2C address %d was already sent - skip", address)
		}
		return nil
	}
	// for go vet false positives, see: https://github.com/golang/go/issues/41205
	if err := d.syscallIoctl(I2C_SLAVE, unsafe.Pointer(uintptr(byte(address))), "Setting address"); err != nil {
		return err
	}
	d.lastAddress = address
	return nil
}

func (d *i2cDevice) syscallIoctl(signal uintptr, payload unsafe.Pointer, sender string) (err error) {
	if err := d.openFileLazy(sender); err != nil {
		return err
	}
	if _, _, errno := d.sys.syscall(Syscall_SYS_IOCTL, d.file, signal, payload); errno != 0 {
		return fmt.Errorf("%s failed with syscall.Errno %v", sender, errno)
	}

	return nil
}

func (d *i2cDevice) openFileLazy(sender string) (err error) { //nolint:unparam // useful for debugging
	// lazy initialization
	// note: "os.ModeExclusive" is undefined without create the file. This means for the existing character device,
	// a second open will not return an error e.g. due to a busy resource, so most likely "os.ModeExclusive" is not really
	// helpful and we drop it to the default "0" used by normal Open().
	if d.file == nil {
		if d.file, err = d.fs.openFile(d.location, os.O_RDWR, 0); err != nil {
			return err
		}
	}
	return nil
}
