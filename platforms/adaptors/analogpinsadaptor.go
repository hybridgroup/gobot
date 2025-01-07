package adaptors

import (
	"fmt"
	"sync"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

type analogPinTranslator func(pin string) (path string, w bool, readBufLen uint16, err error)

// AnalogPinsAdaptor is a adaptor for analog pins, normally used for composition in platforms.
// It is also usable for general sysfs access.
type AnalogPinsAdaptor struct {
	sys       *system.Accesser
	translate analogPinTranslator
	pins      map[string]gobot.AnalogPinner
	mutex     sync.Mutex
}

// NewAnalogPinsAdaptor provides the access to analog pins of the board. Usually sysfs system drivers are used.
// The translator is used to adapt the pin header naming, which is given by user, to the internal file name
// nomenclature. This varies by each platform.
func NewAnalogPinsAdaptor(sys *system.Accesser, t analogPinTranslator) *AnalogPinsAdaptor {
	a := AnalogPinsAdaptor{
		sys:       sys,
		translate: t,
	}

	sys.AddAnalogSupport()

	return &a
}

// Connect prepare new connection to analog pins.
func (a *AnalogPinsAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.pins = make(map[string]gobot.AnalogPinner)
	return nil
}

// Finalize closes connection to analog pins
func (a *AnalogPinsAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.pins = nil
	return nil
}

// AnalogRead returns an analog value from specified pin or identifier, defined by the translation function.
func (a *AnalogPinsAdaptor) AnalogRead(id string) (int, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.analogPin(id)
	if err != nil {
		return 0, err
	}

	return pin.Read()
}

// AnalogWrite writes an analog value to the specified pin or identifier, defined by the translation function.
func (a *AnalogPinsAdaptor) AnalogWrite(id string, val int) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.analogPin(id)
	if err != nil {
		return err
	}

	return pin.Write(val)
}

// analogPin initializes the pin for analog access and returns matched pin for specified identifier.
func (a *AnalogPinsAdaptor) analogPin(id string) (gobot.AnalogPinner, error) {
	if a.pins == nil {
		return nil, fmt.Errorf("not connected for pin %s", id)
	}

	pin := a.pins[id]

	if pin == nil {
		path, w, readBufLen, err := a.translate(id)
		if err != nil {
			return nil, err
		}
		pin = a.sys.NewAnalogPin(path, w, readBufLen)
		a.pins[id] = pin
	}

	return pin, nil
}
