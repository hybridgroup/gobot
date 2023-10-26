package system

import (
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

const dev = "/dev/i2c-1"

func getSyscallFuncImpl(errorMask byte) func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno) {
	// bit 0: error on function query
	// bit 1: error on set address
	// bit 2: error on command
	return func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno) {
		// function query
		if (trap == Syscall_SYS_IOCTL) && (a2 == I2C_FUNCS) {
			if errorMask&0x01 == 0x01 {
				return 0, 0, 1
			}

			var funcPtr *uint64 = (*uint64)(unsafe.Pointer(a3))
			*funcPtr = I2C_FUNC_SMBUS_READ_BYTE | I2C_FUNC_SMBUS_READ_BYTE_DATA |
				I2C_FUNC_SMBUS_READ_WORD_DATA |
				I2C_FUNC_SMBUS_WRITE_BYTE | I2C_FUNC_SMBUS_WRITE_BYTE_DATA |
				I2C_FUNC_SMBUS_WRITE_WORD_DATA
		}
		// set address
		if (trap == Syscall_SYS_IOCTL) && (a2 == I2C_SLAVE) {
			if errorMask&0x02 == 0x02 {
				return 0, 0, 1
			}
		}
		// command
		if (trap == Syscall_SYS_IOCTL) && (a2 == I2C_SMBUS) {
			if errorMask&0x04 == 0x04 {
				return 0, 0, 1
			}
		}
		// Let all operations succeed
		return 0, 0, 0
	}
}

func initTestI2cDeviceWithMockedSys() (*i2cDevice, *mockSyscall) {
	a := NewAccesser()
	msc := a.UseMockSyscall()
	a.UseMockFilesystem([]string{dev})

	d, err := a.NewI2cDevice(dev)
	if err != nil {
		panic(err)
	}

	return d, msc
}

func TestNewI2cDevice(t *testing.T) {
	tests := map[string]struct {
		dev     string
		wantErr string
	}{
		"ok": {
			dev: dev,
		},
		"empty": {
			dev:     "",
			wantErr: "the given character device location is empty",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAccesser()
			// act
			d, err := a.NewI2cDevice(tc.dev)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Equal(t, (*i2cDevice)(nil), d)
			} else {
				var _ gobot.I2cSystemDevicer = d
				assert.NoError(t, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	// act & assert
	assert.NoError(t, d.Close())
}

func TestWriteRead(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	wbuf := []byte{0x01, 0x02, 0x03}
	rbuf := make([]byte, 4)
	// act
	wn, werr := d.Write(1, wbuf)
	rn, rerr := d.Read(1, rbuf)
	// assert
	assert.NoError(t, werr)
	assert.NoError(t, rerr)
	assert.Equal(t, len(wbuf), wn)
	assert.Equal(t, len(wbuf), rn) // will read only the written values
	assert.Equal(t, rbuf[:len(wbuf)], wbuf)
}

func TestReadByte(t *testing.T) {
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"read_byte_ok": {
			funcs: I2C_FUNC_SMBUS_READ_BYTE,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_BYTE,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 1, command: 0, protocol: 1, address: 2 failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus read byte not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const want = byte(5)
			msc.dataSlice = []byte{want}
			// act
			got, err := d.ReadByte(2)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, byte(I2C_SMBUS_READ), msc.smbus.readWrite)
				assert.Equal(t, byte(0), msc.smbus.command) // register is set to 0 in that case
				assert.Equal(t, uint32(I2C_SMBUS_BYTE), msc.smbus.protocol)
			}
		})
	}
}

func TestReadByteData(t *testing.T) {
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"read_byte_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_BYTE_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_BYTE_DATA,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 1, command: 1, protocol: 2, address: 3 failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus read byte data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const (
				reg  = byte(0x01)
				want = byte(0x02)
			)
			msc.dataSlice = []byte{want}
			// act
			got, err := d.ReadByteData(3, reg)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, byte(I2C_SMBUS_READ), msc.smbus.readWrite)
				assert.Equal(t, reg, msc.smbus.command)
				assert.Equal(t, uint32(I2C_SMBUS_BYTE_DATA), msc.smbus.protocol)
			}
		})
	}
}

