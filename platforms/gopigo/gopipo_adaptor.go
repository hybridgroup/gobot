package gopigo

import (
	"strconv"
	"time"

	"github.com/hybridgroup/gobot/platforms/raspi"
)

const (
	ADDRESS = 0x08

	DIGITAL_WRITE = 0x0C // digital write on a port
	DIGITAL_READ  = 0x0D // digital read on a port
	ANALOG_READ   = 0x0E // analog read on a port
	PWM_WRITE     = 0x0F // analog read on a port
	PIN_MODE      = 0x10 // set up the pin mode on a port TODO: function

	PIN_DIGITAL   = "10"
	PIN_ANALOG    = "15"
	PIN_LED_LEFT  = "16"
	PIN_LED_RIGHT = "17"

	PIN_MODE_OUTPUT = 0x01
	PIN_MODE_INPUT  = 0x00
)

var pins = map[string]byte{
	"0":  0,
	"1":  1,
	"10": 10,
	"15": 15,
}

type GoPiGoAdaptor struct {
	name  string
	raspi *raspi.Adaptor
}

func NewAdaptor() (*GoPiGoAdaptor, error) {
	g := &GoPiGoAdaptor{
		name:  "GoPiGo",
		raspi: raspi.NewAdaptor(),
	}
	err := g.raspi.I2cStart(ADDRESS)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GoPiGoAdaptor) Connect() error {
	return g.raspi.Connect()
}

func (g *GoPiGoAdaptor) Finalize() error {
	return g.raspi.Finalize()
}

func (g *GoPiGoAdaptor) I2cStart(address int) error {
	return g.raspi.I2cStart(address)
}

func (g *GoPiGoAdaptor) Name() string {
	return g.name
}

func (g *GoPiGoAdaptor) SetName(name string) {
	g.name = name
}

// DigitalRead reads the from pin
func (g *GoPiGoAdaptor) DigitalRead(pin string) (int, error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return 0, err
	}
	err = g.I2cWrite(ADDRESS, []byte{DIGITAL_READ, byte(p), 0, 0})
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	d, err := g.I2cRead(ADDRESS, 2)
	return int(d[0]), err
}

func (g *GoPiGoAdaptor) DigitalWrite(pin string, val byte) error {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	err = g.I2cWrite(ADDRESS, []byte{DIGITAL_WRITE, byte(p), val, 0})
	return err
}

func (g *GoPiGoAdaptor) AnalogRead(pin string) (int, error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return 0, err
	}
	err = g.I2cWrite(ADDRESS, []byte{ANALOG_READ, byte(p), 0, 0})
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	d, err := g.I2cRead(ADDRESS, 2)
	if err != nil {
		return 0, err
	}
	return int(d[0])*256 + int(d[1]), nil
}

func (g *GoPiGoAdaptor) PwmWrite(pin string, val byte) error {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	err = g.I2cWrite(ADDRESS, []byte{PWM_WRITE, byte(p), val, 0})
	return err
}

func (g *GoPiGoAdaptor) I2cWrite(address int, data []byte) error {
	err := g.raspi.I2cWrite(address, data)
	time.Sleep(5 * time.Millisecond)
	return err
}

func (g *GoPiGoAdaptor) I2cRead(address int, size int) ([]byte, error) {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		d, err := g.raspi.I2cRead(address, size)
		buf[i] = d[0]
		if err != nil {
			return nil, err
		}
		time.Sleep(5 * time.Millisecond)
	}
	return buf, nil
}

func (g *GoPiGoAdaptor) PinMode(pin string, mode byte) error {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	err = g.raspi.I2cWrite(ADDRESS, []byte{PIN_MODE, byte(p), mode, 0})
	time.Sleep(5 * time.Millisecond)
	return err
}
