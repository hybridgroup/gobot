package gpio

import (
	"fmt"

	"gobot.io/x/gobot/v2"
)

// motorOptionApplier needs to be implemented by each configurable option type
type motorOptionApplier interface {
	apply(cfg *motorConfiguration)
}

// motorConfiguration contains all changeable attributes of the driver.
type motorConfiguration struct {
	modeIsAnalog bool
	directionPin string
	forwardPin   string
	backwardPin  string
}

// motorModeIsAnalogOption is the type for applying analog mode to the configuration
type motorModeIsAnalogOption bool

// motorDirectionPinOption is the type for applying a direction pin to the configuration
type motorDirectionPinOption string

// motorForwardPinOption is the type for applying a forward pin to the configuration
type motorForwardPinOption string

// motorBackwardPinOption is the type for applying a backward pin to the configuration
type motorBackwardPinOption string

// MotorDriver Represents a Motor
type MotorDriver struct {
	*driver
	motorCfg         *motorConfiguration
	currentState     byte
	currentSpeed     byte
	currentDirection string
}

// NewMotorDriver return a new MotorDriver given a DigitalWriter and pin. This defaults to digital mode and just switch
// on and off in forward direction. Optional pins can be given, depending on your hardware. So the direction can be
// changed with one pin or by using separated forward and backward pins.
//
// If the given pin supports the PwmWriter the motor can be used/switched to analog mode by writing once to SetSpeed()
// or by calling SetAnalogMode(). The optional pins can be used for direction control.
//
// Supported options:
//
//	"WithName"
//	"WithMotorAnalog"
//	"WithMotorDirectionPin"
//	"WithMotorForwardPin"
//	"WithMotorBackwardPin"
func NewMotorDriver(a DigitalWriter, speedPin string, opts ...interface{}) *MotorDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &MotorDriver{
		driver:           newDriver(a.(gobot.Connection), "Motor", withPin(speedPin)),
		motorCfg:         &motorConfiguration{},
		currentDirection: "forward",
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case motorOptionApplier:
			o.apply(d.motorCfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return d
}

// WithMotorAnalog change the default mode "digital" to analog for the motor.
func WithMotorAnalog() motorOptionApplier {
	return motorModeIsAnalogOption(true)
}

// WithMotorDirectionPin introduces a pin for change the direction of the motor.
func WithMotorDirectionPin(pin string) motorOptionApplier {
	return motorDirectionPinOption(pin)
}

// WithMotorForwardPin introduces a pin for setting the direction to forward.
func WithMotorForwardPin(pin string) motorOptionApplier {
	return motorForwardPinOption(pin)
}

// WithMotorBackwardPin introduces a pin for setting the direction to backward.
func WithMotorBackwardPin(pin string) motorOptionApplier {
	return motorBackwardPinOption(pin)
}

// Off turns the motor off or sets the motor to a 0 speed.
func (d *MotorDriver) Off() error {
	if d.IsDigital() {
		return d.changeState(0)
	}

	return d.SetSpeed(0)
}

// On turns the motor on or sets the motor to a maximum speed.
func (d *MotorDriver) On() error {
	if d.IsDigital() {
		return d.changeState(1)
	}

	if d.currentSpeed == 0 {
		d.currentSpeed = 255
	}

	return d.SetSpeed(d.currentSpeed)
}

// RunMin sets the motor to the minimum speed.
func (d *MotorDriver) RunMin() error {
	return d.Off()
}

// RunMax sets the motor to the maximum speed.
func (d *MotorDriver) RunMax() error {
	return d.SetSpeed(255)
}

// Toggle sets the motor to the opposite of it's current state.
func (d *MotorDriver) Toggle() error {
	if d.IsOn() {
		return d.Off()
	}

	return d.On()
}

// SetSpeed change the speed of the motor, without change the direction.
func (d *MotorDriver) SetSpeed(value byte) error {
	if writer, ok := d.connection.(PwmWriter); ok {
		WithMotorAnalog().apply(d.motorCfg)
		d.currentSpeed = value
		return writer.PwmWrite(d.driverCfg.pin, value)
	}
	return ErrPwmWriteUnsupported
}

// Forward runs the motor forward with the specified speed.
func (d *MotorDriver) Forward(speed byte) error {
	if err := d.SetDirection("forward"); err != nil {
		return err
	}
	if err := d.SetSpeed(speed); err != nil {
		return err
	}

	return nil
}

// Backward runs the motor backward with the specified speed.
func (d *MotorDriver) Backward(speed byte) error {
	if err := d.SetDirection("backward"); err != nil {
		return err
	}
	if err := d.SetSpeed(speed); err != nil {
		return err
	}

	return nil
}

// Direction sets the direction pin to the specified direction.
func (d *MotorDriver) SetDirection(direction string) error {
	d.currentDirection = direction
	if d.motorCfg.directionPin != "" {
		var level byte
		if direction == "forward" {
			level = 1
		} else {
			level = 0
		}
		return d.digitalWrite(d.motorCfg.directionPin, level)
	}

	var forwardLevel, backwardLevel byte
	switch direction {
	case "forward":
		forwardLevel = 1
		backwardLevel = 0
	case "backward":
		forwardLevel = 0
		backwardLevel = 1
	case "none":
		forwardLevel = 0
		backwardLevel = 0
	}

	if d.motorCfg.forwardPin != "" {
		if err := d.digitalWrite(d.motorCfg.forwardPin, forwardLevel); err != nil {
			return err
		}
	}

	if d.motorCfg.backwardPin != "" {
		return d.digitalWrite(d.motorCfg.backwardPin, backwardLevel)
	}

	return nil
}

// IsAnalog returns true if the motor is in analog mode.
func (d *MotorDriver) IsAnalog() bool {
	return d.motorCfg.modeIsAnalog
}

// IsDigital returns true if the motor is in digital mode.
func (d *MotorDriver) IsDigital() bool {
	return !d.motorCfg.modeIsAnalog
}

// IsOn returns true if the motor is on.
func (d *MotorDriver) IsOn() bool {
	if d.IsDigital() {
		return d.currentState == 1
	}
	return d.currentSpeed > 0
}

// IsOff returns true if the motor is off.
func (d *MotorDriver) IsOff() bool {
	return !d.IsOn()
}

// Direction returns the current direction ("forward" or "backward") of the motor.
func (d *MotorDriver) Direction() string {
	return d.currentDirection
}

// Speed returns the current speed of the motor.
func (d *MotorDriver) Speed() byte {
	return d.currentSpeed
}

func (d *MotorDriver) changeState(state byte) error {
	d.currentState = state
	if state == 1 {
		d.currentSpeed = 255
	} else {
		d.currentSpeed = 0
	}

	if d.motorCfg.forwardPin == "" {
		return d.digitalWrite(d.driverCfg.pin, state)
	}

	if state != 1 {
		return d.SetDirection("none")
	}

	if err := d.SetDirection(d.currentDirection); err != nil {
		return err
	}
	if d.driverCfg.pin != "" {
		if err := d.SetSpeed(d.currentSpeed); err != nil {
			return err
		}
	}

	return nil
}

func (o motorModeIsAnalogOption) String() string {
	return "motor mode (analog, digital) option"
}

func (o motorDirectionPinOption) String() string {
	return "direction pin option for motors"
}

func (o motorForwardPinOption) String() string {
	return "forward pin option for motors"
}

func (o motorBackwardPinOption) String() string {
	return "backward pin option for motors"
}

func (o motorModeIsAnalogOption) apply(cfg *motorConfiguration) {
	cfg.modeIsAnalog = bool(o)
}

func (o motorDirectionPinOption) apply(cfg *motorConfiguration) {
	cfg.directionPin = string(o)
}

func (o motorForwardPinOption) apply(cfg *motorConfiguration) {
	cfg.forwardPin = string(o)
}

func (o motorBackwardPinOption) apply(cfg *motorConfiguration) {
	cfg.backwardPin = string(o)
}
