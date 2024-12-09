package edison

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
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
	_ i2c.Connector               = (*Adaptor)(nil)
)

var testPinFiles = []string{
	"/sys/bus/iio/devices/iio:device1/in_voltage0_raw",
	"/sys/kernel/debug/gpio_debug/gpio111/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio115/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio114/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio109/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio131/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio129/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio28/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio27/current_pinmux",
	"/sys/class/pwm/pwmchip0/export",
	"/sys/class/pwm/pwmchip0/unexport",
	"/sys/class/pwm/pwmchip0/pwm1/duty_cycle",
	"/sys/class/pwm/pwmchip0/pwm1/period",
	"/sys/class/pwm/pwmchip0/pwm1/enable",
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio13/value",
	"/sys/class/gpio/gpio13/direction",
	"/sys/class/gpio/gpio40/value",
	"/sys/class/gpio/gpio40/direction",
	"/sys/class/gpio/gpio128/value",
	"/sys/class/gpio/gpio128/direction",
	"/sys/class/gpio/gpio221/value",
	"/sys/class/gpio/gpio221/direction",
	"/sys/class/gpio/gpio243/value",
	"/sys/class/gpio/gpio243/direction",
	"/sys/class/gpio/gpio229/value",
	"/sys/class/gpio/gpio229/direction",
	"/sys/class/gpio/gpio253/value",
	"/sys/class/gpio/gpio253/direction",
	"/sys/class/gpio/gpio261/value",
	"/sys/class/gpio/gpio261/direction",
	"/sys/class/gpio/gpio214/value",
	"/sys/class/gpio/gpio214/direction",
	"/sys/class/gpio/gpio14/direction",
	"/sys/class/gpio/gpio14/value",
	"/sys/class/gpio/gpio165/direction",
	"/sys/class/gpio/gpio165/value",
	"/sys/class/gpio/gpio212/direction",
	"/sys/class/gpio/gpio212/value",
	"/sys/class/gpio/gpio213/direction",
	"/sys/class/gpio/gpio213/value",
	"/sys/class/gpio/gpio236/direction",
	"/sys/class/gpio/gpio236/value",
	"/sys/class/gpio/gpio237/direction",
	"/sys/class/gpio/gpio237/value",
	"/sys/class/gpio/gpio204/direction",
	"/sys/class/gpio/gpio204/value",
	"/sys/class/gpio/gpio205/direction",
	"/sys/class/gpio/gpio205/value",
	"/sys/class/gpio/gpio263/direction",
	"/sys/class/gpio/gpio263/value",
	"/sys/class/gpio/gpio262/direction",
	"/sys/class/gpio/gpio262/value",
	"/sys/class/gpio/gpio240/direction",
	"/sys/class/gpio/gpio240/value",
	"/sys/class/gpio/gpio241/direction",
	"/sys/class/gpio/gpio241/value",
	"/sys/class/gpio/gpio242/direction",
	"/sys/class/gpio/gpio242/value",
	"/sys/class/gpio/gpio218/direction",
	"/sys/class/gpio/gpio218/value",
	"/sys/class/gpio/gpio250/direction",
	"/sys/class/gpio/gpio250/value",
	"/dev/i2c-6",
}

var pwmMockPathsMux13Arduino = []string{
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio109/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio111/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio114/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio115/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio129/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio131/current_pinmux",
	"/sys/class/gpio/gpio13/direction",
	"/sys/class/gpio/gpio13/value",
	"/sys/class/gpio/gpio214/direction",
	"/sys/class/gpio/gpio214/value",
	"/sys/class/gpio/gpio221/direction",
	"/sys/class/gpio/gpio221/value",
	"/sys/class/gpio/gpio240/direction",
	"/sys/class/gpio/gpio240/value",
	"/sys/class/gpio/gpio241/direction",
	"/sys/class/gpio/gpio241/value",
	"/sys/class/gpio/gpio242/direction",
	"/sys/class/gpio/gpio242/value",
	"/sys/class/gpio/gpio243/direction",
	"/sys/class/gpio/gpio243/value",
	"/sys/class/gpio/gpio253/direction",
	"/sys/class/gpio/gpio253/value",
	"/sys/class/gpio/gpio262/direction",
	"/sys/class/gpio/gpio262/value",
	"/sys/class/gpio/gpio263/direction",
	"/sys/class/gpio/gpio263/value",
}

