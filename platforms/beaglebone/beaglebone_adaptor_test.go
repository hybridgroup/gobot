package beaglebone

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gobot.PWMPinnerProvider     = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ aio.AnalogReader            = (*Adaptor)(nil)
	_ gpio.PwmWriter              = (*Adaptor)(nil)
	_ gpio.ServoWriter            = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
	_ spi.Connector               = (*Adaptor)(nil)
)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestPWM(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/ocp/ocp:P9_22_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_21_pinmux/state",
		"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/export",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/unexport",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/enable",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/period",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/polarity",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/enable",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/polarity",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)
	fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle"].Contents = "0"
	fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period"].Contents = "0"

	assert.ErrorContains(t, a.PwmWrite("P9_99", 175), "'P9_99' is not a valid id for a PWM pin")
	_ = a.PwmWrite("P9_21", 175)
	assert.Equal(
		t,
		"500000",
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period"].Contents,
	)
	assert.Equal(
		t,
		"343137",
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle"].Contents,
	)

	_ = a.ServoWrite("P9_21", 100)
	assert.Equal(
		t,
		"500000",
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period"].Contents,
	)
	assert.Equal(
		t,
		"66666",
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle"].Contents,
	)

	assert.NoError(t, a.Finalize())
}

func TestAnalog(t *testing.T) {
	mockPaths := []string{
		"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	fs.Files["/sys/bus/iio/devices/iio:device0/in_voltage1_raw"].Contents = "567\n"
	i, err := a.AnalogRead("P9_40")
	assert.Equal(t, 567, i)
	assert.NoError(t, err)

	_, err = a.AnalogRead("P9_99")
	assert.ErrorContains(t, err, "Not a valid analog pin")

	fs.WithReadError = true
	_, err = a.AnalogRead("P9_40")
	assert.ErrorContains(t, err, "read error")
	fs.WithReadError = false

	assert.NoError(t, a.Finalize())
}

func TestDigitalIO(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/ocp/ocp:P8_07_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_11_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
		"/sys/class/leds/beaglebone:green:usr1/brightness",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/class/gpio/gpio66/value",
		"/sys/class/gpio/gpio66/direction",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
		"/sys/class/gpio/gpio30/value",
		"/sys/class/gpio/gpio30/direction",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	// DigitalIO
	_ = a.DigitalWrite("usr1", 1)
	assert.Equal(t,
		"1",
		fs.Files["/sys/class/leds/beaglebone:green:usr1/brightness"].Contents,
	)

	// no such LED
	err := a.DigitalWrite("usr10101", 1)
	assert.ErrorContains(t, err, " : /sys/class/leds/beaglebone:green:usr10101/brightness: no such file")

	_ = a.DigitalWrite("P9_12", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio60/value"].Contents)

	assert.ErrorContains(t, a.DigitalWrite("P9_99", 1), "'P9_99' is not a valid id for a digital pin")

	_, err = a.DigitalRead("P9_99")
	assert.ErrorContains(t, err, "'P9_99' is not a valid id for a digital pin")

	fs.Files["/sys/class/gpio/gpio66/value"].Contents = "1"
	i, err := a.DigitalRead("P8_07")
	assert.Equal(t, 1, i)
	assert.NoError(t, err)

	assert.NoError(t, a.Finalize())
}

func TestName(t *testing.T) {
	a := NewAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Beaglebone"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestAnalogReadFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/whatever",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	_, err := a.AnalogRead("P9_40")
	assert.Contains(t, err.Error(), "/sys/bus/iio/devices/iio:device0/in_voltage1_raw: no such file")
}

func TestDigitalPinDirectionFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio60/value",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	err := a.DigitalWrite("P9_12", 1)
	assert.Contains(t, err.Error(), "/sys/class/gpio/gpio60/direction: no such file")

	// no pin added after previous problem, so no pin to unexport in finalize
	err = a.Finalize()
	assert.Equal(t, err, nil)
}

func TestDigitalPinFinalizeFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	err := a.DigitalWrite("P9_12", 1)
	assert.NoError(t, err)

	err = a.Finalize()
	assert.Contains(t, err.Error(), "/sys/class/gpio/unexport: no such file")
}

func TestPocketName(t *testing.T) {
	a := NewPocketBeagleAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "PocketBeagle"))
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultChipNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, 8, a.SpiDefaultBitCount())
	assert.Equal(t, int64(500000), a.SpiDefaultMaxSpeed())
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 2, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-2"})
	assert.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 2)
	assert.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	assert.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	assert.Contains(t, err.Error(), "close error")
}

func Test_validateSpiBusNumber(t *testing.T) {
	tests := map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("Bus number -1 out of range"),
		},
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_ok": {
			busNr: 1,
		},
		"number_2_error": {
			busNr:   2,
			wantErr: fmt.Errorf("Bus number 2 out of range"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateSpiBusNumber(tc.busNr)
			// assert
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_validateI2cBusNumber(t *testing.T) {
	tests := map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("Bus number -1 out of range"),
		},
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_error": {
			busNr:   1,
			wantErr: fmt.Errorf("Bus number 1 out of range"),
		},
		"number_2_ok": {
			busNr: 2,
		},
		"number_3_error": {
			busNr:   3,
			wantErr: fmt.Errorf("Bus number 3 out of range"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateI2cBusNumber(tc.busNr)
			// assert
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_translateAndMuxPWMPin(t *testing.T) {
	// arrange
	mockPaths := []string{
		"/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4/",
		"/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2/",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/",
		"/sys/devices/platform/ocp/48300000.epwmss/48300100.ecap/pwm/pwmchip0/",
	}
	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	tests := map[string]struct {
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"P8_13": {
			wantDir:     "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4",
			wantChannel: 1,
		},
		"P8_19": {
			wantDir:     "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4",
			wantChannel: 0,
		},
		"P9_14": {
			wantDir:     "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2",
			wantChannel: 0,
		},
		"P9_16": {
			wantDir:     "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2",
			wantChannel: 1,
		},
		"P9_21": {
			wantDir:     "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0",
			wantChannel: 1,
		},
		"P9_22": {
			wantDir:     "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0",
			wantChannel: 0,
		},
		"P9_42": {
			wantDir:     "/sys/devices/platform/ocp/48300000.epwmss/48300100.ecap/pwm/pwmchip0",
			wantChannel: 0,
		},
		"P9_99": {
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'P9_99' is not a valid id for a PWM pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			muxPath := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", name)
			fs.Add(muxPath)
			// act
			path, channel, err := a.translateAndMuxPWMPin(name)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, path)
			assert.Equal(t, tc.wantChannel, channel)
			if tc.wantErr == nil {
				assert.Equal(t, "pwm", fs.Files[muxPath].Contents)
			}
		})
	}
}
