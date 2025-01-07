package system

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
)

type spiGpioConfig struct {
	pinProvider gobot.DigitalPinnerProvider
	sclkPinID   string
	ncsPinID    string
	sdoPinID    string
	sdiPinID    string
	debug       bool
}

// spiGpio is the implementation of the SPI interface using GPIO's.
type spiGpio struct {
	cfg spiGpioConfig
	// time between clock edges (i.e. half the cycle time)
	tclk    time.Duration
	sclkPin gobot.DigitalPinner
	ncsPin  gobot.DigitalPinner
	sdoPin  gobot.DigitalPinner
	sdiPin  gobot.DigitalPinner
}

// newSpiGpio creates and returns a new SPI connection based on given GPIO's.
func newSpiGpio(cfg spiGpioConfig, maxSpeedHz int64) (*spiGpio, error) {
	spi := &spiGpio{cfg: cfg}
	spi.initializeTime(maxSpeedHz)
	return spi, spi.initializeGpios()
}

func (s *spiGpio) initializeTime(maxSpeedHz int64) {
	// maxSpeedHz is given in Hz, tclk is half the cycle time, tclk=1/(2*f), tclk[ns]=1 000 000 000/(2*maxSpeed)
	// but with gpio's a speed of more than ~15kHz is most likely not possible, so we limit to 10kHz
	if maxSpeedHz > 10000 {
		if s.cfg.debug {
			fmt.Printf("reduce SPI speed for GPIO usage to 10Khz")
		}
		maxSpeedHz = 10000
	}
	s.tclk = time.Duration(1000000000/2/maxSpeedHz) * time.Nanosecond
	if s.cfg.debug {
		fmt.Println("clk", s.tclk)
	}
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

	if err := s.ncsPin.Write(0); err != nil {
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

	return s.ncsPin.Write(1)
}

// Close the SPI connection. Implements gobot.SpiSystemDevicer.
func (s *spiGpio) Close() error {
	var err error
	if s.sclkPin != nil {
		if e := s.sclkPin.Unexport(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	if s.sdoPin != nil {
		if e := s.sdoPin.Unexport(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	if s.sdiPin != nil {
		if e := s.sdiPin.Unexport(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	if s.ncsPin != nil {
		if e := s.ncsPin.Unexport(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	return err
}

func (cfg *spiGpioConfig) String() string {
	return fmt.Sprintf("sclk: %s, ncs: %s, sdo: %s, sdi: %s", cfg.sclkPinID, cfg.ncsPinID, cfg.sdoPinID, cfg.sdiPinID)
}

// transferByte simultaneously transmit and receive a byte
// polarity and phase are assumed to be both 0 (CPOL=0, CPHA=0), so:
// * input data is captured on rising edge of SCLK
// * output data is propagated on falling edge of SCLK
func (s *spiGpio) transferByte(txByte uint8) (uint8, error) {
	rxByte := uint8(0)
	bitMask := uint8(0x80) // start at MSBit

	for i := 0; i < 8; i++ {
		if err := s.sdoPin.Write(int(txByte & bitMask)); err != nil {
			return 0, err
		}

		time.Sleep(s.tclk)
		if err := s.sclkPin.Write(1); err != nil {
			return 0, err
		}

		v, err := s.sdiPin.Read()
		if err != nil {
			return 0, err
		}
		if v != 0 {
			rxByte |= bitMask
		}

		time.Sleep(s.tclk)
		if err := s.sclkPin.Write(0); err != nil {
			return 0, err
		}

		bitMask = bitMask >> 1 // next lower bit
	}

	return rxByte, nil
}

func (s *spiGpio) initializeGpios() error {
	var err error
	// ncs is an output, negated (currently not implemented at pin level)
	s.ncsPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.ncsPinID)
	if err != nil {
		return err
	}
	if err := s.ncsPin.ApplyOptions(WithPinDirectionOutput(1)); err != nil {
		return err
	}
	// sclk is an output, CPOL = 0
	s.sclkPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.sclkPinID)
	if err != nil {
		return err
	}
	if err := s.sclkPin.ApplyOptions(WithPinDirectionOutput(0)); err != nil {
		return err
	}
	// sdi is an input
	s.sdiPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.sdiPinID)
	if err != nil {
		return err
	}
	// sdo is an output
	s.sdoPin, err = s.cfg.pinProvider.DigitalPin(s.cfg.sdoPinID)
	if err != nil {
		return err
	}
	return s.sdoPin.ApplyOptions(WithPinDirectionOutput(0))
}
