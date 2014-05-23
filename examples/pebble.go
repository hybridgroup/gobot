package main

import (
  "github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
  "github.com/hybridgroup/gobot/platforms/pebble"
  "fmt"
)

func main() {
  master := gobot.NewGobot()
  api.NewApi(master).Start()

  pebbleAdaptor := pebble.NewPebbleAdaptor("pebble")
  pebbleDriver  := pebble.NewPebbleDriver(pebbleAdaptor, "pebble")

  work := func() {
    gobot.On(pebbleDriver.Events["button"], func(data interface{}) {
      fmt.Println("Button pushed: " + data.(string))
    })
  }

  robot := gobot.NewRobot("pebble", []gobot.Connection{pebbleAdaptor}, []gobot.Device{pebbleDriver}, work)

  master.Robots = append(master.Robots, robot)
  master.Start()
}
