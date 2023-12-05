package chip

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const defaultI2cBusNumber = 1

// Valids pins are the XIO-P0 through XIO-P7 pins from the
// extender (pins 13-20 on header 14), as well as the SoC pins
// aka all the other pins.
type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for a C.H.I.P.
type Adaptor struct {
	name   string
	sys    *system.Accesser
	mutex  sync.Mutex
	pinmap map[string]sysfsPin
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
}

// NewAdaptor creates a C.H.I.P. Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso):	use GPIO's instead of /dev/spidev#.#
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser()
	a := &Adaptor{
		name: gobot.DefaultName("CHIP"),
		sys:  sys,
	}

	a.pinmap = chipPins
	baseAddr, _ := getXIOBase()
	for i := 0; i < 8; i++ {
		pin := fmt.Sprintf("XIO-P%d", i)
		a.pinmap[pin] = sysfsPin{pin: baseAddr + i, pwmPin: -1}
	}

	var digitalPinsOpts []func(adaptors.DigitalPinsOptioner)
	var pwmPinsOpts []adaptors.PwmPinsOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case func(adaptors.DigitalPinsOptioner):
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.translateDigitalPin, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.translatePWMPin, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, a.validateI2cBusNumber, defaultI2cBusNumber)
	return a
}

// NewProAdaptor creates a C.H.I.P. Pro Adaptor
func NewProAdaptor() *Adaptor {
	a := NewAdaptor()
	a.name = gobot.DefaultName("CHIP Pro")
	a.pinmap = chipProPins
	return a
}

// Name returns the name of the Adaptor
func (a *Adaptor) Name() string { return a.name }

// SetName sets the name of the Adaptor
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect create new connection to board and pins.
func (a *Adaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return a.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (a *Adaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.DigitalPinsAdaptor.Finalize()

	if e := a.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	return err
}

func getXIOBase() (int, error) {
	// Default to original base from 4.3 kernel
	baseAddr := 408
	const expanderID = "pcf8574a"

	labels, err := filepath.Glob("/sys/class/gpio/*/label")
	if err != nil {
		return baseAddr, err
	}

	for _, labelPath := range labels {
		label, err := os.ReadFile(labelPath)
		if err != nil {
			return baseAddr, err
		}
		if strings.HasPrefix(string(label), expanderID) {
			expanderPath, _ := filepath.Split(labelPath)
			basePath := filepath.Join(expanderPath, "base")
			base, err := os.ReadFile(basePath)
			if err != nil {
				return baseAddr, err
			}
			baseAddr, _ = strconv.Atoi(strings.TrimSpace(string(base)))
			break
		}
	}

	return baseAddr, nil
}

func (a *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
	if (busNr < 0) || (busNr > 2) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (a *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if val, ok := a.pinmap[id]; ok {
		return "", val.pin, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}

func (a *Adaptor) translatePWMPin(id string) (string, int, error) {
	sysPin, ok := a.pinmap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a pin", id)
	}
	if sysPin.pwmPin == -1 {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	return "/sys/class/pwm/pwmchip0", sysPin.pwmPin, nil
}
