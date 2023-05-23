package i2c

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*WiichuckDriver)(nil)

func initTestWiichuckDriverWithStubbedAdaptor() *WiichuckDriver {
	d := NewWiichuckDriver(newI2cTestAdaptor())
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d
}

func TestNewWiichuckDriver(t *testing.T) {
	var di interface{} = NewWiichuckDriver(newI2cTestAdaptor())
	d, ok := di.(*WiichuckDriver)
	if !ok {
		t.Errorf("NewWiichuckDriver() should have returned a *WiichuckDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Wiichuck"), true)
	gobottest.Assert(t, d.defaultAddress, 0x52)
	gobottest.Assert(t, d.interval, 10*time.Millisecond)
}

func TestWiichuckDriverStart(t *testing.T) {
	a := newI2cTestAdaptor()
	d := NewWiichuckDriver(a)
	a.Testi2cReadImpl(func(b []byte) (int, error) {
		copy(b, []byte{1, 2, 3, 4, 5, 6})
		return 6, nil
	})
	numberOfCyclesForEvery := 3
	d.interval = 1 * time.Millisecond
	sem := make(chan bool)

	gobottest.Assert(t, d.Start(), nil)

	go func() {
		for {
			time.Sleep(time.Duration(numberOfCyclesForEvery) * time.Millisecond)
			j := d.Joystick()
			if (j["sy_origin"] == float64(44)) &&
				(j["sx_origin"] == float64(45)) {
				sem <- true
				return
			}
		}
	}()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("origin not read correctly")
	}

}

func TestWiichuckDriverHalt(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestWiichuckDriverCanParse(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	d.update(decryptedValue)

	// - This should be done by WiichuckDriver.parse
	gobottest.Assert(t, d.data["sx"], float64(45))
	gobottest.Assert(t, d.data["sy"], float64(44))
	gobottest.Assert(t, d.data["z"], float64(0))
	gobottest.Assert(t, d.data["c"], float64(0))
}

func TestWiichuckDriverCanAdjustOrigins(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	d.update(decryptedValue)

	// - This should be done by WiichuckDriver.adjustOrigins
	gobottest.Assert(t, d.Joystick()["sx_origin"], float64(45))
	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(44))
}

func TestWiichuckDriverCButton(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	d.update(decryptedValue)

	// - This should be done by WiichuckDriver.updateButtons
	done := make(chan bool)

	d.On(d.Event(C), func(data interface{}) {
		gobottest.Assert(t, data, true)
		done <- true
	})

	d.update(decryptedValue)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Errorf("Did not receive 'C' event")
	}
}

func TestWiichuckDriverZButton(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	d.update(decryptedValue)

	done := make(chan bool)

	d.On(d.Event(Z), func(data interface{}) {
		gobottest.Assert(t, data, true)
		done <- true
	})

	d.update(decryptedValue)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Errorf("Did not receive 'Z' event")
	}
}

func TestWiichuckDriverUpdateJoystick(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}

	// - This should be done by WiichuckDriver.updateJoystick
	expectedData := map[string]float64{
		"x": float64(0),
		"y": float64(0),
	}

	done := make(chan bool)

	d.On(d.Event(Joystick), func(data interface{}) {
		gobottest.Assert(t, data, expectedData)
		done <- true
	})

	d.update(decryptedValue)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Errorf("Did not receive 'Joystick' event")
	}
}

func TestWiichuckDriverEncrypted(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()
	encryptedValue := []byte{1, 1, 2, 2, 3, 3}

	d.update(encryptedValue)

	gobottest.Assert(t, d.data["sx"], float64(0))
	gobottest.Assert(t, d.data["sy"], float64(0))
	gobottest.Assert(t, d.data["z"], float64(0))
	gobottest.Assert(t, d.data["c"], float64(0))

	gobottest.Assert(t, d.Joystick()["sx_origin"], float64(-1))
	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(-1))
}

func TestWiichuckDriverSetJoystickDefaultValue(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(-1))

	d.setJoystickDefaultValue("sy_origin", float64(2))

	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(2))

	// when current default value is not -1 it keeps the current value
	d.setJoystickDefaultValue("sy_origin", float64(20))

	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(2))
}

func TestWiichuckDriverCalculateJoystickValue(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.calculateJoystickValue(float64(20), float64(5)), float64(15))
	gobottest.Assert(t, d.calculateJoystickValue(float64(1), float64(2)), float64(-1))
	gobottest.Assert(t, d.calculateJoystickValue(float64(10), float64(5)), float64(5))
	gobottest.Assert(t, d.calculateJoystickValue(float64(5), float64(10)), float64(-5))
}

func TestWiichuckDriverIsEncrypted(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	encryptedValue := []byte{1, 1, 2, 2, 3, 3}
	gobottest.Assert(t, d.isEncrypted(encryptedValue), true)

	encryptedValue = []byte{42, 42, 24, 24, 30, 30}
	gobottest.Assert(t, d.isEncrypted(encryptedValue), true)

	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	gobottest.Assert(t, d.isEncrypted(decryptedValue), false)

	decryptedValue = []byte{1, 1, 2, 2, 5, 6}
	gobottest.Assert(t, d.isEncrypted(decryptedValue), false)

	decryptedValue = []byte{1, 1, 2, 3, 3, 3}
	gobottest.Assert(t, d.isEncrypted(decryptedValue), false)
}

func TestWiichuckDriverDecode(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.decode(byte(0)), float64(46))
	gobottest.Assert(t, d.decode(byte(100)), float64(138))
	gobottest.Assert(t, d.decode(byte(200)), float64(246))
	gobottest.Assert(t, d.decode(byte(254)), float64(0))
}

func TestWiichuckDriverParse(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.data["sx"], float64(0))
	gobottest.Assert(t, d.data["sy"], float64(0))
	gobottest.Assert(t, d.data["z"], float64(0))
	gobottest.Assert(t, d.data["c"], float64(0))

	// First pass
	d.parse([]byte{12, 23, 34, 45, 56, 67})

	gobottest.Assert(t, d.data["sx"], float64(50))
	gobottest.Assert(t, d.data["sy"], float64(23))
	gobottest.Assert(t, d.data["z"], float64(1))
	gobottest.Assert(t, d.data["c"], float64(2))

	// Second pass
	d.parse([]byte{70, 81, 92, 103, 204, 205})

	gobottest.Assert(t, d.data["sx"], float64(104))
	gobottest.Assert(t, d.data["sy"], float64(93))
	gobottest.Assert(t, d.data["z"], float64(1))
	gobottest.Assert(t, d.data["c"], float64(0))
}

func TestWiichuckDriverAdjustOrigins(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(-1))
	gobottest.Assert(t, d.Joystick()["sx_origin"], float64(-1))

	// First pass
	d.parse([]byte{1, 2, 3, 4, 5, 6})
	d.adjustOrigins()

	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(44))
	gobottest.Assert(t, d.Joystick()["sx_origin"], float64(45))

	// Second pass
	d = initTestWiichuckDriverWithStubbedAdaptor()

	d.parse([]byte{61, 72, 83, 94, 105, 206})
	d.adjustOrigins()

	gobottest.Assert(t, d.Joystick()["sy_origin"], float64(118))
	gobottest.Assert(t, d.Joystick()["sx_origin"], float64(65))
}

func TestWiichuckDriverSetName(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestWiichuckDriverOptions(t *testing.T) {
	d := NewWiichuckDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