func TestReadWordData(t *testing.T) {
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"read_word_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_WORD_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_WORD_DATA,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 1, command: 2, protocol: 3, address: 4 failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus read word data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const (
				reg    = byte(0x02)
				msbyte = byte(0xD4)
				lsbyte = byte(0x31)
				want   = uint16(54321)
			)
			// all common drivers read LSByte first
			msc.dataSlice = []byte{lsbyte, msbyte}
			// act
			got, err := d.ReadWordData(4, reg)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, byte(I2C_SMBUS_READ), msc.smbus.readWrite)
				assert.Equal(t, reg, msc.smbus.command)
				assert.Equal(t, uint32(I2C_SMBUS_WORD_DATA), msc.smbus.protocol)
			}
		})
	}
}

func TestReadBlockData(t *testing.T) {
	// arrange
	const (
		reg    = byte(0x03)
		wantB0 = byte(11)
		wantB1 = byte(22)
		wantB2 = byte(33)
		wantB3 = byte(44)
		wantB4 = byte(55)
		wantB5 = byte(66)
		wantB6 = byte(77)
		wantB7 = byte(88)
		wantB8 = byte(99)
		wantB9 = byte(111)
	)
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"read_block_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_I2C_BLOCK,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_I2C_BLOCK,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 1, command: 3, protocol: 8, address: 5 failed with syscall.Errno operation not permitted",
		},
		"error_from_used_fallback_if_not_supported": {
			wantErr: "Read 1 bytes from device by sysfs, expected 10",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			msc.dataSlice = []byte{wantB0, wantB1, wantB2, wantB3, wantB4, wantB5, wantB6, wantB7, wantB8, wantB9}
			buf := []byte{12, 23, 34, 45, 56, 67, 78, 89, 98, 87}
			// act
			err := d.ReadBlockData(5, reg, buf)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, msc.dataSlice, buf)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, uint8(len(buf)+1), msc.sliceSize)
				assert.Equal(t, byte(I2C_SMBUS_READ), msc.smbus.readWrite)
				assert.Equal(t, reg, msc.smbus.command)
				assert.Equal(t, uint32(I2C_SMBUS_I2C_BLOCK_DATA), msc.smbus.protocol)
			}
		})
	}
}

func TestWriteByte(t *testing.T) {
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"write_byte_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_BYTE,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_BYTE,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 0, command: 68, protocol: 1, address: 6 failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus write byte not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const val = byte(0x44)
			// act
			err := d.WriteByte(6, val)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, byte(I2C_SMBUS_WRITE), msc.smbus.readWrite)
				assert.Equal(t, val, msc.smbus.command) // in byte write, the register/command is used for the value
				assert.Equal(t, uint32(I2C_SMBUS_BYTE), msc.smbus.protocol)
			}
		})
	}
}

func TestWriteByteData(t *testing.T) {
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"write_byte_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_BYTE_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_BYTE_DATA,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 0, command: 4, protocol: 2, address: 7 failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus write byte data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const (
				reg = byte(0x04)
				val = byte(0x55)
			)
			// act
			err := d.WriteByteData(7, reg, val)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, byte(I2C_SMBUS_WRITE), msc.smbus.readWrite)
				assert.Equal(t, reg, msc.smbus.command)
				assert.Equal(t, uint32(I2C_SMBUS_BYTE_DATA), msc.smbus.protocol)
				assert.Equal(t, 1, len(msc.dataSlice))
				assert.Equal(t, val, msc.dataSlice[0])
			}
		})
	}
}

