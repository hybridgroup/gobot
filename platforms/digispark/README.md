# Gobot For Digispark

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This repository contains the Gobot adaptor for the [Digispark](http://www.kickstarter.com/projects/digistump/digispark-the-tiny-arduino-enabled-usb-dev-board) ATTiny-based USB development board with the [Little Wire](http://littlewire.cc/) protocol firmware installed.

Want to use Ruby on robots? Check out our sister project Artoo (http://artoo.io)

Want to use Javascript to power your robots? Check out our sister project Cylon.js (http://cylonjs.com/).

For more information about Gobot, check out our repo at https://github.com/hybridgroup/gobot

## Getting Started

Installing gobot-digispark requires the `libusb` package already be installed.

### OSX

To install libusb on OSX using Homebrew:

```
$ brew install libusb
```

### Ubuntu

To install libusb on linux:

```
$ sudo apt-get install libusb-dev
```

Now you can install the library with: `go get github.com/hybridgroup/gobot-digispark`

## Examples
```go
package main

import (
        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot-digispark"
        "github.com/hybridgroup/gobot-gpio"
)

func main() {

        digispark := new(gobotDigispark.DigisparkAdaptor)
        digispark.Name = "Digispark"

        led := gobotGPIO.NewLed(digispark)
        led.Name = "led"
        led.Pin = "0"

        work := func() {
                gobot.Every("0.5s", func() {
                        led.Toggle()
                })
        }

        robot := gobot.Robot{
                Connections: []gobot.Connection{digispark},
                Devices:     []gobot.Device{led},
                Work:        work,
        }

        robot.Start()
}
```
## Connecting to Digispark

If your Digispark (http://www.kickstarter.com/projects/digistump/digispark-the-tiny-arduino-enabled-usb-dev-board) ATTiny-based USB development board already has the Little Wire (http://littlewire.cc/) protocol firmware installed, you can connect right away with Gobot. 

Otherwise, for instructions on how to install Little Wire on a Digispark check out http://digistump.com/board/index.php/topic,160.0.html

### OSX
```
Important: 2012 MBP The USB ports on the 2012 MBPs (Retina and non) cause issues due to their USB3 controllers,
currently the best work around is to use a cheap USB hub (non USB3) - we are working on future solutions. The hub 
on a Cinema display will work as well.
```

The main steps are:
- Plug in the Digispark to the USB port
- Connect to the device via Gobot

First plug the Digispark into your computer via the USB port. Then... (directions go here)

### Ubuntu

The main steps are:
- Add a udev rule to allow access to the Digispark device
- Plug in the Digispark to the USB port
- Connect to the device via Gobot

First, you must add a udev rule, so that Gobot can communicate with the USB device. Ubuntu and other modern Linux distibutions use udev to manage device files when USB devices are added and removed. By default, udev will create a device with read-only permission which will not allow to you download code. You must place the udev rules below into a file named /etc/udev/rules.d/49-micronucleus.rules.

```
# UDEV Rules for Micronucleus boards including the Digispark.
# This file must be placed at:
#
# /etc/udev/rules.d/49-micronucleus.rules    (preferred location)
#   or
# /lib/udev/rules.d/49-micronucleus.rules    (req'd on some broken systems)
#
# After this file is copied, physically unplug and reconnect the board.
#
SUBSYSTEMS=="usb", ATTRS{idVendor}=="1781", ATTRS{idProduct}=="0c9f", MODE:="0666"
KERNEL=="ttyACM*", ATTRS{idVendor}=="1781", ATTRS{idProduct}=="0c9f", MODE:="0666", ENV{ID_MM_DEVICE_IGNORE}="1"

SUBSYSTEMS=="usb", ATTRS{idVendor}=="16d0", ATTRS{idProduct}=="0753", MODE:="0666"
KERNEL=="ttyACM*", ATTRS{idVendor}=="16d0", ATTRS{idProduct}=="0753", MODE:="0666", ENV{ID_MM_DEVICE_IGNORE}="1"
#
# If you share your linux system with other users, or just don't like the
# idea of write permission for everybody, you can replace MODE:="0666" with
# OWNER:="yourusername" to create the device owned by you, or with
# GROUP:="somegroupname" and mange access using standard unix groups.
```

Thanks to [@bluebie](https://github.com/Bluebie) for these instructions! (https://github.com/Bluebie/micronucleus-t85/wiki/Ubuntu-Linux)

Now plug the Digispark into your computer via the USB port.

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
* Add unit tests for any new or changed functionality & Lint and test your code using [Grunt](http://gruntjs.com/).
* All pull requests should be "fast forward"
  * If there are commits after yours use “git rebase -i <new_head_branch>”
  * If you have local changes you may need to use “git stash”
  * For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git

## License
Copyright (c) 2013-2014 The Hybrid Group. Licensed under the Apache 2.0 license.
