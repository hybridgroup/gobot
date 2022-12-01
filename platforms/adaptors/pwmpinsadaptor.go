package adaptors

import (
	"fmt"
	"log"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/system"
)

type pwmPinTranslator func(pin string) (path string, channel int, err error)
type pwmPinCreator func(chip string, line int) gobot.PWMPinner

type pwmPinsOption interface {
	setPWMPinCreator(pwmPinCreator)
}

// PWMPinsAdaptor is a adaptor for PWM pins, normally used for composition in platforms.
type PWMPinsAdaptor struct {
	sys              *system.Accesser
	translate        pwmPinTranslator
	create           pwmPinCreator
	periodDefault    uint32
	polarityNormal   string
	polarityInverted string
	pins             map[string]gobot.PWMPinner
	mutex            sync.Mutex
}

// NewPWMPinsAdaptor provides the access to PWM pins of the board. It uses sysfs system drivers. The translator is used
// to adapt the pin header naming, which is given by user, to the internal file name nomenclature. This varies by each
// platform. If for some reasons the default creator is not suitable, it can be given by the option
// "WithPWMPinCreator()". This is especially needed, if some values needs to be adjusted after the pin was created.
func NewPWMPinsAdaptor(sys *system.Accesser, t pwmPinTranslator, options ...func(pwmPinsOption)) *PWMPinsAdaptor {
	s := &PWMPinsAdaptor{
		translate:        t,
		create:           sys.NewPWMPin,
		periodDefault:    10000000, // 10ms = 100Hz
		polarityNormal:   "normal",
		polarityInverted: "inversed",
	}
	for _, option := range options {
		option(s)
	}
	return s
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
			if errs := pin.Enable(false); errs != nil {
				err = multierror.Append(err, errs)
			}
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
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

// SetPeriod adjusts the period of the specified PWM pin.
// If duty cycle is already set, also this value will be adjusted in the same ratio.
func (a *PWMPinsAdaptor) SetPeriod(id string, period uint32) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.pwmPin(id)
	if err != nil {
		return err
	}
	return setPeriod(pin, period)
}

// PWMPin initializes the pin for PWM and returns matched pwmPin for specified pin number.
// It implements the PWMPinnerProvider interface.
func (a *PWMPinsAdaptor) PWMPin(id string) (gobot.PWMPinner, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.pwmPin(id)
}

func (a *PWMPinsAdaptor) setPWMPinCreator(pc pwmPinCreator) {
	a.create = pc
}

func (a *PWMPinsAdaptor) pwmPin(id string) (gobot.PWMPinner, error) {
	pin := a.pins[id]

	if pin == nil {
		path, channel, err := a.translate(id)
		if err != nil {
			return nil, err
		}
		pin = a.create(path, channel)
		if err := pin.Export(); err != nil {
			return nil, err
		}
		// Make sure pwm is disabled before change anything
		if err := pin.Enable(false); err != nil {
			return nil, err
		}
		if err := setPeriod(pin, a.periodDefault); err != nil {
			return nil, err
		}
		if err := pin.SetPolarity(a.polarityNormal); err != nil {
			return nil, err
		}
		if err := pin.Enable(true); err != nil {
			return nil, err
		}
		a.pins[id] = pin
	}

	return pin, nil
}

// setPeriod adjusts the PWM period of the given pin.
// If duty cycle is already set, also this value will be adjusted in the same ratio.
// The order in which the values are written must be observed, otherwise an error occur "write error: Invalid argument".
func setPeriod(pin gobot.PWMPinner, period uint32) error {
	var errorBase = fmt.Sprintf("setPeriod(%v, %d) failed", pin, period)
	oldDuty, err := pin.DutyCycle()
	if err != nil {
		return fmt.Errorf("%s with '%v'", errorBase, err)
	}

	if oldDuty == 0 {
		if err := pin.SetPeriod(period); err != nil {
			log.Println(1, period)
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
			if err := pin.SetDutyCycle(uint32(duty)); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
		} else {
			if err := pin.SetDutyCycle(uint32(duty)); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
			if err := pin.SetPeriod(period); err != nil {
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
		}
	}
	return nil
}
