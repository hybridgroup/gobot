package i2c

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

// --------- HELPERS
func initTestWiichuckDriver() (driver *WiichuckDriver) {
	driver, _ = initTestWiichuckDriverWithStubbedAdaptor()
	return
}

func initTestWiichuckDriverWithStubbedAdaptor() (*WiichuckDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewWiichuckDriver(adaptor, "bot"), adaptor
}

// --------- TESTS
func TestWiichuckDriver(t *testing.T) {
	// Does it implement gobot.DriverInterface?
	var _ gobot.DriverInterface = (*WiichuckDriver)(nil)

	// Does its adaptor implements the I2cInterface?
	driver := initTestWiichuckDriver()
	var _ I2cInterface = driver.adaptor()
}

func TestNewWiichuckDriver(t *testing.T) {
	// Does it return a pointer to an instance of WiichuckDriver?
	var bm interface{} = NewWiichuckDriver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := bm.(*WiichuckDriver)
	if !ok {
		t.Errorf("NewWiichuckDriver() should have returned a *WiichuckDriver")
	}
}

func TestWiichuckDriverStart(t *testing.T) {
	sem := make(chan bool)
	wii, adaptor := initTestWiichuckDriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() []byte {
		return []byte{1, 2, 3, 4, 5, 6}
	}

	numberOfCyclesForEvery := 3

	wii.SetInterval(1 * time.Millisecond)
	gobot.Assert(t, wii.Start(), true)

	go func() {
		for {
			<-time.After(time.Duration(numberOfCyclesForEvery) * time.Millisecond)
			if (wii.joystick["sy_origin"] == float64(44)) &&
				(wii.joystick["sx_origin"] == float64(45)) {
				sem <- true
			}
		}
	}()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("origin not read correctly")
	}

}

func TestWiichuckDriverInit(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.Init(), true)
}

func TestWiichuckDriverHalt(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.Halt(), true)
}

func TestWiichuckDriverUpdate(t *testing.T) {
	wii := initTestWiichuckDriver()

	// ------ When value is not encrypted
	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	wii.update(decryptedValue)

	// - This should be done by WiichuckDriver.parse
	gobot.Assert(t, wii.data["sx"], float64(45))
	gobot.Assert(t, wii.data["sy"], float64(44))
	gobot.Assert(t, wii.data["z"], float64(0))
	gobot.Assert(t, wii.data["c"], float64(0))

	// - This should be done by WiichuckDriver.adjustOrigins
	gobot.Assert(t, wii.joystick["sx_origin"], float64(45))
	gobot.Assert(t, wii.joystick["sy_origin"], float64(44))

	// - This should be done by WiichuckDriver.updateButtons
	chann := make(chan bool)

	gobot.On(wii.Event("c"), func(data interface{}) {
		gobot.Assert(t, data, true)
		chann <- true
	})
	<-chann

	chann = make(chan bool)
	wii.update(decryptedValue)

	gobot.On(wii.Event("z"), func(data interface{}) {
		gobot.Assert(t, data, true)
		chann <- true
	})
	<-chann

	// - This should be done by WiichuckDriver.updateJoystick
	chann = make(chan bool)
	wii.update(decryptedValue)

	expectedData := map[string]float64{
		"x": float64(0),
		"y": float64(0),
	}

	gobot.On(wii.Event("joystick"), func(data interface{}) {
		gobot.Assert(t, data, expectedData)
		chann <- true
	})
	<-chann

	// ------ When value is encrypted
	wii = initTestWiichuckDriver()
	encryptedValue := []byte{1, 1, 2, 2, 3, 3}

	wii.update(encryptedValue)

	gobot.Assert(t, wii.data["sx"], float64(0))
	gobot.Assert(t, wii.data["sy"], float64(0))
	gobot.Assert(t, wii.data["z"], float64(0))
	gobot.Assert(t, wii.data["c"], float64(0))

	gobot.Assert(t, wii.joystick["sx_origin"], float64(-1))
	gobot.Assert(t, wii.joystick["sy_origin"], float64(-1))
}

func TestWiichuckDriverSetJoystickDefaultValue(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.joystick["sy_origin"], float64(-1))

	wii.setJoystickDefaultValue("sy_origin", float64(2))

	gobot.Assert(t, wii.joystick["sy_origin"], float64(2))

	// when current default value is not -1 it keeps the current value

	wii.setJoystickDefaultValue("sy_origin", float64(20))

	gobot.Assert(t, wii.joystick["sy_origin"], float64(2))

}

func TestWiichuckDriverCalculateJoystickValue(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.calculateJoystickValue(float64(20), float64(5)), float64(15))
	gobot.Assert(t, wii.calculateJoystickValue(float64(1), float64(2)), float64(-1))
	gobot.Assert(t, wii.calculateJoystickValue(float64(10), float64(5)), float64(5))
	gobot.Assert(t, wii.calculateJoystickValue(float64(5), float64(10)), float64(-5))
}

