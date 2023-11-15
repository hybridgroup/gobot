package i2c

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTestYL40DriverWithStubbedAdaptor() (*YL40Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	yl := NewYL40Driver(adaptor, WithPCF8591With400kbitStabilization(0, 2))
	WithPCF8591ForceRefresh(1)(yl.PCF8591Driver)
	_ = yl.Start()
	return yl, adaptor
}

func TestYL40Driver(t *testing.T) {
	// arrange, act
	yl := NewYL40Driver(newI2cTestAdaptor())
	// assert
	assert.NotNil(t, yl.PCF8591Driver)
	assert.Equal(t, time.Duration(0), yl.conf.sensors[YL40Bri].interval)
	assert.NotNil(t, yl.conf.sensors[YL40Bri].scaler)
	assert.Equal(t, time.Duration(0), yl.conf.sensors[YL40Temp].interval)
	assert.NotNil(t, yl.conf.sensors[YL40Temp].scaler)
	assert.Equal(t, time.Duration(0), yl.conf.sensors[YL40AIN2].interval)
	assert.NotNil(t, yl.conf.sensors[YL40AIN2].scaler)
	assert.Equal(t, time.Duration(0), yl.conf.sensors[YL40Poti].interval)
	assert.NotNil(t, yl.conf.sensors[YL40Poti].scaler)
	assert.NotNil(t, yl.conf.aOutScaler)
	assert.NotNil(t, yl.aBri)
	assert.NotNil(t, yl.aTemp)
	assert.NotNil(t, yl.aAIN2)
	assert.NotNil(t, yl.aPoti)
	assert.NotNil(t, yl.aOut)
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
	assert.Equal(t, time.Duration(100), yl.conf.sensors[YL40Bri].interval)
	assert.Equal(t, time.Duration(101), yl.conf.sensors[YL40Temp].interval)
	assert.Equal(t, time.Duration(102), yl.conf.sensors[YL40AIN2].interval)
	assert.Equal(t, time.Duration(103), yl.conf.sensors[YL40Poti].interval)
}

func TestYL40DriverWithYL40InputScaler(t *testing.T) {
	// arrange
	yl := NewYL40Driver(newI2cTestAdaptor())
	f1 := func(input int) float64 { return 0.1 }
	f2 := func(input int) float64 { return 0.2 }
	f3 := func(input int) float64 { return 0.3 }
	f4 := func(input int) float64 { return 0.4 }
	// act
	WithYL40InputScaler(YL40Bri, f1)(yl)
	WithYL40InputScaler(YL40Temp, f2)(yl)
	WithYL40InputScaler(YL40AIN2, f3)(yl)
	WithYL40InputScaler(YL40Poti, f4)(yl)
	// assert
	assert.True(t, fEqual(yl.conf.sensors[YL40Bri].scaler, f1))
	assert.True(t, fEqual(yl.conf.sensors[YL40Temp].scaler, f2))
	assert.True(t, fEqual(yl.conf.sensors[YL40AIN2].scaler, f3))
	assert.True(t, fEqual(yl.conf.sensors[YL40Poti].scaler, f4))
}

func TestYL40DriverWithYL40WithYL40OutputScaler(t *testing.T) {
	// arrange
	yl := NewYL40Driver(newI2cTestAdaptor())
	fo := func(input float64) int { return 123 }
	// act
	WithYL40OutputScaler(fo)(yl)
	// assert
	assert.True(t, fEqual(yl.conf.aOutScaler, fo))
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
	require.NoError(t, err)
	assert.Len(t, adaptor.written, 1)
	assert.Equal(t, ctrlByteOn, adaptor.written[0])
	assert.Equal(t, 2, numCallsRead)
	assert.InDelta(t, want, got, 0.0)
	require.NoError(t, err2)
	assert.InDelta(t, want, got2, 0.0)
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
	require.NoError(t, err)
	assert.Len(t, adaptor.written, 1)
	assert.Equal(t, ctrlByteOn, adaptor.written[0])
	assert.Equal(t, 2, numCallsRead)
	assert.InDelta(t, want, got, 0.0)
	require.NoError(t, err2)
	assert.InDelta(t, want, got2, 0.0)
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
	require.NoError(t, err)
	assert.Len(t, adaptor.written, 1)
	assert.Equal(t, ctrlByteOn, adaptor.written[0])
	assert.Equal(t, 2, numCallsRead)
	assert.InDelta(t, want, got, 0.0)
	require.NoError(t, err2)
	assert.InDelta(t, want, got2, 0.0)
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
	require.NoError(t, err)
	assert.Len(t, adaptor.written, 1)
	assert.Equal(t, ctrlByteOn, adaptor.written[0])
	assert.Equal(t, 2, numCallsRead)
	assert.InDelta(t, want, got, 0.0)
	require.NoError(t, err2)
	assert.InDelta(t, want, got2, 0.0)
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
	require.NoError(t, err)
	assert.Len(t, adaptor.written, 2)
	assert.Equal(t, ctrlByteOn, adaptor.written[0])
	assert.Equal(t, want, adaptor.written[1])
}

func TestYL40DriverStart(t *testing.T) {
	yl := NewYL40Driver(newI2cTestAdaptor())
	require.NoError(t, yl.Start())
}

func TestYL40DriverHalt(t *testing.T) {
	yl := NewYL40Driver(newI2cTestAdaptor())
	require.NoError(t, yl.Halt())
}

func fEqual(want interface{}, got interface{}) bool {
	return fmt.Sprintf("%v", want) == fmt.Sprintf("%v", got)
}
