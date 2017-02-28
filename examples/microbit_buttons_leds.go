package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/microbit"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	buttons := microbit.NewButtonDriver(bleAdaptor)
	leds := microbit.NewLEDDriver(bleAdaptor)

	work := func() {
		buttons.On(microbit.ButtonA, func(data interface{}) {
			fmt.Println("a", data)
			if data.([]byte)[0] == 1 {
				leds.UpLeftArrow()
				return
			}

			leds.Blank()
		})

		buttons.On(microbit.ButtonB, func(data interface{}) {
			fmt.Println("b", data)
			if data.([]byte)[0] == 1 {
				leds.UpRightArrow()
				return
			}

			leds.Blank()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{buttons, leds},
		work,
	)

	robot.Start()
}
