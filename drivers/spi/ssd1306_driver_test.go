package spi

import (
	"image"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on spi.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*SSD1306Driver)(nil)

func initTestSSDDriver() *SSD1306Driver {
	return NewSSD1306Driver(newGpioTestAdaptor())
}

func TestDriverSSDStart(t *testing.T) {
	d := initTestSSDDriver()
	assert.NoError(t, d.Start())
}

func TestDriverSSDHalt(t *testing.T) {
	d := initTestSSDDriver()
	_ = d.Start()
	assert.NoError(t, d.Halt())
}

func TestDriverSSDDisplay(t *testing.T) {
	d := initTestSSDDriver()
	_ = d.Start()
	assert.NoError(t, d.Display())
}

func TestSSD1306DriverShowImage(t *testing.T) {
	d := initTestSSDDriver()
	_ = d.Start()
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	assert.ErrorContains(t, d.ShowImage(img), "Image must match the display width and height")

	img = image.NewRGBA(image.Rect(0, 0, 128, 64))
	assert.NoError(t, d.ShowImage(img))
}

type gpioTestAdaptor struct {
	name string
	port string
	mtx  sync.Mutex
	Connector
	digitalWriteFunc func() (err error)
	servoWriteFunc   func() (err error)
	pwmWriteFunc     func() (err error)
	analogReadFunc   func() (val int, err error)
	digitalReadFunc  func() (val int, err error)
}

func (t *gpioTestAdaptor) ServoWrite(string, byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.servoWriteFunc()
}

func (t *gpioTestAdaptor) PwmWrite(string, byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.pwmWriteFunc()
}

func (t *gpioTestAdaptor) AnalogRead(string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.analogReadFunc()
}

func (t *gpioTestAdaptor) DigitalRead(string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.digitalReadFunc()
}

func (t *gpioTestAdaptor) DigitalWrite(string, byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.digitalWriteFunc()
}
func (t *gpioTestAdaptor) Connect() (err error)  { return }
func (t *gpioTestAdaptor) Finalize() (err error) { return }
func (t *gpioTestAdaptor) Name() string          { return t.name }
func (t *gpioTestAdaptor) SetName(n string)      { t.name = n }
func (t *gpioTestAdaptor) Port() string          { return t.port }

func newGpioTestAdaptor() *gpioTestAdaptor {
	a := newSpiTestAdaptor()
	return &gpioTestAdaptor{
		port: "/dev/null",
		digitalWriteFunc: func() (err error) {
			return nil
		},
		servoWriteFunc: func() (err error) {
			return nil
		},
		pwmWriteFunc: func() (err error) {
			return nil
		},
		analogReadFunc: func() (val int, err error) {
			return 99, nil
		},
		digitalReadFunc: func() (val int, err error) {
			return 1, nil
		},
		Connector: a,
	}
}
