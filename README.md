[![Gobot](https://raw.github.com/hybridgroup/gobot/gh-pages/images/elements/logo.png)](http://gobot.io/)

http://gobot.io/

Gobot is a set of libraries for robotics, physical computing, and the Internet of Things, using the Go programming language (http://golang.org/)

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the same time.

Want to use Ruby or Javascript on robots? Check out our sister projects Artoo (http://artoo.io) and Cylon.js (http://cylonjs.com/)

[![Build Status](https://travis-ci.org/hybridgroup/gobot.png?branch=master)](https://travis-ci.org/hybridgroup/gobot) [![Coverage Status](https://coveralls.io/repos/hybridgroup/gobot/badge.png?branch=coveralls)](https://coveralls.io/r/hybridgroup/gobot?branch=coveralls)

## Examples

### Basic

#### Go with a Sphero

```go
package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot-sphero"
)

func main() {

  spheroAdaptor := new(gobotSphero.SpheroAdaptor)
  spheroAdaptor.Name = "Sphero"
  spheroAdaptor.Port = "/dev/rfcomm0"

  sphero := gobotSphero.NewSphero(spheroAdaptor)
  sphero.Name = "Sphero"

  work := func() {
    gobot.Every("2s", func() {
      sphero.Roll(100, uint16(gobot.Rand(360)))
    })
  }

  robot := gobot.Robot{
    Connections: []gobot.Connection{spheroAdaptor},
    Devices:     []gobot.Device{sphero},
    Work:        work,
  }

  robot.Start()
}
```
#### Go with a Blink

```go
package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot-firmata"
  "github.com/hybridgroup/gobot-gpio"
)

func main() {

  firmata := new(gobotFirmata.FirmataAdaptor)
  firmata.Name = "firmata"
  firmata.Port = "/dev/ttyACM0"

  led := gobotGPIO.NewLed(firmata)
  led.Name = "led"
  led.Pin = "13"

  work := func() {
    gobot.Every("1s", func() {
      led.Toggle()
    })
  }

  robot := gobot.Robot{
    Connections: []gobot.Connection{firmata},
    Devices:     []gobot.Device{led},
    Work:        work,
  }

  robot.Start()
}
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following robotics and physical computing platforms are currently supported:
  
  - [Ardrone](http://ardrone2.parrot.com/) <==> [Library](https://github.com/hybridgroup/gobot-ardrone)
  - [Arduino](http://www.arduino.cc/) <==> [Library](https://github.com/hybridgroup/gobot-firmata)
  - [Beaglebone Black](http://beagleboard.org/Products/BeagleBone+Black/) <=> [Library](https://github.com/hybridgroup/gobot-beaglebone)
  - [Joystick](http://en.wikipedia.org/wiki/Joystick) <=> [Library](https://github.com/hybridgroup/gobot-joystick)
  - [Digispark](http://digistump.com/products/1) <=> [Library](https://github.com/hybridgroup/gobot-digispark)
  - [Firmata](http://firmata.org/wiki/Main_Page) <=> [Library](https://github.com/hybridgroup/gobot-firmata)
  - [Leap Motion](https://www.leapmotion.com/) <=> [Library](https://github.com/hybridgroup/gobot-leap)
  - [OpenCV](http://opencv.org/) <=> [Library](https://github.com/hybridgroup/gobot-opencv)
  - [Spark](https://www.spark.io/) <=> [Library](https://github.com/hybridgroup/gobot-spark)
  - [Sphero](http://www.gosphero.com/) <=> [Library](https://github.com/hybridgroup/gobot-sphero)
  

Support for many devices that use General Purpose Input/Output (GPIO) have
a shared set of drivers provded using the cylon-gpio module:

  - [GPIO](https://en.wikipedia.org/wiki/General_Purpose_Input/Output) <=> [Drivers](https://github.com/hybridgroup/gobot-gpio)
    - Analog Sensor
    - Button
    - Digital Sensor
    - LED
    - Motor
    - Servo

Support for devices that use Inter-Integrated Circuit (I2C) have a shared set of
drivers provded using the gobot-i2c module:

  - [I2C](https://en.wikipedia.org/wiki/I%C2%B2C) <=> [Drivers](https://github.com/hybridgroup/gobot-i2c)
    - BlinkM
    - HMC6352
    - Wii Nunchuck Controller

More platforms and drivers are coming soon...

## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot`

Then install additional libraries for whatever hardware support you want to use from your robot. For example, `go get -u github.com/hybridgroup/gobot-sphero` to use Gobot with a Sphero.

## API:

Gobot includes a RESTful API to query the status of any robot running within a group, including the connection and device status, and execute device commands.

To activate the API, use the `Api` command like this:

```go 
  master := gobot.GobotMaster()
  gobot.Api(master)
```
To specify the api port run your Gobot program with the `PORT` environment variable
```
  $ PORT=8080 go run gobotProgram.go
```

In order to use the [robeaux](https://github.com/hybridgroup/robeaux) AngularJS interface with Gobot you simply clone the robeaux repo and place it in the directory of your Gobot program. The robeaux assets must be in a folder called `robeaux`.

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
