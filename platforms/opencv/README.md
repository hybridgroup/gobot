# OpenCV

This repository contains the Gobot drivers for opencv.

## Getting Started

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
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/platforms/opencv
```

## Using
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
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			window.ShowImage(data.(*cv.IplImage))
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("cameraBot", []gobot.Connection{}, []gobot.Device{window, camera}, work))

	gbot.Start()
}
```
