# OpenCV

OpenCV (Open Source Computer Vision Library) is an open source computer vision and machine learning software library. OpenCV was built to provide a common infrastructure for computer vision applications and to accelerate the use of machine perception in the commercial products. Being a BSD-licensed product, OpenCV makes it easy for businesses to utilize and modify the code.

For more info about OpenCV click [here](http://opencv.org/)

## How to Install

This package requires `opencv` to be installed on your system

### OSX

To install `opencv` on OSX using Homebrew:

```
$ brew tap homebrew/science && brew install opencv
```

### Ubuntu

Follow the official [OpenCV installation guide](http://docs.opencv.org/doc/tutorials/introduction/linux_install/linux_install.html)

### Windows

Follow the official [OpenCV installation guide](http://docs.opencv.org/doc/tutorials/introduction/windows_install/windows_install.html#windows-installation)


Now you can install the package with
```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/opencv
```

## How to Use

Example using the camera.

```go
package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/opencv"
)

func main() {
	gbot := gobot.NewGobot()

	window := opencv.NewWindowDriver("window")
	camera := opencv.NewCameraDriver("camera", 0)

	work := func() {
		gobot.On(camera.Event("frame"), func(data interface{}) {
			window.ShowImage(data.(*cv.IplImage))
		})
	}

	robot := gobot.NewRobot("cameraBot",
		[]gobot.Device{window, camera},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
```
