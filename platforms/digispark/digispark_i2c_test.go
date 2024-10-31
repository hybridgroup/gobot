//go:build libusb
// +build libusb

//nolint:forcetypeassert // ok here
package digispark

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/drivers/i2c"
)

const (
	availableI2cAddress = 0x40
	maxUint8            = ^uint8(0)
)

var (
	_       i2c.Connector = (*Adaptor)(nil)
	i2cData               = []byte{5, 4, 3, 2, 1, 0}
)

type i2cMock struct {
	duration          uint
	direction         uint8
	dataWritten       []byte
	writeStartWasSend bool
	writeStopWasSend  bool
	readStartWasSend  bool
	readStopWasSend   bool
}

func initTestAdaptorI2c() *Adaptor {
	a := NewAdaptor()
	a.connect = func(a *Adaptor) error { return nil }
	a.littleWire = new(i2cMock)
	return a
}

func TestDigisparkAdaptorI2cGetI2cConnection(t *testing.T) {
	// arrange
	var c i2c.Connection
	var err error
	a := initTestAdaptorI2c()

	// act
	c, err = a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// assert
	require.NoError(t, err)
	assert.NotNil(t, c)
}

func TestDigisparkAdaptorI2cGetI2cConnectionFailWithInvalidBus(t *testing.T) {
	// arrange
	a := initTestAdaptorI2c()

	// act
	c, err := a.GetI2cConnection(0x40, 1)

	// assert
	require.ErrorContains(t, err, "Invalid bus number 1, only 0 is supported")
	assert.Nil(t, c)
}

func TestDigisparkAdaptorI2cStartFailWithWrongAddress(t *testing.T) {
	// arrange
	data := []byte{0, 1, 2, 3, 4}
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(0x39, a.DefaultI2cBus())

	// act
	count, err := c.Write(data)

	// assert
	assert.Equal(t, 0, count)
	require.ErrorContains(t, err, fmt.Sprintf("Invalid address, only %d is supported", availableI2cAddress))
	assert.Equal(t, maxUint8, a.littleWire.(*i2cMock).direction)
}

