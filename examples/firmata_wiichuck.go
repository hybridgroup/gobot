// +build example
//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_wiichuck.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	wiichuck := i2c.NewWiichuckDriver(firmataAdaptor)

	work := func() {
		wiichuck.On(wiichuck.Event("joystick"), func(data interface{}) {
			fmt.Println("joystick", data)
		})

		wiichuck.On(wiichuck.Event("c"), func(data interface{}) {
			fmt.Println("c")
		})

		wiichuck.On(wiichuck.Event("z"), func(data interface{}) {
			fmt.Println("z")
		})

		wiichuck.On(wiichuck.Event("error"), func(data interface{}) {
			fmt.Println("Wiichuck error:", data)
		})
	}

	robot := gobot.NewRobot("chuck",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{wiichuck},
		work,
	)

	robot.Start()
}
