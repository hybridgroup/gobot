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
		gobot.On(pebbleDriver.Events["accel"], func(data interface{}) {
			fmt.Println(data.(string))
		})
	}

	robot := gobot.NewRobot("pebble", []gobot.Connection{pebbleAdaptor}, []gobot.Device{pebbleDriver}, work)

	master.Robots = append(master.Robots, robot)
	master.Start()
}
