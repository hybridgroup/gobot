package system

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"unsafe"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

const dev = "/dev/i2c-1"

func syscallFuncsImpl(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	if (trap == syscall.SYS_IOCTL) && (a2 == I2C_FUNCS) {
		var funcPtr *uint64 = (*uint64)(unsafe.Pointer(a3))
		*funcPtr = I2C_FUNC_SMBUS_READ_BYTE | I2C_FUNC_SMBUS_READ_BYTE_DATA |
			I2C_FUNC_SMBUS_READ_WORD_DATA |
			I2C_FUNC_SMBUS_WRITE_BYTE | I2C_FUNC_SMBUS_WRITE_BYTE_DATA |
			I2C_FUNC_SMBUS_WRITE_WORD_DATA
	}
	// Let all operations succeed
	return 0, 0, 0
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
	var tests = map[string]struct {
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
				gobottest.Assert(t, err.Error(), tc.wantErr)
				gobottest.Assert(t, d, (*i2cDevice)(nil))
			} else {
				var _ gobot.I2cSystemDevicer = d
				gobottest.Assert(t, err, nil)
			}
		})
	}
}

func TestClose(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	// act & assert
	gobottest.Assert(t, d.Close(), nil)
}

func TestSetAddress(t *testing.T) {
	// arrange
	d, msc := initTestI2cDeviceWithMockedSys()
	// act
	err := d.SetAddress(0xff)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, msc.devAddress, uintptr(0xff))
}

func TestWriteRead(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	wbuf := []byte{0x01, 0x02, 0x03}
	rbuf := make([]byte, 4)
	// act
	wn, werr := d.Write(wbuf)
	rn, rerr := d.Read(rbuf)
	// assert
	gobottest.Assert(t, werr, nil)
	gobottest.Assert(t, rerr, nil)
	gobottest.Assert(t, wn, len(wbuf))
	gobottest.Assert(t, rn, len(wbuf)) // will read only the written values
	gobottest.Assert(t, wbuf, rbuf[:len(wbuf)])
}

func TestReadByte(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"read_byte_ok": {
			funcs: I2C_FUNC_SMBUS_READ_BYTE,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_BYTE,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			got, err := d.ReadByte()
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, got, want)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, byte(0)) // register is set to 0 in that case
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_BYTE))
			}
		})
	}
}

func TestReadByteData(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"read_byte_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_BYTE_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_BYTE_DATA,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			got, err := d.ReadByteData(reg)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, got, want)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_BYTE_DATA))
			}
		})
	}
}

func TestReadWordData(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"read_word_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_WORD_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_WORD_DATA,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			got, err := d.ReadWordData(reg)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, got, want)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_WORD_DATA))
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
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"read_block_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_I2C_BLOCK,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_I2C_BLOCK,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			err := d.ReadBlockData(reg, buf)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, buf, msc.dataSlice)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.sliceSize, uint8(len(buf)+1))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_I2C_BLOCK_DATA))
			}
		})
	}
}

func TestWriteByte(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"write_byte_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_BYTE,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_BYTE,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			err := d.WriteByte(val)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, val) // in byte write, the register/command is used for the value
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_BYTE))
			}
		})
	}
}

func TestWriteByteData(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"write_byte_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_BYTE_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_BYTE_DATA,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			err := d.WriteByteData(reg, val)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_BYTE_DATA))
				gobottest.Assert(t, len(msc.dataSlice), 1)
				gobottest.Assert(t, msc.dataSlice[0], val)
			}
		})
	}
}

func TestWriteWordData(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"write_word_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_WORD_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_WORD_DATA,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			err := d.WriteWordData(reg, val)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_WORD_DATA))
				gobottest.Assert(t, len(msc.dataSlice), 2)
				// all common drivers write LSByte first
				gobottest.Assert(t, msc.dataSlice[0], wantLSByte)
				gobottest.Assert(t, msc.dataSlice[1], wantMSByte)
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
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"write_block_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_I2C_BLOCK,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_I2C_BLOCK,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "SMBus access failed with syscall.Errno operation not permitted",
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
			err := d.WriteBlockData(reg, data)
			// assert
			if tc.wantErr != "" {
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.sliceSize, uint8(len(data)+1)) // including size element
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.protocol, uint32(I2C_SMBUS_I2C_BLOCK_DATA))
				gobottest.Assert(t, msc.dataSlice[0], uint8(len(data))) // data size
				gobottest.Assert(t, msc.dataSlice[1:], data)
			}
		})
	}
}

func TestWriteBlockDataTooMuch(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	// act
	err := d.WriteBlockData(0x01, make([]byte, 33))
	// assert
	gobottest.Assert(t, err, errors.New("Writing blocks larger than 32 bytes (33) not supported"))
}

func Test_lazyInit(t *testing.T) {
	var tests = map[string]struct {
		requested   uint64
		dev         string
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
		wantFile    bool
		wantFuncs   uint64
	}{
		"ok": {
			requested:   I2C_FUNC_SMBUS_READ_BYTE,
			dev:         dev,
			syscallImpl: syscallFuncsImpl,
			wantFile:    true,
			wantFuncs:   0x7E0000,
		},
		"dev_null_error": {
			dev:         os.DevNull,
			syscallImpl: syscallFuncsImpl,
			wantErr:     " : /dev/null: No such file.",
		},
		"query_funcs_error": {
			dev:         dev,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
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
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
			}
			if tc.wantFile {
				gobottest.Refute(t, d.file, nil)
			} else {
				gobottest.Assert(t, d.file, (*MockFile)(nil))
			}
			gobottest.Assert(t, d.funcs, tc.wantFuncs)
		})
	}
}
