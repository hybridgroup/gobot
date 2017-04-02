package gpio

import "sync"

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
	mtx                     sync.Mutex
	testAdaptorDigitalWrite func() (err error)
	testAdaptorServoWrite   func() (err error)
	testAdaptorPwmWrite     func() (err error)
	testAdaptorAnalogRead   func() (val int, err error)
	testAdaptorDigitalRead  func() (val int, err error)
}

func (t *gpioTestAdaptor) TestAdaptorDigitalWrite(f func() (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalWrite = f
}
func (t *gpioTestAdaptor) TestAdaptorServoWrite(f func() (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorServoWrite = f
}
func (t *gpioTestAdaptor) TestAdaptorPwmWrite(f func() (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorPwmWrite = f
}
func (t *gpioTestAdaptor) TestAdaptorAnalogRead(f func() (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorAnalogRead = f
}
func (t *gpioTestAdaptor) TestAdaptorDigitalRead(f func() (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalRead = f
}

func (t *gpioTestAdaptor) ServoWrite(string, byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorServoWrite()
}
func (t *gpioTestAdaptor) PwmWrite(string, byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorPwmWrite()
}
func (t *gpioTestAdaptor) AnalogRead(string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorAnalogRead()
}
func (t *gpioTestAdaptor) DigitalRead(string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorDigitalRead()
}
func (t *gpioTestAdaptor) DigitalWrite(string, byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorDigitalWrite()
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
