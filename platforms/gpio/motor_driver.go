package gpio

import (
	"github.com/hybridgroup/gobot"
)

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

func NewMotorDriver(a PwmDigitalWriter, name string, pin string) *MotorDriver {
	return &MotorDriver{
		Driver: gobot.Driver{
			Name:    name,
			Pin:     pin,
			Adaptor: a.(gobot.AdaptorInterface),
		},
		CurrentState:     0,
		CurrentSpeed:     0,
		CurrentMode:      "digital",
		CurrentDirection: "forward",
	}
}

func (m *MotorDriver) adaptor() PwmDigitalWriter {
	return m.Driver.Adaptor.(PwmDigitalWriter)
}

func (m *MotorDriver) Start() bool { return true }
func (m *MotorDriver) Halt() bool  { return true }
func (m *MotorDriver) Init() bool  { return true }

func (m *MotorDriver) Off() {
	if m.isDigital() {
		m.changeState(0)
	} else {
		m.Speed(0)
	}
}

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

func (m *MotorDriver) Min() {
	m.Off()
}

func (m *MotorDriver) Max() {
	m.Speed(255)
}

func (m *MotorDriver) IsOn() bool {
	if m.isDigital() {
		return m.CurrentState == 1
	}
	return m.CurrentSpeed > 0
}

func (m *MotorDriver) IsOff() bool {
	return !m.IsOn()
}

func (m *MotorDriver) Toggle() {
	if m.IsOn() {
		m.Off()
	} else {
		m.On()
	}
}

func (m *MotorDriver) Speed(value byte) {
	m.CurrentMode = "analog"
	m.CurrentSpeed = value
	m.adaptor().PwmWrite(m.SpeedPin, value)
}

func (m *MotorDriver) Forward(speed byte) {
	m.Direction("forward")
	m.Speed(speed)
}

func (m *MotorDriver) Backward(speed byte) {
	m.Direction("backward")
	m.Speed(speed)
}

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
