package gpio

import (
	"log"

	"gobot.io/x/gobot/v2"
)

// MotorDriver Represents a Motor
type MotorDriver struct {
	name             string
	connection       DigitalWriter
	SpeedPin         string
	SwitchPin        string
	DirectionPin     string
	ForwardPin       string
	BackwardPin      string
	CurrentState     byte
	CurrentSpeed     byte
	CurrentMode      string
	CurrentDirection string
}

// NewMotorDriver return a new MotorDriver given a DigitalWriter and pin
func NewMotorDriver(a DigitalWriter, speedPin string) *MotorDriver {
	return &MotorDriver{
		name:             gobot.DefaultName("Motor"),
		connection:       a,
		SpeedPin:         speedPin,
		CurrentState:     0,
		CurrentSpeed:     0,
		CurrentMode:      "digital",
		CurrentDirection: "forward",
	}
}

// Name returns the MotorDrivers name
func (d *MotorDriver) Name() string { return d.name }

// SetName sets the MotorDrivers name
func (d *MotorDriver) SetName(n string) { d.name = n }

// Connection returns the MotorDrivers Connection
func (d *MotorDriver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// Start implements the Driver interface
func (d *MotorDriver) Start() error { return nil }

// Halt implements the Driver interface
func (d *MotorDriver) Halt() error { return nil }

// Off turns the motor off or sets the motor to a 0 speed
func (d *MotorDriver) Off() error {
	if d.isDigital() {
		return d.changeState(0)
	}

	return d.Speed(0)
}

// On turns the motor on or sets the motor to a maximum speed
func (d *MotorDriver) On() error {
	if d.isDigital() {
		return d.changeState(1)
	}
	if d.CurrentSpeed == 0 {
		d.CurrentSpeed = 255
	}

	return d.Speed(d.CurrentSpeed)
}

// Min sets the motor to the minimum speed
func (d *MotorDriver) Min() error {
	return d.Off()
}

// Max sets the motor to the maximum speed
func (d *MotorDriver) Max() error {
	return d.Speed(255)
}

// IsOn returns true if the motor is on
func (d *MotorDriver) IsOn() bool {
	if d.isDigital() {
		return d.CurrentState == 1
	}
	return d.CurrentSpeed > 0
}

// IsOff returns true if the motor is off
func (d *MotorDriver) IsOff() bool {
	return !d.IsOn()
}

// Toggle sets the motor to the opposite of it's current state
func (d *MotorDriver) Toggle() error {
	if d.IsOn() {
		return d.Off()
	}

	return d.On()
}

// Speed sets the speed of the motor
func (d *MotorDriver) Speed(value byte) error {
	if writer, ok := d.connection.(PwmWriter); ok {
		d.CurrentMode = "analog"
		d.CurrentSpeed = value
		return writer.PwmWrite(d.SpeedPin, value)
	}
	return ErrPwmWriteUnsupported
}

// Forward sets the forward pin to the specified speed
func (d *MotorDriver) Forward(speed byte) error {
	if err := d.Direction("forward"); err != nil {
		return err
	}
	if err := d.Speed(speed); err != nil {
		return err
	}

	return nil
}

// Backward sets the backward pin to the specified speed
func (d *MotorDriver) Backward(speed byte) error {
	if err := d.Direction("backward"); err != nil {
		return err
	}
	if err := d.Speed(speed); err != nil {
		return err
	}

	return nil
}

// Direction sets the direction pin to the specified speed
func (d *MotorDriver) Direction(direction string) error {
	d.CurrentDirection = direction
	if d.DirectionPin != "" {
		var level byte
		if direction == "forward" {
			level = 1
		} else {
			level = 0
		}
		return d.connection.DigitalWrite(d.DirectionPin, level)
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

	if err := d.connection.DigitalWrite(d.ForwardPin, forwardLevel); err != nil {
		return err
	}

	return d.connection.DigitalWrite(d.BackwardPin, backwardLevel)
}

func (d *MotorDriver) isDigital() bool {
	return d.CurrentMode == "digital"
}

func (d *MotorDriver) changeState(state byte) error {
	d.CurrentState = state
	if state == 1 {
		d.CurrentSpeed = 255
	} else {
		d.CurrentSpeed = 0
	}

	if d.ForwardPin == "" {
		return d.connection.DigitalWrite(d.SpeedPin, state)
	}

	if state != 1 {
		return d.Direction("none")
	}

	if err := d.Direction(d.CurrentDirection); err != nil {
		return err
	}
	if d.SpeedPin != "" {
		if err := d.Speed(d.CurrentSpeed); err != nil {
			return err
		}
	}

	return nil
}
