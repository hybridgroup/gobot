package i2c

import (
	"github.com/hybridgroup/gobot"
	"github.com/maraino/go-mock"
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
}

func (t *i2cTestAdaptor) I2cStart(byte) {}
func (t *i2cTestAdaptor) I2cRead(uint) []byte {
	return []byte{99, 1}
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
	}
}

type I2cInterfaceClient struct {
	gobot.Adaptor
	mock.Mock
}

func (i *I2cInterfaceClient) I2cStart(b byte) { i.Called(b) }
func (i *I2cInterfaceClient) I2cRead(u uint) []byte {
	ret := i.Called(u)
	return ret.Result[0].([]byte)
}
func (i *I2cInterfaceClient) I2cWrite(b []byte) { i.Called(b) }
func (t *I2cInterfaceClient) Connect() bool     { return true }
func (t *I2cInterfaceClient) Finalize() bool    { return true }

func NewI2cInterfaceClient() *I2cInterfaceClient {
	return &I2cInterfaceClient{
		Adaptor: *gobot.NewAdaptor(
			"i2c",
			"I2cInterfaceClient",
		),
	}
}
