// +build example
//
// Do not build by default.

/*
 How to run
 Pass the file name to use to save the raw H264 video from the drone as first param:

	go run examples/tello_video.go "/tmp/tello.h264"
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8890")

	work := func() {
		f, err := os.Create(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		drone.StartVideo()
		gobot.Every(1*time.Second, func() {
			drone.StartVideo()
		})

		drone.On(tello.EvtVideoFrame, func(data interface{}) {
			pkt := data.([]byte)
			if len(pkt) > 6 && pkt[0] == 0x00 && pkt[1] == 0x00 && pkt[2] == 0x00 && pkt[3] == 0x01 {
				fmt.Println("nal type = ", pkt[6]&0x1f)
			}

			fmt.Printf("Writing %d bytes\n", len(pkt))
			_, err := f.Write(pkt)
			if err != nil {
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
