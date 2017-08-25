// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {

	r := raspi.NewAdaptor()
	oled := i2c.NewSSD1306Driver(r)

	stage := false

	work := func() {

		gobot.Every(1*time.Second, func() {
			oled.Clear()
			if stage {
				for x := 0; x < oled.Buffer.Width; x += 5 {
					for y := 0; y < oled.Buffer.Height; y++ {
						oled.Set(x, y, 1)
					}
				}
			}
			stage = !stage
			oled.Display()
		})
	}

	robot := gobot.NewRobot("ssd1306Robot",
		[]gobot.Connection{r},
		[]gobot.Device{oled},
		work,
	)

	robot.Start()
}
