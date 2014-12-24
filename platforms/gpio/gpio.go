package gpio

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var (
	ErrServoWriteUnsupported   = errors.New("ServoWrite is not supported by this platform")
	ErrPwmWriteUnsupported     = errors.New("PwmWrite is not supported by this platform")
	ErrAnalogReadUnsupported   = errors.New("AnalogRead is not supported by this platform")
	ErrDigitalWriteUnsupported = errors.New("DigitalWrite is not supported by this platform")
	ErrDigitalReadUnsupported  = errors.New("DigitalRead is not supported by this platform")
	ErrServoOutOfRange         = errors.New("servo angle must be between 0-180")
)

const (
	Release = "release"
	Push    = "push"
	Error   = "error"
	Data    = "data"
)

type PwmWriter interface {
	gobot.Adaptor
	PwmWrite(string, byte) (err error)
}

type ServoWriter interface {
	gobot.Adaptor
	ServoWrite(string, byte) (err error)
}

type AnalogReader interface {
	gobot.Adaptor
	AnalogRead(string) (val int, err error)
}

type DigitalWriter interface {
	gobot.Adaptor
	DigitalWrite(string, byte) (err error)
}

type DigitalReader interface {
	gobot.Adaptor
	DigitalRead(string) (val int, err error)
}
