package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")
	pin := gpio.NewDirectPinDriver(firmataAdaptor, "pin", "13")
	work := func() {
		level := byte(1)

		gobot.Every(1*time.Second, func() {
			pin.DigitalWrite(level)
			if level == 1 {
				level = 0
			} else {
				level = 1
			}
		})
	}
	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("pinBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{pin}, work))
	gbot.Start()
}
