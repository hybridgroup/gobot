package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
)

// ensure that PCA9685Driver fulfills Gobot Driver interface
var _ gobot.Driver = (*PCA9685Driver)(nil)

// and also the PwmWriter and ServoWriter interfaces
var _ gpio.PwmWriter = (*PCA9685Driver)(nil)
var _ gpio.ServoWriter = (*PCA9685Driver)(nil)

// --------- HELPERS
func initTestPCA9685Driver() (driver *PCA9685Driver) {
	driver, _ = initTestPCA9685DriverWithStubbedAdaptor()
	return
}

func initTestPCA9685DriverWithStubbedAdaptor() (*PCA9685Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewPCA9685Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewPCA9685Driver(t *testing.T) {
	// Does it return a pointer to an instance of PCA9685Driver?
	var pca interface{} = NewPCA9685Driver(newI2cTestAdaptor())
	_, ok := pca.(*PCA9685Driver)
	if !ok {
		t.Errorf("NewPCA9685Driver() should have returned a *PCA9685Driver")
	}
}

func TestPCA9685DriverName(t *testing.T) {
	pca := initTestPCA9685Driver()
	gobottest.Refute(t, pca.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(pca.Name(), "PCA9685"), true)
}

func TestPCA9685DriverOptions(t *testing.T) {
	pca := NewPCA9685Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, pca.GetBusOrDefault(1), 2)
}

// Methods
func TestPCA9685DriverStart(t *testing.T) {
	pca := initTestPCA9685Driver()

	gobottest.Assert(t, pca.Start(), nil)
}

func TestPCA9685DriverStartConnectError(t *testing.T) {
	d, adaptor := initTestPCA9685DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestPCA9685DriverStartWriteError(t *testing.T) {
	pca, adaptor := initTestPCA9685DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, pca.Start(), errors.New("write error"))
}

func TestPCA9685DriverHalt(t *testing.T) {
	pca := initTestPCA9685Driver()
	gobottest.Assert(t, pca.Start(), nil)
	gobottest.Assert(t, pca.Halt(), nil)
}

func TestPCA9685DriverSetPWM(t *testing.T) {
	pca := initTestPCA9685Driver()
	gobottest.Assert(t, pca.Start(), nil)
	gobottest.Assert(t, pca.SetPWM(0, 0, 256), nil)
}

func TestPCA9685DriverSetPWMError(t *testing.T) {
	pca, adaptor := initTestPCA9685DriverWithStubbedAdaptor()
	gobottest.Assert(t, pca.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, pca.SetPWM(0, 0, 256), errors.New("write error"))
}

func TestPCA9685DriverSetPWMFreq(t *testing.T) {
	pca, adaptor := initTestPCA9685DriverWithStubbedAdaptor()
	gobottest.Assert(t, pca.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	gobottest.Assert(t, pca.SetPWMFreq(60), nil)
}

func TestPCA9685DriverSetPWMFreqReadError(t *testing.T) {
	pca, adaptor := initTestPCA9685DriverWithStubbedAdaptor()
	gobottest.Assert(t, pca.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, pca.SetPWMFreq(60), errors.New("read error"))
}

func TestPCA9685DriverSetPWMFreqWriteError(t *testing.T) {
	pca, adaptor := initTestPCA9685DriverWithStubbedAdaptor()
	gobottest.Assert(t, pca.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, pca.SetPWMFreq(60), errors.New("write error"))
}

func TestPCA9685DriverSetName(t *testing.T) {
	pca := initTestPCA9685Driver()
	pca.SetName("TESTME")
	gobottest.Assert(t, pca.Name(), "TESTME")
}
