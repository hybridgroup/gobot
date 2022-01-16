package i2c

import (
	"bytes"
	"fmt"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AS5600Driver)(nil)

func initTestAS5600Driver() *AS5600Driver {

	driver, _ := initTestAS5600DriverWithStubbedAdaptor()

	return driver
}

func initTestAS5600DriverWithStubbedAdaptor() (*AS5600Driver, *i2cTestAdaptor) {

	adaptor := newI2cTestAdaptor()

	return NewAS5600Driver(adaptor), adaptor
}

func TestAS5600Driver(t *testing.T) {

	as := initTestAS5600Driver()

	gobottest.Refute(t, as.Connection(), nil)
}

func TestAS5600DriverStart(t *testing.T) {
	var as *AS5600Driver

	as, _ = initTestAS5600DriverWithStubbedAdaptor()
	gobottest.Assert(t, as.Start(), nil)
}

func TestAS5600DriverHalt(t *testing.T) {
	as := initTestAS5600Driver()

	gobottest.Assert(t, as.Halt(), nil)
}

func TestAS5600DriverSetName(t *testing.T) {

	// Does it change the name of the driver
	as := initTestAS5600Driver()
	as.SetName("TESTME")
	gobottest.Assert(t, as.Name(), "TESTME")
}

func TestAS5600DriverOptions(t *testing.T) {

	as := NewAS5600Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, as.GetBusOrDefault(1), 2)
}

func TestAS5600DriverDetecMagnet(t *testing.T) {

	as, adaptor := initTestAS5600DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{(0x0 | as5600StatusMDBit)})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	as.Start()
	magnet, err := as.DetecMagnet()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, magnet, true)
}

func TestAS5600DriverMagnetTooWeak(t *testing.T) {

	as, adaptor := initTestAS5600DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{(0x0 | as5600StatusMLBit)})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	as.Start()
	magnet, err := as.GetMagnetStrength()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, magnet, as5600MagnetTooWeak)
}

func TestAS5600DriverMagnetTooStrong(t *testing.T) {

	as, adaptor := initTestAS5600DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{(0x0 | as5600StatusMHBit)})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	as.Start()
	magnet, err := as.GetMagnetStrength()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, magnet, as5600MagnetTooStrong)
}

func TestAS5600DriverGetRawAngle(t *testing.T) {
	as, adaptor := initTestAS5600DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		tmp := make([]byte, len(b))
		if adaptor.written[len(adaptor.written)-1] == as5600RAWANGLEMSB {
			tmp[0] = as5600RAWANGLEMSB
		} else if adaptor.written[len(adaptor.written)-1] == as5600RAWANGLELSB {
			tmp[0] = as5600RAWANGLELSB
		}
		for i := 1; i < len(b); i++ {
			tmp[i] = byte(i)
		}
		buf.Write(tmp)
		copy(b, buf.Bytes())

		return buf.Len(), nil
	}
	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		if b[0] == as5600RAWANGLEMSB {
			return 1, nil
		}

		return 0, fmt.Errorf("0x%x wrong register", b[0])
	}
	as.Start()
	angle, err := as.GetRawAngle()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, angle, uint16(0x010c))
}

func TestAS5600DriverGetAngle(t *testing.T) {
	as, adaptor := initTestAS5600DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		tmp := make([]byte, len(b))
		if adaptor.written[len(adaptor.written)-1] == as5600ANGLEMSB {
			tmp[0] = as5600ANGLEMSB
		} else if adaptor.written[len(adaptor.written)-1] == as5600ANGLELSB {
			tmp[0] = as5600ANGLELSB
		}
		for i := 1; i < len(b); i++ {
			tmp[i] = byte(i)
		}
		buf.Write(tmp)
		copy(b, buf.Bytes())

		return buf.Len(), nil
	}
	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		if b[0] == as5600ANGLEMSB {
			return 1, nil
		}

		return 0, fmt.Errorf("0x%x wrong register", b[0])
	}
	as.Start()
	angle, err := as.GetScaledAngle()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, angle, uint16(0x010e))
}

func TestAS5600DriverGetMaxAngle(t *testing.T) {
	as, adaptor := initTestAS5600DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		tmp := make([]byte, len(b))
		if adaptor.written[len(adaptor.written)-1] == as5600MANGMSB {
			tmp[0] = as5600MANGMSB
		} else if adaptor.written[len(adaptor.written)-1] == as5600MANGLSB {
			tmp[0] = as5600MANGLSB
		}
		for i := 1; i < len(b); i++ {
			tmp[i] = byte(i)
		}
		buf.Write(tmp)
		copy(b, buf.Bytes())

		return buf.Len(), nil
	}
	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		if b[0] == as5600MANGMSB {
			return 1, nil
		}

		return 0, fmt.Errorf("0x%x wrong register", b[0])
	}
	as.Start()
	angle, err := as.GetMaxAngle()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, angle, uint16(0x0105))
}
