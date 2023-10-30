package adaptors

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

type (
	digitalPinTranslator  func(pin string) (chip string, line int, err error)
	digitalPinInitializer func(gobot.DigitalPinner) error
)

type digitalPinsOptioner interface {
	setDigitalPinInitializer(digitalPinInitializer)
	setDigitalPinsForSystemGpiod()
	setDigitalPinsForSystemSpi(sclkPin, nssPin, mosiPin, misoPin string)
	prepareDigitalPinsActiveLow(pin string, otherPins ...string)
	prepareDigitalPinsPullDown(pin string, otherPins ...string)
	prepareDigitalPinsPullUp(pin string, otherPins ...string)
	prepareDigitalPinsOpenDrain(pin string, otherPins ...string)
	prepareDigitalPinsOpenSource(pin string, otherPins ...string)
	prepareDigitalPinDebounce(pin string, period time.Duration)
	prepareDigitalPinEventOnFallingEdge(pin string, handler func(lineOffset int, timestamp time.Duration,
		detectedEdge string, seqno uint32, lseqno uint32))
	prepareDigitalPinEventOnRisingEdge(pin string, handler func(lineOffset int, timestamp time.Duration,
		detectedEdge string, seqno uint32, lseqno uint32))
	prepareDigitalPinEventOnBothEdges(pin string, handler func(lineOffset int, timestamp time.Duration,
		detectedEdge string, seqno uint32, lseqno uint32))
	prepareDigitalPinPollForEdgeDetection(pin string, pollInterval time.Duration, pollQuitChan chan struct{})
}

// DigitalPinsAdaptor is a adaptor for digital pins, normally used for composition in platforms.
type DigitalPinsAdaptor struct {
	sys        *system.Accesser
	translate  digitalPinTranslator
	initialize digitalPinInitializer
	pins       map[string]gobot.DigitalPinner
	pinOptions map[string][]func(gobot.DigitalPinOptioner) bool
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

// WithGpiosActiveLow prepares the given pins for inverse reaction on next initialize.
// This is working for inputs and outputs.
func WithGpiosActiveLow(pin string, otherPins ...string) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinsActiveLow(pin, otherPins...)
		}
	}
}

// WithGpiosPullDown prepares the given pins to be pulled down (high impedance to GND) on next initialize.
// This is working for inputs and outputs since Kernel 5.5, but will be ignored with sysfs ABI.
func WithGpiosPullDown(pin string, otherPins ...string) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinsPullDown(pin, otherPins...)
		}
	}
}

// WithGpiosPullUp prepares the given pins to be pulled up (high impedance to VDD) on next initialize.
// This is working for inputs and outputs since Kernel 5.5, but will be ignored with sysfs ABI.
func WithGpiosPullUp(pin string, otherPins ...string) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinsPullUp(pin, otherPins...)
		}
	}
}

// WithGpiosOpenDrain prepares the given output pins to be driven with open drain/collector on next initialize.
// This will be ignored for inputs or with sysfs ABI.
func WithGpiosOpenDrain(pin string, otherPins ...string) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinsOpenDrain(pin, otherPins...)
		}
	}
}

// WithGpiosOpenSource prepares the given output pins to be driven with open source/emitter on next initialize.
// This will be ignored for inputs or with sysfs ABI.
func WithGpiosOpenSource(pin string, otherPins ...string) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinsOpenSource(pin, otherPins...)
		}
	}
}

// WithGpioDebounce prepares the given input pin to be debounced on next initialize.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioDebounce(pin string, period time.Duration) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinDebounce(pin, period)
		}
	}
}

// WithGpioEventOnFallingEdge prepares the given input pin to be generate an event on falling edge.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioEventOnFallingEdge(pin string, handler func(lineOffset int, timestamp time.Duration, detectedEdge string,
	seqno uint32, lseqno uint32),
) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinEventOnFallingEdge(pin, handler)
		}
	}
}

// WithGpioEventOnRisingEdge prepares the given input pin to be generate an event on rising edge.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioEventOnRisingEdge(pin string, handler func(lineOffset int, timestamp time.Duration, detectedEdge string,
	seqno uint32, lseqno uint32),
) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinEventOnRisingEdge(pin, handler)
		}
	}
}

// WithGpioEventOnBothEdges prepares the given input pin to be generate an event on rising and falling edges.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioEventOnBothEdges(pin string, handler func(lineOffset int, timestamp time.Duration, detectedEdge string,
	seqno uint32, lseqno uint32),
) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinEventOnBothEdges(pin, handler)
		}
	}
}

// WithGpioPollForEdgeDetection prepares the given input pin to use a discrete input pin polling function together with
// edge detection.
func WithGpioPollForEdgeDetection(pin string, pollInterval time.Duration, pollQuitChan chan struct{}) func(Optioner) {
	return func(o Optioner) {
		a, ok := o.(digitalPinsOptioner)
		if ok {
			a.prepareDigitalPinPollForEdgeDetection(pin, pollInterval, pollQuitChan)
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
	a.pinOptions = nil
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

	pin, err := a.digitalPin(id, system.WithPinDirectionInput())
	if err != nil {
		return 0, err
	}
	return pin.Read()
}

// DigitalWrite writes digital value to specified pin
func (a *DigitalPinsAdaptor) DigitalWrite(id string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.digitalPin(id, system.WithPinDirectionOutput(int(val)))
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

func (a *DigitalPinsAdaptor) prepareDigitalPinsActiveLow(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinActiveLow())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsPullDown(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinPullDown())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsPullUp(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinPullUp())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsOpenDrain(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinOpenDrain())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinsOpenSource(id string, otherIDs ...string) {
	ids := []string{id}
	ids = append(ids, otherIDs...)

	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	for _, i := range ids {
		a.pinOptions[i] = append(a.pinOptions[i], system.WithPinOpenSource())
	}
}

func (a *DigitalPinsAdaptor) prepareDigitalPinDebounce(id string, period time.Duration) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinDebounce(period))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinEventOnFallingEdge(id string, handler func(int, time.Duration, string,
	uint32, uint32),
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinEventOnFallingEdge(handler))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinEventOnRisingEdge(id string, handler func(int, time.Duration, string,
	uint32, uint32),
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinEventOnRisingEdge(handler))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinEventOnBothEdges(id string, handler func(int, time.Duration, string,
	uint32, uint32),
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinEventOnBothEdges(handler))
}

func (a *DigitalPinsAdaptor) prepareDigitalPinPollForEdgeDetection(
	id string,
	pollInterval time.Duration,
	pollQuitChan chan struct{},
) {
	if a.pinOptions == nil {
		a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)
	}

	a.pinOptions[id] = append(a.pinOptions[id], system.WithPinPollForEdgeDetection(pollInterval, pollQuitChan))
}

func (a *DigitalPinsAdaptor) digitalPin(id string, opts ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	if a.pins == nil {
		return nil, fmt.Errorf("not connected for pin %s", id)
	}

	o := append(a.pinOptions[id], opts...)
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
