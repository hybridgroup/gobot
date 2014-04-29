# Gobot for opencv

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This repository contains the Gobot adaptor for opencv.

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

## Installing
```
go get github.com/hybridgroup/gobot-opencv
```

## Using
```go
package main

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-opencv"
	"path"
	"runtime"
)

func main() {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	opencv := new(gobotOpencv.Opencv)
	opencv.Name = "opencv"

	window := gobotOpencv.NewWindow(opencv)
	window.Name = "window"

	camera := gobotOpencv.NewCamera(opencv)
	camera.Name = "camera"

	work := func() {
		var image *cv.IplImage
		gobot.On(camera.Events["Frame"], func(data interface{}) {
			image = data.(*cv.IplImage)
		})

		go func() {
			for {
				if image != nil {
					i := image.Clone()
					faces := gobotOpencv.DetectFaces(cascade, i)
					i = gobotOpencv.DrawRectangles(i, faces, 0, 255, 0, 5)
					window.ShowImage(i)
				}
			}
		}()
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{opencv},
		Devices:     []gobot.Device{window, camera},
		Work:        work,
	}

	robot.Start()
}
```
## Installing OpenCV and Connecting

In order to use OpenCV you first need to install it and make sure it is working correctly on your computer. You can follow the tutorials in the OpenCV site to install it in your particular OS:

[How to install OpenCV](http://docs.opencv.org/doc/tutorials/introduction/table_of_content_introduction/table_of_content_introduction.html#table-of-content-introduction)

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing

* All patches must be provided under the Apache 2.0 License
* Please use the -s option in git to "sign off" that the commit is your work and you are providing it under the Apache 2.0 License
* Submit a Github Pull Request to the appropriate branch and ideally discuss the changes with us in IRC.
* We will look at the patch, test it out, and give you feedback.
* Avoid doing minor whitespace changes, renamings, etc. along with merged content. These will be done by the maintainers from time to time but they can complicate merges and should be done seperately.
* Take care to maintain the existing coding style.
* Add unit tests for any new or changed functionality.
* All pull requests should be "fast forward"
  * If there are commits after yours use “git rebase -i <new_head_branch>”
  * If you have local changes you may need to use “git stash”
  * For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git


## License
Copyright (c) 2013-2014 The Hybrid Group. Licensed under the Apache 2.0 license.
