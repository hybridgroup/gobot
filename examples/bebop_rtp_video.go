//go:build example
// +build example

//
// Do not build by default.

/*
This example will connect to the Parrot Bebop and streams the drone video
via the RTP protocol.

In order to run this example you will first need to connect to the drone with:

	$ go run bebop_ps3_video.go

then in a separate terminal run this program:

	$ mplayer examples/bebop.sdp

You can view the video feed by navigating to
http://localhost:8090/bebop.mjpeg in a web browser.
*NOTE* firefox works best for viewing the video feed.
*/
package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/parrot/bebop"
)

func main() {
	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		if err := drone.VideoEnable(true); err != nil {
			fmt.Println(err)
		}
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{bebopAdaptor},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
