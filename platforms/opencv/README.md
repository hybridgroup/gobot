# OpenCV

OpenCV (Open Source Computer Vision Library) is an open source computer vision and machine learning software library. OpenCV was built to provide a common infrastructure for computer vision applications and to accelerate the use of machine perception in the commercial products. Being a BSD-licensed product, OpenCV makes it easy for businesses to utilize and modify the code.

For more info about OpenCV click [here](http://opencv.org/)

## How to Install

This package requires OpenCV version 2.4 to be installed on your system. Please note that it is not compatible with OpenCV 3.x at this time.

### OSX

To install OpenCV on OSX using Homebrew:

```
$ brew tap homebrew/science && brew install opencv
```

### Ubuntu

To install OpenCV on Ubuntu 14.04:

```
$ sudo apt-get install libopencv-dev
```

Or, follow the official [OpenCV installation guide](http://docs.opencv.org/doc/tutorials/introduction/linux_install/linux_install.html)

### Windows

Follow the official [OpenCV installation guide](http://docs.opencv.org/doc/tutorials/introduction/windows_install/windows_install.html#windows-installation)


Now you can install the package with
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/opencv
```

## How to Use

Example using the camera.

```go
package main

import (
	cv "github.com/lazywei/go-opencv/opencv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
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

	robot.Start()
}
```
