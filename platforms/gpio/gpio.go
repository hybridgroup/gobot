package gpio

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var (
	// ErrServoWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrServoWriteUnsupported = errors.New("ServoWrite is not supported by this platform")
	// ErrPwmWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrPwmWriteUnsupported = errors.New("PwmWrite is not supported by this platform")
	// ErrAnalogReadUnsupported is error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrAnalogReadUnsupported = errors.New("AnalogRead is not supported by this platform")
	// ErrDigitalWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrDigitalWriteUnsupported = errors.New("DigitalWrite is not supported by this platform")
	// ErrDigitalReadUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrDigitalReadUnsupported = errors.New("DigitalRead is not supported by this platform")
	// ErrServoOutOfRange is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrServoOutOfRange = errors.New("servo angle must be between 0-180")
)

const (
	// Release event
	Release = "release"
	// Push event
	Push = "push"
	// Error event
	Error = "error"
	// Data event
	Data = "data"
	// Vibration event
	Vibration = "vibration"
)

// PwmWriter interface represents an Adaptor which has Pwm capabilities
type PwmWriter interface {
	gobot.Adaptor
	PwmWrite(string, byte) (err error)
}

// ServoWriter interface represents an Adaptor which has Servo capabilities
type ServoWriter interface {
	gobot.Adaptor
	ServoWrite(string, byte) (err error)
}

// AnalogReader interface represents an Adaptor which has Analog capabilities
type AnalogReader interface {
	gobot.Adaptor
	AnalogRead(string) (val int, err error)
}

// DigitalWriter interface represents an Adaptor which has DigitalWrite capabilities
type DigitalWriter interface {
	gobot.Adaptor
	DigitalWrite(string, byte) (err error)
}

// DigitalReader interface represents an Adaptor which has DigitalRead capabilities
type DigitalReader interface {
	gobot.Adaptor
	DigitalRead(string) (val int, err error)
}
