package gpio

import (
	"errors"
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
	// Error event
	Error = "error"
	// ButtonRelease event
	ButtonRelease = "release"
	// ButtonPush event
	ButtonPush = "push"
	// Data event
	Data = "data"
	// Vibration event
	Vibration = "vibration"
	// MotionDetected event
	MotionDetected = "motion-detected"
	// MotionStopped event
	MotionStopped = "motion-stopped"
)

// PwmWriter interface represents an Adaptor which has Pwm capabilities
type PwmWriter interface {
	PwmWrite(string, byte) (err error)
}

// ServoWriter interface represents an Adaptor which has Servo capabilities
type ServoWriter interface {
	ServoWrite(string, byte) (err error)
}

// DigitalWriter interface represents an Adaptor which has DigitalWrite capabilities
type DigitalWriter interface {
	DigitalWrite(string, byte) (err error)
}

// DigitalReader interface represents an Adaptor which has DigitalRead capabilities
type DigitalReader interface {
	DigitalRead(string) (val int, err error)
}
