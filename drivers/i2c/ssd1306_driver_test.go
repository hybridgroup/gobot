//nolint:forcetypeassert // ok here
package i2c

import (
	"fmt"
	"image"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*SSD1306Driver)(nil)

func initTestSSD1306DriverWithStubbedAdaptor(
	width, height int,
	externalVCC bool, //nolint:unparam // keep for tests
) (*SSD1306Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewSSD1306Driver(a, WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height),
		WithSSD1306ExternalVCC(externalVCC))
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewSSD1306Driver(t *testing.T) {
	var di interface{} = NewSSD1306Driver(newI2cTestAdaptor())
	d, ok := di.(*SSD1306Driver)
	if !ok {
		require.Fail(t, "new should have returned a *SSD1306Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "SSD1306"))
	assert.Equal(t, 0x3c, d.defaultAddress)
}

func TestSSD1306StartDefault(t *testing.T) {
	const (
		width       = 128
		height      = 64
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	require.NoError(t, d.Start())
}

func TestSSD1306Start128x32(t *testing.T) {
	const (
		width       = 128
		height      = 32
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	require.NoError(t, d.Start())
}

func TestSSD1306Start96x16(t *testing.T) {
	const (
		width       = 96
		height      = 16
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	require.NoError(t, d.Start())
}

func TestSSD1306StartExternalVCC(t *testing.T) {
	const (
		width       = 128
		height      = 32
		externalVCC = true
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	require.NoError(t, d.Start())
}

func TestSSD1306StartSizeError(t *testing.T) {
	const (
		width       = 128
		height      = 54
		externalVCC = false
	)
	d := NewSSD1306Driver(newI2cTestAdaptor(),
		WithSSD1306DisplayWidth(width), WithSSD1306DisplayHeight(height), WithSSD1306ExternalVCC(externalVCC))
	require.ErrorContains(t, d.Start(), "128x54 resolution is unsupported, supported resolutions: 128x64, 128x32, 96x16")
}

func TestSSD1306Halt(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	require.NoError(t, s.Halt())
}

func TestSSD1306Options(t *testing.T) {
	s := NewSSD1306Driver(newI2cTestAdaptor(), WithBus(2), WithSSD1306DisplayHeight(32), WithSSD1306DisplayWidth(128))
	assert.Equal(t, 2, s.GetBusOrDefault(1))
	assert.Equal(t, 32, s.displayHeight)
	assert.Equal(t, 128, s.displayWidth)
}

func TestSSD1306Display(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(96, 16, false)
	_ = s.Start()
	require.NoError(t, s.Display())
}

func TestSSD1306ShowImage(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	_ = s.Start()
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	require.ErrorContains(t, s.ShowImage(img), "image must match display width and height: 128x64")

	img = image.NewRGBA(image.Rect(0, 0, 128, 64))
	require.NoError(t, s.ShowImage(img))
}

func TestSSD1306Command(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	_ = s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, 0xFF}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.command(0xFF)
	require.NoError(t, err)
}

func TestSSD1306Commands(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	_ = s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, 0x00, 0x80, 0xFF}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.commands([]byte{0x00, 0xFF})
	require.NoError(t, err)
}

func TestSSD1306On(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	_ = s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetDisplayOn}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.On()
	require.NoError(t, err)
}

func TestSSD1306Off(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	_ = s.Start()

	a.i2cWriteImpl = func(got []byte) (int, error) {
		expected := []byte{0x80, ssd1306SetDisplayOff}
		if !reflect.DeepEqual(got, expected) {
			t.Logf("sequence error, got %+v, expected %+v", got, expected)
			return 0, fmt.Errorf("oops")
		}
		return 0, nil
	}
	err := s.Off()
	require.NoError(t, err)
}

func TestSSD1306Reset(t *testing.T) {
	s, a := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	_ = s.Start()

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
	require.NoError(t, err)
}

// COMMANDS

func TestSSD1306CommandsDisplay(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)
	result := s.Command("Display")(map[string]interface{}{})
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestSSD1306CommandsOn(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	result := s.Command("On")(map[string]interface{}{})
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestSSD1306CommandsOff(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	result := s.Command("Off")(map[string]interface{}{})
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestSSD1306CommandsClear(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	result := s.Command("Clear")(map[string]interface{}{})
	assert.Nil(t, result.(map[string]interface{})["err"])
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
	assert.Nil(t, result.(map[string]interface{})["err"])
}

func TestSSD1306CommandsSet(t *testing.T) {
	s, _ := initTestSSD1306DriverWithStubbedAdaptor(128, 64, false)

	assert.Equal(t, byte(0), s.buffer.buffer[0])
	s.Command("Set")(map[string]interface{}{
		"x": int(0),
		"y": int(0),
		"c": int(1),
	})
	assert.Equal(t, byte(1), s.buffer.buffer[0])
}

func TestDisplayBuffer(t *testing.T) {
	width := 128
	height := 64
	size := 1024 // (width*height) / 8
	display := NewDisplayBuffer(width, height, 8)

	if display.Size() != size {
		require.Fail(t, "invalid Size() (%d, expected %d)",
			display.Size(), size)
	}
	if len(display.buffer) != size {
		require.Fail(t, "allocated buffer size invalid (%d, expected %d)",
			len(display.buffer), size)
	}

	assert.Equal(t, byte(0), display.buffer[0])
	assert.Equal(t, byte(0), display.buffer[1])

	display.SetPixel(0, 0, 1)
	display.SetPixel(1, 0, 1)
	display.SetPixel(2, 0, 1)
	display.SetPixel(0, 1, 1)
	assert.Equal(t, byte(3), display.buffer[0])
	assert.Equal(t, byte(1), display.buffer[1])

	display.SetPixel(0, 1, 0)
	assert.Equal(t, byte(1), display.buffer[0])
	assert.Equal(t, byte(1), display.buffer[1])
}
