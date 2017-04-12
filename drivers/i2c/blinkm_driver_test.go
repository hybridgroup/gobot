package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BlinkMDriver)(nil)

// --------- HELPERS
func initTestBlinkMDriver() (driver *BlinkMDriver) {
	driver, _ = initTestBlinkDriverWithStubbedAdaptor()
	return
}

func initTestBlinkDriverWithStubbedAdaptor() (*BlinkMDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBlinkMDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewBlinkMDriver(t *testing.T) {
	// Does it return a pointer to an instance of BlinkMDriver?
	var bm interface{} = NewBlinkMDriver(newI2cTestAdaptor())
	_, ok := bm.(*BlinkMDriver)
	if !ok {
		t.Errorf("NewBlinkMDriver() should have returned a *BlinkMDriver")
	}
}

// Commands
func TestNewBlinkMDriverCommands_Rgb(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobottest.Assert(t, blinkM.Start(), nil)

	result := blinkM.Command("Rgb")(rgb)
	gobottest.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_Fade(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobottest.Assert(t, blinkM.Start(), nil)

	result := blinkM.Command("Fade")(rgb)
	gobottest.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_FirmwareVersion(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, blinkM.Start(), nil)

	param := make(map[string]interface{})

	// When len(data) is 2
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	result := blinkM.Command("FirmwareVersion")(param)

	version, _ := blinkM.FirmwareVersion()
	gobottest.Assert(t, result.(map[string]interface{})["version"].(string), version)

	// When len(data) is not 2
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}
	result = blinkM.Command("FirmwareVersion")(param)

	version, _ = blinkM.FirmwareVersion()
	gobottest.Assert(t, result.(map[string]interface{})["version"].(string), version)
}

func TestNewBlinkMDriverCommands_Color(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobottest.Assert(t, blinkM.Start(), nil)

	param := make(map[string]interface{})

	result := blinkM.Command("Color")(param)

	color, _ := blinkM.Color()
	gobottest.Assert(t, result.(map[string]interface{})["color"].([]byte), color)
}

// Methods
func TestBlinkMDriver(t *testing.T) {
	blinkM := initTestBlinkMDriver()

	gobottest.Assert(t, blinkM.Start(), nil)
	gobottest.Refute(t, blinkM.Connection(), nil)
}

func TestBlinkMDriverStart(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, blinkM.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, blinkM.Start(), errors.New("write error"))
}

func TestBlinkMDriverStartConnectError(t *testing.T) {
	d, adaptor := initTestBlinkDriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestBlinkMDriverHalt(t *testing.T) {
	blinkM := initTestBlinkMDriver()
	gobottest.Assert(t, blinkM.Halt(), nil)
}

func TestBlinkMDriverFirmwareVersion(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, blinkM.Start(), nil)

	// when len(data) is 2
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	version, _ := blinkM.FirmwareVersion()
	gobottest.Assert(t, version, "99.1")

	// when len(data) is not 2
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	version, _ = blinkM.FirmwareVersion()
	gobottest.Assert(t, version, "")

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	version, err := blinkM.FirmwareVersion()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestBlinkMDriverColor(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, blinkM.Start(), nil)

	// when len(data) is 3
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1, 2})
		return 3, nil
	}

	color, _ := blinkM.Color()
	gobottest.Assert(t, color, []byte{99, 1, 2})

	// when len(data) is not 3
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	color, _ = blinkM.Color()
	gobottest.Assert(t, color, []byte{})

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	color, err := blinkM.Color()
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMDriverFade(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, blinkM.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := blinkM.Fade(100, 100, 100)
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMDriverRGB(t *testing.T) {
	blinkM, adaptor := initTestBlinkDriverWithStubbedAdaptor()

	gobottest.Assert(t, blinkM.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := blinkM.Rgb(100, 100, 100)
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMDriverSetName(t *testing.T) {
	d := initTestBlinkMDriver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestBlinkMDriverOptions(t *testing.T) {
	d := NewBlinkMDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
