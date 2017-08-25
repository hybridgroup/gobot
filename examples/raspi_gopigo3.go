// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	g := spi.NewGoPiGo3Driver(r)
	g.Start()
	work := func() {
		g.SetLED(spi.LED_EYE_LEFT, 255, 0, 0)
		g.SetLED(spi.LED_EYE_RIGHT, 0, 255, 0)
		time.Sleep(1 * time.Second)
		g.SetLED(spi.LED_EYE_LEFT, 0, 255, 0)
		g.SetLED(spi.LED_EYE_RIGHT, 255, 0, 0)
	}
	robot := gobot.NewRobot("gopigo3-bot",
		[]gobot.Connection{r},
		[]gobot.Device{g},
		work,
	)
	robot.Start()
}
