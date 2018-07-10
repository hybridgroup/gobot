[![Gobot](https://raw.githubusercontent.com/hybridgroup/gobot-site/master/source/images/elements/gobot-logo-small.png)](http://gobot.io/)

[![GoDoc](https://godoc.org/gobot.io/x/gobot?status.svg)](https://godoc.org/gobot.io/x/gobot)
[![Build Status](https://travis-ci.org/hybridgroup/gobot.png?branch=dev)](https://travis-ci.org/hybridgroup/gobot)
[![Build status](https://ci.appveyor.com/api/projects/status/ix29evnbdrhkr7ud/branch/dev?svg=true)](https://ci.appveyor.com/project/deadprogram/gobot/branch/dev)
[![Coverage Status](https://codecov.io/gh/hybridgroup/gobot/branch/dev/graph/badge.svg)](https://codecov.io/gh/hybridgroup/gobot)
[![Go Report Card](https://goreportcard.com/badge/hybridgroup/gobot)](https://goreportcard.com/report/hybridgroup/gobot)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/hybridgroup/gobot/blob/master/LICENSE.txt)
[![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/hybridgroup/gobot)

Gobot (http://gobot.io/) is a framework using the Go programming language (http://golang.org/) for robotics, physical computing, and the Internet of Things.

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the same time.

Want to use Javascript robotics? Check out our sister project Cylon.js (http://cylonjs.com/)

Want to use Ruby on robots? Check out our sister project Artoo (http://artoo.io)

## Getting Started

Get the Gobot source with: `go get -d -u gobot.io/x/gobot/...`

## Examples

#### Gobot with Arduino

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "13")

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

	robot.Start()
}
```

#### Gobot with Sphero

```go
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero"
)

func main() {
	adaptor := sphero.NewAdaptor("/dev/rfcomm0")
	driver := sphero.NewSpheroDriver(adaptor)

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

	robot.Start()
}
```

#### "Metal" Gobot

You can use the entire Gobot framework as shown in the examples above ("Classic" Gobot), or you can pick and choose from the various Gobot packages to control hardware with nothing but pure idiomatic Golang code ("Metal" Gobot). For example:

```go
package main

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"time"
)

func main() {
	e := edison.NewAdaptor()
	e.Connect()

	led := gpio.NewLedDriver(e, "13")
	led.Start()

	for {
		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
	}
}
```

#### "Master" Gobot

You can also use the full capabilities of the framework aka "Master Gobot" to control swarms of robots or other features such as the built-in API server. For example:

```go
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/platforms/sphero"
)

func NewSwarmBot(port string) *gobot.Robot {
	spheroAdaptor := sphero.NewAdaptor(port)
	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)
	spheroDriver.SetName("Sphero" + port)

	work := func() {
		spheroDriver.Stop()

		spheroDriver.On(sphero.Collision, func(data interface{}) {
			fmt.Println("Collision Detected!")
		})

		gobot.Every(1*time.Second, func() {
			spheroDriver.Roll(100, uint16(gobot.Rand(360)))
		})
		gobot.Every(3*time.Second, func() {
			spheroDriver.SetRGB(uint8(gobot.Rand(255)),
				uint8(gobot.Rand(255)),
				uint8(gobot.Rand(255)),
			)
		})
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{spheroAdaptor},
		[]gobot.Device{spheroDriver},
		work,
	)

	return robot
}

func main() {
	master := gobot.NewMaster()
	api.NewAPI(master).Start()

	spheros := []string{
		"/dev/rfcomm0",
		"/dev/rfcomm1",
		"/dev/rfcomm2",
		"/dev/rfcomm3",
	}

	for _, port := range spheros {
		master.AddRobot(NewSwarmBot(port))
	}

	master.Start()
}
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following robotics and physical computing platforms are currently supported:

- [Arduino](http://www.arduino.cc/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/firmata)
- Audio <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/audio)
- [Beaglebone Black](http://beagleboard.org/boards) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/beaglebone)
- [Beaglebone PocketBeagle](http://beagleboard.org/pocket/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/beaglebone)
- [Bluetooth LE](https://www.bluetooth.com/what-is-bluetooth-technology/bluetooth-technology-basics/low-energy) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/ble)
- [C.H.I.P](http://www.nextthing.co/pages/chip) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/chip)
- [C.H.I.P Pro](https://docs.getchip.com/chip_pro.html) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/chip)
- [Digispark](http://digistump.com/products/1) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/digispark)
- [DJI Tello](https://www.ryzerobotics.com/tello) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/dji/tello)
- [DragonBoard](https://developer.qualcomm.com/hardware/dragonboard-410c) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/dragonboard)
- [ESP8266](http://esp8266.net/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/firmata)
- [GoPiGo 3](https://www.dexterindustries.com/gopigo3/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/dexter/gopigo3)
- [Intel Curie](https://www.intel.com/content/www/us/en/products/boards-kits/curie.html) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/intel-iot/curie)
- [Intel Edison](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/intel-iot/edison)
- [Intel Joule](http://intel.com/joule/getstarted) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/intel-iot/joule)
- [Joystick](http://en.wikipedia.org/wiki/Joystick) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/joystick)
- [Keyboard](https://en.wikipedia.org/wiki/Computer_keyboard) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/keyboard)
- [Leap Motion](https://www.leapmotion.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/leap)
- [MavLink](http://qgroundcontrol.org/mavlink/start) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/mavlink)
- [MegaPi](http://www.makeblock.com/megapi) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/megapi)
- [Microbit](http://microbit.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/microbit)
- [MQTT](http://mqtt.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/mqtt)
- [NATS](http://nats.io/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/nats)
- [Neurosky](http://neurosky.com/products-markets/eeg-biosensors/hardware/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/neurosky)
- [OpenCV](http://opencv.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/opencv)
- [Particle](https://www.particle.io/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/particle)
- [Parrot ARDrone 2.0](http://ardrone2.parrot.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/parrot/ardrone)
- [Parrot Bebop](http://www.parrot.com/usa/products/bebop-drone/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/parrot/bebop)
- [Parrot Minidrone](https://www.parrot.com/us/minidrones) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/parrot/minidrone)
- [Pebble](https://www.getpebble.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/pebble)
- [Raspberry Pi](http://www.raspberrypi.org/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/raspi)
- [Sphero](http://www.sphero.com/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/sphero)
- [Sphero BB-8](http://www.sphero.com/bb8) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/sphero/bb8)
- [Sphero Ollie](http://www.sphero.com/ollie) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/sphero/ollie)
- [Sphero SPRK+](http://www.sphero.com/sprk-plus) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/sphero/sprkplus)
- [Tinker Board](https://www.asus.com/us/Single-Board-Computer/Tinker-Board/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/tinkerboard)
- [UP2](http://www.up-board.org/upsquared/) <=> [Package](https://github.com/hybridgroup/gobot/tree/master/platforms/upboard/up2)

Support for many devices that use General Purpose Input/Output (GPIO) have
a shared set of drivers provided using the `gobot/drivers/gpio` package:

- [GPIO](https://en.wikipedia.org/wiki/General_Purpose_Input/Output) <=> [Drivers](https://github.com/hybridgroup/gobot/tree/master/drivers/gpio)
	- AIP1640 LED
	- Button
	- Buzzer
	- Direct Pin
	- Grove Button
	- Grove Buzzer
	- Grove LED
	- Grove Magnetic Switch
	- Grove Relay
	- Grove Touch Sensor
	- LED
	- Makey Button
	- Motor
	- Proximity Infra Red (PIR) Motion Sensor
	- Relay
	- RGB LED
	- Servo
	- Stepper Motor
	- TM1638 LED Controller

Support for many devices that use Analog Input/Output (AIO) have
a shared set of drivers provided using the `gobot/drivers/aio` package:

- [AIO](https://en.wikipedia.org/wiki/Analog-to-digital_converter) <=> [Drivers](https://github.com/hybridgroup/gobot/tree/master/drivers/aio)
	- Analog Sensor
	- Grove Light Sensor
	- Grove Piezo Vibration Sensor
	- Grove Rotary Dial
	- Grove Sound Sensor
	- Grove Temperature Sensor

Support for devices that use Inter-Integrated Circuit (I2C) have a shared set of
drivers provided using the `gobot/drivers/i2c` package:

- [I2C](https://en.wikipedia.org/wiki/I%C2%B2C) <=> [Drivers](https://github.com/hybridgroup/gobot/tree/master/drivers/i2c)
	- Adafruit Motor Hat
	- ADS1015 Analog to Digital Converter
	- ADS1115 Analog to Digital Converter
	- ADXL345 Digital Accelerometer
	- BH1750 Digital Luminosity/Lux/Light Sensor
	- BlinkM LED
	- BME280 Barometric Pressure/Temperature/Altitude/Humidity Sensor
	- BMP180 Barometric Pressure/Temperature/Altitude Sensor
	- BMP280 Barometric Pressure/Temperature/Altitude Sensor
	- DRV2605L Haptic Controller
	- Grove Digital Accelerometer
	- Grove RGB LCD
	- HMC6352 Compass
	- INA3221 Voltage Monitor
	- JHD1313M1 LCD Display w/RGB Backlight
	- L3GD20H 3-Axis Gyroscope
	- LIDAR-Lite
	- MCP23017 Port Expander
	- MMA7660 3-Axis Accelerometer
	- MPL115A2 Barometer
	- MPU6050 Accelerometer/Gyroscope
	- PCA9685 16-channel 12-bit PWM/Servo Driver
	- SHT3x-D Temperature/Humidity
	- SSD1306 OLED Display Controller
	- TSL2561 Digital Luminosity/Lux/Light Sensor
	- Wii Nunchuck Controller

Support for devices that use Serial Peripheral Interface (SPI) have
a shared set of drivers provided using the `gobot/drivers/spi` package:

- [SPI](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface_Bus) <=> [Drivers](https://github.com/hybridgroup/gobot/tree/master/drivers/spi)
	- APA102 Programmable LEDs
	- MCP3002 Analog/Digital Converter
	- MCP3004 Analog/Digital Converter
	- MCP3008 Analog/Digital Converter
	- MCP3202 Analog/Digital Converter
	- MCP3204 Analog/Digital Converter
	- MCP3208 Analog/Digital Converter
	- MCP3304 Analog/Digital Converter
	- SSD1306 OLED Display Controller

More platforms and drivers are coming soon...

## API:

Gobot includes a RESTful API to query the status of any robot running within a group, including the connection and device status, and execute device commands.

To activate the API, import the `gobot.io/x/gobot/api` package and instantiate the `API` like this:

```go
  master := gobot.NewMaster()
  api.NewAPI(master).Start()
```

You can also specify the api host and port, and turn on authentication:
```go
  master := gobot.NewMaster()
  server := api.NewAPI(master)
  server.Port = "4000"
  server.AddHandler(api.BasicAuth("gort", "klatuu"))
  server.Start()
```

You may access the [robeaux](https://github.com/hybridgroup/robeaux) React.js interface with Gobot by navigating to `http://localhost:3000/index.html`.

## CLI

Gobot uses the Gort [http://gort.io](http://gort.io) Command Line Interface (CLI) so you can access important features right from the command line. We call it "RobotOps", aka "DevOps For Robotics". You can scan, connect, update device firmware, and more!

Gobot also has its own CLI to generate new platforms, adaptors, and drivers. You can check it out in the `/cli` directory.

## Documentation
We're always adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Need help?
* Join our mailing list: https://groups.google.com/forum/#!forum/gobotio
* Issues: https://github.com/hybridgroup/gobot/issues
* twitter: [@gobotio](https://twitter.com/gobotio)

## Contributing
For our contribution guidelines, please go to [https://github.com/hybridgroup/gobot/blob/master/CONTRIBUTING.md
](https://github.com/hybridgroup/gobot/blob/master/CONTRIBUTING.md
).

Gobot is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms. [You can read about it here](https://github.com/hybridgroup/gobot/tree/master/CODE_OF_CONDUCT.md).

## License
Copyright (c) 2013-2018 The Hybrid Group. Licensed under the Apache 2.0 license.

The Contributor Covenant is released under the Creative Commons Attribution 4.0 International Public License, which requires that attribution be included.