var pwmMockPathsMux13ArduinoI2c = []string{
	"/dev/i2c-6",
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio27/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio28/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio109/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio111/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio114/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio115/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio129/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio131/current_pinmux",
	"/sys/class/gpio/gpio13/direction",
	"/sys/class/gpio/gpio13/value",
	"/sys/class/gpio/gpio14/direction",
	"/sys/class/gpio/gpio14/value",
	"/sys/class/gpio/gpio28/direction",
	"/sys/class/gpio/gpio28/value",
	"/sys/class/gpio/gpio165/direction",
	"/sys/class/gpio/gpio165/value",
	"/sys/class/gpio/gpio212/direction",
	"/sys/class/gpio/gpio212/value",
	"/sys/class/gpio/gpio213/direction",
	"/sys/class/gpio/gpio213/value",
	"/sys/class/gpio/gpio214/direction",
	"/sys/class/gpio/gpio214/value",
	"/sys/class/gpio/gpio221/direction",
	"/sys/class/gpio/gpio221/value",
	"/sys/class/gpio/gpio236/direction",
	"/sys/class/gpio/gpio236/value",
	"/sys/class/gpio/gpio237/direction",
	"/sys/class/gpio/gpio237/value",
	"/sys/class/gpio/gpio204/value",
	"/sys/class/gpio/gpio204/direction",
	"/sys/class/gpio/gpio205/value",
	"/sys/class/gpio/gpio205/direction",
	"/sys/class/gpio/gpio240/direction",
	"/sys/class/gpio/gpio240/value",
	"/sys/class/gpio/gpio241/direction",
	"/sys/class/gpio/gpio241/value",
	"/sys/class/gpio/gpio242/direction",
	"/sys/class/gpio/gpio242/value",
	"/sys/class/gpio/gpio243/direction",
	"/sys/class/gpio/gpio243/value",
	"/sys/class/gpio/gpio253/direction",
	"/sys/class/gpio/gpio253/value",
	"/sys/class/gpio/gpio262/direction",
	"/sys/class/gpio/gpio262/value",
	"/sys/class/gpio/gpio263/direction",
	"/sys/class/gpio/gpio263/value",
}

var pwmMockPathsMux13 = []string{
	"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
	"/sys/class/gpio/export",
	"/sys/class/gpio/gpio13/direction",
	"/sys/class/gpio/gpio13/value",
	"/sys/class/gpio/gpio221/direction",
	"/sys/class/gpio/gpio221/value",
	"/sys/class/gpio/gpio253/direction",
	"/sys/class/gpio/gpio253/value",
	"/sys/class/pwm/pwmchip0/export",
	"/sys/class/pwm/pwmchip0/unexport",
	"/sys/class/pwm/pwmchip0/pwm1/duty_cycle",
	"/sys/class/pwm/pwmchip0/pwm1/period",
	"/sys/class/pwm/pwmchip0/pwm1/enable",
}

var pwmMockPathsMux40 = []string{
	"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio40/value",
	"/sys/class/gpio/gpio40/direction",
	"/sys/class/gpio/gpio229/value", // resistor
	"/sys/class/gpio/gpio229/direction",
	"/sys/class/gpio/gpio243/value",
	"/sys/class/gpio/gpio243/direction",
	"/sys/class/gpio/gpio261/value", // level shifter
	"/sys/class/gpio/gpio261/direction",
}

