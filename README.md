[![Gobot](https://raw.githubusercontent.com/hybridgroup/gobot-site/master/source/images/elements/gobot-logo-small.png)](http://gobot.io/)

Gobot (http://gobot.io/) is a framework using the Go programming language (http://golang.org/) for robotics, physical computing, and the Internet of Things.

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the same time.

Want to use Javascript robotics? Check out our sister project Cylon.js (http://cylonjs.com/)

Want to use Ruby on robots? Check out our sister project Artoo (http://artoo.io)

[![GoDoc](https://godoc.org/github.com/hybridgroup/gobot?status.svg)](https://godoc.org/github.com/hybridgroup/gobot)
[![Build Status](https://travis-ci.org/hybridgroup/gobot.png?branch=master)](https://travis-ci.org/hybridgroup/gobot) [![Coverage Status](https://coveralls.io/repos/hybridgroup/gobot/badge.png?branch=master)](https://coveralls.io/r/hybridgroup/gobot?branch=master)

## Getting Started

Get the Gobot source with: `go get -d -u github.com/hybridgroup/gobot/...`

## Examples

#### Gobot with Arduino

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

#### Gobot with Sphero

```go
package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	adaptor := sphero.NewSpheroAdaptor("sphero", "/dev/rfcomm0")
	driver := sphero.NewSpheroDriver(adaptor, "sphero")

	work := func() {
		gobot.Every(3*time.Second, func() {
			driver.Roll(30, uint16(gobot.Rand(360)))
		})
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following robotics and physical computing platforms are currently supported:

- [Ardrone](http://ardrone2.parrot.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/ardrone)
- [Arduino](http://www.arduino.cc/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/firmata)
- [Beaglebone Black](http://beagleboard.org/Products/BeagleBone+Black/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/beaglebone)
- [Bebop](http://www.parrot.com/usa/products/bebop-drone/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/bebop)
- [Digispark](http://digistump.com/products/1) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/digispark)
- [Intel Edison](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/intel-iot/edison)
- [Joystick](http://en.wikipedia.org/wiki/Joystick) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/joystick)
- [Keyboard](https://en.wikipedia.org/wiki/Computer_keyboard) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/keyboard)
- [Leap Motion](https://www.leapmotion.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/leapmotion)
- [MavLink](http://qgroundcontrol.org/mavlink/start) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/mavlinky)
- [MQTT](http://mqtt.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/mqtt)
- [Neurosky](http://neurosky.com/products-markets/eeg-biosensors/hardware/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/neurosky)
- [OpenCV](http://opencv.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/opencv)
- [Pebble](https://www.getpebble.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/pebble)
- [Raspberry Pi](http://www.raspberrypi.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/raspi)
- [Spark](https://www.spark.io/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/spark)
- [Sphero](http://www.gosphero.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/sphero)

Support for many devices that use General Purpose Input/Output (GPIO) have
a shared set of drivers provided using the `gobot/platforms/gpio` package:

- [GPIO](https://en.wikipedia.org/wiki/General_Purpose_Input/Output) <=> [Drivers](https://github.com/hybridgroup/gobot/tree/master/platforms/gpio)
	- Analog Sensor
	- Button
	- Direct Pin
	- Digital Sensor
	- Direct Pin
	- Grove Button
	- Grove Buzzer
	- Grove LED
	- Grove Light Sensor
	- Grove Piezo Vibration Sensor
	- Grove Relay
	- Grove Rotary Dial
	- Grove Sound Sensor
	- Grove Touch Sensor
	- LED
	- Makey Button
	- Motor
	- Servo

Support for devices that use Inter-Integrated Circuit (I2C) have a shared set of
drivers provided using the `gobot/platforms/i2c` package:

- [I2C](https://en.wikipedia.org/wiki/I%C2%B2C) <=> [Drivers](https://github.com/hybridgroup/gobot/tree/master/platforms/i2c)
	- BlinkM
	- Grove Digital Accelerometer
	- Grove RGB LCD
	- HMC6352
	- JHD1313M1
	- LIDAR-Lite
	- MMA7660
	- MPL1150A2
	- MPU6050
	- Wii Nunchuck Controller

More platforms and drivers are coming soon...

## API:

Gobot includes a RESTful API to query the status of any robot running within a group, including the connection and device status, and execute device commands.

To activate the API, require the `github.com/hybridgroup/gobot/api` package and instantiate the `API` like this:

```go
  gbot := gobot.NewGobot()
  api.NewAPI(gbot).Start()
```

You can also specify the api host and port, and turn on authentication:
```go
  gbot := gobot.NewGobot()
  server := api.NewAPI(gbot)
  server.Port = "4000"
  server.Username = "Gort"
  server.Password = "klaatu"
  server.Start()
```

You may access the [robeaux](https://github.com/hybridgroup/robeaux) React.js interface with Gobot by navigating to `http://localhost:3000/index.html`.

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Need help?
* Join our mailing list: https://groups.google.com/forum/#!forum/gobotio
* IRC: `#gobotio @ irc.freenode.net`
* Issues: https://github.com/hybridgroup/gobot/issues
* twitter: [@gobotio](https://twitter.com/gobotio)

## Contributing
For our contribution guidelines, please go to [https://github.com/hybridgroup/gobot/blob/master/CONTRIBUTING.md
](https://github.com/hybridgroup/gobot/blob/master/CONTRIBUTING.md
).

## License
Copyright (c) 2013-2015 The Hybrid Group. Licensed under the Apache 2.0 license.
