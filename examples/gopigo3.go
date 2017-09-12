// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	g "gobot.io/x/gobot/platforms/gopigo3"
)

func main() {
	gopigo3Adaptor := g.NewAdaptor()
	gopigo3 := g.NewGoPiGo3Driver(gopigo3Adaptor)

	work := func() {
		on := uint8(0xFF)
		gobot.Every(1000*time.Millisecond, func() {
			err := gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, on)
			if err != nil {
				fmt.Println(err)
			}
			err = gopigo3.SetLED(g.LED_EYE_LEFT, ^on, 0x00, 0x00)
			if err != nil {
				fmt.Println(err)
			}
			on = ^on
		})
	}

	robot := gobot.NewRobot("gopigo3",
		[]gobot.Connection{gopigo3Adaptor},
		[]gobot.Device{gopigo3},
		work,
	)

	robot.Start()
}
