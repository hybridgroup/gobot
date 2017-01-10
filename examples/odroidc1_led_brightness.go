package main

import (
	//"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/dhart-alldigital/gobot/platforms/odroid/c1"
)

func main() {
	gbot := gobot.NewGobot()
	
	r := c1.NewODroidC1Adaptor("c1")
	led := gpio.NewLedDriver(r, "led", "33")

	work := func() {
		brightness := uint8(0)
		fadeAmount := uint8(20)

		gobot.Every(100*time.Millisecond, func() {
				//log.Printf("[odroidc1_led_brightness] Setting brightness: %v", brightness)
				led.Brightness(brightness)
				brightness = brightness + fadeAmount
				if brightness == 0 || brightness == 255 {
					fadeAmount = -fadeAmount
				}
		})
	}

	robot := gobot.NewRobot("pwmBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
