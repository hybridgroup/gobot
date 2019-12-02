package i2c

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func createReaderHelper(bytesToRead []byte) func(b []byte) (int, error) {
	index := 0
	return func(b []byte) (i int, err error) {
		i, err = 0, nil
		for i < 2 && index < len(bytesToRead) {
			b[i] = bytesToRead[index]
			index++
			i++
		}
		return
	}
}

func createWriterHelper(failAfterNWrites int) func(bytes []byte) (i int, e error) {
	index := 0

	return func(bytes []byte) (i int, e error) {
		index++
		if index >= failAfterNWrites {
			return 0, fmt.Errorf("failed")
		} else {
			return len(bytes), nil
		}
	}
}

func TestNewSparkFunDev15316(t *testing.T) {
	c := newI2cTestAdaptor()
	driver := NewSparkFunDev15316(c)

	if driver.Connector != c {
		t.Errorf("expected to find the same connection but didn't")
	}
}

func TestSparkFunDev15316_Init(t *testing.T) {
	c := newI2cTestAdaptor()
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	if driver.channels == nil {
		t.Fatalf("expected channels to be setup by now")
	}

	tmp := []byte{0, 0x20, 0, 0x10, 0xfe, 0x1e, 0, 0x20}
	if !reflect.DeepEqual(c.written, tmp) {
		t.Fatalf("expected %v but found %v", c.written, tmp)
	}

	if len(c.written) != 8 {
		t.Fatalf("expected to find 4 bytes written instead found %v", len(c.written))
	}

	if driver.oscClock != defaultClockFreq {
		t.Fatalf("expected oscClock to be %v but found %v", defaultClockFreq, driver.oscClock)
	}

	if driver.updateRate != defaultUpdateRate {
		t.Fatalf("expected updateRate to be %v but found %v", defaultUpdateRate, driver.updateRate)
	}

	err = driver.InitWithRate(minPWMUpdateRate)
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	if driver.updateRate != minPWMUpdateRate {
		t.Fatalf("expected updateRate to be %v but found %v", minPWMUpdateRate, driver.updateRate)
	}

	err = driver.InitWithRate(minPWMUpdateRate - 1)
	if err == nil {
		t.Errorf("expected an error but found none")
	}

	err = driver.InitWithRate(maxPWMUpdateRate + 1)
	if err == nil {
		t.Errorf("expected an error but found none")
	}
}

func TestSparkFunDev15316_Init2(t *testing.T) {
	var nilDriver *SparkFunDev15316
	err := nilDriver.Init()
	if err == nil {
		t.Fatalf("expected an error but found none")
	}

	c := newI2cTestAdaptor()
	c.i2cConnectErr = true
	driver := NewSparkFunDev15316(c)
	err = driver.Init()
	if err == nil {
		t.Errorf("expected an error on Init()")
	}

	c = newI2cTestAdaptor()
	driver = NewSparkFunDev15316(c)
	err = driver.InitWithRateAndClock(defaultUpdateRate, maxOSCClock+1)
	if err == nil {
		t.Errorf("expected an error on Init()")
	}

	////////////////////////////////////
	c = newI2cTestAdaptor()
	driver = NewSparkFunDev15316(c)
	c.i2cWriteImpl = createWriterHelper(1)

	err = driver.Init()
	if err == nil {
		t.Errorf("expected an error on Init()")
	}

	////////////////////////////////////
	c = newI2cTestAdaptor()
	driver = NewSparkFunDev15316(c)
	c.i2cWriteImpl = createWriterHelper(2)

	err = driver.Init()
	if err == nil {
		t.Errorf("expected an error on Init()")
	}
	////////////////////////////////////
	c = newI2cTestAdaptor()
	driver = NewSparkFunDev15316(c)
	c.i2cWriteImpl = createWriterHelper(3)

	err = driver.Init()
	if err == nil {
		t.Errorf("expected an error on Init()")
	}
	////////////////////////////////////
	c = newI2cTestAdaptor()
	driver = NewSparkFunDev15316(c)
	c.i2cWriteImpl = createWriterHelper(4)

	err = driver.Init()
	if err == nil {
		t.Errorf("expected an error on Init()")
	}

}

func TestSparkFunDev15316_Configure(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0, InstantPower(true))
	if err != nil {
		t.Fatalf("expected no error for configure instead found %v", err)
	}

	if _, has := driver.channels[Ch0]; !has {
		t.Fatalf("expected to find a channel")
	}

	channel := driver.channels[Ch0]
	if channel.calMaxPos != 1638 {
		t.Fatalf("expected calMaxPos to be %v but found %v", 1638, channel.calMaxPos)
	}

	if channel.currentPos != 1228 {
		t.Fatalf("expected the current position to be 1228")
	}

	if channel.calMinPos != 819 {
		t.Fatalf("expected calMinPos to be %v but found %v", 819, channel.calMinPos)
	}

	if channel.calCenterPos != 1228 {
		t.Fatalf("expected calCenterPos to be %v but found %v", 1228, channel.calCenterPos)
	}

	if math.Abs(channel.calPosFactor-8.19) > 1e-6 {
		t.Fatalf("expected calPosFactor to be %v but found %v", 8.19, channel.calPosFactor)
	}

	if !channel.instantPower {
		t.Fatalf("expected to find instant power")
	}

	err = driver.Configure(Channel(16))
	if err == nil {
		t.Fatalf("expected an error for configure")
	}
}

func TestSparkFunDev15316_Configure2(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0, CenterOffset(-1200))
	if err == nil {
		t.Fatalf("expected to find an error")
	}

	err = driver.Configure(Ch0, CenterOffset(+1200))
	if err == nil {
		t.Fatalf("expected to find an error")
	}

	driver.channels = nil
	err = driver.Configure(Ch0)
	if err == nil {
		t.Fatalf("expected to find an error")
	}

}

