/*
	This example will connect to the Bebop and stream it's video to a webpage
	via ffserver. This requires you to have both ffmpeg and ffserver installed
	on your computer.

	In order to run this example you will first need to  start ffserver with:
		$ ffserver -f ff.conf

	then in a seperate terminal run this program:
		$ go run video.go

	You will then be able to view the video feed by navigatin to
	http://localhost:8090/bebop.mjpeg in a web browser. *NOTE* firefox works
	best for viewing the video feed.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/hybridgroup/gobot/platforms/bebop/client"
)

func main() {
	bebop := client.New()

	if err := bebop.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "http://localhost:8090/bebop.ffm")

	ffmpegErr, err := ffmpeg.StderrPipe()

	if err != nil {
		fmt.Println(err)
		return
	}

	ffmpegIn, err := ffmpeg.StdinPipe()

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ffmpeg.Start(); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			buf, err := ioutil.ReadAll(ffmpegErr)
			if err != nil {
				fmt.Println(err)
			}
			if len(buf) > 0 {
				fmt.Println(string(buf))
			}
		}
	}()

	go func() {
		for {
			if _, err := ffmpegIn.Write(<-bebop.Video()); err != nil {
				fmt.Println(err)
			}
		}
	}()

	<-time.After(99 * time.Second)
}
