package gpio

import "sync"

type gpioTestBareAdaptor struct{}

func (t *gpioTestBareAdaptor) Connect() (err error)  { return }
func (t *gpioTestBareAdaptor) Finalize() (err error) { return }
func (t *gpioTestBareAdaptor) Name() string          { return "" }
func (t *gpioTestBareAdaptor) SetName(n string)      {}

type gpioTestAdaptor struct {
	name             string
	port             string
	mtx              sync.Mutex
	digitalReadFunc  func(ping string) (val int, err error)
	digitalWriteFunc func(pin string, val byte) (err error)
	pwmWriteFunc     func(pin string, val byte) (err error)
	servoWriteFunc   func(pin string, val byte) (err error)
}

func newGpioTestAdaptor() *gpioTestAdaptor {
	t := gpioTestAdaptor{
		name: "gpio_test_adaptor",
		port: "/dev/null",
		digitalWriteFunc: func(pin string, val byte) (err error) {
			return nil
		},
		servoWriteFunc: func(pin string, val byte) (err error) {
			return nil
		},
		pwmWriteFunc: func(pin string, val byte) (err error) {
			return nil
		},
		digitalReadFunc: func(pin string) (val int, err error) {
			return 1, nil
		},
	}

	return &t
}

// DigitalRead capabilities (interface DigitalReader)
func (t *gpioTestAdaptor) DigitalRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.digitalReadFunc(pin)
}

// DigitalWrite capabilities (interface DigitalWriter)
func (t *gpioTestAdaptor) DigitalWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.digitalWriteFunc(pin, val)
}

// PwmWrite capabilities (interface PwmWriter)
func (t *gpioTestAdaptor) PwmWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.pwmWriteFunc(pin, val)
}

// ServoWrite capabilities (interface ServoWriter)
func (t *gpioTestAdaptor) ServoWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.servoWriteFunc(pin, val)
}

func (t *gpioTestAdaptor) Connect() (err error)  { return }
func (t *gpioTestAdaptor) Finalize() (err error) { return }
func (t *gpioTestAdaptor) Name() string          { return t.name }
func (t *gpioTestAdaptor) SetName(n string)      { t.name = n }
func (t *gpioTestAdaptor) Port() string          { return t.port }
