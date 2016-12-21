package gpio

type gpioTestBareAdaptor struct{}

func (t *gpioTestBareAdaptor) Connect() (err error)  { return }
func (t *gpioTestBareAdaptor) Finalize() (err error) { return }
func (t *gpioTestBareAdaptor) Name() string          { return "" }
func (t *gpioTestBareAdaptor) SetName(n string)      {}

type gpioTestDigitalWriter struct {
	gpioTestBareAdaptor
}

func (t *gpioTestDigitalWriter) DigitalWrite(string, byte) (err error) { return }

type gpioTestAdaptor struct {
	name string
	port string
}

var testAdaptorDigitalWrite = func() (err error) {
	return nil
}
var testAdaptorServoWrite = func() (err error) {
	return nil
}
var testAdaptorPwmWrite = func() (err error) {
	return nil
}
var testAdaptorAnalogRead = func() (val int, err error) {
	return 99, nil
}
var testAdaptorDigitalRead = func() (val int, err error) {
	return 1, nil
}

func (t *gpioTestAdaptor) DigitalWrite(string, byte) (err error) {
	return testAdaptorDigitalWrite()
}
func (t *gpioTestAdaptor) ServoWrite(string, byte) (err error) {
	return testAdaptorServoWrite()
}
func (t *gpioTestAdaptor) PwmWrite(string, byte) (err error) {
	return testAdaptorPwmWrite()
}
func (t *gpioTestAdaptor) AnalogRead(string) (val int, err error) {
	return testAdaptorAnalogRead()
}
func (t *gpioTestAdaptor) DigitalRead(string) (val int, err error) {
	return testAdaptorDigitalRead()
}
func (t *gpioTestAdaptor) Connect() (err error)  { return }
func (t *gpioTestAdaptor) Finalize() (err error) { return }
func (t *gpioTestAdaptor) Name() string          { return t.name }
func (t *gpioTestAdaptor) SetName(n string)      { t.name = n }
func (t *gpioTestAdaptor) Port() string          { return t.port }

func newGpioTestAdaptor() *gpioTestAdaptor {
	return &gpioTestAdaptor{
		port: "/dev/null",
	}
}
