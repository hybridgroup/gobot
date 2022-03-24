package i2c

import (
	"fmt"
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
)

func initTestYL40DriverWithStubbedAdaptor() (*YL40Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	yl := NewYL40Driver(adaptor, WithPCF8591With400kbitStabilization(0, 2))
	WithPCF8591ForceRefresh(1)(yl.PCF8591Driver)
	yl.Start()
	return yl, adaptor
}

func TestYL40Driver(t *testing.T) {
	// arrange, act
	yl := NewYL40Driver(newI2cTestAdaptor())
	//assert
	gobottest.Refute(t, yl.PCF8591Driver, nil)
	gobottest.Assert(t, yl.conf.sensors[YL40Bri].interval, time.Duration(0))
	gobottest.Refute(t, yl.conf.sensors[YL40Bri].scaler, nil)
	gobottest.Assert(t, yl.conf.sensors[YL40Temp].interval, time.Duration(0))
	gobottest.Refute(t, yl.conf.sensors[YL40Temp].scaler, nil)
	gobottest.Assert(t, yl.conf.sensors[YL40AIN2].interval, time.Duration(0))
	gobottest.Refute(t, yl.conf.sensors[YL40AIN2].scaler, nil)
	gobottest.Assert(t, yl.conf.sensors[YL40Poti].interval, time.Duration(0))
	gobottest.Refute(t, yl.conf.sensors[YL40Poti].scaler, nil)
	gobottest.Refute(t, yl.conf.aOutScaler, nil)
	gobottest.Refute(t, yl.aBri, nil)
	gobottest.Refute(t, yl.aTemp, nil)
	gobottest.Refute(t, yl.aAIN2, nil)
	gobottest.Refute(t, yl.aPoti, nil)
	gobottest.Refute(t, yl.aOut, nil)
}

func TestYL40DriverWithYL40Interval(t *testing.T) {
	// arrange, act
	yl := NewYL40Driver(newI2cTestAdaptor(),
		WithYL40Interval(YL40Bri, 100),
		WithYL40Interval(YL40Temp, 101),
		WithYL40Interval(YL40AIN2, 102),
		WithYL40Interval(YL40Poti, 103),
	)
	// assert
	gobottest.Assert(t, yl.conf.sensors[YL40Bri].interval, time.Duration(100))
	gobottest.Assert(t, yl.conf.sensors[YL40Temp].interval, time.Duration(101))
	gobottest.Assert(t, yl.conf.sensors[YL40AIN2].interval, time.Duration(102))
	gobottest.Assert(t, yl.conf.sensors[YL40Poti].interval, time.Duration(103))
}

func TestYL40DriverWithYL40InputScaler(t *testing.T) {
	// arrange
	yl := NewYL40Driver(newI2cTestAdaptor())
	f1 := func(input int) (value float64) { return 0.1 }
	f2 := func(input int) (value float64) { return 0.2 }
	f3 := func(input int) (value float64) { return 0.3 }
	f4 := func(input int) (value float64) { return 0.4 }
	//act
	WithYL40InputScaler(YL40Bri, f1)(yl)
	WithYL40InputScaler(YL40Temp, f2)(yl)
	WithYL40InputScaler(YL40AIN2, f3)(yl)
	WithYL40InputScaler(YL40Poti, f4)(yl)
	// assert
	gobottest.Assert(t, fEqual(yl.conf.sensors[YL40Bri].scaler, f1), true)
	gobottest.Assert(t, fEqual(yl.conf.sensors[YL40Temp].scaler, f2), true)
	gobottest.Assert(t, fEqual(yl.conf.sensors[YL40AIN2].scaler, f3), true)
	gobottest.Assert(t, fEqual(yl.conf.sensors[YL40Poti].scaler, f4), true)
}

func TestYL40DriverWithYL40WithYL40OutputScaler(t *testing.T) {
	// arrange
	yl := NewYL40Driver(newI2cTestAdaptor())
	fo := func(input float64) (value int) { return 123 }
	//act
	WithYL40OutputScaler(fo)(yl)
	// assert
	gobottest.Assert(t, fEqual(yl.conf.aOutScaler, fo), true)
}

func TestYL40DriverReadBrightness(t *testing.T) {
	// sequence to read the input with PCF8591, see there tests
	// arrange
	yl, adaptor := initTestYL40DriverWithStubbedAdaptor()
	adaptor.written = []byte{} // reset writes of Start() and former test
	// ANAOUT was switched on by Start()
	ctrlByteOn := uint8(pcf8591_ANAON) | uint8(pcf8591_ALLSINGLE) | uint8(pcf8591_CHAN0)
	returnRead := []uint8{0x01, 0x02, 0x03, 73}
	// scaler for brightness is 255..0 => 0..1000
	want := float64(255-returnRead[3]) * 1000 / 255
	// arrange reads
	numCallsRead := 0
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b = returnRead[0:len(b)]
		}
		if numCallsRead == 2 {
			b[0] = returnRead[len(returnRead)-1]
		}
		return len(b), nil
	}
	// act
	got, err := yl.ReadBrightness()
	got2, err2 := yl.Brightness()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(adaptor.written), 1)
	gobottest.Assert(t, adaptor.written[0], ctrlByteOn)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, got, want)
	gobottest.Assert(t, err2, nil)
	gobottest.Assert(t, got2, want)
}

