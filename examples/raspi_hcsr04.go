//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio/hcsr04"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	hcsr04 := hcsr04.NewHCSR04(r, "11", "13")

	work := func() {
		hsrc04.StartDistanceMonitor()

		gobot.Every(1*time.Second, func() {
			fmt.Println(hsrc04.GetDistance())
			time.Sleep(hcsr04.MonitorUpdate)
		})
	}

	robot := gobot.NewRobot("distanceBot",
		[]gobot.Connection{r},
		[]gobot.Device{hsrc04},
		work,
	)

	robot.Start()
}
