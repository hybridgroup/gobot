# Bebop

The Parrot Bebop and Parrot Bebop 2 are inexpensive quadcopters that can be controlled using their built-in API commands via a WiFi connection. They include a built-in front-facing HD video camera, as well as a second lower resolution bottom-facing video camera.

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

By default, the Parrot Bebop is a WiFi access point, so there is no additional work to establish a connection to a single drone, you just connect to it.

Once connected, if you want to stream drone video you can either:

	- Use the RTP protocol with an external player such as mplayer or VLC.
	- Grab the video frames from the drone's data frames, and work with them directly.
