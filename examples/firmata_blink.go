package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	robot := gbot.Robots().Add(gobot.NewRobot("bot"))
	adaptor := robot.Connections().Add(firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")).(gpio.PwmDigitalWriter)
	driver := robot.Devices().Add(gpio.NewLedDriver("myLed", adaptor, "13")).(*gpio.LedDriver)
	robot.Work = func() {
		gobot.Every(1*time.Second, func() {
			driver.Toggle()
		})
	}
	gbot.Start()
}