func TestSparkFunDev15316_SetRadialPos(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0, InstantPower(true), Mode(RotationPositionalServo))
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	c.written = make([]byte, 0)

	err = driver.SetRadialPos(Ch0, 50)
	if err != nil {
		t.Fatalf("expected no error but found: %v", err)
	}

	tmp := []byte{0x08, 0xcc, 0x04}
	if !reflect.DeepEqual(c.written, tmp) {
		t.Fatalf("expected %x but found: %x", tmp, c.written)
	}
}

func TestSparkFunDev15316_SetRadialPos2(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0,
		InstantPower(false),
		Mode(RotationPositionalServo),
		RampingUpdateFreq(LinearRamping),
	)
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	c.written = make([]byte, 0)

	err = driver.SetRadialPos(Ch0, 0)
	if err != nil {
		t.Fatalf("expected no error but found: %v", err)
	}

	tmp := []byte{0x08, 0xcc, 0x04, 0x08, 0xcb, 0x04}
	for i := 0; i < len(tmp); i++ {
		if c.written[i] != tmp[i] {
			t.Fatalf("expected %x but found: %x", tmp, c.written)
		}
	}
}

func TestSparkFunDev15316_SetRadialPos3(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0,
		InstantPower(false),
		Mode(RotationPositionalServo),
		RampingUpdateFreq(LinearRamping),
	)
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	c.written = make([]byte, 0)

	err = driver.SetRadialPos(Ch0, 100)
	if err != nil {
		t.Fatalf("expected no error but found: %v", err)
	}

	tmp := []byte{0x08, 0xcc, 0x04, 0x08, 0xcd, 0x04}
	for i := 0; i < len(tmp); i++ {
		if c.written[i] != tmp[i] {
			t.Fatalf("expected %x but found: %x", tmp, c.written)
		}
	}
}

func TestSparkFunDev15316_SetRadialPos4(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0)
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	////////////////////////

	err = driver.SetRadialPos(Ch15+1, 50)
	if err == nil {
		t.Errorf("expected an error on SetRadialPos()")
	}

	////////////////////////
	tmp := driver.channels
	driver.channels = nil

	err = driver.SetRadialPos(Ch0, 50)
	if err == nil {
		t.Errorf("expected an error on SetRadialPos()")
	}
	driver.channels = tmp
	///////////////////////

	err = driver.SetRadialPos(Ch1, 50)
	if err == nil {
		t.Errorf("expected an error on SetRadialPos()")
	}

	////////////////////////
	err = driver.SetRadialPos(Ch0, 50)
	if err == nil {
		t.Errorf("expected an error on SetRadialPos()")
	}
}

func TestSparkFunDev15316_SetContinuous(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0, InstantPower(true), Mode(RotationContinuousServo))
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	c.written = make([]byte, 0)

	err = driver.SetContinuous(Ch0, Stop)
	if err != nil {
		t.Fatalf("expected no error but found: %v", err)
	}

	tmp := []byte{0x08, 0xcc, 0x04}
	if !reflect.DeepEqual(c.written, tmp) {
		t.Fatalf("expected %x but found: %x", tmp, c.written)
	}

	err = driver.SetContinuous(Ch0, CCW)
	if err != nil {
		t.Fatalf("expected no error but found: %v", err)
	}

	err = driver.SetContinuous(Ch0, CW)
	if err != nil {
		t.Fatalf("expected no error but found: %v", err)
	}
}

func TestSparkFunDev15316_SetContinuous2(t *testing.T) {
	c := newI2cTestAdaptor()
	c.i2cReadImpl = createReaderHelper([]byte{0xcc, 0x04, 0xcc, 0x04})
	driver := NewSparkFunDev15316(c)

	err := driver.Init()
	if err != nil {
		t.Errorf("expected no error on Init()")
	}

	err = driver.Configure(Ch0, InstantPower(true), Mode(RotationContinuousServo))
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	////////////////////////
	err = driver.SetContinuous(Ch15+1, Stop)
	if err == nil {
		t.Fatalf("expected an error")
	}
	////////////////////////
	err = driver.SetContinuous(Ch1, Stop)
	if err == nil {
		t.Fatalf("expected an error")
	}

	////////////////////////
	err = driver.Configure(Ch1, InstantPower(true), Mode(RotationPositionalServo))
	if err != nil {
		t.Errorf("expected no error on Configure()")
	}

	err = driver.SetContinuous(Ch1, Stop)
	if err == nil {
		t.Fatalf("expected an error")
	}

	////////////////////////
	err = driver.SetContinuous(Ch0, CW+1)
	if err == nil {
		t.Fatalf("expected an error")
	}

	////////////////////////
	driver.channels = nil

	err = driver.SetContinuous(Ch0, CW)
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestStringer(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{fmt.Sprint(RotationContinuousServo), "RotationContinuous"},
		{fmt.Sprint(RotationPositionalServo), "RotationPositional"},
		{fmt.Sprint(LinearServo), "Linear"},

		{fmt.Sprint(Stop), "Stop"},
		{fmt.Sprint(CCW), "CCW"},
		{fmt.Sprint(CW), "CW"},

		{fmt.Sprint(LogRamping), "Log"},
		{fmt.Sprint(LinearRamping), "Linear"},
		{fmt.Sprint(ExponentialRamping), "Exponential"},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.input != test.expected {
				t.Errorf("expected '%v' but found '%v'", test.expected, test.input)
			}
		})
	}
}
