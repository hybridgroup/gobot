// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	raspiAdaptor := raspi.NewAdaptor()
	gpg3 := gopigo3.NewDriver(raspiAdaptor)

	work := func() {
		on := uint8(0xFF)
		gobot.Every(1000*time.Millisecond, func() {
			err := gpg3.SetLED(gopigo3.LED_EYE_RIGHT, 0x00, 0x00, on)
			if err != nil {
				fmt.Println(err)
			}
			err = gpg3.SetLED(gopigo3.LED_EYE_LEFT, ^on, 0x00, 0x00)
			if err != nil {
				fmt.Println(err)
			}
			on = ^on
		})
	}

	robot := gobot.NewRobot("gopigo3eyeleds",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gpg3},
		work,
	)

	robot.Start()
}
