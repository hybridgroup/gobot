package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

// --------- HELPERS
func initTestBlinkMDriver() (driver *BlinkMDriver) {
	driver, _ = initTestBlinkDriverWithStubbedAdaptor()
	return
}

func initTestBlinkDriverWithStubbedAdaptor() (*BlinkMDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewBlinkMDriver(adaptor, "bot"), adaptor
}

// --------- TESTS

func TestBlinkMDriver(t *testing.T) {
	// Does it implement gobot.DriverInterface?
	var _ gobot.DriverInterface = (*BlinkMDriver)(nil)

	// Does its adaptor implements the I2cInterface?
	driver := initTestBlinkMDriver()
	var _ I2cInterface = driver.adaptor()
}

func TestNewBlinkMDriver(t *testing.T) {
	// Does it return a pointer to an instance of BlinkMDriver?
	var bm interface{} = NewBlinkMDriver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := bm.(*BlinkMDriver)
	if !ok {
		t.Errorf("NewBlinkMDriver() should have returned a *BlinkMDriver")
	}
}

// Commands
func TestNewBlinkMDriverCommands_Rgb(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	result := blinkM.Driver.Command("Rgb")(rgb)
	gobot.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_Fade(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	result := blinkM.Driver.Command("Fade")(rgb)
	gobot.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_FirmwareVersion(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	param := make(map[string]interface{})

	// When len(data) is 2
	adaptor.i2cReadImpl = func() []byte {
		return []byte{99, 1}
	}

	result := blinkM.Driver.Command("FirmwareVersion")(param)

	gobot.Assert(t, result, blinkM.FirmwareVersion())

	// When len(data) is not 2
	adaptor.i2cReadImpl = func() []byte {
		return []byte{99}
	}
	result = blinkM.Driver.Command("FirmwareVersion")(param)

	gobot.Assert(t, result, blinkM.FirmwareVersion())
}

func TestNewBlinkMDriverCommands_Color(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	param := make(map[string]interface{})

	result := blinkM.Driver.Command("Color")(param)

	gobot.Assert(t, result, blinkM.Color())
}

// Methods
func TestBlinkMDriverStart(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobot.Assert(t, blinkM.Start(), true)
}

func TestBlinkMDriverInit(t *testing.T) {
	blinkM := initTestBlinkMDriver()
	gobot.Assert(t, blinkM.Init(), true)
}

func TestBlinkMDriverHalt(t *testing.T) {
	blinkM := initTestBlinkMDriver()
	gobot.Assert(t, blinkM.Halt(), true)
}

func TestBlinkMDriverFirmwareVersion(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	// when len(data) is 2
	adaptor.i2cReadImpl = func() []byte {
		return []byte{99, 1}
	}

	gobot.Assert(t, blinkM.FirmwareVersion(), "99.1")

	// when len(data) is not 2
	adaptor.i2cReadImpl = func() []byte {
		return []byte{99}
	}

	gobot.Assert(t, blinkM.FirmwareVersion(), "")
}

func TestBlinkMDriverColor(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	// when len(data) is 3
	adaptor.i2cReadImpl = func() []byte {
		return []byte{99, 1, 2}
	}

	gobot.Assert(t, blinkM.Color(), []byte{99, 1, 2})

	// when len(data) is not 3
	adaptor.i2cReadImpl = func() []byte {
		return []byte{99}
	}

	gobot.Assert(t, blinkM.Color(), []byte{})

}
