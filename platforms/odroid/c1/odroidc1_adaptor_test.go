package c1

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}
func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}
func (NullReadWriteCloser) Close() error {
	return nil
}

func initTestODroidC1Adaptor() *ODroidC1Adaptor {
	a := NewODroidC1Adaptor("myAdaptor")
	a.Connect()
	return a
}

func TestODroidC1Adaptor(t *testing.T) {
	a := NewODroidC1Adaptor("myAdaptor")
	gobot.Assert(t, a.Name(), "myAdaptor")
	gobot.Assert(t, a.i2cLocation, "/sys/bus/i2c")
}

func TestODroidC1AdaptorFinalize(t *testing.T) {
	a := initTestODroidC1Adaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		
		"/sys/class/gpio/gpio74/value",
		"/sys/class/gpio/gpio74/direction",
	})

	sysfs.SetFilesystem(fs)
	a.DigitalWrite("3", 1)
	a.i2cDevice = new(NullReadWriteCloser)
	gobot.Assert(t, len(a.Finalize()), 0)
}

func TestODroidC1AdaptorDigitalIO(t *testing.T) {
	a := initTestODroidC1Adaptor()
	fs := sysfs.NewMockFilesystem([]string{		
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio83/value",
		"/sys/class/gpio/gpio83/direction",
		"/sys/class/gpio/gpio116/value",
		"/sys/class/gpio/gpio116/direction",
	})

	sysfs.SetFilesystem(fs)

	a.DigitalWrite("7", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio83/value"].Contents, "1")

	a.DigitalWrite("13", 1)
	i, _ := a.DigitalRead("13")
	gobot.Assert(t, i, 1)
}

func TestODroidC1AdaptorPWM(t *testing.T) {
	a := initTestODroidC1Adaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/devices/platform/pwm-ctrl/duty0",
		"/sys/devices/platform/pwm-ctrl/freq0",
		"/sys/devices/platform/pwm-ctrl/enable0",
	})

	sysfs.SetFilesystem(fs)

	// PWM
	a.PwmWrite("33", 175)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/platform/pwm-ctrl/freq0"].Contents,
		"500000",
	)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/platform/pwm-ctrl/duty0"].Contents,
		"343137",
	)

	a.ServoWrite("33", 100)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/platform/pwm-ctrl/freq0"].Contents,
		"16666666",
	)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/platform/pwm-ctrl/duty0"].Contents,
		"1898148",
	)
}

func TestODroidC1AdaptorAnalog(t *testing.T) {
	a := initTestODroidC1Adaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/saradc/saradc_ch0",
		"/sys/class/saradc/saradc_ch1",
	})

	sysfs.SetFilesystem(fs)

	// Analog
	fs.Files["/sys/class/saradc/saradc_ch0"].Contents = "567\n"
	i, _ := a.AnalogRead("40")
	gobot.Assert(t, i, 567)
	
	fs.Files["/sys/class/saradc/saradc_ch1"].Contents = "123\n"
	i, _ = a.AnalogRead("37")
	gobot.Assert(t, i, 123)

	i, err := a.AnalogRead("9999")
	gobot.Assert(t, err, errors.New("Not a valid pin"))

}

func TestODroidC1AdaptorI2c(t *testing.T) {
	a := initTestODroidC1Adaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/bus/i2c",
	})
	sysfs.SetFilesystem(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.I2cWrite([]byte{0x00, 0x01})
	data, _ := a.I2cRead(2)
	gobot.Assert(t, data, []byte{0x00, 0x01})
}
