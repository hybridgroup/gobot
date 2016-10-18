package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/digispark"
)

func main() {
	master := gobot.NewMaster()
	api.NewAPI(master).Start()

	digisparkAdaptor := digispark.NewAdaptor()
	led := gpio.NewLedDriver(digisparkAdaptor, "0")

	robot := gobot.NewRobot("digispark",
		[]gobot.Connection{digisparkAdaptor},
		[]gobot.Device{led},
	)

	master.AddRobot(robot)

	master.Start()
}
