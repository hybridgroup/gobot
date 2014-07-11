# Ardrone

This package  provides the Gobot adaptor and driver for the [Parrot Ardrone](http://ardrone2.parrot.com).

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

## Getting Started

## Installing
```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/ardrone
```
## Using
```go
package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
)

func main() {
	gbot := gobot.NewGobot()

	ardroneAdaptor := ardrone.NewArdroneAdaptor("Drone")
	drone := ardrone.NewArdroneDriver(ardroneAdaptor, "Drone")

	work := func() {
		drone.TakeOff()
		gobot.On(drone.Event("flying"), func(data interface{}) {
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{drone},
		work,
	)
	gbot.AddRobot(robot)

	gbot.Start()
}
```
