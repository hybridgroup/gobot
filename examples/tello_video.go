// +build example
//
// Do not build by default.

/*
You must have ffmpeg and ffplay installed in order to run this code.

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
		// launch ffplay and get it ready to display video
		ffplay := exec.Command("ffplay", "-i", "pipe:0")
		ffplayIn, _ := ffplay.StdinPipe()
		if err := ffplay.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			gobot.Every(1*time.Second, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			// if len(pkt) > 6 && pkt[0] == 0x00 && pkt[1] == 0x00 && pkt[2] == 0x00 && pkt[3] == 0x01 {
			// 	nalType := pkt[6] & 0x1f
			// 	//fmt.Println("nal type = ", nalType)
			// }

			//fmt.Printf("Writing %d bytes\n", len(pkt))
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
