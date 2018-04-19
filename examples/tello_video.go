// +build example
//
// Do not build by default.

/*
You must have ffmpeg and ffplay installed in order to run this code. it will connect to the Tello
and then open a window using ffplay showing the streaming video.

How to run

	go run examples/tello_video.go
*/

package main

import (
	"fmt"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8890")

	work := func() {
		ffplay := exec.Command("ffplay", "-fast", "-i", "pipe:0")
		ffplayIn, _ := ffplay.StdinPipe()
		if err := ffplay.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(4)
			gobot.Every(250*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffplayIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
