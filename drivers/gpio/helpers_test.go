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
	name                    string
	port                    string
	testAdaptorDigitalWrite func() (err error)
	testAdaptorServoWrite   func() (err error)
	testAdaptorPwmWrite     func() (err error)
	testAdaptorAnalogRead   func() (val int, err error)
	testAdaptorDigitalRead  func() (val int, err error)
}

func (t *gpioTestAdaptor) DigitalWrite(string, byte) (err error) {
	return t.testAdaptorDigitalWrite()
}
func (t *gpioTestAdaptor) ServoWrite(string, byte) (err error) {
	return t.testAdaptorServoWrite()
}
func (t *gpioTestAdaptor) PwmWrite(string, byte) (err error) {
	return t.testAdaptorPwmWrite()
}
func (t *gpioTestAdaptor) AnalogRead(string) (val int, err error) {
	return t.testAdaptorAnalogRead()
}
func (t *gpioTestAdaptor) DigitalRead(string) (val int, err error) {
	return t.testAdaptorDigitalRead()
}
func (t *gpioTestAdaptor) Connect() (err error)  { return }
func (t *gpioTestAdaptor) Finalize() (err error) { return }
func (t *gpioTestAdaptor) Name() string          { return t.name }
func (t *gpioTestAdaptor) SetName(n string)      { t.name = n }
func (t *gpioTestAdaptor) Port() string          { return t.port }

func newGpioTestAdaptor() *gpioTestAdaptor {
	return &gpioTestAdaptor{
		port: "/dev/null",
		testAdaptorDigitalWrite: func() (err error) {
			return nil
		},
		testAdaptorServoWrite: func() (err error) {
			return nil
		},
		testAdaptorPwmWrite: func() (err error) {
			return nil
		},
		testAdaptorAnalogRead: func() (val int, err error) {
			return 99, nil
		},
		testAdaptorDigitalRead: func() (val int, err error) {
			return 1, nil
		},
	}
}
