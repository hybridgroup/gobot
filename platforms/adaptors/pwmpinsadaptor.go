package adaptors

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

// note for period in nano seconds:
// 100000000ns = 100ms = 10Hz, 10000000ns = 10ms = 100Hz,  1000000ns = 1ms = 1kHz,
// 100000ns = 100us = 10kHz, 10000ns = 10us = 100kHz, 1000ns = 1us = 1MHz,
// 100ns = 10MHz, 10ns = 100MHz, 1ns = 1GHz
const pwmPeriodDefault = 10000000 // 10 ms = 100 Hz

type (
	pwmPinTranslator  func(pin string) (path string, channel int, err error)
	pwmPinInitializer func(gobot.PWMPinner) error
)

type pwmPinsOption interface {
	setInitializer(pwmPinInitializer)
	setDefaultPeriod(uint32)
	setPolarityInvertedIdentifier(string)
}

// PWMPinsAdaptor is a adaptor for PWM pins, normally used for composition in platforms.
type PWMPinsAdaptor struct {
	sys                        *system.Accesser
	translate                  pwmPinTranslator
	initialize                 pwmPinInitializer
	periodDefault              uint32
	polarityNormalIdentifier   string
	polarityInvertedIdentifier string
	adjustDutyOnSetPeriod      bool
	pins                       map[string]gobot.PWMPinner
	mutex                      sync.Mutex
}

// NewPWMPinsAdaptor provides the access to PWM pins of the board. It uses sysfs system drivers. The translator is used
// to adapt the pin header naming, which is given by user, to the internal file name nomenclature. This varies by each
// platform. If for some reasons the default initializer is not suitable, it can be given by the option
// "WithPWMPinInitializer()".
func NewPWMPinsAdaptor(sys *system.Accesser, t pwmPinTranslator, options ...func(pwmPinsOption)) *PWMPinsAdaptor {
	a := &PWMPinsAdaptor{
		sys:                        sys,
		translate:                  t,
		periodDefault:              pwmPeriodDefault,
		polarityNormalIdentifier:   "normal",
		polarityInvertedIdentifier: "inverted",
		adjustDutyOnSetPeriod:      true,
	}
	a.initialize = a.getDefaultInitializer()
	for _, option := range options {
		option(a)
	}
	return a
}

// WithPWMPinInitializer substitute the default initializer.
func WithPWMPinInitializer(pc pwmPinInitializer) func(pwmPinsOption) {
	return func(a pwmPinsOption) {
		a.setInitializer(pc)
	}
}

// WithPWMPinDefaultPeriod substitute the default period of 10 ms (100 Hz) for all created pins.
func WithPWMPinDefaultPeriod(periodNanoSec uint32) func(pwmPinsOption) {
	return func(a pwmPinsOption) {
		a.setDefaultPeriod(periodNanoSec)
	}
}

// WithPolarityInvertedIdentifier use the given identifier, which will replace the default "inverted".
func WithPolarityInvertedIdentifier(identifier string) func(pwmPinsOption) {
	return func(a pwmPinsOption) {
		a.setPolarityInvertedIdentifier(identifier)
	}
}

// Connect prepare new connection to PWM pins.
func (a *PWMPinsAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.pins = make(map[string]gobot.PWMPinner)
	return nil
}

// Finalize closes connection to PWM pins.
func (a *PWMPinsAdaptor) Finalize() (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, pin := range a.pins {
		if pin != nil {
			if errs := pin.SetEnabled(false); errs != nil {
				err = multierror.Append(err, errs)
			}
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	a.pins = nil
	return err
}

// PwmWrite writes a PWM signal to the specified pin.
func (a *PWMPinsAdaptor) PwmWrite(id string, val byte) (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.pwmPin(id)
	if err != nil {
		return
	}
	period, err := pin.Period()
	if err != nil {
		return err
	}
	duty := gobot.FromScale(float64(val), 0, 255.0)
	return pin.SetDutyCycle(uint32(float64(period) * duty))
}

// ServoWrite writes a servo signal to the specified pin.
func (a *PWMPinsAdaptor) ServoWrite(id string, angle byte) (err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.pwmPin(id)
	if err != nil {
		return
	}
	period, err := pin.Period()
	if err != nil {
		return err
	}

	// 0.5 ms => -90
	// 1.5 ms =>   0
	// 2.0 ms =>  90
	minDuty := 100 * 0.0005 * float64(period)
	maxDuty := 100 * 0.0020 * float64(period)
	duty := uint32(gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty))
	return pin.SetDutyCycle(duty)
}

// SetPeriod adjusts the period of the specified PWM pin immediately.
// If duty cycle is already set, also this value will be adjusted in the same ratio.
func (a *PWMPinsAdaptor) SetPeriod(id string, period uint32) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.pwmPin(id)
	if err != nil {
		return err
	}
	return setPeriod(pin, period, a.adjustDutyOnSetPeriod)
}

