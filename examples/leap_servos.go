// +build example
//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/leap"
)

// Video: https://www.youtube.com/watch?v=ayNMyUfdAqc
func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/tty.usbmodem1451")
	servo1 := gpio.NewServoDriver(firmataAdaptor, "3")
	servo2 := gpio.NewServoDriver(firmataAdaptor, "4")
	servo3 := gpio.NewServoDriver(firmataAdaptor, "5")
	servo4 := gpio.NewServoDriver(firmataAdaptor, "6")
	servo5 := gpio.NewServoDriver(firmataAdaptor, "7")

	leapMotionAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	l := leap.NewDriver(leapMotionAdaptor)

	work := func() {
		fist := false
		l.On(leap.MessageEvent, func(data interface{}) {
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

	robot.Start()
}
