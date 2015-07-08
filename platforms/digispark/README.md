# Digispark

The Digispark is an Attiny85 based microcontroller development board similar to the Arduino line, only cheaper, smaller, and a bit less powerful. With a whole host of shields to extend its functionality and the ability to use the familiar Arduino IDE the Digispark is a great way to jump into electronics, or perfect for when an Arduino is too big or too much.

This package provides the Gobot adaptor for the [Digispark](http://digistump.com/products/1) ATTiny-based USB development board with the [Little Wire](http://littlewire.cc/) protocol firmware installed.

## How to Install

This package requires `libusb`.

### OSX

To install `libusb` on OSX using Homebrew:

```
$ brew install libusb && brew install libusb-compat
```

### Ubuntu

To install libusb on linux:

```
$ sudo apt-get install libusb-dev
```

Now you can install the package with

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/digispark
```

## How to Use

```go
package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/digispark"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	digisparkAdaptor := digispark.NewDigisparkAdaptor("Digispark")
	led := gpio.NewLedDriver(digisparkAdaptor, "led", "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{digisparkAdaptor},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
```
## How to Connect

If your Digispark already has the Little Wire protocol firmware installed, you can connect right away with Gobot.

Otherwise, for instructions on how to install Little Wire on a Digispark check out http://digistump.com/board/index.php/topic,160.0.html

### OSX

```
Important: 2012 MBP The USB ports on the 2012 MBPs (Retina and non) cause issues due to their USB3 controllers,
currently the best work around is to use a cheap USB hub (non USB3) - we are working on future solutions. The hub on a Cinema display will work as well.
```
Plug the Digispark into your computer via the USB port and you're ready to go!

### Ubuntu

Ubuntu requires a few extra steps to set up the digispark for communication with Gobot:
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
