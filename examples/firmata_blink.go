package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")

	led := gpio.NewLedDriver(firmataAdaptor, "myLed", "13")

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
		})
	}

	gbot.Robots = append(gbot.Robots, gobot.NewRobot("Jerry", []gobot.Connection{firmataAdaptor}, []gobot.Device{led}, work))
	gbot.Start()
}
