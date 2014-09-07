package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/leap"
)

// Video: https://www.youtube.com/watch?v=ayNMyUfdAqc
func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/tty.usbmodem1451")
	servo1 := gpio.NewServoDriver(firmataAdaptor, "servo", "3")
	servo2 := gpio.NewServoDriver(firmataAdaptor, "servo", "4")
	servo3 := gpio.NewServoDriver(firmataAdaptor, "servo", "5")
	servo4 := gpio.NewServoDriver(firmataAdaptor, "servo", "6")
	servo5 := gpio.NewServoDriver(firmataAdaptor, "servo", "7")

	leapMotionAdaptor := leap.NewLeapMotionAdaptor("leap", "127.0.0.1:6437")
	l := leap.NewLeapMotionDriver(leapMotionAdaptor, "leap")

	work := func() {
		fist := false
		gobot.On(l.Event("message"), func(data interface{}) {
			handIsOpen := len(data.(leap.Frame).Pointables) > 0
			if handIsOpen && fist {
				servo1.Move(0)
				servo2.Move(0)
				servo3.Move(0)
				servo4.Move(0)
				servo5.Move(0)
				fist = false
			} else if !handIsOpen && !fist {
				servo1.Move(120)
				servo2.Move(120)
				servo3.Move(120)
				servo4.Move(120)
				servo5.Move(120)
				fist = true
			}
		})
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{firmataAdaptor, leapMotionAdaptor},
		[]gobot.Device{servo1, servo2, servo3, servo4, servo5, l},
		work,
	)

	gbot.AddRobot(robot)
	gbot.Start()
}
