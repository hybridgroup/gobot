package dragonboard

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
)

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
)

func initTestAdaptor() *Adaptor {
	a := NewAdaptor()
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestName(t *testing.T) {
	a := initTestAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "DragonBoard"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestDigitalIO(t *testing.T) {
	a := initTestAdaptor()
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)

	_ = a.DigitalWrite("GPIO_B", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio12/value"].Contents)

	fs.Files["/sys/class/gpio/gpio36/value"].Contents = "1"
	i, _ := a.DigitalRead("GPIO_A")
	assert.Equal(t, 1, i)

	require.ErrorContains(t, a.DigitalWrite("GPIO_M", 1), "'GPIO_M' is not a valid id for a digital pin")
	require.NoError(t, a.Finalize())
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a := initTestAdaptor()
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)

	require.NoError(t, a.Connect())
	require.NoError(t, a.DigitalWrite("GPIO_B", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	require.ErrorContains(t, err, "write error")
}

func TestI2cDefaultBus(t *testing.T) {
	a := initTestAdaptor()
	assert.Equal(t, 0, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 1)
	require.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	require.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	require.ErrorContains(t, err, "close error")
}
