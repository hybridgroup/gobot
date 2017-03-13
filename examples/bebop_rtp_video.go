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
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/bebop"
)

func main() {
	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		drone.VideoEnable(true)
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{bebopAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
