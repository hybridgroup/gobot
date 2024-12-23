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

type digitalPinGpiosForSPI struct {
	sclkPin string
	ncsPin  string
	sdoPin  string
	sdiPin  string
}

type digitalPinsDebouncePin struct {
	id     string
	period time.Duration
}

type digitalPinsEventOnEdgePin struct {
	id      string
	handler func(lineOffset int, timestamp time.Duration, detectedEdge string, seqno uint32, lseqno uint32)
}

type digitalPinsPollForEdgeDetectionPin struct {
	id           string
	pollInterval time.Duration
	pollQuitChan chan struct{}
}

// digitalPinsConfiguration contains all changeable attributes of the adaptor.
type digitalPinsConfiguration struct {
	initialize               digitalPinInitializer
	useSysfs                 *bool
	gpiosForSPI              *digitalPinGpiosForSPI
	activeLowPins            []string
	pullDownPins             []string
	pullUpPins               []string
	openDrainPins            []string
	openSourcePins           []string
	debouncePins             []digitalPinsDebouncePin
	eventOnFallingEdgePins   []digitalPinsEventOnEdgePin
	eventOnRisingEdgePins    []digitalPinsEventOnEdgePin
	eventOnBothEdgesPins     []digitalPinsEventOnEdgePin
	pollForEdgeDetectionPins []digitalPinsPollForEdgeDetectionPin
}

// DigitalPinsAdaptor is a adaptor for digital pins, normally used for composition in platforms.
type DigitalPinsAdaptor struct {
	sys            *system.Accesser
	digitalPinsCfg *digitalPinsConfiguration
	translate      digitalPinTranslator
	pins           map[string]gobot.DigitalPinner
	pinOptions     map[string][]func(gobot.DigitalPinOptioner) bool
	mutex          sync.Mutex
}

// NewDigitalPinsAdaptor provides the access to digital pins of the board. It supports sysfs and gpiod system drivers.
// This is decided by the given accesser. The translator is used to adapt the pin header naming, which is given by user,
// to the internal file name or chip/line nomenclature. This varies by each platform. If for some reasons the default
// initializer is not suitable, it can be given by the option "WithDigitalPinInitializer()". This is especially needed,
// if some values needs to be adjusted after the pin was created but before the pin is exported.
func NewDigitalPinsAdaptor(
	sys *system.Accesser,
	t digitalPinTranslator,
	opts ...DigitalPinsOptionApplier,
) *DigitalPinsAdaptor {
	a := &DigitalPinsAdaptor{
		sys:            sys,
		digitalPinsCfg: &digitalPinsConfiguration{initialize: func(pin gobot.DigitalPinner) error { return pin.Export() }},
		translate:      t,
	}
	for _, o := range opts {
		o.apply(a.digitalPinsCfg)
	}
	return a
}

// WithDigitalPinInitializer can be used to substitute the default initializer.
func WithDigitalPinInitializer(pc digitalPinInitializer) digitalPinsInitializeOption {
	return digitalPinsInitializeOption(pc)
}

// WithGpiodAccess can be used to change the default legacy sysfs implementation to the character device Kernel ABI,
// provided by the gpiod package.
func WithGpiodAccess() digitalPinsSystemSysfsOption {
	return digitalPinsSystemSysfsOption(false)
}

// WithSysfsAccess can be used to change the default character device implementation, provided by the gpiod package, to
// the legacy sysfs Kernel ABI.
func WithSysfsAccess() digitalPinsSystemSysfsOption {
	return digitalPinsSystemSysfsOption(true)
}

// WithSpiGpioAccess can be used to switch the default SPI implementation to GPIO usage.
func WithSpiGpioAccess(sclkPin, ncsPin, sdoPin, sdiPin string) digitalPinsForSystemSpiOption {
	o := digitalPinsForSystemSpiOption{
		sclkPin: sclkPin,
		ncsPin:  ncsPin,
		sdoPin:  sdoPin,
		sdiPin:  sdiPin,
	}

	return o
}

// WithGpiosActiveLow prepares the given pins for inverse reaction on next initialize.
// This is working for inputs and outputs.
func WithGpiosActiveLow(pin string, otherPins ...string) digitalPinsActiveLowOption {
	pins := append([]string{pin}, otherPins...)
	return digitalPinsActiveLowOption(pins)
}

// WithGpiosPullDown prepares the given pins to be pulled down (high impedance to GND) on next initialize.
// This is working for inputs and outputs since Kernel 5.5, but will be ignored with sysfs ABI.
func WithGpiosPullDown(pin string, otherPins ...string) digitalPinsPullDownOption {
	pins := append([]string{pin}, otherPins...)
	return digitalPinsPullDownOption(pins)
}

// WithGpiosPullUp prepares the given pins to be pulled up (high impedance to VDD) on next initialize.
// This is working for inputs and outputs since Kernel 5.5, but will be ignored with sysfs ABI.
func WithGpiosPullUp(pin string, otherPins ...string) digitalPinsPullUpOption {
	pins := append([]string{pin}, otherPins...)
	return digitalPinsPullUpOption(pins)
}

// WithGpiosOpenDrain prepares the given output pins to be driven with open drain/collector on next initialize.
// This will be ignored for inputs or with sysfs ABI.
func WithGpiosOpenDrain(pin string, otherPins ...string) digitalPinsOpenDrainOption {
	pins := append([]string{pin}, otherPins...)
	return digitalPinsOpenDrainOption(pins)
}

