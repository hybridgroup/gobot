package adaptors

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/system"
)

type digitalPinTranslator func(pin string) (chip string, line int, err error)
type digitalPinInitializer func(gobot.DigitalPinner) error

type digitalPinsOptioner interface {
	setDigitalPinInitializer(digitalPinInitializer)
	setDigitalPinsForSystemGpiod()
	setDigitalPinsForSystemSpi(sclkPin, nssPin, mosiPin, misoPin string)
}

// DigitalPinsAdaptor is a adaptor for digital pins, normally used for composition in platforms.
type DigitalPinsAdaptor struct {
	sys        *system.Accesser
	translate  digitalPinTranslator
	initialize digitalPinInitializer
	pins       map[string]gobot.DigitalPinner
	mutex      sync.Mutex
}

// NewDigitalPinsAdaptor provides the access to digital pins of the board. It supports sysfs and gpiod system drivers.
// This is decided by the given accesser. The translator is used to adapt the pin header naming, which is given by user,
// to the internal file name or chip/line nomenclature. This varies by each platform. If for some reasons the default
// initializer is not suitable, it can be given by the option "WithDigitalPinInitializer()". This is especially needed,
// if some values needs to be adjusted after the pin was created but before the pin is exported.
func NewDigitalPinsAdaptor(sys *system.Accesser, t digitalPinTranslator, options ...func(Optioner)) *DigitalPinsAdaptor {
	a := &DigitalPinsAdaptor{
		sys:        sys,
		translate:  t,
		initialize: func(pin gobot.DigitalPinner) error { return pin.Export() },
	}
	for _, option := range options {
		option(a)
	}
	return a
}

// WithDigitalPinInitializer can be used to substitute the default initializer.
func WithDigitalPinInitializer(pc digitalPinInitializer) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.setDigitalPinInitializer(pc)
		}
	}
}

// WithGpiodAccess can be used to change the default sysfs implementation to the character device Kernel ABI.
// The access is provided by the gpiod package.
func WithGpiodAccess() func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.setDigitalPinsForSystemGpiod()
		}
	}
}

// WithSpiGpioAccess can be used to switch the default SPI implementation to GPIO usage.
func WithSpiGpioAccess(sclkPin, nssPin, mosiPin, misoPin string) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.setDigitalPinsForSystemSpi(sclkPin, nssPin, mosiPin, misoPin)
		}
	}
}

// Connect prepare new connection to digital pins.
func (a *DigitalPinsAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.pins = make(map[string]gobot.DigitalPinner)
	return nil
}

// Finalize closes connection to digital pins
func (a *DigitalPinsAdaptor) Finalize() (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, pin := range a.pins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	a.pins = nil
	return
}

// DigitalPin returns a digital pin. If the pin is initially acquired, it is an input.
// Pin direction and other options can be changed afterwards by pin.ApplyOptions() at any time.
func (a *DigitalPinsAdaptor) DigitalPin(id string) (gobot.DigitalPinner, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.digitalPin(id)
}

// DigitalRead reads digital value from pin
func (a *DigitalPinsAdaptor) DigitalRead(id string) (int, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.digitalPin(id, system.WithDirectionInput())
	if err != nil {
		return 0, err
	}
	return pin.Read()
}

// DigitalWrite writes digital value to specified pin
func (a *DigitalPinsAdaptor) DigitalWrite(id string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.digitalPin(id, system.WithDirectionOutput(int(val)))
	if err != nil {
		return err
	}
	return pin.Write(int(val))
}

func (a *DigitalPinsAdaptor) setDigitalPinInitializer(pinInit digitalPinInitializer) {
	a.initialize = pinInit
}

func (a *DigitalPinsAdaptor) setDigitalPinsForSystemGpiod() {
	system.WithDigitalPinGpiodAccess()(a.sys)
}

func (a *DigitalPinsAdaptor) setDigitalPinsForSystemSpi(sclkPin, nssPin, mosiPin, misoPin string) {
	system.WithSpiGpioAccess(a, sclkPin, nssPin, mosiPin, misoPin)(a.sys)
}

func (a *DigitalPinsAdaptor) digitalPin(id string, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	if a.pins == nil {
		return nil, fmt.Errorf("not connected")
	}

	pin := a.pins[id]

	if pin == nil {
		chip, line, err := a.translate(id)
		if err != nil {
			return nil, err
		}
		pin = a.sys.NewDigitalPin(chip, line, o...)
		if err = a.initialize(pin); err != nil {
			return nil, err
		}
		a.pins[id] = pin
	} else {
		if err := pin.ApplyOptions(o...); err != nil {
			return nil, err
		}
	}

	return pin, nil
}
