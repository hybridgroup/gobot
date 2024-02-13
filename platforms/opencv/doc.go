/*
Package opencv contains the Gobot drivers for opencv.

Installing:

# This package requires `opencv` to be installed on your system

Then you can install the package with:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

	package main

	import (
		cv "gobot.io/x/go-opencv/opencv"
		"gobot.io/x/gobot/v2"
		"gobot.io/x/gobot/v2/platforms/opencv"
	)

	func main() {
		window := opencv.NewWindowDriver()
		camera := opencv.NewCameraDriver(0)

		work := func() {
			camera.On(camera.Event("frame"), func(data interface{}) {
				window.ShowImage(data.(*cv.IplImage))
			})
		}

		robot := gobot.NewRobot("cameraBot",
			[]gobot.Device{window, camera},
			work,
		)

		if err := robot.Start(); err != nil {
			panic(err)
		}
	}

For further information refer to opencv README:
https://github.com/hybridgroup/gobot/blob/release/platforms/opencv/README.md
*/
package opencv // import "gobot.io/x/gobot/v2/platforms/opencv"
