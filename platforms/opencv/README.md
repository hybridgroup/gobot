# OpenCV

OpenCV (Open Source Computer Vision Library) is an open source computer vision and machine learning software library.
OpenCV was built to provide a common infrastructure for computer vision applications and to accelerate the use of machine
perception in the commercial products. Being a BSD-licensed product, OpenCV makes it easy for businesses to utilize and
modify the code.

For more info about OpenCV click [here](http://opencv.org/)

## How to Install

This package requires OpenCV 3.4+ be installed on your system, along with GoCV, which is the Go programming language
wrapper used by Gobot. The best way is to follow the installation instructions on the GoCV website at [https://gocv.io](https://gocv.io).

The instructions should automatically install OpenCV 4+

### macOS

To install on macOS follow the instructions here:

<https://gocv.io/getting-started/macos/>

### Ubuntu

To install on Ubuntu follow the instructions here:

<https://gocv.io/getting-started/linux/>

### Windows

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

To install on Windows follow the instructions here:

<https://gocv.io/getting-started/windows/>

## How to Use

Here is an example using the camera:

```go
package main

import (
  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/opencv"
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

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

Build your application with build tag "gocv".
