package gopigo3

import (
	"testing"
	"time"

	"log"

	"github.com/hybridgroup/gobot"
)

type GoPiGoTestAdaptor struct {
	t       *testing.T
	asserts []byte
}

func initTestGoPiGo(t *testing.T, asserts []byte) *GoPiGo3Driver {
	return &GoPiGoTestAdaptor{
		t:       t,
		asserts: asserts,
	}
}

func ExampleNewGoPiGoDriver() {
	adapter, err := NewAdaptor()
	if err != nil {
		panic(err)
	}

	driver := NewGoPiGo3Driver(adapter)
	work := func() {
		gobot.Every(10*time.Second, func() {
			driver.SetLed(LED_EYE_LEFT, 10, 10, 10)
			time.Sleep(1 * time.Second)
			driver.SetLed(LED_EYE_LEFT, 0, 0, 0)
			driver.SetLed(LED_EYE_RIGHT, 10, 10, 10)
			time.Sleep(1 * time.Second)
			driver.SetLed(LED_EYE_RIGHT, 0, 0, 0)
			driver.SetLed(LED_BLINKER_LEFT, 10, 10, 10)
			time.Sleep(1 * time.Second)
			driver.SetLed(LED_BLINKER_LEFT, 0, 0, 0)
			driver.SetLed(LED_BLINKER_RIGHT, 10, 10, 10)
			time.Sleep(1 * time.Second)
			driver.SetLed(LED_BLINKER_RIGHT, 0, 0, 0)
			driver.SetLed(LED_WIFI, 10, 10, 10)
			time.Sleep(1 * time.Second)
			driver.SetLed(LED_WIFI, 0, 0, 0)
		})

		gobot.Every(10* time.Second, func() {
			driver.SetMotorPower(MOTOR_LEFT, 50)
			driver.SetMotorPower(MOTOR_RIGHT, 50)
			time.Sleep(3 * time.Second)

			driver.SetMotorPower(MOTOR_LEFT, 0)
			driver.SetMotorPower(MOTOR_RIGHT, 0)
		})
	}

	robot := gobot.NewRobot("GoPiGoBot",
		[]gobot.Connection{adapter},
		[]gobot.Device{driver},
		work)
	err = robot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
