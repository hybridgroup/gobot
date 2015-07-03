# Firmata

Arduino is an open-source electronics prototyping platform based on flexible, easy-to-use hardware and software. It's intended for artists, designers, hobbyists and anyone interested in creating interactive objects or environments.

This package provides the adaptor for microcontrollers such as Arduino that support the [Firmata](http://firmata.org/wiki/Main_Page) protocol

For more info about the arduino platform click [here](http://arduino.cc/).

## How to Install

```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/firmata
```

## How to Use

```go
package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "13")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
```

## How to Connect

### Upload the Firmata Firmware to the Arduino

This section assumes you're using an Arduino Uno or another compatible board, and a UNIX operating system (OS X or Linux). If you already have the Firmata sketch installed, you can skip straight to the examples.

### OS X

First plug the Arduino into your computer via the USB/serial port.
A dialog box will appear telling you that a new network interface has been detected.
Click "Network Preferences...", and when it opens, simply click "Apply".

Once plugged in, use [Gort](http://gort.io)'s `gort scan serial` command to find out your connection info and serial port address:

```
$ gort scan serial
```

Use the `gort arduino install` command to install `avrdude`, this will allow you to upload firmata to the arduino:

```
$ gort arduino install
```

Once the avrdude uploader is installed we upload the firmata protocol to the arduino, use the arduino serial port address found when you ran `gort scan serial`:

```
$ gort arduino upload firmata /dev/tty.usbmodem1421
```

Now you are ready to connect and communicate with the Arduino using serial port connection

### Ubuntu

First plug the Arduino into your computer via the USB/serial port.

Once plugged in, use [Gort](http://gort.io)'s `gort scan serial` command to find out your connection info and serial port address:

```
$ gort scan serial
```

Use the `gort arduino install` command to install `avrdude`, this will allow you to upload firmata to the arduino:

```
$ gort arduino install
```

Once the avrdude uploader is installed we upload the firmata protocol to the arduino, use the arduino serial port address found when you ran `gort scan serial`, or leave it blank to use the default address `ttyACM0`:

```
$ gort arduino upload firmata /dev/ttyACM0
```

Now you are ready to connect and communicate with the Arduino using serial port connection

## Hardware Support
The following firmata devices have been tested and are currently supported:

  - [Arduino uno r3](http://arduino.cc/en/Main/arduinoBoardUno)
  - [Teensy 3.0](http://www.pjrc.com/store/teensy3.html)

More devices are coming soon...