func initTestAdaptorWithMockedFilesystem(boardType string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor(boardType)
	fs := a.sys.UseMockFilesystem(testPinFiles)
	fs.Files["/sys/class/pwm/pwmchip0/pwm1/period"].Contents = "5000"
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()

	assert.True(t, strings.HasPrefix(a.Name(), "Edison"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestConnect(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem("arduino")

	assert.Equal(t, 6, a.DefaultI2cBus())
	assert.Equal(t, "arduino", a.board)
	require.NoError(t, a.Connect())
}

func TestArduinoSetupFail263(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	delete(fs.Files, "/sys/class/gpio/gpio263/direction")

	err := a.arduinoSetup()
	require.ErrorContains(t, err, "/sys/class/gpio/gpio263/direction: no such file")
}

func TestArduinoSetupFail240(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	delete(fs.Files, "/sys/class/gpio/gpio240/direction")

	err := a.arduinoSetup()
	require.ErrorContains(t, err, "/sys/class/gpio/gpio240/direction: no such file")
}

func TestArduinoSetupFail111(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio111/current_pinmux")

	err := a.arduinoSetup()
	require.ErrorContains(t, err, "/sys/kernel/debug/gpio_debug/gpio111/current_pinmux: no such file")
}

func TestArduinoSetupFail131(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio131/current_pinmux")

	err := a.arduinoSetup()
	require.ErrorContains(t, err, "/sys/kernel/debug/gpio_debug/gpio131/current_pinmux: no such file")
}

func TestArduinoI2CSetupFailTristate(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	require.NoError(t, a.arduinoSetup())

	fs.WithWriteError = true
	err := a.arduinoI2CSetup()
	require.ErrorContains(t, err, "write error")
}

func TestArduinoI2CSetupFail14(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	require.NoError(t, a.arduinoSetup())
	delete(fs.Files, "/sys/class/gpio/gpio14/direction")

	err := a.arduinoI2CSetup()
	require.ErrorContains(t, err, "/sys/class/gpio/gpio14/direction: no such file")
}

func TestArduinoI2CSetupUnexportFail(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	require.NoError(t, a.arduinoSetup())
	delete(fs.Files, "/sys/class/gpio/unexport")

	err := a.arduinoI2CSetup()
	require.ErrorContains(t, err, "/sys/class/gpio/unexport: no such file")
}

func TestArduinoI2CSetupFail236(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	require.NoError(t, a.arduinoSetup())
	delete(fs.Files, "/sys/class/gpio/gpio236/direction")

	err := a.arduinoI2CSetup()
	require.ErrorContains(t, err, "/sys/class/gpio/gpio236/direction: no such file")
}

func TestArduinoI2CSetupFail28(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	require.NoError(t, a.arduinoSetup())
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio28/current_pinmux")

	err := a.arduinoI2CSetup()
	require.ErrorContains(t, err, "/sys/kernel/debug/gpio_debug/gpio28/current_pinmux: no such file")
}

func TestConnectArduinoError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithWriteError = true

	err := a.Connect()
	require.ErrorContains(t, err, "write error")
}

func TestConnectArduinoWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithWriteError = true

	err := a.Connect()
	require.ErrorContains(t, err, "write error")
}

func TestConnectSparkfun(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem("sparkfun")

	require.NoError(t, a.Connect())
	assert.Equal(t, 1, a.DefaultI2cBus())
	assert.Equal(t, "sparkfun", a.board)
}

func TestConnectMiniboard(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem("miniboard")

	require.NoError(t, a.Connect())
	assert.Equal(t, 1, a.DefaultI2cBus())
	assert.Equal(t, "miniboard", a.board)
}

func TestConnectUnknown(t *testing.T) {
	a := NewAdaptor("wha")

	err := a.Connect()
	require.ErrorContains(t, err, "Unknown board type: wha")
}

func TestFinalize(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	_ = a.DigitalWrite("3", 1)
	require.NoError(t, a.PwmWrite("5", 100))

	_, _ = a.GetI2cConnection(0xff, 6)
	require.NoError(t, a.Finalize())

	// assert that finalize after finalize is working
	require.NoError(t, a.Finalize())

	// assert that re-connect is working
	require.NoError(t, a.Connect())
	// remove one file to force Finalize error
	delete(fs.Files, "/sys/class/gpio/unexport")
	err := a.Finalize()
	require.ErrorContains(t, err, "1 error occurred")
	require.ErrorContains(t, err, "/sys/class/gpio/unexport")
}

func TestFinalizeError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	require.NoError(t, a.PwmWrite("5", 100))

	fs.WithWriteError = true
	err := a.Finalize()
	require.ErrorContains(t, err, "6 errors occurred")
	require.ErrorContains(t, err, "write error")
	require.ErrorContains(t, err, "SetEnabled(false) failed for id 1 with write error")
	require.ErrorContains(t, err, "Unexport() failed for id 1 with write error")
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	require.NoError(t, a.DigitalWrite("13", 1))
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio40/value"].Contents)

	require.NoError(t, a.DigitalWrite("2", 0))
	i, err := a.DigitalRead("2")
	require.NoError(t, err)
	assert.Equal(t, 0, i)
}

func TestDigitalPinInFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio40/value")
	delete(fs.Files, "/sys/class/gpio/gpio40/direction")
	_ = a.Connect() // we drop this error

	_, err := a.DigitalPin("13")
	require.ErrorContains(t, err, "no such file")
}

func TestDigitalPinInResistorFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio229/value")
	delete(fs.Files, "/sys/class/gpio/gpio229/direction")
	_ = a.Connect() // we drop this error

	_, err := a.DigitalPin("13")
	require.ErrorContains(t, err, "no such file")
}

func TestDigitalPinInLevelShifterFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio261/value")
	delete(fs.Files, "/sys/class/gpio/gpio261/direction")
	_ = a.Connect() // we drop this error

	_, err := a.DigitalPin("13")
	require.ErrorContains(t, err, "no such file")
}

func TestDigitalPinInMuxFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio243/value")
	delete(fs.Files, "/sys/class/gpio/gpio243/direction")
	_ = a.Connect() // we drop this error

	_, err := a.DigitalPin("13")
	require.ErrorContains(t, err, "no such file")
}

func TestDigitalWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithWriteError = true

	err := a.DigitalWrite("13", 1)
	require.ErrorContains(t, err, "write error")
}

func TestDigitalReadWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithWriteError = true

	_, err := a.DigitalRead("13")
	require.ErrorContains(t, err, "write error")
}

func TestPwm(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")

	err := a.PwmWrite("5", 100)
	require.NoError(t, err)
	assert.Equal(t, "1960", fs.Files["/sys/class/pwm/pwmchip0/pwm1/duty_cycle"].Contents)

	err = a.PwmWrite("7", 100)
	require.ErrorContains(t, err, "'7' is not a valid id for a PWM pin")
}

func TestPwmExportError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux13Arduino)
	delete(fs.Files, "/sys/class/pwm/pwmchip0/export")
	err := a.Connect()
	require.NoError(t, err)

	err = a.PwmWrite("5", 100)
	require.ErrorContains(t, err, "/sys/class/pwm/pwmchip0/export: no such file")
}

func TestPwmEnableError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux13)
	delete(fs.Files, "/sys/class/pwm/pwmchip0/pwm1/enable")
	_ = a.Connect() // we drop this error

	err := a.PwmWrite("5", 100)
	require.ErrorContains(t, err, "/sys/class/pwm/pwmchip0/pwm1/enable: no such file")
}

func TestPwmWritePinError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithWriteError = true

	err := a.PwmWrite("5", 100)
	require.ErrorContains(t, err, "write error")
}

func TestPwmWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithWriteError = true

	err := a.PwmWrite("5", 100)
	require.ErrorContains(t, err, "write error")
}

func TestPwmReadError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithReadError = true

	err := a.PwmWrite("5", 100)
	require.ErrorContains(t, err, "read error")
}

func TestAnalog(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.Files["/sys/bus/iio/devices/iio:device1/in_voltage0_raw"].Contents = "1000\n"

	i, err := a.AnalogRead("0")
	require.NoError(t, err)
	assert.Equal(t, 250, i)
}

func TestAnalogError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem("arduino")
	fs.WithReadError = true

	_, err := a.AnalogRead("0")
	require.ErrorContains(t, err, "read error")
}

func TestI2cWorkflow(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem("arduino")
	a.sys.UseMockSyscall()

	con, err := a.GetI2cConnection(0xff, 6)
	require.NoError(t, err)

	_, err = con.Write([]byte{0x00, 0x01})
	require.NoError(t, err)

	data := []byte{42, 42}
	_, err = con.Read(data)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x01}, data)

	require.NoError(t, a.Finalize())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux13ArduinoI2c)
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 6)
	require.NoError(t, err)
	_, err = con.Write([]byte{0x0A})
	require.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	require.Error(t, err)
	require.ErrorContains(t, err, "close error")
}

func Test_validateAndSetupI2cBusNumber(t *testing.T) {
	tests := map[string]struct {
		board   string
		busNr   int
		wantErr string
	}{
		"arduino_number_negative_error": {
			busNr:   -1,
			wantErr: "Unsupported I2C bus '-1'",
		},
		"arduino_number_1_error": {
			busNr:   1,
			wantErr: "Unsupported I2C bus '1'",
		},
		"arduino_number_6_ok": {
			busNr: 6,
		},
		"sparkfun_number_negative_error": {
			board:   "sparkfun",
			busNr:   -1,
			wantErr: "Unsupported I2C bus '-1'",
		},
		"sparkfun_number_1_ok": {
			board: "sparkfun",
			busNr: 1,
		},
		"miniboard_number_6_error": {
			board:   "miniboard",
			busNr:   6,
			wantErr: "Unsupported I2C bus '6'",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor(tc.board)
			a.sys.UseMockFilesystem(pwmMockPathsMux13ArduinoI2c)
			_ = a.Connect()
			// act
			err := a.validateAndSetupI2cBusNumber(tc.busNr)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
