package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*JHD1313M1Driver)(nil)

// --------- HELPERS
func initTestJHD1313M1Driver() *JHD1313M1Driver {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	return d
}

func initTestJHD1313M1DriverWithStubbedAdaptor() (*JHD1313M1Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewJHD1313M1Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewJHD1313M1Driver(t *testing.T) {
	// Does it return a pointer to an instance of JHD1313M1Driver?
	var mpl interface{} = NewJHD1313M1Driver(newI2cTestAdaptor())
	_, ok := mpl.(*JHD1313M1Driver)
	if !ok {
		require.Fail(t, "NewJHD1313M1Driver() should have returned a *JHD1313M1Driver")
	}
}

// Methods
func TestJHD1313M1Driver(t *testing.T) {
	jhd := initTestJHD1313M1Driver()

	assert.NotNil(t, jhd.Connection())
	assert.True(t, strings.HasPrefix(jhd.Name(), "JHD1313M1"))
}

func TestJHD1313MDriverSetName(t *testing.T) {
	d := initTestJHD1313M1Driver()
	d.SetName("TESTME")
	assert.Equal(t, "TESTME", d.Name())
}

func TestJHD1313MDriverOptions(t *testing.T) {
	d := NewJHD1313M1Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestJHD1313MDriverStart(t *testing.T) {
	d := initTestJHD1313M1Driver()
	require.NoError(t, d.Start())
}

func TestJHD1313MStartConnectError(t *testing.T) {
	d, adaptor := initTestJHD1313M1DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	require.ErrorContains(t, d.Start(), "Invalid i2c connection")
}

func TestJHD1313MDriverStartWriteError(t *testing.T) {
	d, adaptor := initTestJHD1313M1DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.Start(), "write error")
}

func TestJHD1313MDriverHalt(t *testing.T) {
	d := initTestJHD1313M1Driver()
	_ = d.Start()
	require.NoError(t, d.Halt())
}

func TestJHD1313MDriverSetRgb(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.SetRGB(0x00, 0x00, 0x00))
}

func TestJHD1313MDriverSetRgbError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.SetRGB(0x00, 0x00, 0x00), "write error")
}

func TestJHD1313MDriverClear(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Clear())
}

func TestJHD1313MDriverClearError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.Clear(), "write error")
}

func TestJHD1313MDriverHome(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Home())
}

func TestJHD1313MDriverWrite(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Write("Hello"))
}

func TestJHD1313MDriverWriteError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	require.ErrorContains(t, d.Write("Hello"), "write error")
}

func TestJHD1313MDriverWriteTwoLines(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Write("Hello\nthere"))
}

func TestJHD1313MDriverWriteTwoLinesError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.Write("Hello\nthere"), "write error")
}

func TestJHD1313MDriverSetPosition(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.SetPosition(2))
}

func TestJHD1313MDriverSetSecondLinePosition(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.SetPosition(18))
}

func TestJHD1313MDriverSetPositionInvalid(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	assert.Equal(t, jhd1313m1ErrInvalidPosition, d.SetPosition(-1))
	assert.Equal(t, jhd1313m1ErrInvalidPosition, d.SetPosition(32))
}

func TestJHD1313MDriverScroll(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Scroll(true))
}

func TestJHD1313MDriverReverseScroll(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Scroll(false))
}

func TestJHD1313MDriverSetCustomChar(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_ = d.Start()
	require.NoError(t, d.SetCustomChar(0, data))
}

func TestJHD1313MDriverSetCustomCharError(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_ = d.Start()
	require.ErrorContains(t, d.SetCustomChar(10, data), "can't set a custom character at a position greater than 7")
}

func TestJHD1313MDriverSetCustomCharWriteError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	require.ErrorContains(t, d.SetCustomChar(0, data), "write error")
}

func TestJHD1313MDriverCommands(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	_ = d.Start()

	err := d.Command("SetRGB")(map[string]interface{}{"r": "1", "g": "1", "b": "1"})
	assert.Nil(t, err)

	err = d.Command("Clear")(map[string]interface{}{})
	assert.Nil(t, err)

	err = d.Command("Home")(map[string]interface{}{})
	assert.Nil(t, err)

	err = d.Command("Write")(map[string]interface{}{"msg": "Hello"})
	assert.Nil(t, err)

	err = d.Command("SetPosition")(map[string]interface{}{"pos": "1"})
	assert.Nil(t, err)

	err = d.Command("Scroll")(map[string]interface{}{"lr": "true"})
	assert.Nil(t, err)
}
