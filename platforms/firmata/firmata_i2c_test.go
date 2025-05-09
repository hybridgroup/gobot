//go:build !windows
// +build !windows

//nolint:forcetypeassert,gosec // ok here
package firmata

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/firmata/client"
)

// make sure that this Adaptor fulfills all required I2C interfaces
var _ i2c.Connector = (*Adaptor)(nil)

type i2cMockFirmataBoard struct {
	gobot.Eventer
	i2cDataForRead []byte
	numBytesToRead int
	i2cWritten     []byte
}

// setup mock for i2c tests
func (t *i2cMockFirmataBoard) I2cRead(address int, numBytes int) error {
	t.numBytesToRead = numBytes
	i2cReply := client.I2cReply{Data: t.i2cDataForRead}
	go func() {
		<-time.After(10 * time.Millisecond)
		t.Publish(t.Event("I2cReply"), i2cReply)
	}()
	return nil
}

func (t *i2cMockFirmataBoard) I2cWrite(address int, data []byte) error {
	t.i2cWritten = append(t.i2cWritten, data...)
	return nil
}
func (*i2cMockFirmataBoard) I2cConfig(int) error { return nil }

// GPIO, PWM and servo functions unused in this test scenarios
func (*i2cMockFirmataBoard) Connect(io.ReadWriteCloser) error { return nil }
func (*i2cMockFirmataBoard) Disconnect() error                { return nil }
func (*i2cMockFirmataBoard) Pins() []client.Pin               { return nil }
func (*i2cMockFirmataBoard) AnalogWrite(int, int) error       { return nil }
func (*i2cMockFirmataBoard) SetPinMode(int, int) error        { return nil }
func (*i2cMockFirmataBoard) ReportAnalog(int, int) error      { return nil }
func (*i2cMockFirmataBoard) ReportDigital(int, int) error     { return nil }
func (*i2cMockFirmataBoard) DigitalWrite(int, int) error      { return nil }
func (*i2cMockFirmataBoard) ServoConfig(int, int, int) error  { return nil }

// WriteSysex of the client implementation not tested here
func (*i2cMockFirmataBoard) WriteSysex([]byte) error { return nil }

func newI2cMockFirmataBoard() *i2cMockFirmataBoard {
	m := &i2cMockFirmataBoard{
		Eventer: gobot.NewEventer(),
	}
	m.AddEvent("I2cReply")
	return m
}

func initTestTestAdaptorWithI2cConnection() (i2c.Connection, *i2cMockFirmataBoard) {
	a := NewAdaptor()
	a.Board = newI2cMockFirmataBoard()
	con, err := a.GetI2cConnection(0, 0)
	if err != nil {
		panic(err)
	}
	return con, a.Board.(*i2cMockFirmataBoard)
}

func TestClose(t *testing.T) {
	i2c, _ := initTestTestAdaptorWithI2cConnection()
	require.NoError(t, i2c.Close())
}

func TestRead(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{111}
	buf := []byte{0}
	// act
	countRead, err := con.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 1, countRead)
	assert.Equal(t, 1, brd.numBytesToRead)
	assert.Equal(t, brd.i2cDataForRead, buf)
	assert.Empty(t, brd.i2cWritten)
}

func TestReadByte(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{222}
	// act
	val, err := con.ReadByte()
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, brd.numBytesToRead)
	assert.Equal(t, brd.i2cDataForRead[0], val)
	assert.Empty(t, brd.i2cWritten)
}

func TestReadByteData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{100}
	reg := uint8(0x01)
	// act
	val, err := con.ReadByteData(reg)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, brd.numBytesToRead)
	assert.Equal(t, brd.i2cDataForRead[0], val)
	assert.Len(t, brd.i2cWritten, 1)
	assert.Equal(t, reg, brd.i2cWritten[0])
}

func TestReadWordData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	lsb := uint8(0x11)
	msb := uint8(0xff)
	brd.i2cDataForRead = []byte{lsb, msb}
	reg := uint8(0x22)
	// act
	val, err := con.ReadWordData(reg)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 2, brd.numBytesToRead)
	assert.Equal(t, uint16(lsb)|uint16(msb)<<8, val)
	assert.Len(t, brd.i2cWritten, 1)
	assert.Equal(t, reg, brd.i2cWritten[0])
}

func TestReadBlockData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{50, 40, 30, 20, 10}
	reg := uint8(0x33)
	buf := []byte{1, 2, 3, 4, 5}
	// act
	err := con.ReadBlockData(reg, buf)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 5, brd.numBytesToRead)
	assert.Equal(t, brd.i2cDataForRead, buf)
	assert.Len(t, brd.i2cWritten, 1)
	assert.Equal(t, reg, brd.i2cWritten[0])
}

func TestWrite(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	want := []byte{0x00, 0x01}
	wantLen := len(want)
	// act
	written, err := con.Write(want)
	// assert
	require.NoError(t, err)
	assert.Equal(t, wantLen, written)
	assert.Equal(t, want, brd.i2cWritten)
}

func TestWrite20bytes(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	want := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	wantLen := len(want)
	// act
	written, err := con.Write(want)
	// assert
	require.NoError(t, err)
	assert.Equal(t, wantLen, written)
	assert.Equal(t, want, brd.i2cWritten)
}

func TestWriteByte(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	want := uint8(0x11)
	// act
	err := con.WriteByte(want)
	// assert
	require.NoError(t, err)
	assert.Len(t, brd.i2cWritten, 1)
	assert.Equal(t, want, brd.i2cWritten[0])
}

func TestWriteByteData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	reg := uint8(0x12)
	val := uint8(0x22)
	// act
	err := con.WriteByteData(reg, val)
	// assert
	require.NoError(t, err)
	assert.Len(t, brd.i2cWritten, 2)
	assert.Equal(t, reg, brd.i2cWritten[0])
	assert.Equal(t, val, brd.i2cWritten[1])
}

func TestWriteWordData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	reg := uint8(0x13)
	val := uint16(0x8002)
	// act
	err := con.WriteWordData(reg, val)
	// assert
	require.NoError(t, err)
	assert.Len(t, brd.i2cWritten, 3)
	assert.Equal(t, reg, brd.i2cWritten[0])
	assert.Equal(t, uint8(val&0x00FF), brd.i2cWritten[1])
	assert.Equal(t, uint8(val>>8), brd.i2cWritten[2])
}

func TestWriteBlockData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	reg := uint8(0x14)
	val := []byte{}
	// we prepare more than 32 bytes, because the call has to drop it
	for i := uint8(0); i < 40; i++ {
		val = append(val, i)
	}
	// act
	err := con.WriteBlockData(reg, val)
	// assert
	require.NoError(t, err)
	assert.Len(t, brd.i2cWritten, 33)
	assert.Equal(t, reg, brd.i2cWritten[0])
	assert.Equal(t, val[0:32], brd.i2cWritten[1:])
}

func TestDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 0, a.DefaultI2cBus())
}

func TestGetI2cConnectionInvalidBus(t *testing.T) {
	a := NewAdaptor()
	_, err := a.GetI2cConnection(0x01, 99)
	require.ErrorContains(t, err, "Invalid bus number 99, only 0 is supported")
}
