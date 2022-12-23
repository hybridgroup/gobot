package system

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
	xspi "github.com/warthog618/gpiod/spi"
)

// spiGpiod is the implementation of the SPI interface using the periph.io sysfs implementation for Linux.
type spiGpiod struct {
	xs *xspi.SPI
}

// newSpiGpiod creates and returns a new SPI connection based on given GPIO's.
func newSpiGpiod(chipName string, sclk, ssz, mosi, miso int, tclk time.Duration) (*spiGpiod, error) {
	c, err := gpiod.NewChip(chipName, gpiod.WithConsumer("spi_emulation"))
	xs, err := xspi.New(c, sclk, ssz, mosi, miso)
	xspi.WithTclk(tclk)
	if err != nil {
		return nil, err
	}
	return &spiGpiod{xs: xs}, nil
}

// TxRx uses the SPI device to send/receive data. Implements gobot.SpiSystemDevicer.
func (c *spiGpiod) TxRx(tx []byte, rx []byte) error {
	dataLen := len(rx)
	if len(tx) != len(rx) {
		return fmt.Errorf("length of tx (%d) must be the same as length of rx (%d)", len(tx), len(rx))
	}

	for idx, b := range tx {
		if err := c.writeByte(b); err != nil {
			return err
		}

		val, err := c.readByte()
		if err != nil {
			return err
		}
		rx[idx] = val
	}

	if len(rx) != dataLen {
		return fmt.Errorf("Read length (%d) differ to expected (%d)", len(rx), dataLen)
	}
	return nil
}

// Close the SPI connection. Implements gobot.SpiSystemDevicer.
func (c *spiGpiod) Close() error {
	c.xs.Close()
	return nil
}

func (c *spiGpiod) writeByte(b byte) error {
	// bit wise clock out the given byte
	for j := 0; j < 8; j++ {
		mask := byte(1 << uint(j))
		if (b & mask) == 0 {
			if err := c.xs.ClockOut(0); err != nil {
				return err
			}
			continue
		}
		if err := c.xs.ClockOut(1); err != nil {
			return err
		}
	}
	return nil
}

func (c *spiGpiod) readByte() (uint8, error) {
	// bit wise clock in a byte
	var b uint8
	for i := uint(0); i < 8; i++ {
		v, err := c.xs.ClockIn()
		if err != nil {
			return 0, err
		}
		b = b << 1
		if v != 0 {
			b = b | 0x01
		}
	}
	return b, nil
}