// PWMPin initializes the pin for PWM and returns matched pwmPin for specified pin number.
// It implements the PWMPinnerProvider interface.
func (a *PWMPinsAdaptor) PWMPin(id string) (gobot.PWMPinner, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.pwmPin(id)
}

func (a *PWMPinsAdaptor) setInitializer(pinInit pwmPinInitializer) {
	a.initialize = pinInit
}

func (a *PWMPinsAdaptor) setDefaultPeriod(periodNanoSec uint32) {
	a.periodDefault = periodNanoSec
}

func (a *PWMPinsAdaptor) setPolarityInvertedIdentifier(identifier string) {
	a.polarityInvertedIdentifier = identifier
}

func (a *PWMPinsAdaptor) getDefaultInitializer() func(gobot.PWMPinner) error {
	return func(pin gobot.PWMPinner) error {
		if err := pin.Export(); err != nil {
			return err
		}
		// Make sure PWM is disabled before change anything (period needs to be >0 for this check)
		if period, _ := pin.Period(); period > 0 {
			if err := pin.SetEnabled(false); err != nil {
				return err
			}
		}
		if err := setPeriod(pin, a.periodDefault, a.adjustDutyOnSetPeriod); err != nil {
			return err
		}
		// period needs to be set >1 before all next statements
		if err := pin.SetPolarity(true); err != nil {
			return err
		}
		return pin.SetEnabled(true)
	}
}

func (a *PWMPinsAdaptor) pwmPin(id string) (gobot.PWMPinner, error) {
	if a.pins == nil {
		return nil, fmt.Errorf("not connected")
	}

	pin := a.pins[id]

	if pin == nil {
		path, channel, err := a.translate(id)
		if err != nil {
			return nil, err
		}
		pin = a.sys.NewPWMPin(path, channel, a.polarityNormalIdentifier, a.polarityInvertedIdentifier)
		if err := a.initialize(pin); err != nil {
			return nil, err
		}
		a.pins[id] = pin
	}

	return pin, nil
}

// setPeriod adjusts the PWM period of the given pin. If duty cycle is already set and this feature is not suppressed,
// also this value will be adjusted in the same ratio. The order of writing the values must be observed, otherwise an
// error occur "write error: Invalid argument".
func setPeriod(pin gobot.PWMPinner, period uint32, adjustDuty bool) error {
	errorBase := fmt.Sprintf("setPeriod(%v, %d) failed", pin, period)

	var oldDuty uint32
	var err error
	if adjustDuty {
		if oldDuty, err = pin.DutyCycle(); err != nil {
			return fmt.Errorf("%s with '%v'", errorBase, err)
		}
	}

	if oldDuty == 0 {
		if err := pin.SetPeriod(period); err != nil {
			return fmt.Errorf("%s with '%v'", errorBase, err)
		}
	} else {
		// adjust duty cycle in the same ratio
		oldPeriod, err := pin.Period()
		if err != nil {
			return fmt.Errorf("%s with '%v'", errorBase, err)
		}
		duty := uint32(uint64(oldDuty) * uint64(period) / uint64(oldPeriod))

		// the order depends on value (duty must not be bigger than period in any situation)
		if duty > oldPeriod {
			if err := pin.SetPeriod(period); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
			if err := pin.SetDutyCycle(duty); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
		} else {
			if err := pin.SetDutyCycle(duty); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
			if err := pin.SetPeriod(period); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
		}
	}
	return nil
}
