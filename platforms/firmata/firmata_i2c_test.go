package firmata

import (
	"errors"
	"io"
	"testing"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/gobottest"
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
func (i2cMockFirmataBoard) I2cConfig(int) error { return nil }

// GPIO, PWM and servo functions unused in this test scenarios
func (i2cMockFirmataBoard) Connect(io.ReadWriteCloser) error { return nil }
func (i2cMockFirmataBoard) Disconnect() error                { return nil }
func (i2cMockFirmataBoard) Pins() []client.Pin               { return nil }
func (i2cMockFirmataBoard) AnalogWrite(int, int) error       { return nil }
func (i2cMockFirmataBoard) SetPinMode(int, int) error        { return nil }
func (i2cMockFirmataBoard) ReportAnalog(int, int) error      { return nil }
func (i2cMockFirmataBoard) ReportDigital(int, int) error     { return nil }
func (i2cMockFirmataBoard) DigitalWrite(int, int) error      { return nil }
func (i2cMockFirmataBoard) ServoConfig(int, int, int) error  { return nil }

// WriteSysex of the client implementation not tested here
func (i2cMockFirmataBoard) WriteSysex([]byte) error { return nil }

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
	gobottest.Assert(t, i2c.Close(), nil)
}

func TestRead(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{111}
	buf := []byte{0}
	// act
	countRead, err := con.Read(buf)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, countRead, 1)
	gobottest.Assert(t, brd.numBytesToRead, 1)
	gobottest.Assert(t, buf, brd.i2cDataForRead)
	gobottest.Assert(t, len(brd.i2cWritten), 0)
}

func TestReadByte(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{222}
	// act
	val, err := con.ReadByte()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, brd.numBytesToRead, 1)
	gobottest.Assert(t, val, brd.i2cDataForRead[0])
	gobottest.Assert(t, len(brd.i2cWritten), 0)
}

func TestReadByteData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	brd.i2cDataForRead = []byte{100}
	reg := uint8(0x01)
	// act
	val, err := con.ReadByteData(reg)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, brd.numBytesToRead, 1)
	gobottest.Assert(t, val, brd.i2cDataForRead[0])
	gobottest.Assert(t, len(brd.i2cWritten), 1)
	gobottest.Assert(t, brd.i2cWritten[0], reg)
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, brd.numBytesToRead, 2)
	gobottest.Assert(t, val, uint16(lsb)|uint16(msb)<<8)
	gobottest.Assert(t, len(brd.i2cWritten), 1)
	gobottest.Assert(t, brd.i2cWritten[0], reg)
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, brd.numBytesToRead, 5)
	gobottest.Assert(t, buf, brd.i2cDataForRead)
	gobottest.Assert(t, len(brd.i2cWritten), 1)
	gobottest.Assert(t, brd.i2cWritten[0], reg)
}

func TestWrite(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	want := []byte{0x00, 0x01}
	wantLen := len(want)
	// act
	written, err := con.Write(want)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, written, wantLen)
	gobottest.Assert(t, brd.i2cWritten, want)
}

func TestWrite20bytes(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	want := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	wantLen := len(want)
	// act
	written, err := con.Write(want)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, written, wantLen)
	gobottest.Assert(t, brd.i2cWritten, want)
}

func TestWriteByte(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	want := uint8(0x11)
	// act
	err := con.WriteByte(want)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(brd.i2cWritten), 1)
	gobottest.Assert(t, brd.i2cWritten[0], want)
}

func TestWriteByteData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	reg := uint8(0x12)
	val := uint8(0x22)
	// act
	err := con.WriteByteData(reg, val)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(brd.i2cWritten), 2)
	gobottest.Assert(t, brd.i2cWritten[0], reg)
	gobottest.Assert(t, brd.i2cWritten[1], val)
}

func TestWriteWordData(t *testing.T) {
	// arrange
	con, brd := initTestTestAdaptorWithI2cConnection()
	reg := uint8(0x13)
	val := uint16(0x8002)
	// act
	err := con.WriteWordData(reg, val)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(brd.i2cWritten), 3)
	gobottest.Assert(t, brd.i2cWritten[0], reg)
	gobottest.Assert(t, brd.i2cWritten[1], uint8(val&0x00FF))
	gobottest.Assert(t, brd.i2cWritten[2], uint8(val>>8))
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(brd.i2cWritten), 33)
	gobottest.Assert(t, brd.i2cWritten[0], reg)
	gobottest.Assert(t, brd.i2cWritten[1:], val[0:32])
}

func TestDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.DefaultI2cBus(), 0)
}

func TestGetI2cConnectionInvalidBus(t *testing.T) {
	a := NewAdaptor()
	_, err := a.GetI2cConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Invalid bus number 99, only 0 is supported"))
}
