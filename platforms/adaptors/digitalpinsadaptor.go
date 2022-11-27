package adaptors

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/system"
)

type translator func(pin string) (chip string, line int, err error)
type creator func(chip string, line int, o ...func(gobot.DigitalPinOptioner) bool) gobot.DigitalPinner

type digitalPinsOption interface {
	setCreator(creator) // needed e.g. by Beaglebone adaptor
}

// DigitalPinsAdaptor is a adaptor for digital pins, normally used for composition in platforms.
type DigitalPinsAdaptor struct {
	sys       *system.Accesser
	translate translator
	create    creator
	pins      map[string]gobot.DigitalPinner
	mutex     sync.Mutex
}

// NewDigitalPinsAdaptor provides the access to digital pins of the board. It supports sysfs and gpiod system drivers.
// This is decided by the given accesser. The translator is used to adapt the header naming, which is given by user, to
// the internal file name or chip/line nomenclature. This varies by each platform. If for some reasons the default
// creator is not suitable, it can be given by the option "WithPinCreator()". This is especially needed, if some values
// needs to be adjusted after the pin was created. E.g. for Beaglebone platform.
func NewDigitalPinsAdaptor(sys *system.Accesser, t translator, options ...func(digitalPinsOption)) *DigitalPinsAdaptor {
	s := &DigitalPinsAdaptor{
		translate: t,
		create:    sys.NewDigitalPin,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// WithPinCreator can be used to substitute the default creator.
func WithPinCreator(pc creator) func(digitalPinsOption) {
	return func(a digitalPinsOption) {
		a.setCreator(pc)
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

func (a *DigitalPinsAdaptor) setCreator(pc creator) {
	a.create = pc
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
		pin = a.create(chip, line, o...)
		if err = pin.Export(); err != nil {
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
