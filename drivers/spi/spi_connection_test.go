package spi

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

var _ gobot.SpiOperations = (*spiConnection)(nil)

func initTestConnectionWithMockedSystem() (Connection, *system.MockSpiAccess) {
	a := system.NewAccesser()
	sysdev := a.UseMockSpi()
	const (
		busNum   = 15
		chipNum  = 14
		mode     = 13
		bits     = 12
		maxSpeed = int64(11)
	)
	d, err := a.NewSpiDevice(busNum, chipNum, mode, bits, maxSpeed)
	if err != nil {
		panic(err)
	}
	c := NewConnection(d)
	return c, sysdev
}

func TestReadCommandData(t *testing.T) {
	// arrange
	command := []byte{0x11, 0x12}
	want := []byte{0x31, 0x32}
	c, sysdev := initTestConnectionWithMockedSystem()
	sysdev.SetSimRead(want)
	// act
	got := []byte{0x01, 0x02}
	err := c.ReadCommandData(command, got)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), command)
	gobottest.Assert(t, got, want)
}

func TestReadByteData(t *testing.T) {
	// arrange
	const (
		reg  = 0x15
		want = uint8(0x41)
	)
	c, sysdev := initTestConnectionWithMockedSystem()
	sysdev.SetSimRead([]byte{0x00, want}) // the answer is one cycle behind
	// act
	got, err := c.ReadByteData(reg)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), []byte{reg, 0x00}) // for read register we need n+1 bytes
	gobottest.Assert(t, got, want)
}

func TestReadBlockData(t *testing.T) {
	// arrange
	const (
		reg = 0x16
	)
	want := []byte{42, 24, 56, 65}
	c, sysdev := initTestConnectionWithMockedSystem()
	sysdev.SetSimRead(append([]byte{0x00}, want...)) // the answer is one cycle behind
	// act
	got := make([]byte, 4)
	err := c.ReadBlockData(reg, got)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), []byte{reg, 0x00, 0x00, 0x00, 0x00}) // for read registers we need n+1 bytes
	gobottest.Assert(t, got, want)
}

func TestWriteByte(t *testing.T) {
	// arrange
	const want = 0x02
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteByte(want)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), []byte{want})
}

func TestWriteByteData(t *testing.T) {
	// arrange
	const (
		reg = 0x22
		val = 0x33
	)
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteByteData(reg, val)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), []byte{reg, val})
}

func TestWriteBlockData(t *testing.T) {
	// arrange
	const reg = 0x33
	data := []byte{0x22, 0x11}
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteBlockData(reg, data)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), append([]byte{reg}, data...))
}

func TestWriteBytes(t *testing.T) {
	// arrange
	want := []byte{0x03}
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteBytes(want)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sysdev.Written(), want)
}
