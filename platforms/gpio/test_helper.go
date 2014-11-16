package gpio

import "github.com/hybridgroup/gobot"

type gpioTestAdaptor struct {
	gobot.Adaptor
}

func (t *gpioTestAdaptor) AnalogWrite(string, byte)  {}
func (t *gpioTestAdaptor) DigitalWrite(string, byte) {}
func (t *gpioTestAdaptor) ServoWrite(string, byte)   {}
func (t *gpioTestAdaptor) PwmWrite(string, byte)     {}
func (t *gpioTestAdaptor) InitServo()                {}
func (t *gpioTestAdaptor) AnalogRead(string) int {
	return 99
}
func (t *gpioTestAdaptor) DigitalRead(string) int {
	return 1
}
func (t *gpioTestAdaptor) Connect() error  { return nil }
func (t *gpioTestAdaptor) Finalize() error { return nil }

func newGpioTestAdaptor(name string) *gpioTestAdaptor {
	return &gpioTestAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"/dev/null",
		),
	}
}
