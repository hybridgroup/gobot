package i2c

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*TH02Driver)(nil)

// // --------- HELPERS
func initTestTH02Driver() *SHT3xDriver {
	driver, _ := initTestSHT3xDriverWithStubbedAdaptor()
	return driver
}

func initTestTH02DriverWithStubbedAdaptor() (*TH02Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewTH02Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewTH02Driver(t *testing.T) {
	i2cd := newI2cTestAdaptor()
	defer i2cd.Close()
	// Does it return a pointer to an instance of SHT3xDriver?
	var iface interface{} = NewTH02Driver(i2cd)
	_, ok := iface.(*TH02Driver)
	if !ok {
		t.Errorf("NewTH02Driver() should have returned a *NewTH02Driver")
	}
	b := NewTH02Driver(i2cd, func(Config) {})
	gobottest.Refute(t, b.Connection(), nil)

	//cover some basically useless protions the Interface demands
	if name := b.Name(); name != b.name {
		t.Errorf("Didnt return the proper name.  Got %q wanted %q", name, b.name)
	}

	if b.SetName("42"); b.name != "42" {
		t.Errorf("yikes - didnt set name.")
	}
}

func TestTH02Driver_Accuracy(t *testing.T) {
	i2cd := newI2cTestAdaptor()
	defer i2cd.Close()
	b := NewTH02Driver(i2cd)

	if b.SetAddress(0x42); b.addr != 0x42 {
		t.Error("Didnt set address as expected")
	}

	if b.SetAccuracy(0x42); b.accuracy != TH02HighAccuracy {
		t.Error("Setting an invalid accuracy should resolve to TH02HighAccuracy")
	}

	if b.SetAccuracy(TH02LowAccuracy); b.accuracy != TH02LowAccuracy {
		t.Error("Expected setting low accuracy to actually set to low accuracy")
	}

	if acc := b.Accuracy(); acc != TH02LowAccuracy {
		t.Errorf("Accuract() didnt return what was expected")
	}
}

func TestTH022DriverStart(t *testing.T) {
	b, _ := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, b.Start(), nil)
}

func TestTH02StartConnectError(t *testing.T) {
	d, adaptor := initTestTH02DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestTH02DriverHalt(t *testing.T) {
	sht3x := initTestTH02Driver()
	gobottest.Assert(t, sht3x.Halt(), nil)
}

func TestTH02DriverOptions(t *testing.T) {
	d := NewTH02Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	d.Halt()
}

func TestTH02Driver_ReadData(t *testing.T) {
	d, i2cd := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	type x struct {
		rd, wr func([]byte) (int, error)
		rtn    uint16
		errNil bool
	}

	tests := map[string]x{
		"example RH": x{
			rd: func(b []byte) (int, error) {
				copy(b, []byte{0x00, 0x07, 0xC0})
				return 3, nil
			},
			wr: func([]byte) (int, error) {
				return 1, nil
			},
			errNil: true,
			rtn:    1984,
		},
		"example T": x{
			rd: func(b []byte) (int, error) {
				copy(b, []byte{0x00, 0x12, 0xC0})
				return 3, nil
			},
			wr: func([]byte) (int, error) {
				return 1, nil
			},
			errNil: true,
			rtn:    4800,
		},
		"timeout - no wait for ready": x{
			rd: func(b []byte) (int, error) {
				time.Sleep(200 * time.Millisecond)
				copy(b, []byte{0x01})
				return 1, fmt.Errorf("nope")
			},
			wr: func([]byte) (int, error) {
				return 1, nil
			},
			errNil: false,
			rtn:    0,
		},
		"unable to write status register": x{
			rd: func(b []byte) (int, error) {
				copy(b, []byte{0x00})
				return 0, nil
			},
			wr: func([]byte) (int, error) {
				return 0, fmt.Errorf("Nope")
			},
			errNil: false,
			rtn:    0,
		},
		"unable to read doesnt provide enought data": x{
			rd: func(b []byte) (int, error) {
				copy(b, []byte{0x00, 0x01})
				return 2, nil
			},
			wr: func([]byte) (int, error) {
				return 1, nil
			},
			errNil: false,
			rtn:    0,
		},
	}

	for name, x := range tests {
		t.Log("Running", name)
		i2cd.i2cReadImpl = x.rd
		i2cd.i2cWriteImpl = x.wr
		got, err := d.readData()
		gobottest.Assert(t, err == nil, x.errNil)
		gobottest.Assert(t, got, x.rtn)
	}
}

func TestTH02Driver_waitForReady(t *testing.T) {
	d, i2cd := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	i2cd.i2cReadImpl = func(b []byte) (int, error) {
		time.Sleep(50 * time.Millisecond)
		copy(b, []byte{0x01, 0x00})
		return 3, nil
	}

	i2cd.i2cWriteImpl = func([]byte) (int, error) {
		return 1, nil
	}

	timeout := 10 * time.Microsecond
	if err := d.waitForReady(&timeout); err == nil {
		t.Error("Expected a timeout error")
	}
}

func TestTH02Driver_WriteRegister(t *testing.T) {
	d, i2cd := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	i2cd.i2cWriteImpl = func([]byte) (int, error) {
		return 1, nil
	}

	if err := d.writeRegister(0x00, 0x00); err != nil {
		t.Errorf("expected a nil error write")
	}
}

func TestTH02Driver_Heater(t *testing.T) {
	d, i2cd := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	i2cd.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xff})
		return 1, nil
	}

	i2cd.i2cWriteImpl = func([]byte) (int, error) {
		return 1, nil
	}

	on, err := d.Heater()
	gobottest.Assert(t, on, true)
	gobottest.Assert(t, err, nil)
}
func TestTH02Driver_SerialNumber(t *testing.T) {
	d, i2cd := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	i2cd.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x42})
		return 1, nil
	}

	i2cd.i2cWriteImpl = func([]byte) (int, error) {
		return 1, nil
	}

	sn, err := d.SerialNumber()

	gobottest.Assert(t, sn, uint32((0x42)>>4))
	gobottest.Assert(t, err, nil)
}

func TestTH02Driver_ApplySettings(t *testing.T) {
	d := &TH02Driver{}

	type x struct {
		acc, base, out byte
		heating        bool
	}

	tests := map[string]x{
		"low acc, heating":     x{acc: TH02LowAccuracy, base: 0x00, heating: true, out: 0x01},
		"high acc, no heating": x{acc: TH02HighAccuracy, base: 0x00, heating: false, out: 0x23},
	}

	for name, x := range tests {
		t.Log(name)
		d.accuracy = x.acc
		d.heating = x.heating
		got := d.applysettings(x.base)
		gobottest.Assert(t, x.out, got)
	}
}

func TestTH02Driver_Sample(t *testing.T) {
	d, i2cd := initTestTH02DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)

	i2cd.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x07, 0xC0})
		return 4, nil
	}

	i2cd.i2cWriteImpl = func([]byte) (int, error) {
		return 1, nil
	}

	temp, rh, _ := d.Sample()

	gobottest.Assert(t, temp, float32(0))
	gobottest.Assert(t, rh, float32(0))

}
