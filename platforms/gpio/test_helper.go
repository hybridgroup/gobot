package gpio

type gpioTestAdaptor struct {
	name string
	port string
}

func (t *gpioTestAdaptor) AnalogWrite(string, byte) (err error)  { return nil }
func (t *gpioTestAdaptor) DigitalWrite(string, byte) (err error) { return nil }
func (t *gpioTestAdaptor) ServoWrite(string, byte) (err error)   { return nil }
func (t *gpioTestAdaptor) PwmWrite(string, byte) (err error)     { return nil }
func (t *gpioTestAdaptor) InitServo() (err error)                { return nil }
func (t *gpioTestAdaptor) AnalogRead(string) (val int, err error) {
	return 99, nil
}
func (t *gpioTestAdaptor) DigitalRead(string) (val int, err error) {
	return 1, nil
}
func (t *gpioTestAdaptor) Connect() (errs []error)  { return }
func (t *gpioTestAdaptor) Finalize() (errs []error) { return }
func (t *gpioTestAdaptor) Name() string             { return t.name }
func (t *gpioTestAdaptor) Port() string             { return t.port }

func newGpioTestAdaptor(name string) *gpioTestAdaptor {
	return &gpioTestAdaptor{
		name: name,
		port: "/dev/null",
	}
}
