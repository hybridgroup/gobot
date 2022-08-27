package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*ADS1x15Driver)(nil)

// that supports the AnalogReader interface
var _ aio.AnalogReader = (*ADS1x15Driver)(nil)

func initTestADS1x15DriverWithStubbedAdaptor() (*ADS1x15Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	const defaultDataRate = 3
	dataRates := map[int]uint16{defaultDataRate: 0x0003}
	d := newADS1x15Driver(a, "ADS1x15", dataRates, defaultDataRate)
	noConversion := []uint8{0x80, 0x00} // no conversion in progress
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, noConversion)
		return 2, nil
	}
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

var ads1x15TestChannel = map[string]interface{}{
	"channel": int(2),
}

var ads1x15TestChannelGainDataRate = map[string]interface{}{
	"channel":  int(1),
	"gain":     int(2),
	"dataRate": int(3),
}

func TestADS1x15CommandsReadDifferenceWithDefaults(t *testing.T) {
	// arrange
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()
	// act
	result := d.Command("ReadDifferenceWithDefaults")(ads1x15TestChannel)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
	gobottest.Assert(t, result.(map[string]interface{})["val"], -4.096)
}

func TestADS1x15CommandsReadDifference(t *testing.T) {
	// arrange
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()
	// act
	result := d.Command("ReadDifference")(ads1x15TestChannelGainDataRate)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
	gobottest.Assert(t, result.(map[string]interface{})["val"], -2.048)
}

func TestADS1x15CommandsReadWithDefaults(t *testing.T) {
	// arrange
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()
	// act
	result := d.Command("ReadWithDefaults")(ads1x15TestChannel)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
	gobottest.Assert(t, result.(map[string]interface{})["val"], -4.096)
}

func TestADS1x15CommandsRead(t *testing.T) {
	// arrange
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()
	// act
	result := d.Command("Read")(ads1x15TestChannelGainDataRate)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
	gobottest.Assert(t, result.(map[string]interface{})["val"], -2.048)
}

func TestADS1x15CommandsAnalogRead(t *testing.T) {
	// arrange
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()
	ads1x15TestPin := map[string]interface{}{
		"pin": string("2"),
	}
	// act
	result := d.Command("AnalogRead")(ads1x15TestPin)
	// assert
	gobottest.Assert(t, result.(map[string]interface{})["err"], nil)
	gobottest.Assert(t, result.(map[string]interface{})["val"], -32768)
}

func TestADS1x15_ads1x15BestGainForVoltage(t *testing.T) {
	g, err := ads1x15BestGainForVoltage(1.5)
	gobottest.Assert(t, g, 2)

	g, err = ads1x15BestGainForVoltage(20.0)
	gobottest.Assert(t, err, errors.New("The maximum voltage which can be read is 6.144000"))
}