func TestDigisparkAdaptorI2cWrite(t *testing.T) {
	// arrange
	data := []byte{0, 1, 2, 3, 4}
	dataLen := len(data)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	count, err := c.Write(data)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(0), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, dataLen, count)
	assert.Equal(t, data, a.littleWire.(*i2cMock).dataWritten)
	assert.True(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cWriteByte(t *testing.T) {
	// arrange
	data := byte(0x02)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	err := c.WriteByte(data)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(0), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, []byte{data}, a.littleWire.(*i2cMock).dataWritten)
	assert.True(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cWriteByteData(t *testing.T) {
	// arrange
	reg := uint8(0x03)
	data := byte(0x09)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	err := c.WriteByteData(reg, data)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(0), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, []byte{reg, data}, a.littleWire.(*i2cMock).dataWritten)
	assert.True(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cWriteWordData(t *testing.T) {
	// arrange
	reg := uint8(0x04)
	data := uint16(0x0508)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	err := c.WriteWordData(reg, data)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(0), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, []byte{reg, 0x08, 0x05}, a.littleWire.(*i2cMock).dataWritten)
	assert.True(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cWriteBlockData(t *testing.T) {
	// arrange
	reg := uint8(0x05)
	data := []byte{0x80, 0x81, 0x82}
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	err := c.WriteBlockData(reg, data)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(0), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, append([]byte{reg}, data...), a.littleWire.(*i2cMock).dataWritten)
	assert.True(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.False(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cRead(t *testing.T) {
	// arrange
	data := []byte{0, 1, 2, 3, 4}
	dataLen := len(data)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	count, err := c.Read(data)

	// assert
	assert.Equal(t, dataLen, count)
	require.NoError(t, err)
	assert.False(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(1), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, i2cData[:dataLen], data)
	assert.False(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cReadByte(t *testing.T) {
	// arrange
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	data, err := c.ReadByte()

	// assert
	require.NoError(t, err)
	assert.False(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(1), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, i2cData[0], data)
	assert.False(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cReadByteData(t *testing.T) {
	// arrange
	reg := uint8(0x04)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	data, err := c.ReadByteData(reg)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(1), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, []byte{reg}, a.littleWire.(*i2cMock).dataWritten)
	assert.Equal(t, i2cData[0], data)
	assert.False(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cReadWordData(t *testing.T) {
	// arrange
	reg := uint8(0x05)
	// 2 bytes of i2cData are used swapped
	expectedValue := uint16(0x0405)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	data, err := c.ReadWordData(reg)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(1), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, []byte{reg}, a.littleWire.(*i2cMock).dataWritten)
	assert.Equal(t, expectedValue, data)
	assert.False(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cReadBlockData(t *testing.T) {
	// arrange
	reg := uint8(0x05)
	data := []byte{0, 0, 0, 0, 0}
	dataLen := len(data)
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	err := c.ReadBlockData(reg, data)

	// assert
	require.NoError(t, err)
	assert.True(t, a.littleWire.(*i2cMock).writeStartWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStartWasSend)
	assert.Equal(t, uint8(1), a.littleWire.(*i2cMock).direction)
	assert.Equal(t, []byte{reg}, a.littleWire.(*i2cMock).dataWritten)
	assert.Equal(t, i2cData[:dataLen], data)
	assert.False(t, a.littleWire.(*i2cMock).writeStopWasSend)
	assert.True(t, a.littleWire.(*i2cMock).readStopWasSend)
}

func TestDigisparkAdaptorI2cUpdateDelay(t *testing.T) {
	// arrange
	a := initTestAdaptorI2c()
	c, _ := a.GetI2cConnection(availableI2cAddress, a.DefaultI2cBus())

	// act
	err := c.(*digisparkI2cConnection).UpdateDelay(uint(100))

	// assert
	require.NoError(t, err)
	assert.Equal(t, uint(100), a.littleWire.(*i2cMock).duration)
}

// setup mock for i2c tests
func (l *i2cMock) i2cInit() error {
	l.direction = maxUint8
	return l.error()
}

func (l *i2cMock) i2cStart(address7bit uint8, direction uint8) error {
	if address7bit != availableI2cAddress {
		return fmt.Errorf("Invalid address, only %d is supported", availableI2cAddress)
	}
	if err := l.error(); err != nil {
		return err
	}
	l.direction = direction
	if direction == 1 {
		l.readStartWasSend = true
	} else {
		l.writeStartWasSend = true
	}
	return nil
}

func (l *i2cMock) i2cWrite(sendBuffer []byte, length int, endWithStop uint8) error {
	l.dataWritten = append(l.dataWritten, sendBuffer...)
	if endWithStop > 0 {
		l.writeStopWasSend = true
	}
	return l.error()
}

func (l *i2cMock) i2cRead(readBuffer []byte, length int, endWithStop uint8) error {
	if len(readBuffer) < length {
		length = len(readBuffer)
	}
	if len(i2cData) < length {
		length = len(i2cData)
	}
	copy(readBuffer[:length], i2cData[:length])

	if endWithStop > 0 {
		l.readStopWasSend = true
	}
	return l.error()
}

func (l *i2cMock) i2cUpdateDelay(duration uint) error {
	l.duration = duration
	return l.error()
}

// GPIO, PWM and servo functions unused in this test scenarios
func (l *i2cMock) digitalWrite(pin uint8, state uint8) error                  { return nil }
func (l *i2cMock) pinMode(pin uint8, mode uint8) error                        { return nil }
func (l *i2cMock) pwmInit() error                                             { return nil }
func (l *i2cMock) pwmStop() error                                             { return nil }
func (l *i2cMock) pwmUpdateCompare(channelA uint8, channelB uint8) error      { return nil }
func (l *i2cMock) pwmUpdatePrescaler(value uint) error                        { return nil }
func (l *i2cMock) servoInit() error                                           { return nil }
func (l *i2cMock) servoUpdateLocation(locationA uint8, locationB uint8) error { return nil }
func (l *i2cMock) error() error                                               { return nil }
