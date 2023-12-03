package adaptors

import "time"

// pwmPinOptionApplier needs to be implemented by each configurable option type
type PwmPinsOptionApplier interface {
	apply(cfg *pwmPinsConfiguration)
}

// pwmPinInitializeOption is the type for applying another than the default initializer.
type pwmPinsInitializeOption pwmPinInitializer

// pwmPinPeriodDefaultOption is the type for applying another than the default period of 10 ms (100 Hz) for all
// created pins.
type pwmPinsPeriodDefaultOption uint32

// pwmPinPolarityInvertedIdentifierOption is the type for applying another identifier, which will replace the default
// "inverted".
type pwmPinsPolarityInvertedIdentifierOption string

// pwmPinsAdjustDutyOnSetPeriodOption is the type for applying the automatic adjustment of duty cycle on setting
// the period to on/off.
type pwmPinsAdjustDutyOnSetPeriodOption bool

// pwmPinsDefaultPeriodForPinOption is the type for applying another than the default period of 10 ms (100 Hz) only for
// the given pin id.
type pwmPinsDefaultPeriodForPinOption struct {
	id     string
	period uint32
}

// pwmPinsServoDutyScaleForPinOption is the type for applying another than the default 0.5-2.5 ms range of duty cycle
// for servo calls on the specified pin id.
type pwmPinsServoDutyScaleForPinOption struct {
	id  string
	min time.Duration
	max time.Duration
}

// pwmPinsServoAngleScaleForPinOption is the type for applying another than the default 0.0-180.0° range of angle for
// servo calls on the specified pin id.
type pwmPinsServoAngleScaleForPinOption struct {
	id        string
	minDegree float64
	maxDegree float64
}

func (o pwmPinsInitializeOption) String() string {
	return "pin initializer option for PWM's"
}

func (o pwmPinsPeriodDefaultOption) String() string {
	return "default period option for PWM's"
}

func (o pwmPinsPolarityInvertedIdentifierOption) String() string {
	return "inverted identifier option for PWM's"
}

func (o pwmPinsAdjustDutyOnSetPeriodOption) String() string {
	return "adjust duty cycle on set period option for PWM's"
}

func (o pwmPinsDefaultPeriodForPinOption) String() string {
	return "default period for the pin option for PWM's"
}

func (o pwmPinsServoDutyScaleForPinOption) String() string {
	return "duty cycle min-max range for a servo pin option for PWM's"
}

func (o pwmPinsServoAngleScaleForPinOption) String() string {
	return "angle min-max range for a servo pin option for PWM's"
}

func (o pwmPinsInitializeOption) apply(cfg *pwmPinsConfiguration) {
	cfg.initialize = pwmPinInitializer(o)
}

func (o pwmPinsPeriodDefaultOption) apply(cfg *pwmPinsConfiguration) {
	cfg.periodDefault = uint32(o)
}

func (o pwmPinsPolarityInvertedIdentifierOption) apply(cfg *pwmPinsConfiguration) {
	cfg.polarityInvertedIdentifier = string(o)
}

func (o pwmPinsAdjustDutyOnSetPeriodOption) apply(cfg *pwmPinsConfiguration) {
	cfg.adjustDutyOnSetPeriod = bool(o)
}

func (o pwmPinsDefaultPeriodForPinOption) apply(cfg *pwmPinsConfiguration) {
	cfg.pinsDefaultPeriod[o.id] = o.period
}

func (o pwmPinsServoDutyScaleForPinOption) apply(cfg *pwmPinsConfiguration) {
	scale, ok := cfg.pinsServoScale[o.id]
	if !ok {
		scale = pwmPinServoScale{minDegree: 0, maxDegree: 180}
	}

	scale.minDuty = o.min
	scale.maxDuty = o.max

	cfg.pinsServoScale[o.id] = scale
}

func (o pwmPinsServoAngleScaleForPinOption) apply(cfg *pwmPinsConfiguration) {
	scale, ok := cfg.pinsServoScale[o.id]
	if !ok {
		scale = pwmPinServoScale{} // default values for duty cycle will be set on initialize, if zero
	}

	scale.minDegree = o.minDegree
	scale.maxDegree = o.maxDegree

	cfg.pinsServoScale[o.id] = scale
}
