package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/audio"
)

func main() {
	gbot := gobot.NewGobot()

	e := audio.NewAudioAdaptor("sound")
	laser := audio.NewAudioDriver(e, "laser", "./examples/laser.mp3")

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

	gbot.AddRobot(robot)

	gbot.Start()
}
