package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

// --------- HELPERS
func initTestBlinkMDriver() *BlinkMDriver {
	return NewBlinkMDriver(newI2cTestAdaptor("adaptor"), "bot")
}

func initTestMockedBlinkMDriver() (*BlinkMDriver, *I2cInterfaceClient) {
	inter := NewI2cInterfaceClient()
	return NewBlinkMDriver(inter, "bot"), inter
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
func TestNewBlinkMDriverCommandsDefinition(t *testing.T) {

	blinkM := initTestBlinkMDriver()

	errorMessage := func(command string) string {
		return command + " command should exist."
	}

	assertCommandExists := func(command string) {
		if blinkM.Driver.Command(command) == nil {
			t.Errorf(errorMessage(command))
		}
	}

	assertCommandExists("Rgb")
	assertCommandExists("Fade")
	assertCommandExists("FirmwareVersion")
	assertCommandExists("Color")
}

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
	blinkM, inter := initTestMockedBlinkMDriver()

	// Expectations
	inter.When("I2cWrite", []byte("Z")).Times(1)
	inter.When("I2cRead", uint(2)).Return([]byte{99, 1}).Times(1)

	param := make(map[string]interface{})

	result := blinkM.Driver.Command("FirmwareVersion")(param)

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}

	gobot.Assert(t, result, blinkM.FirmwareVersion())
}

func TestNewBlinkMDriverCommands_Color(t *testing.T) {
	blinkM, inter := initTestMockedBlinkMDriver()

	// Expectations
	inter.When("I2cWrite", []byte("g")).Times(1)
	inter.When("I2cRead", uint(3)).Return([]byte{99, 1, 2}).Times(1)

	param := make(map[string]interface{})

	result := blinkM.Driver.Command("Color")(param)

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}

	gobot.Assert(t, result, blinkM.Color())
}

// Methods
func TestBlinkMAdaptor(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobot.Assert(t, blinkM.adaptor(), blinkM.Adaptor())
}

func TestBlinkMDriverStart(t *testing.T) {
	blinkM, inter := initTestMockedBlinkMDriver()

	// Expectations
	inter.When("I2cStart", uint8(0x9)).Times(1)
	inter.When("I2cWrite", []byte("o")).Times(1)

	//// call to Rgb
	inter.When("I2cWrite", []byte("n")).Times(1)
	inter.When("I2cWrite", []byte{0, 0, 0}).Times(1)

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

func TestBlinkMDriverRgb(t *testing.T) {
	blinkM, inter := initTestMockedBlinkMDriver()

	// Expectations
	inter.When("I2cWrite", []byte("n")).Times(1)
	inter.When("I2cWrite", []byte{red, green, blue}).Times(1)

	blinkM.Rgb(red, green, blue)

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}
}

func TestBlinkMDriverFade(t *testing.T) {
	blinkM, inter := initTestMockedBlinkMDriver()

	inter.When("I2cWrite", []byte("c")).Times(1)
	inter.When("I2cWrite", []byte{red, green, blue}).Times(1)

	blinkM.Fade(red, green, blue)

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}
}

func TestBlinkMDriverFirmwareVersion(t *testing.T) {
	blinkM, inter := initTestMockedBlinkMDriver()

	// Expectations
	inter.When("I2cWrite", []byte("Z")).Times(1)
	inter.When("I2cRead", uint(2)).Return([]byte{99, 1}).Times(1)

	result := blinkM.FirmwareVersion()

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}

	gobot.Assert(t, result, "99.1")

	//// when len(data) is not 2
	blinkM, inter = initTestMockedBlinkMDriver()

	inter.When("I2cWrite", []byte("Z"))
	inter.When("I2cRead", uint(2)).Return([]byte{99}) // data length not 2 but 1

	result = blinkM.FirmwareVersion()

	gobot.Assert(t, result, "")
}

func TestBlinkMDriverColor(t *testing.T) {
	blinkM, inter := initTestMockedBlinkMDriver()

	// Expectations
	inter.When("I2cWrite", []byte("g")).Times(1)
	inter.When("I2cRead", uint(3)).Return([]byte{99, 1, 2}).Times(1)

	result := blinkM.Color()

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}

	gobot.Assert(t, result, []byte{99, 1, 2})

	//// when len(data) is not 3
	blinkM, inter = initTestMockedBlinkMDriver()

	inter.When("I2cWrite", []byte("g"))
	inter.When("I2cRead", uint(3)).Return([]byte{99, 2}) // data length not 3 but 2

	result = blinkM.Color()

	gobot.Assert(t, result, []byte{})

}
