package beaglebone

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

type NullReadWriteCloser struct {
	contents []byte
}

func (n *NullReadWriteCloser) SetAddress(int) error {
	return nil
}

func (n *NullReadWriteCloser) Write(b []byte) (int, error) {
	n.contents = make([]byte, len(b))
	copy(n.contents[:], b[:])

	return len(b), nil
}

func (n *NullReadWriteCloser) Read(b []byte) (int, error) {
	copy(b, n.contents)
	return len(b), nil
}

var closeErr error

func (n *NullReadWriteCloser) Close() error {
	return closeErr
}

func TestBeagleboneAdaptor(t *testing.T) {
	glob = func(pattern string) (matches []string, err error) {
		return make([]string, 2), nil
	}
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-2",
		"/sys/devices/platform/bone_capemgr",
		"/sys/devices/platform/ocp/ocp4",
		"/sys/class/leds/beaglebone:green:usr1/brightness",
		"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
		"/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5",
		"/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/run",
		"/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/period",
		"/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/polarity",
		"/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/duty",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	sysfs.SetFilesystem(fs)
	a := NewAdaptor()
	a.slots = "/sys/devices/platform/bone_capemgr"
	a.ocp = "/sys/devices/platform/ocp/ocp4"
	a.kernel = "4.4"

	a.Connect()

	a.analogPath = "/sys/bus/iio/devices/iio:device0"

	// PWM
	glob = func(pattern string) (matches []string, err error) {
		pattern = strings.TrimSuffix(pattern, "*")
		return []string{pattern + "5"}, nil
	}

	gobottest.Assert(t, a.PwmWrite("P9_99", 175), errors.New("Not a valid pin"))
	a.PwmWrite("P9_14", 175)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/period"].Contents,
		"500000",
	)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/duty"].Contents,
		"343137",
	)

	a.ServoWrite("P9_14", 100)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/period"].Contents,
		"16666666",
	)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/ocp4/pwm_test_P9_14.5/duty"].Contents,
		"1898148",
	)

	// Analog
	fs.Files["/sys/bus/iio/devices/iio:device0/in_voltage1_raw"].Contents = "567\n"
	i, _ := a.AnalogRead("P9_40")
	gobottest.Assert(t, i, 567)

	_, err := a.AnalogRead("P9_99")
	gobottest.Assert(t, err, errors.New("Not a valid pin"))

	// DigitalIO
	a.DigitalWrite("usr1", 1)
	gobottest.Assert(t,
		fs.Files["/sys/class/leds/beaglebone:green:usr1/brightness"].Contents,
		"1",
	)

	a.DigitalWrite("P9_12", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio60/value"].Contents, "1")

	gobottest.Assert(t, a.DigitalWrite("P9_99", 1), errors.New("Not a valid pin"))

	fs.Files["/sys/class/gpio/gpio10/value"].Contents = "1"
	i, _ = a.DigitalRead("P8_31")
	gobottest.Assert(t, i, 1)

	// I2c
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 2)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}
