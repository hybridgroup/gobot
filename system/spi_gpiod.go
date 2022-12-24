package system

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot"
)

type spiGpioConfig struct {
	pinProvider gobot.DigitalPinnerProvider
	sclkPinId   string
	nssPinId    string
	mosiPinId   string
	misoPinId   string
}

// spiGpiod is the implementation of the SPI interface using the periph.io sysfs implementation for Linux.
type spiGpio struct {
	cfg spiGpioConfig
	// time between clock edges (i.e. half the cycle time)
	tclk    time.Duration
	sclkPin gobot.DigitalPinner
	nssPin  gobot.DigitalPinner
	mosiPin gobot.DigitalPinner
	misoPin gobot.DigitalPinner
}

// newSpiGpio creates and returns a new SPI connection based on given GPIO's.
func newSpiGpio(cfg spiGpioConfig, maxSpeed int64) (*spiGpio, error) {
	spi := &spiGpio{cfg: cfg}
	spi.initializeTime(maxSpeed)
	return spi, spi.initializeGpios()
}

func (s *spiGpio) initializeTime(maxSpeed int64) {
	// maxSpeed is given in Hz, tclk is half the cycle time, tclk=1/(2*f), tclk[ns]=1 000 000 000/(2*maxSpeed)
	// but with gpio's a speed of more than ~15kHz is most likely not possible, so we limit to 10kHz
	if maxSpeed > 10000 {
		maxSpeed = 10000
	}
	tclk := time.Duration(1000000000/2/maxSpeed) * time.Nanosecond
	log.Println("clk", tclk)
}

// TxRx uses the SPI device to send/receive data. Implements gobot.SpiSystemDevicer.
func (s *spiGpio) TxRx(tx []byte, rx []byte) error {
	var doRx bool
	if rx != nil {
		doRx = true
		if len(tx) != len(rx) {
			return fmt.Errorf("length of tx (%d) must be the same as length of rx (%d)", len(tx), len(rx))
		}
	}

	if err := s.nssPin.Write(0); err != nil {
		return err
	}

	for idx, b := range tx {
		val, err := s.transferByte(b)
		if err != nil {
			return err
		}
		if doRx {
			rx[idx] = val
		}
	}

	return s.nssPin.Write(1)
}

// Close the SPI connection. Implements gobot.SpiSystemDevicer.
func (s *spiGpio) Close() error {
	if s.sclkPin != nil {
		s.sclkPin.Unexport()
	}
	if s.mosiPin != nil {
		s.mosiPin.Unexport()
	}
	if s.misoPin != nil {
		s.misoPin.Unexport()
	}
	if s.nssPin != nil {
		s.nssPin.Unexport()
	}
	return nil
}

// transferByte simultaneously transmit and receive a byte
// polarity and phase are assumed to be both 0 (CPOL=0, CPHA=0), so:
// * input data is captured on rising edge of SCLK
// * output data is propagated on falling edge of SCLK
func (c *spiGpio) transferByte(txByte uint8) (uint8, error) {
	rxByte := uint8(0)
	bitMask := uint8(0x80) // start at MSBit

	for i := 0; i < 8; i++ {
		if err := c.mosiPin.Write(int(txByte & bitMask)); err != nil {
			return 0, err
		}

		time.Sleep(c.tclk)
		if err := c.sclkPin.Write(1); err != nil {
			return 0, err
		}

		v, err := c.misoPin.Read()
		if err != nil {
			return 0, err
		}
		if v != 0 {
			rxByte |= bitMask
		}

		time.Sleep(c.tclk)
		if err := c.sclkPin.Write(0); err != nil {
			return 0, err
		}

		bitMask = bitMask >> 1 // next lower bit
	}

	return rxByte, nil
}

func (cfg *spiGpioConfig) String() string {
	return fmt.Sprintf("sclk: %s, nss: %s, mosi: %s, miso: %s", cfg.sclkPinId, cfg.nssPinId, cfg.mosiPinId, cfg.misoPinId)
}

func (s *spiGpio) initializeGpios() error {
	var err error
	// nss is an output, negotiated (currently not implemented at pin level)
	s.nssPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.nssPinId)
	if err != nil {
		return err
	}
	if err := s.nssPin.ApplyOptions(WithDirectionOutput(1)); err != nil {
		return err
	}
	// sclk is an output, CPOL = 0
	s.sclkPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.sclkPinId)
	if err != nil {
		return err
	}
	if err := s.sclkPin.ApplyOptions(WithDirectionOutput(0)); err != nil {
		return err
	}
	// miso is an input
	s.misoPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.misoPinId)
	if err != nil {
		return err
	}
	// mosi is an output
	s.mosiPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.mosiPinId)
	if err != nil {
		return err
	}
	return s.mosiPin.ApplyOptions(WithDirectionOutput(0))
}
