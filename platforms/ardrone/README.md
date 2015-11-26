# Ardrone

The ARDrone from Parrot is an inexpensive quadcopter that is controlled using WiFi. It includes a built-in front-facing HD video camera, as well as a second lower resolution bottom facing video camera.

For more info about the ARDrone platform click [here](http://ardrone2.parrot.com/).

## How to Install
```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/ardrone
```
## How to Use
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
		gobot.On(drone.Event("flying"), func(data interface{}) {
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
	gbot.AddRobot(robot)

	gbot.Start()
}
```
## How to Connect

The ARDrone is a WiFi device, so there is no additional work to establish a connection to a single drone. However, in order to connect to multiple drones, you need to perform some configuration steps on each drone via SSH.
