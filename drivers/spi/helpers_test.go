package spi

import (
	"fmt"
)

type spiTestAdaptor struct {
	busNum        int
	spiConnectErr bool
	device        *spiTestDevice
	simRead       []byte
}

func (a *spiTestAdaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (device Connection, err error) {
	if a.spiConnectErr {
		return nil, fmt.Errorf("Invalid SPI connection in helper")
	}
	a.busNum = busNum
	a.device = &spiTestDevice{simRead: a.simRead}
	return a.device, nil
}

func (a *spiTestAdaptor) SpiDefaultBusNumber() int  { return a.busNum }
func (a *spiTestAdaptor) SpiDefaultChipNumber() int { return 0 }
func (a *spiTestAdaptor) SpiDefaultMode() int       { return 0 }
func (a *spiTestAdaptor) SpiDefaultBitCount() int   { return 0 }
func (a *spiTestAdaptor) SpiDefaultMaxSpeed() int64 { return 0 }

type spiTestDevice struct {
	spiReadErr  bool
	spiWriteErr bool
	written     []byte
	simRead     []byte
}

func (t *spiTestDevice) ReadData(command, data []byte) error {
	if t.spiReadErr {
		return fmt.Errorf("Error on SPI read in helper")
	}
	t.written = append(t.written, command...)
	copy(data, t.simRead)
	return nil
}

func (t *spiTestDevice) WriteData(data []byte) error {
	if t.spiWriteErr {
		return fmt.Errorf("Error on SPI write in helper")
	}
	t.written = append(t.written, data...)
	return nil
}

func (t *spiTestDevice) Close() error { return nil }

func newSpiTestAdaptor(simRead []byte) *spiTestAdaptor {
	return &spiTestAdaptor{simRead: simRead}
}
