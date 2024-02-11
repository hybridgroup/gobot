//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	width := 128
	height := 32
	r := raspi.NewAdaptor()
	oled := i2c.NewSSD1306Driver(r, i2c.WithSSD1306DisplayWidth(width), i2c.WithSSD1306DisplayHeight(height))

	stage := false

	work := func() {
		gobot.Every(1*time.Second, func() {
			oled.Clear()
			if stage {
				for x := 0; x < width; x += 5 {
					for y := 0; y < height; y++ {
						oled.Set(x, y, 1)
					}
				}
			}
			stage = !stage
			if err := oled.Display(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("ssd1306Robot",
		[]gobot.Connection{r},
		[]gobot.Device{oled},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
