package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/pebble"
)

func main() {
	master := gobot.NewGobot()
	api.NewAPI(master).Start()

	pebbleAdaptor := pebble.NewPebbleAdaptor("pebble")
	pebbleDriver := pebble.NewPebbleDriver(pebbleAdaptor, "pebble")

	work := func() {
		pebbleDriver.SendNotification("Hello Pebble!")
		gobot.On(pebbleDriver.Events["button"], func(data interface{}) {
			fmt.Println("Button pushed: " + data.(string))
		})

		gobot.On(pebbleDriver.Events["tap"], func(data interface{}) {
			fmt.Println("Tap event detected")
		})
	}

	robot := gobot.NewRobot("pebble", []gobot.Connection{pebbleAdaptor}, []gobot.Device{pebbleDriver}, work)

	master.Robots = append(master.Robots, robot)
	master.Start()
}
