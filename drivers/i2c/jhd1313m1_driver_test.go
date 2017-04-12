package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*JHD1313M1Driver)(nil)

// --------- HELPERS
func initTestJHD1313M1Driver() (driver *JHD1313M1Driver) {
	driver, _ = initTestJHD1313M1DriverWithStubbedAdaptor()
	return
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
		t.Errorf("NewJHD1313M1Driver() should have returned a *JHD1313M1Driver")
	}
}

// Methods
func TestJHD1313M1Driver(t *testing.T) {
	jhd := initTestJHD1313M1Driver()

	gobottest.Refute(t, jhd.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(jhd.Name(), "JHD1313M1"), true)
}

func TestJHD1313MDriverSetName(t *testing.T) {
	d := initTestJHD1313M1Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestJHD1313MDriverOptions(t *testing.T) {
	d := NewJHD1313M1Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestJHD1313MDriverStart(t *testing.T) {
	d := initTestJHD1313M1Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestJHD1313MStartConnectError(t *testing.T) {
	d, adaptor := initTestJHD1313M1DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestJHD1313MDriverStartWriteError(t *testing.T) {
	d, adaptor := initTestJHD1313M1DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestJHD1313MDriverHalt(t *testing.T) {
	d := initTestJHD1313M1Driver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestJHD1313MDriverSetRgb(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.SetRGB(0x00, 0x00, 0x00), nil)
}

func TestJHD1313MDriverSetRgbError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.SetRGB(0x00, 0x00, 0x00), errors.New("write error"))
}

func TestJHD1313MDriverClear(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Clear(), nil)
}

func TestJHD1313MDriverClearError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Clear(), errors.New("write error"))
}

func TestJHD1313MDriverHome(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Home(), nil)
}

func TestJHD1313MDriverWrite(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Write("Hello"), nil)
}

func TestJHD1313MDriverWriteError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	gobottest.Assert(t, d.Write("Hello"), errors.New("write error"))
}

func TestJHD1313MDriverWriteTwoLines(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Write("Hello\nthere"), nil)
}

func TestJHD1313MDriverWriteTwoLinesError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Write("Hello\nthere"), errors.New("write error"))
}

func TestJHD1313MDriverSetPosition(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.SetPosition(2), nil)
}

func TestJHD1313MDriverSetSecondLinePosition(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.SetPosition(18), nil)
}

func TestJHD1313MDriverSetPositionInvalid(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.SetPosition(-1), ErrInvalidPosition)
	gobottest.Assert(t, d.SetPosition(32), ErrInvalidPosition)
}

func TestJHD1313MDriverScroll(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Scroll(true), nil)
}

func TestJHD1313MDriverReverseScroll(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Scroll(false), nil)
}

func TestJHD1313MDriverSetCustomChar(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	d.Start()
	gobottest.Assert(t, d.SetCustomChar(0, data), nil)
}

func TestJHD1313MDriverSetCustomCharError(t *testing.T) {
	d, _ := initTestJHD1313M1DriverWithStubbedAdaptor()
	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	d.Start()
	gobottest.Assert(t, d.SetCustomChar(10, data), errors.New("can't set a custom character at a position greater than 7"))
}

func TestJHD1313MDriverSetCustomCharWriteError(t *testing.T) {
	d, a := initTestJHD1313M1DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	data := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	gobottest.Assert(t, d.SetCustomChar(0, data), errors.New("write error"))
}
