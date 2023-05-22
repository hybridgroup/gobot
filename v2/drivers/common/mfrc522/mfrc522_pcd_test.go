package mfrc522

import (
	"testing"

	"gobot.io/x/gobot/v2/gobottest"
)

type busConnMock struct {
	written []byte
	readIdx int
	simRead []byte
	fifoIdx int
	simFifo []byte
}

func (c *busConnMock) ReadByteData(reg uint8) (uint8, error) {
	c.written = append(c.written, reg)

	switch reg {
	case regFIFOLevel:
		return uint8(len(c.simFifo)), nil
	case regFIFOData:
		c.fifoIdx++
		return c.simFifo[c.fifoIdx-1], nil
	default:
		if len(c.simRead) > 0 {
			c.readIdx++
			return c.simRead[c.readIdx-1], nil
		}
		return 0, nil
	}
}

func (c *busConnMock) WriteByteData(reg uint8, data byte) error {
	c.written = append(c.written, reg)
	c.written = append(c.written, data)
	return nil
}

func initTestMFRC522CommonWithStubbedConnector() (*MFRC522Common, *busConnMock) {
	c := &busConnMock{}
	d := NewMFRC522Common()
	d.connection = c
	return d, c
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
	d, c := initTestMFRC522CommonWithStubbedConnector()
	wantWritten := []byte{0x37}
	const want = uint8(5)
	c.simRead = []byte{want}
	// act
	got, err := d.getVersion()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, got, want)
	gobottest.Assert(t, c.written, wantWritten)
}

func Test_switchAntenna(t *testing.T) {
	var tests = map[string]struct {
		target      bool
		simRead     byte
		wantWritten []byte
	}{
		"switch_on": {
			target:      true,
			simRead:     0xFD,
			wantWritten: []byte{0x14, 0x14, 0xFF},
		},
		"is_already_on": {
			target:      true,
			simRead:     0x03,
			wantWritten: []byte{0x14},
		},
		"switch_off": {
			target:      false,
			simRead:     0x03,
			wantWritten: []byte{0x14, 0x14, 0x00},
		},
		"is_already_off": {
			target:      false,
			simRead:     0xFD,
			wantWritten: []byte{0x14},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, c := initTestMFRC522CommonWithStubbedConnector()
			c.simRead = []byte{tc.simRead}
			// act
			err := d.switchAntenna(tc.target)
			// assert
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, c.written, tc.wantWritten)
		})
	}
}

func Test_stopCrypto1(t *testing.T) {
	// arrange
	d, c := initTestMFRC522CommonWithStubbedConnector()
	c.simRead = []byte{0xFF}
	wantWritten := []byte{0x08, 0x08, 0xF7}
	// act
	err := d.stopCrypto1()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, c.written, wantWritten)
}

func Test_communicateWithPICC(t *testing.T) {
	// arrange
	d, c := initTestMFRC522CommonWithStubbedConnector()
	dataToFifo := []byte{0xF1, 0xF2}
	writtenPrepare := []byte{0x02, 0xF7, 0x04, 0x7F, 0x0A, 0x80, 0x01, 0x00}
	writtenWriteFifo := []byte{0x09, 0xF1, 0x09, 0xF2}
	writtenTransceive := []byte{0x0D, 0x00, 0x01, 0x0C}
	writtenBitFramingStart := []byte{0x0D, 0x0D, 0x80}
	writtenWaitAndFinish := []byte{0x04, 0x0D, 0x0D, 0x7F, 0x06}
	writtenReadFifo := []byte{0x0A, 0x09, 0x09, 0x0C}
	// read: bit framing register set, simulate calculation done, bit framing register clear, simulate no error,
	// simulate control register 0x00
	c.simRead = []byte{0x00, 0x30, 0xFF, 0x00, 0x00}
	c.simFifo = []byte{0x11, 0x22}

	backData := []byte{0x00, 0x00}
	// act
	// transceive, all 8 bits, no CRC
	err := d.communicateWithPICC(0x0C, dataToFifo, backData, 0x00, false)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, c.written[:8], writtenPrepare)
	gobottest.Assert(t, c.written[8:12], writtenWriteFifo)
	gobottest.Assert(t, c.written[12:16], writtenTransceive)
	gobottest.Assert(t, c.written[16:19], writtenBitFramingStart)
	gobottest.Assert(t, c.written[19:24], writtenWaitAndFinish)
	gobottest.Assert(t, c.written[24:], writtenReadFifo)
	gobottest.Assert(t, backData, []byte{0x11, 0x22})
}

func Test_calculateCRC(t *testing.T) {
	// arrange
	d, c := initTestMFRC522CommonWithStubbedConnector()
	dataToFifo := []byte{0x02, 0x03}
	writtenPrepare := []byte{0x01, 0x00, 0x05, 0x04, 0x0A, 0x80}
	writtenFifo := []byte{0x09, 0x02, 0x09, 0x03}
	writtenCalc := []byte{0x01, 0x03, 0x05, 0x01, 0x00}
	writtenGetResult := []byte{0x22, 0x21}
	c.simRead = []byte{0x04, 0x11, 0x22} // calculation done, crcL, crcH
	gotCrcBack := []byte{0x00, 0x00}
	// act
	err := d.calculateCRC(dataToFifo, gotCrcBack)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, c.written[:6], writtenPrepare)
	gobottest.Assert(t, c.written[6:10], writtenFifo)
	gobottest.Assert(t, c.written[10:15], writtenCalc)
	gobottest.Assert(t, c.written[15:], writtenGetResult)
	gobottest.Assert(t, gotCrcBack, []byte{0x11, 0x22})
}

func Test_writeFifo(t *testing.T) {
	// arrange
	d, c := initTestMFRC522CommonWithStubbedConnector()
	dataToFifo := []byte{0x11, 0x22, 0x33}
	wantWritten := []byte{0x09, 0x11, 0x09, 0x22, 0x09, 0x33}
	// act
	err := d.writeFifo(dataToFifo)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, c.written, wantWritten)
}

func Test_readFifo(t *testing.T) {
	// arrange
	d, c := initTestMFRC522CommonWithStubbedConnector()
	c.simFifo = []byte{0x30, 0x20, 0x10}
	wantWritten := []byte{0x0A, 0x09, 0x09, 0x09, 0x0C}
	backData := []byte{0x00, 0x00, 0x00}
	// act
	_, err := d.readFifo(backData)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, c.written, wantWritten)
	gobottest.Assert(t, backData, c.simFifo)
}
