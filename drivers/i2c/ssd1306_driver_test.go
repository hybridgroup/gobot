package i2c

import (
	"errors"
	"fmt"
	"image"
	"reflect"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*SSD1306Driver)(nil)

func initTestSSD1306DriverWithStubbedAdaptor(width, height int, externalVCC bool) (*SSD1306Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewSSD1306Driver(a, WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewSSD1306Driver(t *testing.T) {
	var di interface{} = NewSSD1306Driver(newI2cTestAdaptor())
	d, ok := di.(*SSD1306Driver)
	if !ok {
		t.Errorf("new should have returned a *SSD1306Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "SSD1306"), true)
	gobottest.Assert(t, d.defaultAddress, 0x3c)
}

func TestSSD1306StartDefault(t *testing.T) {
	const (
		width       = 128
		height      = 64
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	gobottest.Assert(t, d.Start(), nil)
}

func TestSSD1306Start128x32(t *testing.T) {
	const (
		width       = 128
		height      = 32
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	gobottest.Assert(t, d.Start(), nil)
}

func TestSSD1306Start96x16(t *testing.T) {
	const (
		width       = 96
		height      = 16
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	gobottest.Assert(t, d.Start(), nil)
}

func TestSSD1306StartExternalVCC(t *testing.T) {
	const (
		width       = 128
		height      = 32
		externalVCC = true
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	gobottest.Assert(t, d.Start(), nil)
}

func TestSSD1306StartSizeError(t *testing.T) {
	const (
		width       = 128
		height      = 54
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	gobottest.Assert(t, d.Start(), errors.New("128x54 resolution is unsupported, supported resolutions: 128x64, 128x32, 96x16"))
}

func TestSSD1306Halt(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	gobottest.Assert(t, s.Halt(), nil)
}

func TestSSD1306Options(t *testing.T) {
	s := NewSSD1306Driver(newI2cTestAdaptor(), WithBus(2), WithSSD1306DisplayHeight(32), WithSSD1306DisplayWidth(128))
	gobottest.Assert(t, s.GetBusOrDefault(1), 2)
	gobottest.Assert(t, s.displayHeight, 32)
	gobottest.Assert(t, s.displayWidth, 128)
}

func TestSSD1306Display(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(96, 16, false)
	s.Start()
	gobottest.Assert(t, s.Display(), nil)
}

func TestSSD1306ShowImage(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	gobottest.Assert(t, s.ShowImage(img), errors.New("image must match display width and height: 128x64"))

	img = image.NewRGBA(image.Rect(0, 0, 128, 64))
	gobottest.Assert(t, s.ShowImage(img), nil)
}

func TestSSD1306Command(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
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

func TestSSD1306Commands(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
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

func TestSSD1306On(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
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

func TestSSD1306Off(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
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
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
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

func TestSSD1306CommandsDisplay(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	result := s.Command("Display")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306CommandsOn(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	result := s.Command("On")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306CommandsOff(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	result := s.Command("Off")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306CommandsClear(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	result := s.Command("Clear")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306CommandsSetContrast(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	a.i2cWriteImpl = func(got []byte) (int, error) {
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

func TestSSD1306CommandsSet(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	gobottest.Assert(t, s.buffer.buffer[0], byte(0))
	s.Command("Set")(map[string]interface{}{
		"x": int(0),
		"y": int(0),
		"c": int(1),
	})
	gobottest.Assert(t, s.buffer.buffer[0], byte(1))
}

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
