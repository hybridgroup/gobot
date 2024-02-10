//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	raspiAdaptor := raspi.NewAdaptor()
	gpg3 := gopigo3.NewDriver(raspiAdaptor)
	screen := i2c.NewGroveLcdDriver(raspiAdaptor)

	work := func() {
		manufacturerName, _ := gpg3.GetManufacturerName()
		boardName, _ := gpg3.GetBoardName()
		hardwareVersion, _ := gpg3.GetHardwareVersion()
		if err := screen.Write(manufacturerName[0:15]); err != nil {
			fmt.Println(err)
		}
		if err := screen.SetPosition(16); err != nil {
			fmt.Println(err)
		}
		if err := screen.Write(boardName + " " + hardwareVersion); err != nil {
			fmt.Println(err)
		}
		if err := screen.SetRGB(0, 0, 255); err != nil {
			fmt.Println(err)
		}
		if err := screen.Home(); err != nil {
			fmt.Println(err)
		}
	}

	robot := gobot.NewRobot("gopigo3lcd",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gpg3, screen},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
