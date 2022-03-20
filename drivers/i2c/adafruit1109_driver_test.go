package i2c

import (
	"errors"
	"fmt"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func initTestAdafruit1109DriverWithStubbedAdaptor() (*Adafruit1109Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewAdafruit1109Driver(adaptor), adaptor
}

func TestAdafruit1109DriverStart(t *testing.T) {
	ada, _ := initTestAdafruit1109DriverWithStubbedAdaptor()
	gobottest.Assert(t, ada.Start(), nil)
}

func TestAdafruit1109DriverStartWriteErr(t *testing.T) {
	d, adaptor := initTestAdafruit1109DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestAdafruit1109DriverStartReadErr(t *testing.T) {
	d, adaptor := initTestAdafruit1109DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, d.Start(), errors.New("read error"))
}

func TestAdafruit1109Driver_parseId(t *testing.T) {
	// arrange
	ports := []string{"A", "B"}
	for _, port := range ports {
		for pin := uint8(0); pin <= 7; pin++ {
			id := fmt.Sprintf("%s_%d", port, pin)
			t.Run(id, func(t *testing.T) {
				// act
				got := adafruit1109ParseId(id)
				// assert
				gobottest.Assert(t, got, adafruit1109PortPin{port, pin})
			})
		}
	}
}
