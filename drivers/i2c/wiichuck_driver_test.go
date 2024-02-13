package i2c

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
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
		require.Fail(t, "NewWiichuckDriver() should have returned a *WiichuckDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "Wiichuck"))
	assert.Equal(t, 0x52, d.defaultAddress)
	assert.Equal(t, 10*time.Millisecond, d.interval)
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

	require.NoError(t, d.Start())

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
		require.Fail(t, "origin not read correctly")
	}
}

func TestWiichuckDriverHalt(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestWiichuckDriverCanParse(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	_ = d.update(decryptedValue)

	// - This should be done by WiichuckDriver.parse
	assert.InDelta(t, float64(45), d.data["sx"], 0.0)
	assert.InDelta(t, float64(44), d.data["sy"], 0.0)
	assert.InDelta(t, float64(0), d.data["z"], 0.0)
	assert.InDelta(t, float64(0), d.data["c"], 0.0)
}

func TestWiichuckDriverCanAdjustOrigins(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	_ = d.update(decryptedValue)

	// - This should be done by WiichuckDriver.adjustOrigins
	assert.InDelta(t, float64(45), d.Joystick()["sx_origin"], 0.0)
	assert.InDelta(t, float64(44), d.Joystick()["sy_origin"], 0.0)
}

func TestWiichuckDriverCButton(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	_ = d.update(decryptedValue)

	// - This should be done by WiichuckDriver.updateButtons
	done := make(chan bool)

	_ = d.On(d.Event(C), func(data interface{}) {
		assert.Equal(t, true, data)
		done <- true
	})

	_ = d.update(decryptedValue)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		require.Fail(t, "Did not receive 'C' event")
	}
}

func TestWiichuckDriverZButton(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	_ = d.update(decryptedValue)

	done := make(chan bool)

	_ = d.On(d.Event(Z), func(data interface{}) {
		assert.Equal(t, true, data)
		done <- true
	})

	_ = d.update(decryptedValue)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		require.Fail(t, "Did not receive 'Z' event")
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

	_ = d.On(d.Event(Joystick), func(data interface{}) {
		assert.Equal(t, expectedData, data)
		done <- true
	})

	_ = d.update(decryptedValue)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		require.Fail(t, "Did not receive 'Joystick' event")
	}
}

func TestWiichuckDriverEncrypted(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()
	encryptedValue := []byte{1, 1, 2, 2, 3, 3}

	_ = d.update(encryptedValue)

	assert.InDelta(t, float64(0), d.data["sx"], 0.0)
	assert.InDelta(t, float64(0), d.data["sy"], 0.0)
	assert.InDelta(t, float64(0), d.data["z"], 0.0)
	assert.InDelta(t, float64(0), d.data["c"], 0.0)

	assert.InDelta(t, float64(-1), d.Joystick()["sx_origin"], 0.0)
	assert.InDelta(t, float64(-1), d.Joystick()["sy_origin"], 0.0)
}

func TestWiichuckDriverSetJoystickDefaultValue(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	assert.InDelta(t, float64(-1), d.Joystick()["sy_origin"], 0.0)

	d.setJoystickDefaultValue("sy_origin", float64(2))

	assert.InDelta(t, float64(2), d.Joystick()["sy_origin"], 0.0)

	// when current default value is not -1 it keeps the current value
	d.setJoystickDefaultValue("sy_origin", float64(20))

	assert.InDelta(t, float64(2), d.Joystick()["sy_origin"], 0.0)
}

func TestWiichuckDriverCalculateJoystickValue(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	assert.InDelta(t, float64(15), d.calculateJoystickValue(float64(20), float64(5)), 0.0)
	assert.InDelta(t, float64(-1), d.calculateJoystickValue(float64(1), float64(2)), 0.0)
	assert.InDelta(t, float64(5), d.calculateJoystickValue(float64(10), float64(5)), 0.0)
	assert.InDelta(t, float64(-5), d.calculateJoystickValue(float64(5), float64(10)), 0.0)
}

func TestWiichuckDriverIsEncrypted(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	encryptedValue := []byte{1, 1, 2, 2, 3, 3}
	assert.True(t, d.isEncrypted(encryptedValue))

	encryptedValue = []byte{42, 42, 24, 24, 30, 30}
	assert.True(t, d.isEncrypted(encryptedValue))

	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	assert.False(t, d.isEncrypted(decryptedValue))

	decryptedValue = []byte{1, 1, 2, 2, 5, 6}
	assert.False(t, d.isEncrypted(decryptedValue))

	decryptedValue = []byte{1, 1, 2, 3, 3, 3}
	assert.False(t, d.isEncrypted(decryptedValue))
}

func TestWiichuckDriverDecode(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	assert.InDelta(t, float64(46), d.decode(byte(0)), 0.0)
	assert.InDelta(t, float64(138), d.decode(byte(100)), 0.0)
	assert.InDelta(t, float64(246), d.decode(byte(200)), 0.0)
	assert.InDelta(t, float64(0), d.decode(byte(254)), 0.0)
}

func TestWiichuckDriverParse(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	assert.InDelta(t, float64(0), d.data["sx"], 0.0)
	assert.InDelta(t, float64(0), d.data["sy"], 0.0)
	assert.InDelta(t, float64(0), d.data["z"], 0.0)
	assert.InDelta(t, float64(0), d.data["c"], 0.0)

	// First pass
	d.parse([]byte{12, 23, 34, 45, 56, 67})

	assert.InDelta(t, float64(50), d.data["sx"], 0.0)
	assert.InDelta(t, float64(23), d.data["sy"], 0.0)
	assert.InDelta(t, float64(1), d.data["z"], 0.0)
	assert.InDelta(t, float64(2), d.data["c"], 0.0)

	// Second pass
	d.parse([]byte{70, 81, 92, 103, 204, 205})

	assert.InDelta(t, float64(104), d.data["sx"], 0.0)
	assert.InDelta(t, float64(93), d.data["sy"], 0.0)
	assert.InDelta(t, float64(1), d.data["z"], 0.0)
	assert.InDelta(t, float64(0), d.data["c"], 0.0)
}

func TestWiichuckDriverAdjustOrigins(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()

	assert.InDelta(t, float64(-1), d.Joystick()["sy_origin"], 0.0)
	assert.InDelta(t, float64(-1), d.Joystick()["sx_origin"], 0.0)

	// First pass
	d.parse([]byte{1, 2, 3, 4, 5, 6})
	d.adjustOrigins()

	assert.InDelta(t, float64(44), d.Joystick()["sy_origin"], 0.0)
	assert.InDelta(t, float64(45), d.Joystick()["sx_origin"], 0.0)

	// Second pass
	d = initTestWiichuckDriverWithStubbedAdaptor()

	d.parse([]byte{61, 72, 83, 94, 105, 206})
	d.adjustOrigins()

	assert.InDelta(t, float64(118), d.Joystick()["sy_origin"], 0.0)
	assert.InDelta(t, float64(65), d.Joystick()["sx_origin"], 0.0)
}

func TestWiichuckDriverSetName(t *testing.T) {
	d := initTestWiichuckDriverWithStubbedAdaptor()
	d.SetName("TESTME")
	assert.Equal(t, "TESTME", d.Name())
}

func TestWiichuckDriverOptions(t *testing.T) {
	d := NewWiichuckDriver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}
