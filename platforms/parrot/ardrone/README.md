# Ardrone

The ARDrone from Parrot is an inexpensive quadcopter that is controlled using WiFi. It includes a built-in front-facing HD video camera, as well as a second lower resolution bottom facing video camera.

For more info about the ARDrone platform click [here](http://ardrone2.parrot.com/).

## How to Install
```
go get -d -u gobot.io/x/gobot/...
```

## How to Use
```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/ardrone"
)

func main() {
	ardroneAdaptor := ardrone.NewAdaptor("Drone")
	drone := ardrone.NewDriver(ardroneAdaptor, "Drone")

	work := func() {
		drone.On(drone.Event("flying"), func(data interface{}) {
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
		drone.TakeOff()
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{ardroneAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
```

## How to Connect

The ARDrone is a WiFi device, so there is no additional work to establish a connection to a single drone. However, in order to connect to multiple drones, you need to perform some configuration steps on each drone via SSH.
