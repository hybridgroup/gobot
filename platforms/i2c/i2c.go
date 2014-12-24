package i2c

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var (
	ErrEncryptedBytes = errors.New("Encrypted bytes")
	ErrNotEnoughBytes = errors.New("Not enough bytes read")
)

const (
	Error    = "error"
	Joystick = "joystick"
	C        = "c"
	Z        = "z"
)

type I2c interface {
	gobot.Adaptor
	I2cStart(address byte) (err error)
	I2cRead(len uint) (data []byte, err error)
	I2cWrite(buf []byte) (err error)
}
