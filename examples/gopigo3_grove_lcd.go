// +build example
//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	raspiAdaptor := raspi.NewAdaptor()
	gpg3 := gopigo3.NewDriver(raspiAdaptor)
	screen := i2c.NewGroveLcdDriver(raspiAdaptor)

	work := func() {
		manufacturerName := ""
		boardName := ""
		hardwareVersion := ""
		manufacturerName, _ = gpg3.GetManufacturerName()
		boardName, _ = gpg3.GetBoardName()
		hardwareVersion, _ = gpg3.GetHardwareVersion()
		screen.Write(manufacturerName[0:15])
		screen.SetPosition(16)
		screen.Write(boardName + " " + hardwareVersion)
		screen.SetRGB(0, 0, 255)
		screen.Home()
	}

	robot := gobot.NewRobot("gopigo3lcd",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gpg3, screen},
		work,
	)

	robot.Start()
}
