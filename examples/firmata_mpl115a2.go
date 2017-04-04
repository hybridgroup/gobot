// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	mpl115a2 := i2c.NewMPL115A2Driver(firmataAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			press, _ := mpl115a2.Pressure()
			fmt.Println("Pressure", press)

			temp, _ := mpl115a2.Temperature()
			fmt.Println("Temperature", temp)
		})
	}

	robot := gobot.NewRobot("mpl115a2Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{mpl115a2},
		work,
	)

	robot.Start()
}
