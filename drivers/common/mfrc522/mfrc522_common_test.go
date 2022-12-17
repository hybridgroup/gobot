package mfrc522

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

type busConnMock struct {
	written []byte
}

func (c *busConnMock) WriteByteData(reg byte, data byte) error {
	c.written = append(c.written, reg)
	c.written = append(c.written, data)
	return nil
}
func (c *busConnMock) ReadByteData(reg byte) (byte, error) {
	c.written = append(c.written, reg)
	return 0, nil
}

func TestNewMFRC522Common(t *testing.T) {
	// act
	d := NewMFRC522Common()
	// assert
	gobottest.Refute(t, d, nil)
}

func TestConnect(t *testing.T) {
	// arrange
	c := &busConnMock{}
	d := NewMFRC522Common()
	// act
	d.Connect(c)
	// assert
	gobottest.Refute(t, d, nil)
	gobottest.Assert(t, d.connection, c)
}

func TestInitialize(t *testing.T) {
	// arrange
	wantSoftReset := []byte{0x01, 0x0F, 0x01}
	wantInit := []byte{0x12, 0x00, 0x13, 0x00, 0x24, 0x26, 0x2A, 0x8F, 0x2B, 0xFF, 0x2D, 0xE8, 0x2C, 0x03, 0x15, 0x40, 0x11, 0xA1}
	wantAntenna := []byte{0x14, 0x14, 0x03}
	wantGetVersion := []byte{0x37}
	c := &busConnMock{}
	d := NewMFRC522Common()
	d.connection = c
	// act
	err := d.Initialize()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, c.written[:3], wantSoftReset)
	gobottest.Assert(t, c.written[3:21], wantInit)
	gobottest.Assert(t, c.written[21:24], wantAntenna)
	gobottest.Assert(t, c.written[24:], wantGetVersion)
}
