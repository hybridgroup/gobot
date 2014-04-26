package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type MotorInterface interface {
	PwmWrite(string, byte)
	DigitalWrite(string, byte)
}

type Motor struct {
	gobot.Driver
	Adaptor          MotorInterface
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

func NewMotor(a MotorInterface) *Motor {
	m := new(Motor)
	m.Adaptor = a
	m.Commands = []string{
		"OffC",
		"OnC",
		"IsOnC",
		"IsOffC",
		"ToggleC",
		"SpeedC",
		"MinC",
		"MaxC",
		"ForwardC",
		"BackwardC",
		"CurrentSpeedC",
	}
	m.CurrentState = 0
	m.CurrentSpeed = 0
	m.CurrentMode = "digital"
	m.CurrentDirection = "forward"
	return m
}

func (m *Motor) Start() bool { return true }
func (m *Motor) Halt() bool  { return true }
func (m *Motor) Init() bool  { return true }

func (m *Motor) Off() {
	if m.isDigital() {
		m.changeState(0)
	} else {
		m.Speed(0)
	}
}

func (m *Motor) On() {
	if m.isDigital() {
		m.changeState(1)
	} else {
		if m.CurrentSpeed == 0 {
			m.CurrentSpeed = 255
		}
		m.Speed(m.CurrentSpeed)
	}
}

func (m *Motor) Min() {
	m.Off()
}

func (m *Motor) Max() {
	m.Speed(255)
}

func (m *Motor) IsOn() bool {
	if m.isDigital() {
		return m.CurrentState == 1
	} else {
		return m.CurrentSpeed > 0
	}
}

func (m *Motor) IsOff() bool {
	return !m.IsOn()
}

func (m *Motor) Toggle() {
	if m.IsOn() {
		m.Off()
	} else {
		m.On()
	}
}

func (m *Motor) Speed(value byte) {
	m.CurrentMode = "analog"
	m.CurrentSpeed = value
	m.Adaptor.PwmWrite(m.SpeedPin, value)
}

func (m *Motor) Forward(speed byte) {
	m.Direction("forward")
	m.Speed(speed)
}

func (m *Motor) Backward(speed byte) {
	m.Direction("backward")
	m.Speed(speed)
}

func (m *Motor) Direction(direction string) {
	m.CurrentDirection = direction
	if m.DirectionPin != "" {
		var level byte
		if direction == "forward" {
			level = 1
		} else {
			level = 0
		}
		m.Adaptor.DigitalWrite(m.DirectionPin, level)
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
		m.Adaptor.DigitalWrite(m.ForwardPin, forwardLevel)
		m.Adaptor.DigitalWrite(m.BackwardPin, backwardLevel)
	}
}

func (m *Motor) isDigital() bool {
	if m.CurrentMode == "digital" {
		return true
	}
	return false
}

func (m *Motor) changeState(state byte) {
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
		m.Adaptor.DigitalWrite(m.SpeedPin, state)
	}
}
