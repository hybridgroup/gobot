package gpio

import (
	"fmt"
	"sync"

	"gobot.io/x/gobot/v2"
)

type gpioTestBareAdaptor struct{}

func (t *gpioTestBareAdaptor) Connect() (err error)  { return }
func (t *gpioTestBareAdaptor) Finalize() (err error) { return }
func (t *gpioTestBareAdaptor) Name() string          { return "" }
func (t *gpioTestBareAdaptor) SetName(n string)      {}

type digitalPinMock struct {
	writeFunc func(val int) (err error)
}

type gpioTestAdaptor struct {
	name             string
	pinMap           map[string]gobot.DigitalPinner
	port             string
	mtx              sync.Mutex
	digitalReadFunc  func(ping string) (val int, err error)
	digitalWriteFunc func(pin string, val byte) (err error)
	pwmWriteFunc     func(pin string, val byte) (err error)
	servoWriteFunc   func(pin string, val byte) (err error)
}

func newGpioTestAdaptor() *gpioTestAdaptor {
	t := gpioTestAdaptor{
		name:   "gpio_test_adaptor",
		pinMap: make(map[string]gobot.DigitalPinner),
		port:   "/dev/null",
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

// DigitalPin (interface DigitalPinnerProvider) return a pin object
func (t *gpioTestAdaptor) DigitalPin(id string) (gobot.DigitalPinner, error) {
	if pin, ok := t.pinMap[id]; ok {
		return pin, nil
	}
	return nil, fmt.Errorf("pin '%s' not found in '%s'", id, t.name)
}

// ApplyOptions (interface DigitalPinOptionApplier by DigitalPinner) apply all given options to the pin immediately
func (d *digitalPinMock) ApplyOptions(options ...func(gobot.DigitalPinOptioner) bool) error {
	return nil
}

// Export (interface DigitalPinner) exports the pin for use by the adaptor
func (d *digitalPinMock) Export() error {
	return nil
}

// Unexport (interface DigitalPinner) releases the pin from the adaptor, so it is free for the operating system
func (d *digitalPinMock) Unexport() error {
	return nil
}

// Read (interface DigitalPinner) reads the current value of the pin
func (d *digitalPinMock) Read() (n int, err error) {
	return 0, err
}

// Write (interface DigitalPinner) writes to the pin
func (d *digitalPinMock) Write(b int) error {
	return d.writeFunc(b)
}

func (t *gpioTestAdaptor) addDigitalPin(id string) *digitalPinMock {
	dpm := &digitalPinMock{
		writeFunc: func(val int) (err error) { return nil },
	}
	t.pinMap[id] = dpm
	return dpm
}