func TestYL40DriverReadTemperature(t *testing.T) {
	// sequence to read the input with PCF8591, see there tests
	// arrange
	yl, adaptor := initTestYL40DriverWithStubbedAdaptor()
	adaptor.written = []byte{} // reset writes of Start() and former test
	// ANAOUT was switched on by Start()
	ctrlByteOn := uint8(pcf8591_ANAON) | uint8(pcf8591_ALLSINGLE) | uint8(pcf8591_CHAN1)
	returnRead := []uint8{0x01, 0x02, 0x03, 232}
	// scaler for temperature is 255..0 => NTC °C, 232 relates to nearly 25°C
	// in TestTemperatureSensorDriverNtcScaling we have already used this NTC values
	want := 24.805280460718336
	// arrange reads
	numCallsRead := 0
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b = returnRead[0:len(b)]
		}
		if numCallsRead == 2 {
			b[0] = returnRead[len(returnRead)-1]
		}
		return len(b), nil
	}
	// act
	got, err := yl.ReadTemperature()
	got2, err2 := yl.Temperature()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(adaptor.written), 1)
	gobottest.Assert(t, adaptor.written[0], ctrlByteOn)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, got, want)
	gobottest.Assert(t, err2, nil)
	gobottest.Assert(t, got2, want)
}

func TestYL40DriverReadAIN2(t *testing.T) {
	// sequence to read the input with PCF8591, see there tests
	// arrange
	yl, adaptor := initTestYL40DriverWithStubbedAdaptor()
	adaptor.written = []byte{} // reset writes of Start() and former test
	// ANAOUT was switched on by Start()
	ctrlByteOn := uint8(pcf8591_ANAON) | uint8(pcf8591_ALLSINGLE) | uint8(pcf8591_CHAN2)
	returnRead := []uint8{0x01, 0x02, 0x03, 72}
	// scaler for analog input 2 is 0..255 => 0..3.3
	want := float64(returnRead[3]) * 33 / 2550
	// arrange reads
	numCallsRead := 0
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b = returnRead[0:len(b)]
		}
		if numCallsRead == 2 {
			b[0] = returnRead[len(returnRead)-1]
		}
		return len(b), nil
	}
	// act
	got, err := yl.ReadAIN2()
	got2, err2 := yl.AIN2()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(adaptor.written), 1)
	gobottest.Assert(t, adaptor.written[0], ctrlByteOn)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, got, want)
	gobottest.Assert(t, err2, nil)
	gobottest.Assert(t, got2, want)
}

func TestYL40DriverReadPotentiometer(t *testing.T) {
	// sequence to read the input with PCF8591, see there tests
	// arrange
	yl, adaptor := initTestYL40DriverWithStubbedAdaptor()
	adaptor.written = []byte{} // reset writes of Start() and former test
	// ANAOUT was switched on by Start()
	ctrlByteOn := uint8(pcf8591_ANAON) | uint8(pcf8591_ALLSINGLE) | uint8(pcf8591_CHAN3)
	returnRead := []uint8{0x01, 0x02, 0x03, 63}
	// scaler for potentiometer is 255..0 => -100..100
	want := float64(returnRead[3])*-200/255 + 100
	// arrange reads
	numCallsRead := 0
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b = returnRead[0:len(b)]
		}
		if numCallsRead == 2 {
			b[0] = returnRead[len(returnRead)-1]
		}
		return len(b), nil
	}
	// act
	got, err := yl.ReadPotentiometer()
	got2, err2 := yl.Potentiometer()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(adaptor.written), 1)
	gobottest.Assert(t, adaptor.written[0], ctrlByteOn)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, got, want)
	gobottest.Assert(t, err2, nil)
	gobottest.Assert(t, got2, want)
}

func TestYL40DriverAnalogWrite(t *testing.T) {
	// sequence to write the output of PCF8591, see there
	// arrange
	pcf, adaptor := initTestYL40DriverWithStubbedAdaptor()
	adaptor.written = []byte{} // reset writes of Start() and former test
	ctrlByteOn := uint8(pcf8591_ANAON)
	want := uint8(175)
	// write is scaled by 0..3.3V => 0..255
	write := float64(want) * 33 / 2550
	// arrange writes
	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		return len(b), nil
	}
	// act
	err := pcf.Write(write)
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(adaptor.written), 2)
	gobottest.Assert(t, adaptor.written[0], ctrlByteOn)
	gobottest.Assert(t, adaptor.written[1], want)
}

func TestYL40DriverStart(t *testing.T) {
	yl := NewYL40Driver(newI2cTestAdaptor())
	gobottest.Assert(t, yl.Start(), nil)
}

func TestYL40DriverHalt(t *testing.T) {
	yl := NewYL40Driver(newI2cTestAdaptor())
	gobottest.Assert(t, yl.Halt(), nil)
}

func fEqual(want interface{}, got interface{}) bool {
	return fmt.Sprintf("%v", want) == fmt.Sprintf("%v", got)
}
