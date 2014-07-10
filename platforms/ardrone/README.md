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
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/ardrone"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	adaptor := ardrone.NewArdroneAdaptor("Drone")
	drone := ardrone.NewArdroneDriver(adaptor, "Drone")

	work := func() {
		drone.TakeOff()
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("drone", []gobot.Connection{adaptor}, []gobot.Device{drone}, work))

	gbot.Start()
}
```