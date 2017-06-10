# Digispark

The Digispark is an Attiny85 based microcontroller development board similar to the Arduino line, only cheaper, smaller, and a bit less powerful. With a whole host of shields to extend its functionality and the ability to use the familiar Arduino IDE the Digispark is a great way to jump into electronics, or perfect for when an Arduino is too big or too much.

This package provides the Gobot adaptor for the [Digispark](http://digistump.com/products/1) ATTiny-based USB development board with the [Little Wire](http://littlewire.github.io/) protocol firmware installed.

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
go get -d -u gobot.io/x/gobot/...
```

## How to Use

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/digispark"
)

func main() {
	digisparkAdaptor := digispark.NewAdaptor()
	led := gpio.NewLedDriver(digisparkAdaptor, "0")

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

	robot.Start()
}
```

## How to Connect

If your Digispark already has the Little Wire protocol firmware installed, you can connect right away with Gobot.

Otherwise, you must first flash your Digispark with the Little Wire firmware.

The easiest way to flash your Digispark is to use Gort [https://gort.io](https://gort.io).

Download and install Gort, and then use the following commands:

Then, install the needed Digispark firmware.

```
gort digispark install
```

### OSX

**Important**: 2012 MBP The USB ports on the 2012 MBPs (Retina and non) cause issues due to their USB3 controllers,
currently the best work around is to use a cheap USB hub (non USB3).

Plug the Digispark into your computer via the USB port and run:

```
gort digispark upload littlewire
```

### Ubuntu

Ubuntu requires an extra one-time step to set up the Digispark for communication with Gobot. Run the following command:

```
gort digispark set-udev-rules
```

You might need to enter your administrative password. This steps adds a udev rule to allow access to the Digispark device.

Once this is done, you can upload Little Wire to your Digispark:

```
gort digispark upload littlewire
```

### Windows

We need instructions here, because it supposedly works.

### Manual instructions

For manual instructions on how to install Little Wire on a Digispark check out http://digistump.com/board/index.php/topic,160.0.html

Thanks to [@bluebie](https://github.com/Bluebie) for these instructions! (https://github.com/Bluebie/micronucleus-t85/wiki/Ubuntu-Linux)

Now plug the Digispark into your computer via the USB port.
