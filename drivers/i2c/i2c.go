package i2c

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var (
	ErrEncryptedBytes  = errors.New("Encrypted bytes")
	ErrNotEnoughBytes  = errors.New("Not enough bytes read")
	ErrNotReady        = errors.New("Device is not ready")
	ErrInvalidPosition = errors.New("Invalid position value")
)

const (
	Error    = "error"
	Joystick = "joystick"
	C        = "c"
	Z        = "z"
)

type I2cStarter interface {
	I2cStart(address int) (err error)
}

type I2cReader interface {
	I2cRead(address int, len int) (data []byte, err error)
}

type I2cWriter interface {
	I2cWrite(address int, buf []byte) (err error)
}

type I2c interface {
	gobot.Adaptor
	I2cStarter
	I2cReader
	I2cWriter
}
