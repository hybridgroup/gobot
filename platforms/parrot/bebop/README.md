# Bebop

The Bebop from Parrot is an inexpensive quadcopter that is controlled using WiFi. It includes a built-in front-facing HD video camera, as well as a second lower resolution bottom facing video camera.


## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/parrot/bebop
```

## How to Use
```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/bebop"
)

func main() {
	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
    drone.HullProtection(true)
		drone.TakeOff()
		gobot.On(drone.Event("flying"), func(data interface{}) {
			gobot.After(3*time.Second, func() {
				drone.Land()
			})
		})
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{bebopAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
```

## How to Connect

The Bebop is a WiFi device, so there is no additional work to establish a connection to a single drone. However, in order to connect to multiple drones, you need to perform some configuration steps on each drone via SSH.
