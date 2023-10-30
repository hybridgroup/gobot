package spi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
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
	assert.NoError(t, err)
	assert.Equal(t, command, sysdev.Written())
	assert.Equal(t, want, got)
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
	assert.NoError(t, err)
	assert.Equal(t, []byte{reg, 0x00}, sysdev.Written()) // for read register we need n+1 bytes
	assert.Equal(t, want, got)
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
	assert.NoError(t, err)
	assert.Equal(t, []byte{reg, 0x00, 0x00, 0x00, 0x00}, sysdev.Written()) // for read registers we need n+1 bytes
	assert.Equal(t, want, got)
}

func TestWriteByte(t *testing.T) {
	// arrange
	const want = 0x02
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteByte(want)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, []byte{want}, sysdev.Written())
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
	assert.NoError(t, err)
	assert.Equal(t, []byte{reg, val}, sysdev.Written())
}

func TestWriteBlockData(t *testing.T) {
	// arrange
	const reg = 0x33
	data := []byte{0x22, 0x11}
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteBlockData(reg, data)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, append([]byte{reg}, data...), sysdev.Written())
}

func TestWriteBytes(t *testing.T) {
	// arrange
	want := []byte{0x03}
	c, sysdev := initTestConnectionWithMockedSystem()
	// act
	err := c.WriteBytes(want)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, want, sysdev.Written())
}
