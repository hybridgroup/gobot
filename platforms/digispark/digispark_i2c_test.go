package digispark

import (
	"errors"
	"fmt"
	"testing"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
)

const availableI2cAddress = 0x40
const maxUint8 = ^uint8(0)

var _ i2c.Connector = (*Adaptor)(nil)
var i2cData = []byte{5, 4, 3, 2, 1, 0}

type i2cMock struct {
	address      int
	duration     uint
	direction    uint8
	dataWritten  []byte
	startWasSend bool
	stopWasSend  bool
}

func initTestAdaptorI2c() *Adaptor {
	a := NewAdaptor()
	a.connect = func(a *Adaptor) (err error) { return nil }
	a.littleWire = new(i2cMock)
	return a
}

func TestDigisparkAdaptorI2cGetConnection(t *testing.T) {
	// arrange
	var c i2c.Connection
	var err error
	a := initTestAdaptorI2c()

	// act
	c, err = a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, c, nil)
}

func TestDigisparkAdaptorI2cGetConnectionFailWithInvalidBus(t *testing.T) {
	// arrange
	a := initTestAdaptorI2c()

	// act
	c, err := a.GetConnection(0x40, 1)

	// assert
	gobottest.Assert(t, err, errors.New("Invalid bus number 1, only 0 is supported"))
	gobottest.Assert(t, c, nil)
}

func TestDigisparkAdaptorI2cStartFailWithWrongAddress(t *testing.T) {
	// arrange
	data := []byte{0, 1, 2, 3, 4}
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(0x39, a.GetDefaultBus())

	// act
	count, err := c.Write(data)

	// assert
	gobottest.Assert(t, count, 0)
	gobottest.Assert(t, err, fmt.Errorf("Invalid address, only %d is supported", availableI2cAddress))
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, maxUint8)
}

func TestDigisparkAdaptorI2cWrite(t *testing.T) {
	// arrange
	data := []byte{0, 1, 2, 3, 4}
	dataLen := len(data)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	count, err := c.Write(data)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(0))
	gobottest.Assert(t, count, dataLen)
	gobottest.Assert(t, a.littleWire.(*i2cMock).dataWritten, data)
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cWriteByte(t *testing.T) {
	// arrange
	data := byte(0x02)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	err = c.WriteByte(data)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(0))
	gobottest.Assert(t, a.littleWire.(*i2cMock).dataWritten, []byte{data})
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cWriteByteData(t *testing.T) {
	// arrange
	reg := uint8(0x03)
	data := byte(0x09)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	err = c.WriteByteData(reg, data)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(0))
	gobottest.Assert(t, a.littleWire.(*i2cMock).dataWritten, []byte{reg, data})
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cWriteWordData(t *testing.T) {
	// arrange
	reg := uint8(0x04)
	data := uint16(0x0508)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	err = c.WriteWordData(reg, data)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(0))
	gobottest.Assert(t, a.littleWire.(*i2cMock).dataWritten, []byte{reg, 0x08, 0x05})
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cRead(t *testing.T) {
	// arrange
	data := []byte{0, 1, 2, 3, 4}
	dataLen := len(data)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	count, err := c.Read(data)

	// assert
	gobottest.Assert(t, count, dataLen)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(1))
	gobottest.Assert(t, data, i2cData[:dataLen])
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cReadByte(t *testing.T) {
	// arrange
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	data, err := c.ReadByte()

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(1))
	gobottest.Assert(t, data, i2cData[0])
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cReadByteData(t *testing.T) {
	// arrange
	reg := uint8(0x04)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	data, err := c.ReadByteData(reg)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(1))
	gobottest.Assert(t, a.littleWire.(*i2cMock).dataWritten, []byte{reg})
	gobottest.Assert(t, data, i2cData[0])
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cReadWordData(t *testing.T) {
	// arrange
	reg := uint8(0x05)
	// 2 bytes of i2cData are used swapped
	expectedValue := uint16(0x0405)
	a := initTestAdaptorI2c()
	c, err := a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	data, err := c.ReadWordData(reg)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).startWasSend, true)
	gobottest.Assert(t, a.littleWire.(*i2cMock).direction, uint8(1))
	gobottest.Assert(t, a.littleWire.(*i2cMock).dataWritten, []byte{reg})
	gobottest.Assert(t, data, expectedValue)
	gobottest.Assert(t, a.littleWire.(*i2cMock).stopWasSend, true)
}

func TestDigisparkAdaptorI2cUpdateDelay(t *testing.T) {
	// arrange
	var c i2c.Connection
	var err error
	a := initTestAdaptorI2c()
	c, err = a.GetConnection(availableI2cAddress, a.GetDefaultBus())

	// act
	err = c.(*digisparkI2cConnection).UpdateDelay(uint(100))

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*i2cMock).duration, uint(100))
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
	l.startWasSend = true
	return nil
}

func (l *i2cMock) i2cWrite(sendBuffer []byte, length int, endWithStop uint8) error {
	l.dataWritten = append(l.dataWritten, sendBuffer...)
	if endWithStop > 0 {
		l.stopWasSend = true
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
		l.stopWasSend = true
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
