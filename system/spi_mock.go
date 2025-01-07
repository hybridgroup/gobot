package system

import (
	"fmt"

	"gobot.io/x/gobot/v2"
)

type MockSpiAccess struct {
	underlyingSpiAccess spiAccesser
	CreateError         bool
	busNum              int
	chipNum             int
	mode                int
	bits                int
	maxSpeed            int64
	sysdev              *spiMock
}

func newMockSpiAccess(underlyingSpiAccess spiAccesser) *MockSpiAccess {
	return &MockSpiAccess{underlyingSpiAccess: underlyingSpiAccess}
}

func (spi *MockSpiAccess) createDevice(
	busNum, chipNum, mode, bits int,
	maxSpeed int64,
) (gobot.SpiSystemDevicer, error) {
	spi.busNum = busNum
	spi.chipNum = chipNum
	spi.mode = mode
	spi.bits = bits
	spi.maxSpeed = maxSpeed
	spi.sysdev = newSpiMock(busNum, chipNum, mode, bits, maxSpeed)
	var err error
	if spi.CreateError {
		err = fmt.Errorf("error while create SPI connection in mock")
	}
	return spi.sysdev, err
}

func (spi *MockSpiAccess) isType(accesserType spiBusAccesserType) bool {
	return spi.underlyingSpiAccess.isType(accesserType)
}

func (*MockSpiAccess) isSupported() bool {
	return true
}

// SetReadError can be used to simulate a read error.
func (spi *MockSpiAccess) SetReadError(val bool) {
	spi.sysdev.simReadErr = val
}

// SetWriteError can be used to simulate a write error.
func (spi *MockSpiAccess) SetWriteError(val bool) {
	spi.sysdev.simWriteErr = val
}

// SetCloseError can be used to simulate a error on Close().
func (spi *MockSpiAccess) SetCloseError(val bool) {
	spi.sysdev.simCloseErr = val
}

// SetSimRead is used to set the byte stream for next read.
func (spi *MockSpiAccess) SetSimRead(data []byte) {
	spi.sysdev.simRead = make([]byte, len(data))
	copy(spi.sysdev.simRead, data)
}

// Written returns the byte stream which was last written.
func (spi *MockSpiAccess) Written() []byte {
	return spi.sysdev.written
}

// Reset resets the last written values.
func (spi *MockSpiAccess) Reset() {
	spi.sysdev.written = []byte{}
}

// spiMock is the a mock implementation, used in tests
type spiMock struct {
	id          string
	simReadErr  bool
	simWriteErr bool
	simCloseErr bool
	written     []byte
	simRead     []byte
}

// newSpiMock creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface.
func newSpiMock(busNum, chipNum, mode, bits int, maxSpeed int64) *spiMock {
	return &spiMock{id: fmt.Sprintf("bu:%d, c:%d, m:%d, bi:%d, s:%d", busNum, chipNum, mode, bits, maxSpeed)}
}

// Close the SPI connection to the device. Implements gobot.SpiSystemDevicer.
func (c *spiMock) Close() error {
	if c.simCloseErr {
		return fmt.Errorf("error while SPI close in mock")
	}
	return nil
}

// TxRx uses the SPI device TX to send/receive data. gobot.SpiSystemDevicer.
func (c *spiMock) TxRx(tx []byte, rx []byte) error {
	if c.simReadErr {
		return fmt.Errorf("error while SPI read in mock")
	}
	c.written = append(c.written, tx...)
	// the answer can be one cycle behind, this must be considered in test setup
	copy(rx, c.simRead)
	return nil
}
