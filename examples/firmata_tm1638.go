//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_tm1638.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])

	// Even thought the modules are connected among them, they are independent (not chained)
	modules := make([]*gpio.TM1638Driver, 4)
	modules[0] = gpio.NewTM1638Driver(firmataAdaptor, "9", "8", "7")
	modules[1] = gpio.NewTM1638Driver(firmataAdaptor, "9", "8", "6")
	modules[2] = gpio.NewTM1638Driver(firmataAdaptor, "9", "8", "5")
	modules[3] = gpio.NewTM1638Driver(firmataAdaptor, "9", "8", "4")

	var ledInt byte
	var color byte
	var offset int

	// Repeat and concat strings until long enough that with scroll it still shows text
	const showText = "  HELLO WORLD    -    gobot.io    -    TM1638  "
	text := showText
	for len(text)-len(showText) < len(modules)*8 {
		text += showText
	}

	work := func() {
		gobot.Every(400*time.Millisecond, func() {
			// Enable and change the color of the LEDs
			if err := modules[0].SetLED(color, ledInt); err != nil {
				fmt.Println(err)
			}
			if err := modules[1].SetLED(color, ledInt); err != nil {
				fmt.Println(err)
			}
			if err := modules[2].SetLED(color, ledInt); err != nil {
				fmt.Println(err)
			}
			if err := modules[3].SetLED(color, ledInt); err != nil {
				fmt.Println(err)
			}

			ledInt++
			if ledInt > 7 {
				ledInt = 0
				color++
				if color > 2 {
					color = 0
				}
			}

			// Scroll the text
			for i := 0; i < 4; i++ {
				if err := modules[i].SetDisplayText(text[offset+8*i : offset+8*i+8]); err != nil {
					fmt.Println(err)
				}
			}
			offset++
			if offset >= len(showText) {
				offset = 0
			}
		})
	}

	robot := gobot.NewRobot("tm1638Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{modules[0], modules[1], modules[2], modules[3]},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
