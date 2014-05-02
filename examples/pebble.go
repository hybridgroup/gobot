package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/pebble"
  "fmt"
)

func main() {
  pebbleAdaptor := new(gobotPebble.PebbleAdaptor)
  pebbleAdaptor.Name = "Pebble"

  pebble := gobotPebble.NewPebble(pebbleAdaptor)
  pebble.Name = "pebble"

  master := gobot.GobotMaster()
  api    := gobot.Api(master)
  api.Port = "8080"

  work := func() {
    gobot.On(pebble.Events["button"], func(data interface{}) {
      fmt.Println("Button pushed: " + data.(string))
    })
  }

  robot := gobot.Robot{
    Connections: []gobot.Connection{pebbleAdaptor},
    Devices:     []gobot.Device{pebble},
    Work:        work,
  }

  robot.Name = "pebble"

  master.Robots = append(master.Robots, &robot)
  master.Start()
}
