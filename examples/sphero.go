package main

import (
	"github.com/hybridgroup/gobot"
	//"github.com/hybridgroup/gobot-sphero"
)

func main() {

	spheroAdaptor := new(gobotSphero.SpheroAdaptor)
	spheroAdaptor.Name = "Sphero"
	spheroAdaptor.Port = "127.0.0.1:4560"

	sphero := gobotSphero.NewSphero(spheroAdaptor)
	sphero.Name = "Sphero"

	connections := []interface{}{
		spheroAdaptor,
	}
	devices := []interface{}{
		sphero,
	}

	work := func() {
		gobot.Every("2s", func() {
			sphero.Roll(100, uint16(gobot.Random(0, 360)))
		})
	}

	robot := gobot.Robot{
		Connections: connections,
		Devices:     devices,
		Work:        work,
	}

	robot.Start()
}
