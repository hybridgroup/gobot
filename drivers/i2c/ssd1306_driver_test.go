package i2c

import (
	"errors"
	"fmt"
	"image"
	"reflect"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*SSD1306Driver)(nil)

func TestDisplayBuffer(t *testing.T) {
	width := 128
	height := 64
	size := 1024 // (width*height) / 8
	display := NewDisplayBuffer(width, height, 8)

	if display.Size() != size {
		t.Errorf("invalid Size() (%d, expected %d)",
			display.Size(), size)
	}
	if len(display.buffer) != size {
		t.Errorf("allocated buffer size invalid (%d, expected %d)",
			len(display.buffer), size)
	}

	gobottest.Assert(t, display.buffer[0], byte(0))
	gobottest.Assert(t, display.buffer[1], byte(0))

	display.SetPixel(0, 0, 1)
	display.SetPixel(1, 0, 1)
	display.SetPixel(2, 0, 1)
	display.SetPixel(0, 1, 1)
	gobottest.Assert(t, display.buffer[0], byte(3))
	gobottest.Assert(t, display.buffer[1], byte(1))

	display.SetPixel(0, 1, 0)
	gobottest.Assert(t, display.buffer[0], byte(1))
	gobottest.Assert(t, display.buffer[1], byte(1))
}

// --------- HELPERS
func initTestSSD1306Driver(width, height int, externalVCC bool) (driver *SSD1306Driver) {
	driver, _ = initTestSSD1306DriverWithStubbedAdaptor(width, height, externalVCC)
	return
}

func initTestSSD1306DriverWithStubbedAdaptor(width, height int, externalVCC bool) (*SSD1306Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewSSD1306Driver(adaptor, WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC)), adaptor
}

// --------- TESTS

func TestNewSSD1306Driver(t *testing.T) {
	// Does it return a pointer to an instance of SSD1306Driver?
	var bm interface{} = NewSSD1306Driver(newI2cTestAdaptor())
	_, ok := bm.(*SSD1306Driver)
	if !ok {
		t.Errorf("new should have returned a *SSD1306Driver")
	}

	b := NewSSD1306Driver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

// Methods

func TestSSD1306DriverStartDefaul(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	gobottest.Assert(t, s.Start(), nil)
}

func TestSSD1306DriverStart128x32(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 32, false)
	gobottest.Assert(t, s.Start(), nil)
}

func TestSSD1306DriverStart96x16(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(96, 16, false)
	gobottest.Assert(t, s.Start(), nil)
}

func TestSSD1306DriverStartExternalVCC(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 32, true)
	gobottest.Assert(t, s.Start(), nil)
}

func TestSSD1306StartConnectError(t *testing.T) {
	d, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestSSD1306StartSizeError(t *testing.T) {
	d, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 54, false)
	//adaptor.Testi2cConnectErr(false)
	gobottest.Assert(t, d.Start(), errors.New("128x54 resolution is unsupported, supported resolutions: 128x64, 128x32, 96x16"))
}

func TestSSD1306DriverHalt(t *testing.T) {
	s := initTestSSD1306Driver(128, 64, false)

	gobottest.Assert(t, s.Halt(), nil)
}

// Test Name & SetName
func TestSSD1306DriverName(t *testing.T) {
	s := initTestSSD1306Driver(96, 16, false)

	gobottest.Assert(t, strings.HasPrefix(s.Name(), "SSD1306"), true)
	s.SetName("Ole Oled")
	gobottest.Assert(t, s.Name(), "Ole Oled")
}

func TestSSD1306DriverOptions(t *testing.T) {
	s := NewSSD1306Driver(newI2cTestAdaptor(), WithBus(2), WithSSD1306DisplayHeight(32), WithSSD1306DisplayWidth(128))
	gobottest.Assert(t, s.GetBusOrDefault(1), 2)
	gobottest.Assert(t, s.displayHeight, 32)
	gobottest.Assert(t, s.displayWidth, 128)
}

func TestSSD1306DriverDisplay(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(96, 16, false)
	s.Start()
	gobottest.Assert(t, s.Display(), nil)
}

func TestSSD1306DriverShowImage(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	gobottest.Assert(t, s.ShowImage(img), errors.New("image must match display width and height: 128x64"))

	img = image.NewRGBA(image.Rect(0, 0, 128, 64))
	gobottest.Assert(t, s.ShowImage(img), nil)
}

func TestSSD1306DriverCommand(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, 0xFF}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.command(0xFF)
	gobottest.Assert(t, err, nil)
}

func TestSSD1306DriverCommands(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, 0x00, 0x80, 0xFF}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.commands([]byte{0x00, 0xFF})
	gobottest.Assert(t, err, nil)
}

func TestSSD1306DriverOn(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetDisplayOn}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.On()
	gobottest.Assert(t, err, nil)
}

func TestSSD1306DriverOff(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetDisplayOff}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.Off()
	gobottest.Assert(t, err, nil)
}

func TestSSD1306Reset(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expectedOff := []byte{0x80, ssd1306SetDisplayOff}
		expectedOn := []byte{0x80, ssd1306SetDisplayOn}
		if !reflect.DeepEqual(got, expectedOff) && !reflect.DeepEqual(got, expectedOn) {
			t.Logf("sequence error, got %+v, expected: %+v or %+v", got, expectedOff, expectedOn)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.Reset()
	gobottest.Assert(t, err, nil)
}

// COMMANDS

func TestSSD1306DriverCommandsDisplay(t *testing.T) {
	s := initTestSSD1306Driver(128, 64, false)
	s.Start()

	result := s.Command("Display")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsOn(t *testing.T) {
	s := initTestSSD1306Driver(128, 64, false)
	s.Start()

	result := s.Command("On")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsOff(t *testing.T) {
	s := initTestSSD1306Driver(128, 64, false)
	s.Start()

	result := s.Command("Off")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsClear(t *testing.T) {
	s := initTestSSD1306Driver(128, 64, false)
	s.Start()

	result := s.Command("Clear")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsSetContrast(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetContrast, 0x80, 0x10}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}

	result := s.Command("SetContrast")(map[string]interface{}{
		"contrast": byte(0x10),
	})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsSet(t *testing.T) {
	s := initTestSSD1306Driver(128, 64, false)
	s.Start()

	gobottest.Assert(t, s.buffer.buffer[0], byte(0))
	s.Command("Set")(map[string]interface{}{
		"x": int(0),
		"y": int(0),
		"c": int(1),
	})
	gobottest.Assert(t, s.buffer.buffer[0], byte(1))
}
