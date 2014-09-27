package gpio

import (
	"github.com/hybridgroup/gobot"
)

// Represents a Motor
type MotorDriver struct {
	gobot.Driver
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

// NewMotorDriver return a new MotorDriver given a PwmDigitalWriter, name and pin
func NewMotorDriver(a PwmDigitalWriter, name string, pin string) *MotorDriver {
	return &MotorDriver{
		Driver: *gobot.NewDriver(
			name,
			"MotorDriver",
			a.(gobot.AdaptorInterface),
		),
		CurrentState:     0,
		CurrentSpeed:     0,
		CurrentMode:      "digital",
		CurrentDirection: "forward",
	}
}

func (m *MotorDriver) adaptor() PwmDigitalWriter {
	return m.Adaptor().(PwmDigitalWriter)
}

// Start starts the MotorDriver. Returns true on successful start of the driver
func (m *MotorDriver) Start() bool { return true }

// Halt halts the MotorDriver. Returns true on successful halt of the driver
func (m *MotorDriver) Halt() bool { return true }

// Off turns the motor off or sets the motor to a 0 speed
func (m *MotorDriver) Off() {
	if m.isDigital() {
		m.changeState(0)
	} else {
		m.Speed(0)
	}
}

// On turns the motor on or sets the motor to a maximum speed
func (m *MotorDriver) On() {
	if m.isDigital() {
		m.changeState(1)
	} else {
		if m.CurrentSpeed == 0 {
			m.CurrentSpeed = 255
		}
		m.Speed(m.CurrentSpeed)
	}
}

// Min sets the motor to the minimum speed
func (m *MotorDriver) Min() {
	m.Off()
}

// Max sets the motor to the maximum speed
func (m *MotorDriver) Max() {
	m.Speed(255)
}

// InOn returns true if the motor is on
func (m *MotorDriver) IsOn() bool {
	if m.isDigital() {
		return m.CurrentState == 1
	}
	return m.CurrentSpeed > 0
}

// InOff returns true if the motor is off
func (m *MotorDriver) IsOff() bool {
	return !m.IsOn()
}

// Toggle sets the motor to the opposite of it's current state
func (m *MotorDriver) Toggle() {
	if m.IsOn() {
		m.Off()
	} else {
		m.On()
	}
}

// Speed sets the speed of the motor
func (m *MotorDriver) Speed(value byte) {
	m.CurrentMode = "analog"
	m.CurrentSpeed = value
	m.adaptor().PwmWrite(m.SpeedPin, value)
}

// Forward sets the forward pin to the specified speed
func (m *MotorDriver) Forward(speed byte) {
	m.Direction("forward")
	m.Speed(speed)
}

// Backward sets the backward pin to the specified speed
func (m *MotorDriver) Backward(speed byte) {
	m.Direction("backward")
	m.Speed(speed)
}

// Direction sets the direction pin to the specified speed
func (m *MotorDriver) Direction(direction string) {
	m.CurrentDirection = direction
	if m.DirectionPin != "" {
		var level byte
		if direction == "forward" {
			level = 1
		} else {
			level = 0
		}
		m.adaptor().DigitalWrite(m.DirectionPin, level)
	} else {
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
		m.adaptor().DigitalWrite(m.ForwardPin, forwardLevel)
		m.adaptor().DigitalWrite(m.BackwardPin, backwardLevel)
	}
}

func (m *MotorDriver) isDigital() bool {
	if m.CurrentMode == "digital" {
		return true
	}
	return false
}

func (m *MotorDriver) changeState(state byte) {
	m.CurrentState = state
	if state == 1 {
		m.CurrentSpeed = 0
	} else {
		m.CurrentSpeed = 255
	}
	if m.ForwardPin != "" {
		if state == 0 {
			m.Direction(m.CurrentDirection)
			if m.SpeedPin != "" {
				m.Speed(m.CurrentSpeed)
			}
		} else {
			m.Direction("none")
		}
	} else {
		m.adaptor().DigitalWrite(m.SpeedPin, state)
	}
}
