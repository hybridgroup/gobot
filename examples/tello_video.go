//go:build example
// +build example

//
// Do not build by default.

/*
You must have MPlayer (https://mplayerhq.hu) installed in order to run this code. it will connect to the Tello
and then open a window using MPlayer showing the streaming video.

How to run

	go run examples/tello_video.go
*/

package main

import (
	"fmt"
	"os/exec"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8890")

	mplayer := exec.Command("mplayer", "-fps", "60", "-")
	mplayerIn, _ := mplayer.StdinPipe()
	if err := mplayer.Start(); err != nil {
		fmt.Println(err)
		return
	}

	work := func() {
		_ = drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			if err := drone.StartVideo(); err != nil {
				fmt.Println(err)
			}
			if err := drone.SetVideoEncoderRate(tello.VideoBitRateAuto); err != nil {
				fmt.Println(err)
			}
			gobot.Every(100*time.Millisecond, func() {
				if err := drone.StartVideo(); err != nil {
					fmt.Println(err)
				}
			})
		})

		_ = drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := mplayerIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
