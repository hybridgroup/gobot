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
	led := gpio.NewLedDriver(firmataAdaptor, "myLed", "13")
	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}
	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("blinkBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{led}, work))
	gbot.Start()
}
