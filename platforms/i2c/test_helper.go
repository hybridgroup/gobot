package i2c

import (
	"github.com/hybridgroup/gobot"
)

var rgb = map[string]interface{}{
	"red":   1.0,
	"green": 1.0,
	"blue":  1.0,
}

func castColor(color string) byte {
	return byte(rgb[color].(float64))
}

var red = castColor("red")
var green = castColor("green")
var blue = castColor("blue")

type i2cTestAdaptor struct {
	gobot.Adaptor
	i2cReadImpl func() []byte
}

func (t *i2cTestAdaptor) I2cStart(byte) {}
func (t *i2cTestAdaptor) I2cRead(uint) []byte {
	return t.i2cReadImpl()
}
func (t *i2cTestAdaptor) I2cWrite([]byte) {}
func (t *i2cTestAdaptor) Connect() bool   { return true }
func (t *i2cTestAdaptor) Finalize() bool  { return true }

func newI2cTestAdaptor(name string) *i2cTestAdaptor {
	return &i2cTestAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"I2cTestAdaptor",
		),
		i2cReadImpl: func() []byte {
			return []byte{}
		},
	}
}
