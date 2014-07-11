package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/leap"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	servo1 := gpio.NewServoDriver(firmataAdaptor, "servo", "5")
	servo2 := gpio.NewServoDriver(firmataAdaptor, "servo", "3")

	leapAdaptor := leap.NewLeapMotionAdaptor("leap", "127.0.0.1:6437")
	leapDriver := leap.NewLeapMotionDriver(leapAdaptor, "leap")

	work := func() {
		x := 90.0
		z := 90.0
		gobot.On(leapDriver.Event("message"), func(data interface{}) {
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

	gbot.AddRobot(robot)

	gbot.Start()
}
