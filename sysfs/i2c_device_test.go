package sysfs

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
)

const dev = "/dev/i2c-1"

func initTestI2cDeviceWithMockedSys() (*i2cDevice, *MockSyscall) {
	SetFilesystem(NewMockFilesystem([]string{dev}))
	msc := &MockSyscall{}
	SetSyscall(msc)
	d, err := NewI2cDevice(dev)
	if err != nil {
		panic(err)
	}
	return d, msc
}

func cleanTestI2cDevice() {
	defer SetFilesystem(&NativeFilesystem{})
	defer SetSyscall(&NativeSyscall{})
}

func TestNewI2cDevice(t *testing.T) {
	var tests = map[string]struct {
		dev     string
		wantErr string
	}{
		"ok": {
			dev: dev,
		},
		"null": {
			dev:     os.DevNull,
			wantErr: " : /dev/null: No such file.",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			SetFilesystem(NewMockFilesystem([]string{dev}))
			defer SetFilesystem(&NativeFilesystem{})
			SetSyscall(&MockSyscall{})
			defer SetSyscall(&NativeSyscall{})
			// act
			i, err := NewI2cDevice(tc.dev)
			var _ i2c.I2cDevice = i
			if tc.wantErr != "" {
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
			}
		})
	}
}

func TestNewI2cDeviceQueryFuncError(t *testing.T) {
	// arrange
	SetFilesystem(NewMockFilesystem([]string{dev}))
	defer SetFilesystem(&NativeFilesystem{})
	SetSyscall(&MockSyscall{Impl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 }})
	defer SetSyscall(&NativeSyscall{})
	// act
	_, err := NewI2cDevice(dev)
	// assert
	gobottest.Assert(t, err, errors.New("Querying functionality failed with syscall.Errno operation not permitted"))
}

func TestClose(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	defer cleanTestI2cDevice()
	// act & assert
	gobottest.Assert(t, d.Close(), nil)
}

func TestSetAddress(t *testing.T) {
	// arrange
	d, msc := initTestI2cDeviceWithMockedSys()
	defer cleanTestI2cDevice()
	// act
	err := d.SetAddress(0xff)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, msc.devAddress, uintptr(0xff))
}

func TestWriteRead(t *testing.T) {
	// arrange
	d, _ := initTestI2cDeviceWithMockedSys()
	defer cleanTestI2cDevice()
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
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus read byte not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const want = byte(5)
			msc.dataSlice = []byte{want}
			// act
			got, err := d.ReadByte()
			// assert
			if tc.wantErr != "" {
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, got, want)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, byte(0)) // register is set to 0 in that case
				gobottest.Assert(t, msc.smbus.size, uint32(I2C_SMBUS_BYTE))
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
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus read byte data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
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
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, got, want)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.size, uint32(I2C_SMBUS_BYTE_DATA))
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
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus read word data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
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
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, got, want)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.size, uint32(I2C_SMBUS_WORD_DATA))
			}
		})
	}
}

func TestReadBlockData(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"read_block_data_ok": {
			funcs: I2C_FUNC_SMBUS_READ_BLOCK_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_READ_BLOCK_DATA,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_from_used_fallback_if_not_supported": {
			wantErr: "Read 1 bytes from device by sysfs, expected 3",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const (
				reg       = byte(0x03)
				wantByte0 = byte(0x11)
				wantByte1 = byte(0x22)
				wantByte2 = byte(0x33)
			)
			msc.dataSlice = []byte{wantByte0, wantByte1, wantByte2}
			wantSize := uint32(len(msc.dataSlice) + 1) // register is also part of send data
			buf := []byte{17, 28, 39}
			// act
			err := d.ReadBlockData(reg, buf)
			// assert
			if tc.wantErr != "" {
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, buf, msc.dataSlice)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_READ))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.size, wantSize)
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
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus write byte not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const val = byte(0x44)
			// act
			err := d.WriteByte(val)
			// assert
			if tc.wantErr != "" {
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, val) // in byte write, the register/command is used for the value
				gobottest.Assert(t, msc.smbus.size, uint32(I2C_SMBUS_BYTE))
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
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus write byte data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
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
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.size, uint32(I2C_SMBUS_BYTE_DATA))
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
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
		"error_not_supported": {
			wantErr: "SMBus write word data not supported",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
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
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.size, uint32(I2C_SMBUS_WORD_DATA))
				gobottest.Assert(t, len(msc.dataSlice), 2)
				// all common drivers write LSByte first
				gobottest.Assert(t, msc.dataSlice[0], wantLSByte)
				gobottest.Assert(t, msc.dataSlice[1], wantMSByte)
			}
		})
	}
}

func TestWriteBlockData(t *testing.T) {
	var tests = map[string]struct {
		funcs       uint64
		syscallImpl func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
		wantErr     string
	}{
		"write_word_data_ok": {
			funcs: I2C_FUNC_SMBUS_WRITE_BLOCK_DATA,
		},
		"error_syscall": {
			funcs:       I2C_FUNC_SMBUS_WRITE_BLOCK_DATA,
			syscallImpl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) { return 0, 0, 1 },
			wantErr:     "Failed with syscall.Errno operation not permitted",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, msc := initTestI2cDeviceWithMockedSys()
			defer cleanTestI2cDevice()
			msc.Impl = tc.syscallImpl
			d.funcs = tc.funcs
			const (
				reg   = byte(0x06)
				byte0 = byte(0x66)
				byte1 = byte(0x77)
				byte2 = byte(0x88)
			)
			data := []byte{byte0, byte1, byte2}
			wantSize := uint32(len(data) + 1) // register is also part of send data
			// act
			err := d.WriteBlockData(reg, data)
			// assert
			if tc.wantErr != "" {
				gobottest.Refute(t, err, nil)
				gobottest.Assert(t, err.Error(), tc.wantErr)
			} else {
				gobottest.Assert(t, err, nil)
				gobottest.Assert(t, msc.lastFile, d.file)
				gobottest.Assert(t, msc.lastSignal, uintptr(I2C_SMBUS))
				gobottest.Assert(t, msc.smbus.readWrite, byte(I2C_SMBUS_WRITE))
				gobottest.Assert(t, msc.smbus.command, reg)
				gobottest.Assert(t, msc.smbus.size, wantSize)
				gobottest.Assert(t, msc.dataSlice, data)
			}
		})
	}
}

func TestWriteBlockDataTooMuch(t *testing.T) {
	// arrange
	SetFilesystem(NewMockFilesystem([]string{dev}))
	defer SetFilesystem(&NativeFilesystem{})
	d, _ := NewI2cDevice(dev)
	// act
	err := d.WriteBlockData(0x01, make([]byte, 33))
	// assert
	gobottest.Assert(t, err, errors.New("Writing blocks larger than 32 bytes (33) not supported"))
}
