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
}

func (spi *MockSpiAccess) createConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiOperations, error) {
	spi.busNum = busNum
	spi.chipNum = chipNum
	spi.mode = mode
	spi.bits = bits
	spi.maxSpeed = maxSpeed
	var err error
	if spi.CreateError {
		err = fmt.Errorf("error while create SPI connection in mock")
	}
	return newSpiConnectionMock(busNum, chipNum, mode, bits, maxSpeed), err
}

// spiConnectionMock is the a mock implementation, used in tests
type spiConnectionMock struct {
	id string
}

// newspiConnectionMock creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface.
func newSpiConnectionMock(busNum, chipNum, mode, bits int, maxSpeed int64) *spiConnectionMock {
	return &spiConnectionMock{id: fmt.Sprintf("bu:%d, c:%d, m:%d, bi:%d, s:%d", busNum, chipNum, mode, bits, maxSpeed)}
}

// Close the SPI connection.
func (c *spiConnectionMock) Close() error {
	return nil
}

// ReadData uses the SPI device TX to send/receive data.
func (c *spiConnectionMock) ReadData(command []byte, data []byte) error {
	return nil
}

// WriteData uses the SPI device TX to send data.
func (c *spiConnectionMock) WriteData(data []byte) error {
	return nil
}