func TestWiichuckDriverIsEncrypted(t *testing.T) {
	wii := initTestWiichuckDriver()

	encryptedValue := []byte{1, 1, 2, 2, 3, 3}
	gobot.Assert(t, wii.isEncrypted(encryptedValue), true)

	encryptedValue = []byte{42, 42, 24, 24, 30, 30}
	gobot.Assert(t, wii.isEncrypted(encryptedValue), true)

	decryptedValue := []byte{1, 2, 3, 4, 5, 6}
	gobot.Assert(t, wii.isEncrypted(decryptedValue), false)

	decryptedValue = []byte{1, 1, 2, 2, 5, 6}
	gobot.Assert(t, wii.isEncrypted(decryptedValue), false)

	decryptedValue = []byte{1, 1, 2, 3, 3, 3}
	gobot.Assert(t, wii.isEncrypted(decryptedValue), false)
}

func TestWiichuckDriverDecode(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.decode(byte(0)), float64(46))
	gobot.Assert(t, wii.decode(byte(100)), float64(138))
	gobot.Assert(t, wii.decode(byte(200)), float64(246))
	gobot.Assert(t, wii.decode(byte(254)), float64(0))
}

func TestWiichuckDriverParse(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.data["sx"], float64(0))
	gobot.Assert(t, wii.data["sy"], float64(0))
	gobot.Assert(t, wii.data["z"], float64(0))
	gobot.Assert(t, wii.data["c"], float64(0))

	// First pass
	wii.parse([]byte{12, 23, 34, 45, 56, 67})

	gobot.Assert(t, wii.data["sx"], float64(50))
	gobot.Assert(t, wii.data["sy"], float64(23))
	gobot.Assert(t, wii.data["z"], float64(1))
	gobot.Assert(t, wii.data["c"], float64(2))

	// Second pass
	wii.parse([]byte{70, 81, 92, 103, 204, 205})

	gobot.Assert(t, wii.data["sx"], float64(104))
	gobot.Assert(t, wii.data["sy"], float64(93))
	gobot.Assert(t, wii.data["z"], float64(1))
	gobot.Assert(t, wii.data["c"], float64(0))
}

func TestWiichuckDriverAdjustOrigins(t *testing.T) {
	wii := initTestWiichuckDriver()

	gobot.Assert(t, wii.joystick["sy_origin"], float64(-1))
	gobot.Assert(t, wii.joystick["sx_origin"], float64(-1))

	// First pass
	wii.parse([]byte{1, 2, 3, 4, 5, 6})
	wii.adjustOrigins()

	gobot.Assert(t, wii.joystick["sy_origin"], float64(44))
	gobot.Assert(t, wii.joystick["sx_origin"], float64(45))

	// Second pass
	wii = initTestWiichuckDriver()

	wii.parse([]byte{61, 72, 83, 94, 105, 206})
	wii.adjustOrigins()

	gobot.Assert(t, wii.joystick["sy_origin"], float64(118))
	gobot.Assert(t, wii.joystick["sx_origin"], float64(65))
}

func TestWiichuckDriverUpdateButtons(t *testing.T) {
	//when data["c"] is 0
	chann := make(chan bool)
	wii := initTestWiichuckDriver()

	wii.data["c"] = 0

	wii.updateButtons()

	gobot.On(wii.Event("c"), func(data interface{}) {
		gobot.Assert(t, true, data)
		chann <- true
	})
	<-chann

	//when data["z"] is 0
	chann = make(chan bool)
	wii = initTestWiichuckDriver()

	wii.data["z"] = 0

	wii.updateButtons()

	gobot.On(wii.Event("z"), func(data interface{}) {
		gobot.Assert(t, true, data)
		chann <- true
	})
	<-chann
}

func TestWiichuckDriverUpdateJoystick(t *testing.T) {
	chann := make(chan bool)
	wii := initTestWiichuckDriver()

	// First pass
	wii.data["sx"] = 40
	wii.data["sy"] = 55
	wii.joystick["sx_origin"] = 1
	wii.joystick["sy_origin"] = 5

	wii.updateJoystick()

	expectedData := map[string]float64{
		"x": float64(39),
		"y": float64(50),
	}

	gobot.On(wii.Event("joystick"), func(data interface{}) {
		gobot.Assert(t, data, expectedData)
		chann <- true
	})
	<-chann

	//// Second pass
	chann = make(chan bool)
	wii = initTestWiichuckDriver()

	wii.data["sx"] = 178
	wii.data["sy"] = 34
	wii.joystick["sx_origin"] = 14
	wii.joystick["sy_origin"] = 27

	wii.updateJoystick()

	expectedData = map[string]float64{
		"x": float64(164),
		"y": float64(7),
	}

	gobot.On(wii.Event("joystick"), func(data interface{}) {
		gobot.Assert(t, data, expectedData)
		chann <- true
	})
	<-chann
}
