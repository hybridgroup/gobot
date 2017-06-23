// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	pca9685 := i2c.NewPCA9685Driver(firmataAdaptor)
	servo := gpio.NewServoDriver(pca9685, "15")

	work := func() {
		pca9685.SetPWMFreq(60)

		for i := 10; i < 150; i += 10 {
			fmt.Println("Turning", i)
			servo.Move(uint8(i))
			time.Sleep(1 * time.Second)
		}

		for i := 150; i > 10; i -= 10 {
			fmt.Println("Turning", i)
			servo.Move(uint8(i))
			time.Sleep(1 * time.Second)
		}
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{pca9685, servo},
		work,
	)

	robot.Start()
}
