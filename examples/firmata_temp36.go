package main

import (
	"time"

	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			val, err := firmataAdaptor.AnalogRead("0")
			if err != nil {
				fmt.Println(err)
				return
			}

			voltage := (float64(val) * 5) / 1024 // if using 3.3V replace 5 with 3.3
			tempC := (voltage - 0.5) * 100
			tempF := (tempC * 9 / 5) + 32

			fmt.Printf("%.2f°C\n", tempC)
			fmt.Printf("%.2f°F\n", tempF)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
