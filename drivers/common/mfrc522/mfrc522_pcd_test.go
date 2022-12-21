package mfrc522

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

type busConnMock struct {
	written []byte
	simRead byte
}

func (c *busConnMock) ReadByteData(reg uint8) (uint8, error) {
	c.written = append(c.written, reg)
	return c.simRead, nil
}

func (c *busConnMock) WriteByteData(reg uint8, data byte) error {
	c.written = append(c.written, reg)
	c.written = append(c.written, data)
	return nil
}

func TestNewMFRC522Common(t *testing.T) {
	// act
	d := NewMFRC522Common()
	// assert
	gobottest.Refute(t, d, nil)
}

func TestInitialize(t *testing.T) {
	// arrange
	wantSoftReset := []byte{0x01, 0x0F, 0x01}
	wantInit := []byte{0x12, 0x00, 0x13, 0x00, 0x24, 0x26, 0x2A, 0x8F, 0x2B, 0xFF, 0x2D, 0xE8, 0x2C, 0x03, 0x15, 0x40, 0x11, 0x29}
	wantAntennaOn := []byte{0x14, 0x14, 0x03}
	wantGain := []byte{0x26, 0x50}
	c := &busConnMock{}
	d := NewMFRC522Common()
	// act
	err := d.Initialize(c)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, d.connection, c)
	gobottest.Assert(t, c.written[:3], wantSoftReset)
	gobottest.Assert(t, c.written[3:21], wantInit)
	gobottest.Assert(t, c.written[21:24], wantAntennaOn)
	gobottest.Assert(t, c.written[24:], wantGain)
}

func Test_getVersion(t *testing.T) {
	// arrange
	c := &busConnMock{}
	d := NewMFRC522Common()
	d.connection = c
	wantWritten := []byte{0x37}
	const want = uint8(5)
	c.simRead = want
	// act
	got, err := d.getVersion()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, got, want)
	gobottest.Assert(t, c.written, wantWritten)
}
