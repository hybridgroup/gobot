package beaglebone

import (
	"strings"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

func TestBeagleboneAdaptor(t *testing.T) {
	glob = func(pattern string) (matches []string, err error) {
		return make([]string, 2), nil
	}
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
		"/sys/devices/bone_capemgr.4",
		"/sys/devices/ocp.3",
		"/sys/devices/ocp.3/gpio-leds.8/leds/beaglebone:green:usr1/brightness",
		"/sys/devices/ocp.3/helper.5",
		"/sys/devices/ocp.3/helper.5/AIN1",
		"/sys/devices/ocp.3/pwm_test_P9_14.5",
		"/sys/devices/ocp.3/pwm_test_P9_14.5/run",
		"/sys/devices/ocp.3/pwm_test_P9_14.5/period",
		"/sys/devices/ocp.3/pwm_test_P9_14.5/duty",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	sysfs.SetFilesystem(fs)
	a := NewBeagleboneAdaptor("myAdaptor")
	a.slots = "/sys/devices/bone_capemgr.4"
	a.ocp = "/sys/devices/ocp.3"

	a.Connect()

	a.helper = "/sys/devices/ocp.3/helper.5"

	// PWM
	glob = func(pattern string) (matches []string, err error) {
		pattern = strings.TrimSuffix(pattern, "*")
		return []string{pattern + "5"}, nil
	}

	a.PwmWrite("P9_14", 175)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/ocp.3/pwm_test_P9_14.5/period"].Contents,
		"500000",
	)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/ocp.3/pwm_test_P9_14.5/duty"].Contents,
		"156862",
	)

	a.AnalogWrite("P9_14", 175)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/ocp.3/pwm_test_P9_14.5/period"].Contents,
		"500000",
	)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/ocp.3/pwm_test_P9_14.5/duty"].Contents,
		"156862",
	)

	a.ServoWrite("P9_14", 100)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/ocp.3/pwm_test_P9_14.5/period"].Contents,
		"20000000",
	)
	gobot.Assert(
		t,
		fs.Files["/sys/devices/ocp.3/pwm_test_P9_14.5/duty"].Contents,
		"17222222",
	)

	// Analog
	fs.Files["/sys/devices/ocp.3/helper.5/AIN1"].Contents = "567\n"
	i, _ := a.AnalogRead("P9_40")
	gobot.Assert(t, i, 567)

	// DigitalIO
	a.DigitalWrite("usr1", 1)
	gobot.Assert(t,
		fs.Files["/sys/devices/ocp.3/gpio-leds.8/leds/beaglebone:green:usr1/brightness"].Contents,
		"1",
	)

	a.DigitalWrite("P9_12", 1)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio60/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio10/value"].Contents = "1"
	i, _ = a.DigitalRead("P8_31")
	gobot.Assert(t, i, 1)

	// I2c
	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.I2cStart(0xff)

	a.I2cWrite([]byte{0x00, 0x01})
	data, _ := a.I2cRead(2)
	gobot.Assert(t, data, []byte{0x00, 0x01})

	gobot.Assert(t, a.Finalize(), nil)
}