// WithGpiosOpenSource prepares the given output pins to be driven with open source/emitter on next initialize.
// This will be ignored for inputs or with sysfs ABI.
func WithGpiosOpenSource(pin string, otherPins ...string) digitalPinsOpenSourceOption {
	pins := append([]string{pin}, otherPins...)
	return digitalPinsOpenSourceOption(pins)
}

// WithGpioDebounce prepares the given input pin to be debounced on next initialize.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioDebounce(pin string, period time.Duration) digitalPinsDebounceOption {
	return digitalPinsDebounceOption{id: pin, period: period}
}

// WithGpioEventOnFallingEdge prepares the given input pin to be generate an event on falling edge.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioEventOnFallingEdge(pin string, handler func(lineOffset int, timestamp time.Duration, detectedEdge string,
	seqno uint32, lseqno uint32),
) digitalPinsEventOnFallingEdgeOption {
	return digitalPinsEventOnFallingEdgeOption{id: pin, handler: handler}
}

// WithGpioEventOnRisingEdge prepares the given input pin to be generate an event on rising edge.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioEventOnRisingEdge(pin string, handler func(lineOffset int, timestamp time.Duration, detectedEdge string,
	seqno uint32, lseqno uint32),
) digitalPinsEventOnRisingEdgeOption {
	return digitalPinsEventOnRisingEdgeOption{id: pin, handler: handler}
}

// WithGpioEventOnBothEdges prepares the given input pin to be generate an event on rising and falling edges.
// This is working for inputs since Kernel 5.10, but will be ignored for outputs or with sysfs ABI.
func WithGpioEventOnBothEdges(pin string, handler func(lineOffset int, timestamp time.Duration, detectedEdge string,
	seqno uint32, lseqno uint32),
) digitalPinsEventOnBothEdgesOption {
	return digitalPinsEventOnBothEdgesOption{id: pin, handler: handler}
}

// WithGpioPollForEdgeDetection prepares the given input pin to use a discrete input pin polling function together with
// edge detection.
func WithGpioPollForEdgeDetection(
	pin string,
	pollInterval time.Duration,
	pollQuitChan chan struct{},
) digitalPinsPollForEdgeDetectionOption {
	return digitalPinsPollForEdgeDetectionOption{id: pin, pollInterval: pollInterval, pollQuitChan: pollQuitChan}
}

// Connect prepare new connection to digital pins.
func (a *DigitalPinsAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.pinOptions != nil {
		return fmt.Errorf("digital pin adaptor already connected, please call Finalize() for re-connect")
	}

	cfg := a.digitalPinsCfg

	if cfg.useSysfs != nil {
		if *cfg.useSysfs {
			system.WithDigitalPinSysfsAccess()(a.sys)
		} else {
			system.WithDigitalPinGpiodAccess()(a.sys)
		}
	}

	if cfg.gpiosForSPI != nil {
		system.WithSpiGpioAccess(a, cfg.gpiosForSPI.sclkPin, cfg.gpiosForSPI.ncsPin, cfg.gpiosForSPI.sdoPin,
			cfg.gpiosForSPI.sdiPin)(a.sys)
	}

	a.pinOptions = make(map[string][]func(gobot.DigitalPinOptioner) bool)

	for _, pin := range cfg.activeLowPins {
		a.pinOptions[pin] = append(a.pinOptions[pin], system.WithPinActiveLow())
	}

	for _, pin := range cfg.pullDownPins {
		a.pinOptions[pin] = append(a.pinOptions[pin], system.WithPinPullDown())
	}

	for _, pin := range cfg.pullUpPins {
		a.pinOptions[pin] = append(a.pinOptions[pin], system.WithPinPullUp())
	}

	for _, pin := range cfg.openDrainPins {
		a.pinOptions[pin] = append(a.pinOptions[pin], system.WithPinOpenDrain())
	}

	for _, pin := range cfg.openSourcePins {
		a.pinOptions[pin] = append(a.pinOptions[pin], system.WithPinOpenSource())
	}

	for _, pin := range cfg.debouncePins {
		a.pinOptions[pin.id] = append(a.pinOptions[pin.id], system.WithPinDebounce(pin.period))
	}

	for _, pin := range cfg.eventOnFallingEdgePins {
		a.pinOptions[pin.id] = append(a.pinOptions[pin.id], system.WithPinEventOnFallingEdge(pin.handler))
	}

	for _, pin := range cfg.eventOnRisingEdgePins {
		a.pinOptions[pin.id] = append(a.pinOptions[pin.id], system.WithPinEventOnRisingEdge(pin.handler))
	}

	for _, pin := range cfg.eventOnBothEdgesPins {
		a.pinOptions[pin.id] = append(a.pinOptions[pin.id], system.WithPinEventOnBothEdges(pin.handler))
	}

	for _, pin := range cfg.pollForEdgeDetectionPins {
		a.pinOptions[pin.id] = append(a.pinOptions[pin.id],
			system.WithPinPollForEdgeDetection(pin.pollInterval, pin.pollQuitChan))
	}

	a.pins = make(map[string]gobot.DigitalPinner)
	return nil
}

// Finalize closes connection to digital pins
func (a *DigitalPinsAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var err error
	for _, pin := range a.pins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	a.pins = nil
	a.pinOptions = nil
	return err
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

func (a *DigitalPinsAdaptor) digitalPin(
	id string,
	opts ...func(gobot.DigitalPinOptioner) bool,
) (gobot.DigitalPinner, error) {
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
		if err = a.digitalPinsCfg.initialize(pin); err != nil {
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
