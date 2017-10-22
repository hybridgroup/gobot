package i2c

import (
	"errors"
	"fmt"
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

	display := NewDisplayBuffer(width, height)

	if display.Size() != size {
		t.Errorf("invalid Size() (%d, expected %d)",
			display.Size(), size)
	}
	if len(display.buffer) != size {
		t.Errorf("Allocated buffer size invalid (%d, expected %d)",
			len(display.buffer), size)
	}

	gobottest.Assert(t, display.buffer[0], byte(0))
	gobottest.Assert(t, display.buffer[1], byte(0))

	display.Set(0, 0, 1)
	display.Set(1, 0, 1)
	display.Set(2, 0, 1)
	display.Set(0, 1, 1)
	gobottest.Assert(t, display.buffer[0], byte(3))
	gobottest.Assert(t, display.buffer[1], byte(1))

	display.Set(0, 1, 0)
	gobottest.Assert(t, display.buffer[0], byte(1))
	gobottest.Assert(t, display.buffer[1], byte(1))
}

// --------- HELPERS
func initTestSSD1306Driver() (driver *SSD1306Driver) {
	driver, _ = initTestSSD1306DriverWithStubbedAdaptor()
	return
}

func initTestSSD1306DriverWithStubbedAdaptor() (*SSD1306Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewSSD1306Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewSSD1306Driver(t *testing.T) {
	// Does it return a pointer to an instance of SSD1306Driver?
	var bm interface{} = NewSSD1306Driver(newI2cTestAdaptor())
	_, ok := bm.(*SSD1306Driver)
	if !ok {
		t.Errorf("NewSSD1306Driver() should have returned a *SSD1306Driver")
	}

	b := NewSSD1306Driver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

// Methods

func TestSSD1306DriverStart(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor()

	gobottest.Assert(t, s.Start(), nil)
}

func TestSSD1306StartConnectError(t *testing.T) {
	d, adaptor := initTestSSD1306DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestSSD1306DriverHalt(t *testing.T) {
	s := initTestSSD1306Driver()

	gobottest.Assert(t, s.Halt(), nil)
}

// Test Name & SetName
func TestSSD1306DriverName(t *testing.T) {
	s := initTestSSD1306Driver()

	gobottest.Assert(t, strings.HasPrefix(s.Name(), "SSD1306"), true)
	s.SetName("Ole Oled")
	gobottest.Assert(t, s.Name(), "Ole Oled")
}

func TestSSD1306DriverOptions(t *testing.T) {
	s := NewSSD1306Driver(newI2cTestAdaptor(), WithBus(2), WithDisplayHeight(32), WithDisplayWidth(128))
	gobottest.Assert(t, s.GetBusOrDefault(1), 2)
	gobottest.Assert(t, s.DisplayHeight, 32)
	gobottest.Assert(t, s.DisplayWidth, 128)
}

func TestSSD1306DriverDisplay(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor()
	s.Start()
	gobottest.Assert(t, s.Display(), nil)
}

func TestSSD1306DriverCommand(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor()
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, 0xFF}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("Sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("Woops!")
		}
		return 0, nil
	}
	err := s.command(0xFF)
	gobottest.Assert(t, err, nil)
}

func TestSSD1306DriverCommands(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor()
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, 0x00, 0x80, 0xFF}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("Sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("Woops!")
		}
		return 0, nil
	}
	err := s.commands([]byte{0x00, 0xFF})
	gobottest.Assert(t, err, nil)
}

func TestSSD1306DriverOn(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor()
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetDisplayOn}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("Sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("Woops!")
		}
		return 0, nil
	}
	err := s.On()
	gobottest.Assert(t, err, nil)
}

func TestSSD1306DriverOff(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor()
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetDisplayOff}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("Sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("Woops!")
		}
		return 0, nil
	}
	err := s.Off()
	gobottest.Assert(t, err, nil)
}

// COMMANDS

func TestSSD1306DriverCommandsDisplay(t *testing.T) {
	s := initTestSSD1306Driver()
	s.Start()

	result := s.Command("Display")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsOn(t *testing.T) {
	s := initTestSSD1306Driver()
	s.Start()

	result := s.Command("On")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsOff(t *testing.T) {
	s := initTestSSD1306Driver()
	s.Start()

	result := s.Command("Off")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsClear(t *testing.T) {
	s := initTestSSD1306Driver()
	s.Start()

	result := s.Command("Clear")(map[string]interface{}{})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsSetContrast(t *testing.T) {
	s, adaptor := initTestSSD1306DriverWithStubbedAdaptor()
	s.Start()

	adaptor.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetContrast, 0x80, 0x10}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("Sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("Woops!")
		}
		return 0, nil
	}

	result := s.Command("SetContrast")(map[string]interface{}{
		"contrast": byte(0x10),
	})
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
}

func TestSSD1306DriverCommandsSet(t *testing.T) {
	s := initTestSSD1306Driver()
	s.Start()

	gobottest.Assert(t, s.Buffer.buffer[0], byte(0))
	s.Command("Set")(map[string]interface{}{
		"x": int(0),
		"y": int(0),
		"c": int(1),
	})
	gobottest.Assert(t, s.Buffer.buffer[0], byte(1))
}