func TestWriteWordData(t *testing.T) {
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"write_word_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_WORD_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_WORD_DATA,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 0, command: 5, protocol: 3, address: 8 failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus write word data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const (
				reg        = byte(0x05)
				val        = uint16(54321)
				wantLSByte = byte(0x31)
				wantMSByte = byte(0xD4)
			)
			// act
			err := d.WriteWordData(8, reg, val)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, byte(I2C_SMBUS_WRITE), msc.smbus.readWrite)
				assert.Equal(t, reg, msc.smbus.command)
				assert.Equal(t, uint32(I2C_SMBUS_WORD_DATA), msc.smbus.protocol)
				assert.Equal(t, 2, len(msc.dataSlice))
				// all common drivers write LSByte first
				assert.Equal(t, wantLSByte, msc.dataSlice[0])
				assert.Equal(t, wantMSByte, msc.dataSlice[1])
			}
		})
	}
}

func TestWriteBlockData(t *testing.T) {
	// arrange
	const (
		reg = byte(0x06)
		b0  = byte(0x09)
		b1  = byte(0x11)
		b2  = byte(0x22)
		b3  = byte(0x33)
		b4  = byte(0x44)
		b5  = byte(0x55)
		b6  = byte(0x66)
		b7  = byte(0x77)
		b8  = byte(0x88)
		b9  = byte(0x99)
	)
	tests := map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
	}{
		"write_block_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_I2C_BLOCK,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_I2C_BLOCK,
			syscallImpl: getSyscallFuncImpl(0x04),
			wantErr:     "SMBus access r/w: 0, command: 6, protocol: 8, address: 9 failed with syscall.Errno operation not permitted",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			data := []byte{b0, b1, b2, b3, b4, b5, b6, b7, b8, b9}
			// act
			err := d.WriteBlockData(9, reg, data)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, d.file, msc.lastFile)
				assert.Equal(t, uintptr(I2C_SMBUS), msc.lastSignal)
				assert.Equal(t, uint8(len(data)+1), msc.sliceSize) // including size element
				assert.Equal(t, byte(I2C_SMBUS_WRITE), msc.smbus.readWrite)
				assert.Equal(t, reg, msc.smbus.command)
				assert.Equal(t, uint32(I2C_SMBUS_I2C_BLOCK_DATA), msc.smbus.protocol)
				assert.Equal(t, uint8(len(data)), msc.dataSlice[0]) // data size
				assert.Equal(t, data, msc.dataSlice[1:])
			}
		})
	}
}

func TestWriteBlockDataTooMuch(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	// act
	err := d.WriteBlockData(10, 0x01, make([]byte, 33))
	// assert
	assert.ErrorContains(t, err, "Writing blocks larger than 32 bytes (33) not supported")
}

func Test_setAddress(t *testing.T) {
	// arrange
	d, msc := initTestI2cDeviceWithMockedSys()
	// act
	err := d.setAddress(0xff)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, uintptr(0xff), msc.devAddress)
}

func Test_queryFunctionality(t *testing.T) {
	tests := map[string]struct {
		requested   uint64
		dev         string
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err SyscallErrno)
		wantErr     string
		wantFile    bool
		wantFuncs   uint64
	}{
		"ok": {
			requested:   I2C_FUNC_SMBUS_READ_BYTE,
			dev:         dev,
			syscallImpl: getSyscallFuncImpl(0x00),
			wantFile:    true,
			wantFuncs:   0x7E0000,
		},
		"dev_null_error": {
			dev:         os.DevNull,
			syscallImpl: getSyscallFuncImpl(0x00),
			wantErr:     " : /dev/null: no such file",
		},
		"query_funcs_error": {
			dev:         dev,
			syscallImpl: getSyscallFuncImpl(0x01),
			wantErr:     "Querying functionality failed with syscall.Errno operation not permitted",
			wantFile:    true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			d.location = tc.dev
			msc.Impl = tc.syscallImpl
			// act
			err := d.queryFunctionality(tc.requested, "test"+name)
			// assert
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			if tc.wantFile {
				assert.NotNil(t, d.file)
			} else {
				assert.Equal(t, (*MockFile)(nil), d.file)
			}
			assert.Equal(t, tc.wantFuncs, d.funcs)
		})
	}
}
