//go:build example
// +build example

//
// Do not build by default.

package main

//
// before to run
// 1. check Jetson io pwm pin configure `sudo /opt/nvidia/jetson-io/jetson-io.py`
// 2. if end pin configure, reboot Jetson nano.
// 3. run gobot
import (
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/jetson"
)

func main() {
	jetsonAdaptor := jetson.NewAdaptor()
	servo := gpio.NewServoDriver(jetsonAdaptor, "32")

	counter := 0
	flg := true
	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			log.Println("Turning", counter)
			servo.Move(uint8(counter))
			if counter == 140 {
				flg = false
			} else if counter == 30 {
				flg = true
			}

			if flg {
				counter = counter + 1
			} else {
				counter = counter - 1
			}
		})
	}

	robot := gobot.NewRobot("Jetsonservo",
		[]gobot.Connection{jetsonAdaptor},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()
}
