package adaptors

import (
	"fmt"
	"log"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

// note for period in nano seconds:
// 100000000ns = 100ms = 10Hz, 10000000ns = 10ms = 100Hz,  1000000ns = 1ms = 1kHz,
// 100000ns = 100us = 10kHz, 10000ns = 10us = 100kHz, 1000ns = 1us = 1MHz,
// 100ns = 10MHz, 10ns = 100MHz, 1ns = 1GHz
const pwmPeriodDefault = 10000000 // 10 ms = 100 Hz

// 50Hz = 0.02 sec = 20 ms
const fiftyHzNanos = 20 * 1000 * 1000

type (
	pwmPinTranslator  func(pin string) (path string, channel int, err error)
	pwmPinInitializer func(id string, pin gobot.PWMPinner) error
)

type pwmPinServoScale struct {
	minDegree, maxDegree float64
	minDuty, maxDuty     time.Duration
}

// pwmPinsConfiguration contains all changeable attributes of the adaptor.
type pwmPinsConfiguration struct {
	initialize                 pwmPinInitializer
	usePiBlasterPin            bool
	periodDefault              uint32
	periodMinimum              uint32
	dutyRateMinimum            float64 // is the minimal relation of duty/period (except 0.0)
	polarityNormalIdentifier   string
	polarityInvertedIdentifier string
	adjustDutyOnSetPeriod      bool
	pinsDefaultPeriod          map[string]uint32           // the key is the pin id
	pinsServoScale             map[string]pwmPinServoScale // the key is the pin id
}

// PWMPinsAdaptor is a adaptor for PWM pins, normally used for composition in platforms.
type PWMPinsAdaptor struct {
	sys        *system.Accesser
	translate  pwmPinTranslator
	pwmPinsCfg *pwmPinsConfiguration
	pins       map[string]gobot.PWMPinner
	mutex      sync.Mutex
}

// NewPWMPinsAdaptor provides the access to PWM pins of the board. It uses sysfs system drivers. The translator is used
// to adapt the pin header naming, which is given by user, to the internal file name nomenclature. This varies by each
// platform. If for some reasons the default initializer is not suitable, it can be given by the option
// "WithPWMPinInitializer()".
//
// Further options:
//
//	"WithPWMDefaultPeriod"
//	"WithPWMPolarityInvertedIdentifier"
//	"WithPWMNoDutyCycleAdjustment"
//	"WithPWMDefaultPeriodForPin"
//	"WithPWMServoDutyCycleRangeForPin"
//	"WithPWMServoAngleRangeForPin"
func NewPWMPinsAdaptor(sys *system.Accesser, t pwmPinTranslator, opts ...PwmPinsOptionApplier) *PWMPinsAdaptor {
	a := &PWMPinsAdaptor{
		sys:       sys,
		translate: t,
		pwmPinsCfg: &pwmPinsConfiguration{
			periodDefault:              pwmPeriodDefault,
			pinsDefaultPeriod:          make(map[string]uint32),
			pinsServoScale:             make(map[string]pwmPinServoScale),
			polarityNormalIdentifier:   "normal",
			polarityInvertedIdentifier: "inversed",
			adjustDutyOnSetPeriod:      true,
		},
	}
	a.pwmPinsCfg.initialize = a.getDefaultInitializer()

	for _, o := range opts {
		o.apply(a.pwmPinsCfg)
	}

	return a
}

// WithPWMPinInitializer substitute the default initializer.
func WithPWMPinInitializer(pc pwmPinInitializer) pwmPinsInitializeOption {
	return pwmPinsInitializeOption(pc)
}

// WithPWMUsePiBlaster substitute the default sysfs-implementation for PWM-pins by the implementation for pi-blaster.
func WithPWMUsePiBlaster() pwmPinsUsePiBlasterPinOption {
	return pwmPinsUsePiBlasterPinOption(true)
}

// WithPWMDefaultPeriod substitute the default period of 10 ms (100 Hz) for all created pins.
func WithPWMDefaultPeriod(periodNanoSec uint32) pwmPinsPeriodDefaultOption {
	return pwmPinsPeriodDefaultOption(periodNanoSec)
}

// WithPWMMinimumPeriod substitute the default minimum period limit of 0 nanoseconds.
func WithPWMMinimumPeriod(periodNanoSec uint32) pwmPinsPeriodMinimumOption {
	return pwmPinsPeriodMinimumOption(periodNanoSec)
}

// WithPWMMinimumDutyRate substitute the default minimum duty rate of 1/period. The given limit only come into effect,
// if the rate is > 0, because a rate of 0.0 is always allowed.
func WithPWMMinimumDutyRate(dutyRate float64) pwmPinsDutyRateMinimumOption {
	return pwmPinsDutyRateMinimumOption(dutyRate)
}

// WithPWMPolarityInvertedIdentifier use the given identifier, which will replace the default "inversed".
func WithPWMPolarityInvertedIdentifier(identifier string) pwmPinsPolarityInvertedIdentifierOption {
	return pwmPinsPolarityInvertedIdentifierOption(identifier)
}

// WithPWMNoDutyCycleAdjustment switch off the automatic adjustment of duty cycle on setting the period.
func WithPWMNoDutyCycleAdjustment() pwmPinsAdjustDutyOnSetPeriodOption {
	return pwmPinsAdjustDutyOnSetPeriodOption(false)
}

// WithPWMDefaultPeriodForPin substitute the default period of 10 ms (100 Hz) for the given pin.
// This option also overrides a default period given by the WithPWMDefaultPeriod() option.
// This is often needed for servo applications, where the default period is 50Hz (20.000.000 ns).
func WithPWMDefaultPeriodForPin(pin string, periodNanoSec uint32) pwmPinsDefaultPeriodForPinOption {
	o := pwmPinsDefaultPeriodForPinOption{id: pin, period: periodNanoSec}
	return o
}

// WithPWMServoDutyCycleRangeForPin set new values for range of duty cycle for servo calls, which replaces the default
// 0.5-2.5 ms range. The given duration values will be internally converted to nanoseconds.
func WithPWMServoDutyCycleRangeForPin(pin string, minimum, maximum time.Duration) pwmPinsServoDutyScaleForPinOption {
	return pwmPinsServoDutyScaleForPinOption{id: pin, min: minimum, max: maximum}
}

// WithPWMServoAngleRangeForPin set new values for range of angle for servo calls, which replaces
// the default 0.0-180.0° range.
func WithPWMServoAngleRangeForPin(pin string, minimum, maximum float64) pwmPinsServoAngleScaleForPinOption {
	return pwmPinsServoAngleScaleForPinOption{id: pin, minDegree: minimum, maxDegree: maximum}
}

// Connect prepare new connection to PWM pins.
func (a *PWMPinsAdaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.pins = make(map[string]gobot.PWMPinner)

	if a.pwmPinsCfg.dutyRateMinimum == 0 && a.pwmPinsCfg.periodDefault > 0 {
		a.pwmPinsCfg.dutyRateMinimum = 1 / float64(a.pwmPinsCfg.periodDefault)
	}

	return nil
}

// Finalize closes connection to PWM pins.
func (a *PWMPinsAdaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var err error
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

// PwmWrite writes a PWM signal to the specified pin. The given value is between 0 and 255.
func (a *PWMPinsAdaptor) PwmWrite(id string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.pwmPin(id)
	if err != nil {
		return err
	}
	periodNanos, err := pin.Period()
	if err != nil {
		return err
	}

	dutyNanos := float64(periodNanos) * gobot.FromScale(float64(val), 0, 255.0)

	if err := a.validateDutyCycle(id, dutyNanos, float64(periodNanos)); err != nil {
		return err
	}

	return pin.SetDutyCycle(uint32(dutyNanos))
}

// ServoWrite writes a servo signal to the specified pin. The given angle is between 0 and 180°.
func (a *PWMPinsAdaptor) ServoWrite(id string, angle byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	pin, err := a.pwmPin(id)
	if err != nil {
		return err
	}
	periodNanos, err := pin.Period()
	if err != nil {
		return err
	}

	if periodNanos != fiftyHzNanos {
		log.Printf("WARNING: the PWM acts with a period of %d, but should use %d (50Hz) for servos\n",
			periodNanos, fiftyHzNanos)
	}

	scale, ok := a.pwmPinsCfg.pinsServoScale[id]
	if !ok {
		return fmt.Errorf("no scaler found for servo pin '%s'", id)
	}

	dutyNanos := gobot.ToScale(gobot.FromScale(float64(angle),
		scale.minDegree, scale.maxDegree),
		float64(scale.minDuty), float64(scale.maxDuty))

	if err := a.validateDutyCycle(id, dutyNanos, float64(periodNanos)); err != nil {
		return err
	}

	return pin.SetDutyCycle(uint32(dutyNanos))
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
	return setPeriod(pin, period, a.pwmPinsCfg.adjustDutyOnSetPeriod)
}

// PWMPin initializes the pin for PWM and returns matched pwmPin for specified pin number.
// It implements the PWMPinnerProvider interface.
func (a *PWMPinsAdaptor) PWMPin(id string) (gobot.PWMPinner, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.pwmPin(id)
}

func (a *PWMPinsAdaptor) getDefaultInitializer() func(string, gobot.PWMPinner) error {
	return func(id string, pin gobot.PWMPinner) error {
		if err := pin.Export(); err != nil {
			return err
		}
		// Make sure PWM is disabled before change anything (period needs to be >0 for this check)
		if period, _ := pin.Period(); period > 0 {
			if err := pin.SetEnabled(false); err != nil {
				return err
			}
		}

		// looking for a pin specific period
		defaultPeriod, ok := a.pwmPinsCfg.pinsDefaultPeriod[id]
		if !ok {
			defaultPeriod = a.pwmPinsCfg.periodDefault
		}

		if err := setPeriod(pin, defaultPeriod, a.pwmPinsCfg.adjustDutyOnSetPeriod); err != nil {
			return err
		}

		// ensure servo scaler is present
		//
		// usually for the most servos (at 50Hz) for 180° (SG90, AD002)
		// 0.5 ms =>   0 (1/40 part of period at 50Hz)
		// 1.5 ms =>  90
		// 2.5 ms => 180 (1/8 part of period at 50Hz)
		scale, ok := a.pwmPinsCfg.pinsServoScale[id]
		if !ok {
			scale = pwmPinServoScale{
				minDegree: 0,
				maxDegree: 180,
			}
		}
		if scale.minDuty == 0 {
			scale.minDuty = time.Duration(defaultPeriod / 40)
		}
		if scale.maxDuty == 0 {
			scale.maxDuty = time.Duration(defaultPeriod / 8)
		}
		a.pwmPinsCfg.pinsServoScale[id] = scale

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

		if a.pwmPinsCfg.usePiBlasterPin {
			pin = newPiBlasterPWMPin(a.sys, channel)
		} else {
			pin = a.sys.NewPWMPin(path, channel, a.pwmPinsCfg.polarityNormalIdentifier,
				a.pwmPinsCfg.polarityInvertedIdentifier)
		}
		if err := a.pwmPinsCfg.initialize(id, pin); err != nil {
			return nil, err
		}
		a.pins[id] = pin
	}

	return pin, nil
}

func (a *PWMPinsAdaptor) validateDutyCycle(id string, dutyNanos, periodNanos float64) error {
	if periodNanos == 0 {
		return nil
	}

	if dutyNanos > periodNanos {
		return fmt.Errorf("duty cycle (%d) exceeds period (%d) for PWM pin id '%s'",
			uint32(dutyNanos), uint32(periodNanos), id)
	}

	if dutyNanos == 0 {
		return nil
	}

	rate := dutyNanos / periodNanos
	if rate < a.pwmPinsCfg.dutyRateMinimum {
		return fmt.Errorf("duty rate (%.8f) is lower than allowed (%.8f) for PWM pin id '%s'",
			rate, a.pwmPinsCfg.dutyRateMinimum, id)
	}
	return nil
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

		//nolint:gosec // TODO: fix later
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
