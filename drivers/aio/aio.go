package aio

import (
	"errors"
)

var (
	// ErrAnalogReadUnsupported is error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrAnalogReadUnsupported = errors.New("AnalogRead is not supported by this platform")
)

const (
	// Error event
	Error = "error"
	// Data event
	Data = "data"
	// Value event
	Value = "value"
	// Vibration event
	Vibration = "vibration"
)

// AnalogReader interface represents an Adaptor which has AnalogRead capabilities
type AnalogReader interface {
	//gobot.Adaptor
	AnalogRead(pin string) (val int, err error)
}

// AnalogWriter interface represents an Adaptor which has AnalogWrite capabilities
type AnalogWriter interface {
	//gobot.Adaptor
	AnalogWrite(pin string, val int) (err error)
}
