# OpenCV

OpenCV (Open Source Computer Vision Library) is an open source computer vision and machine learning software library. OpenCV was built to provide a common infrastructure for computer vision applications and to accelerate the use of machine perception in the commercial products. Being a BSD-licensed product, OpenCV makes it easy for businesses to utilize and modify the code.

For more info about OpenCV click [here](http://opencv.org/)

## How to Install

This package requires OpenCV version 3.4 be installed on your system, along with GoCV, which is the Go programming language wrapper used by Gobot. The best way is to follow the installation instructions on the GoCV website at [https://gocv.io](https://gocv.io).

### macOS

To install on macOS follow the instructions here:

https://gocv.io/getting-started/macos/

### Ubuntu

To install on Ubuntu follow the instructions here:

https://gocv.io/getting-started/linux/

### Windows

To install on Windows follow the instructions here:

https://gocv.io/getting-started/windows/


Now you can install the Gobot package itself with

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

When you run code that uses OpenCV, you must setup some environment variables first. The best way to do this, is to first run the `env.sh` script that comes with GoCV, like this:

```
source $GOPATH/src/gocv.io/x/gocv/env.sh
```

Once you have run this script you can use `go run` or `go build` on your Gobot code that uses OpenCV as you normally would.

Here is an example using the camera:

```go
package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gocv.io/x/gocv"
)

func main() {
	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver(0)

	work := func() {
		camera.On(opencv.Frame, func(data interface{}) {
			img := data.(gocv.Mat)
			window.ShowImage(img)
			window.WaitKey(1)
		})
	}

	robot := gobot.NewRobot("cameraBot",
		[]gobot.Device{window, camera},
		work,
	)

	robot.Start()
}
```
