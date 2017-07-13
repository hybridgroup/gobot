// +build example
//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_cat_toy.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/leap"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	servo1 := gpio.NewServoDriver(firmataAdaptor, "5")
	servo2 := gpio.NewServoDriver(firmataAdaptor, "3")

	leapAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	leapDriver := leap.NewDriver(leapAdaptor)

	work := func() {
		x := 90.0
		z := 90.0
		leapDriver.On(leap.MessageEvent, func(data interface{}) {
			if len(data.(leap.Frame).Hands) > 0 {
				hand := data.(leap.Frame).Hands[0]
				x = gobot.ToScale(gobot.FromScale(hand.X(), -300, 300), 30, 150)
				z = gobot.ToScale(gobot.FromScale(hand.Z(), -300, 300), 30, 150)
			}
		})
		gobot.Every(10*time.Millisecond, func() {
			servo1.Move(uint8(x))
			servo2.Move(uint8(z))
			fmt.Println("Current Angle: ", servo1.CurrentAngle, ",", servo2.CurrentAngle)
		})
	}

	robot := gobot.NewRobot("pwmBot",
		[]gobot.Connection{firmataAdaptor, leapAdaptor},
		[]gobot.Device{servo1, servo2, leapDriver},
		work,
	)

	robot.Start()
}
