package i2c

import "github.com/hybridgroup/gobot"

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
