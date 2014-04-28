package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/firmata"
	"github.com/hybridgroup/gobot/gpio"
	"github.com/hybridgroup/gobot/leap"
)

func main() {
	firmataAdaptor := firmata.NewFirmataAdaptor()
	firmataAdaptor.Name = "firmata"
	firmataAdaptor.Port = "/dev/ttyACM0"

	servo1 := gpio.NewServoDriver(firmataAdaptor)
	servo1.Name = "servo"
	servo1.Pin = "5"

	servo2 := gpio.NewServoDriver(firmataAdaptor)
	servo2.Name = "servo"
	servo2.Pin = "3"

	leapAdaptor := leap.NewLeapMotionAdaptor()
	leapAdaptor.Name = "leap"
	leapAdaptor.Port = "127.0.0.1:6437"

	leapDriver := leap.NewLeapMotionDriver(leapAdaptor)
	leapDriver.Name = "leap"

	work := func() {
		x := 90.0
		z := 90.0
		gobot.On(leapDriver.Events["Message"], func(data interface{}) {
			if len(data.(leap.Frame).Hands) > 0 {
				hand := data.(leap.Frame).Hands[0]
				x = gobot.ToScale(gobot.FromScale(hand.X(), -300, 300), 30, 150)
				z = gobot.ToScale(gobot.FromScale(hand.Z(), -300, 300), 30, 150)
			}
		})
		gobot.Every("0.01s", func() {
			servo1.Move(uint8(x))
			servo2.Move(uint8(z))
			fmt.Println("Current Angle: ", servo1.CurrentAngle, ",", servo2.CurrentAngle)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmataAdaptor, leapAdaptor},
		Devices:     []gobot.Device{servo1, servo2, leapDriver},
		Work:        work,
	}

	robot.Start()
}
