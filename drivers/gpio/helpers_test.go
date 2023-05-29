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
	testAdaptorDigitalWrite func(pin string, val byte) (err error)
	testAdaptorServoWrite   func(pin string, val byte) (err error)
	testAdaptorPwmWrite     func(pin string, val byte) (err error)
	testAdaptorAnalogRead   func(ping string) (val int, err error)
	testAdaptorDigitalRead  func(ping string) (val int, err error)
}

func (t *gpioTestAdaptor) TestAdaptorDigitalWrite(f func(pin string, val byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalWrite = f
}
func (t *gpioTestAdaptor) TestAdaptorServoWrite(f func(pin string, val byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorServoWrite = f
}
func (t *gpioTestAdaptor) TestAdaptorPwmWrite(f func(pin string, val byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorPwmWrite = f
}
func (t *gpioTestAdaptor) TestAdaptorAnalogRead(f func(pin string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorAnalogRead = f
}
func (t *gpioTestAdaptor) TestAdaptorDigitalRead(f func(pin string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalRead = f
}

func (t *gpioTestAdaptor) ServoWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorServoWrite(pin, val)
}
func (t *gpioTestAdaptor) PwmWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorPwmWrite(pin, val)
}
func (t *gpioTestAdaptor) AnalogRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorAnalogRead(pin)
}
func (t *gpioTestAdaptor) DigitalRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorDigitalRead(pin)
}
func (t *gpioTestAdaptor) DigitalWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorDigitalWrite(pin, val)
}
func (t *gpioTestAdaptor) Connect() (err error)  { return }
func (t *gpioTestAdaptor) Finalize() (err error) { return }
func (t *gpioTestAdaptor) Name() string          { return t.name }
func (t *gpioTestAdaptor) SetName(n string)      { t.name = n }
func (t *gpioTestAdaptor) Port() string          { return t.port }

func newGpioTestAdaptor() *gpioTestAdaptor {
	return &gpioTestAdaptor{
		port: "/dev/null",
		testAdaptorDigitalWrite: func(pin string, val byte) (err error) {
			return nil
		},
		testAdaptorServoWrite: func(pin string, val byte) (err error) {
			return nil
		},
		testAdaptorPwmWrite: func(pin string, val byte) (err error) {
			return nil
		},
		testAdaptorAnalogRead: func(pin string) (val int, err error) {
			return 99, nil
		},
		testAdaptorDigitalRead: func(pin string) (val int, err error) {
			return 1, nil
		},
	}
}
