package i2c

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
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

	result := blinkM.Command("Rgb")(rgb)
	gobottest.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_Fade(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	result := blinkM.Command("Fade")(rgb)
	gobottest.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_FirmwareVersion(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	param := make(map[string]interface{})

	// When len(data) is 2
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99, 1}, nil
	}

	result := blinkM.Command("FirmwareVersion")(param)

	version, _ := blinkM.FirmwareVersion()
	gobottest.Assert(t, result.(map[string]interface{})["version"].(string), version)

	// When len(data) is not 2
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99}, nil
	}
	result = blinkM.Command("FirmwareVersion")(param)

	version, _ = blinkM.FirmwareVersion()
	gobottest.Assert(t, result.(map[string]interface{})["version"].(string), version)
}

func TestNewBlinkMDriverCommands_Color(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	param := make(map[string]interface{})

	result := blinkM.Command("Color")(param)

	color, _ := blinkM.Color()
	gobottest.Assert(t, result.(map[string]interface{})["color"].([]byte), color)
}

// Methods
func TestBlinkMDriver(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobottest.Assert(t, blinkM.Name(), "bot")
	gobottest.Assert(t, blinkM.Connection().Name(), "adaptor")
}

func TestBlinkMDriverStart(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, len(blinkM.Start()), 0)

	adaptor.i2cStartImpl = func() error {
		return errors.New("start error")
	}

	gobottest.Assert(t, blinkM.Start()[0], errors.New("start error"))
	adaptor.i2cStartImpl = func() error {
		return nil
	}
	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	gobottest.Assert(t, blinkM.Start()[0], errors.New("write error"))
}

func TestBlinkMDriverHalt(t *testing.T) {
	blinkM := initTestBlinkMDriver()
	gobottest.Assert(t, len(blinkM.Halt()), 0)
}

func TestBlinkMDriverFirmwareVersion(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	// when len(data) is 2
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99, 1}, nil
	}

	version, _ := blinkM.FirmwareVersion()
	gobottest.Assert(t, version, "99.1")

	// when len(data) is not 2
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99}, nil
	}

	version, _ = blinkM.FirmwareVersion()
	gobottest.Assert(t, version, "")

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}

	version, err := blinkM.FirmwareVersion()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestBlinkMDriverColor(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	// when len(data) is 3
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99, 1, 2}, nil
	}

	color, _ := blinkM.Color()
	gobottest.Assert(t, color, []byte{99, 1, 2})

	// when len(data) is not 3
	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99}, nil
	}

	color, _ = blinkM.Color()
	gobottest.Assert(t, color, []byte{})

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}

	color, err := blinkM.Color()
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMDriverFade(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}

	err := blinkM.Fade(100, 100, 100)
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMDriverRGB(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}

	err := blinkM.Rgb(100, 100, 100)
	gobottest.Assert(t, err, errors.New("write error"))

}
