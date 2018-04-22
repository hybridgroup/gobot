// +build example
//
// Do not build by default.

/*
You must have ffmpeg and OpenCV installed in order to run this code. It will connect to the Tello
and then open a window using OpenCV showing the streaming video.

How to run

	go run examples/tello_opencv.go
*/

package main

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/opencv"
	"gocv.io/x/gocv"
)

const (
	frameSize = 960 * 720 * 3
)

func main() {
	drone := tello.NewDriver("8890")
	window := opencv.NewWindowDriver()

	work := func() {
		ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
			"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
		ffmpegIn, _ := ffmpeg.StdinPipe()
		ffmpegOut, _ := ffmpeg.StdoutPipe()
		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			for {
				buf := make([]byte, frameSize)
				if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
					fmt.Println(err)
					continue
				}

				img := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
				if img.Empty() {
					continue
				}
				window.ShowImage(img)
				window.WaitKey(1)
			}
		}()

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			drone.SetExposure(0)

			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone, window},
		work,
	)

	robot.Start()
}
