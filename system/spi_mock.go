package system

import (
	"fmt"

	"gobot.io/x/gobot"
)

// MockSpiAccess contains parameters of mocked SPI access
type MockSpiAccess struct {
	CreateError bool
	busNum      int
	chipNum     int
	mode        int
	bits        int
	maxSpeed    int64
	connection  *spiConnectionMock
}

func (spi *MockSpiAccess) createConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiOperations, error) {
	spi.busNum = busNum
	spi.chipNum = chipNum
	spi.mode = mode
	spi.bits = bits
	spi.maxSpeed = maxSpeed
	spi.connection = newSpiConnectionMock(busNum, chipNum, mode, bits, maxSpeed)
	var err error
	if spi.CreateError {
		err = fmt.Errorf("error while create SPI connection in mock")
	}
	return spi.connection, err
}

// SetReadError can be used to simulate a read error.
func (spi *MockSpiAccess) SetReadError(val bool) {
	spi.connection.simReadErr = val
}

// SetWriteError can be used to simulate a write error.
func (spi *MockSpiAccess) SetWriteError(val bool) {
	spi.connection.simWriteErr = val
}

// SetCloseError can be used to simulate a error on Close().
func (spi *MockSpiAccess) SetCloseError(val bool) {
	spi.connection.simCloseErr = val
}

// SetSimRead is used to set the byte stream for next read.
func (spi *MockSpiAccess) SetSimRead(val []byte) {
	spi.connection.simRead = val
}

// Written returns the byte stream which was last written.
func (spi *MockSpiAccess) Written() []byte {
	return spi.connection.written
}

// spiConnectionMock is the a mock implementation, used in tests
type spiConnectionMock struct {
	id          string
	simReadErr  bool
	simWriteErr bool
	simCloseErr bool
	written     []byte
	simRead     []byte
}

// newspiConnectionMock creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface.
func newSpiConnectionMock(busNum, chipNum, mode, bits int, maxSpeed int64) *spiConnectionMock {
	return &spiConnectionMock{id: fmt.Sprintf("bu:%d, c:%d, m:%d, bi:%d, s:%d", busNum, chipNum, mode, bits, maxSpeed)}
}

// Close the SPI connection.
func (c *spiConnectionMock) Close() error {
	if c.simCloseErr {
		return fmt.Errorf("error while SPI close in mock")
	}
	return nil
}

// ReadData uses the SPI device TX to send/receive data.
func (c *spiConnectionMock) ReadData(command []byte, data []byte) error {
	if c.simReadErr {
		return fmt.Errorf("error while SPI read in mock")
	}
	c.written = append(c.written, command...)
	copy(data, c.simRead)
	return nil
}

// WriteData uses the SPI device TX to send data.
func (c *spiConnectionMock) WriteData(data []byte) error {
	if c.simWriteErr {
		return fmt.Errorf("error while SPI write in mock")
	}
	c.written = append(c.written, data...)
	return nil
}
