//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/audio"
)

func main() {
	e := audio.NewAdaptor()
	laser := audio.NewDriver(e, "./examples/laser.mp3")

	work := func() {
		gobot.Every(2*time.Second, func() {
			laser.Play()
		})
	}

	robot := gobot.NewRobot("soundBot",
		[]gobot.Connection{e},
		[]gobot.Device{laser},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
