//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_pca9685.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
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
