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
	// Vibration event
	Vibration = "vibration"
)

// AnalogReader interface represents an Adaptor which has Analog capabilities
type AnalogReader interface {
	//gobot.Adaptor
	AnalogRead(string) (val int, err error)
}

// DHTReader interface represents an Adaptor which has capabilities to
// read digital humidity and temperature for the GrovePi DHT sensors.
type DHTReader interface {
	ReadDHT(pin string) (temperature, humidity float32, err error)
}
