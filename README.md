[![Gobot](https://raw.githubusercontent.com/hybridgroup/gobot-site/master/source/images/elements/gobot-logo-small.png)](http://gobot.io/)

[![GoDoc](https://godoc.org/gobot.io/x/gobot/v2?status.svg)](https://godoc.org/gobot.io/x/gobot/v2)
[![CircleCI Build status](https://circleci.com/gh/hybridgroup/gobot/tree/dev.svg?style=svg)](https://circleci.com/gh/hybridgroup/gobot/tree/dev)
[![Appveyor Build status](https://ci.appveyor.com/api/projects/status/ix29evnbdrhkr7ud/branch/dev?svg=true)](https://ci.appveyor.com/project/deadprogram/gobot/branch/dev)
[![codecov](https://codecov.io/gh/hybridgroup/gobot/branch/dev/graph/badge.svg)](https://codecov.io/gh/hybridgroup/gobot)
[![Go Report Card](https://goreportcard.com/badge/hybridgroup/gobot)](https://goreportcard.com/report/hybridgroup/gobot)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/hybridgroup/gobot/blob/release/LICENSE.txt)

Gobot (<https://gobot.io/>) is a framework using the Go programming language (<https://golang.org/>) for robotics, physical
computing, and the Internet of Things.

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the
same time.

Want to run Go directly on microcontrollers? Check out our sister project TinyGo (<https://tinygo.org/>)

## Getting Started

### Get in touch

Get the Gobot source code by running this commands:

```sh
git clone https://github.com/hybridgroup/gobot.git
git checkout release
```

Afterwards have a look at the [examples directory](./examples). You need to find an example matching your platform for your
first test (e.g. "raspi_blink.go"). Than build the binary (cross compile), transfer it to your target and run it.

`env GOOS=linux GOARCH=arm GOARM=5 go build -o ./output/my_raspi_bink examples/raspi_blink.go`

> Building the code on your local machine with the example code above will create a binary for ARMv5. This is probably not
> what you need for your specific target platform. Please read also the platform specific documentation in the platform
> subfolders.

### Create your first project

Create a new folder and a new Go module project.

```sh
mkdir ~/my_gobot_example
cd ~/my_gobot_example
go mod init my.gobot.example.com
```

Copy your example file besides the go.mod file, import the requirements and build.

```sh
cp /<path to gobot folder>/examples/raspi_blink.go ~/my_gobot_example/
go mod tidy
env GOOS=linux GOARCH=arm GOARM=5 go build -o ./output/my_raspi_bink raspi_blink.go
```

Now you are ready to modify the example and test your changes. Start by removing the build directives at the beginning
of the file.

## Examples

### Gobot with Arduino

```go
package main

import (
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/drivers/gpio"
  "gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
  firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
  led := gpio.NewLedDriver(firmataAdaptor, "13")

  work := func() {
    gobot.Every(1*time.Second, func() {
      if err := led.Toggle(); err != nil {
        fmt.Println(err)
      }
    })
  }

  robot := gobot.NewRobot("bot",
    []gobot.Connection{firmataAdaptor},
    []gobot.Device{led},
    work,
  )

  if err := robot.Start(); err != nil {
    panic(err)
  }
}
```

### Gobot with Sphero

```go
package main

import (
  "fmt"
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/drivers/serial"
  "gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
  adaptor := serialport.NewAdaptor("/dev/rfcomm0")
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

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

### "Metal" Gobot

You can use the entire Gobot framework as shown in the examples above ("Classic" Gobot), or you can pick and choose from
the various Gobot packages to control hardware with nothing but pure idiomatic Golang code ("Metal" Gobot). For example:

```go
package main

import (
  "gobot.io/x/gobot/v2/drivers/gpio"
  "gobot.io/x/gobot/v2/platforms/intel-iot/edison"
  "time"
)

func main() {
  e := edison.NewAdaptor()
  if err := e.Connect(); err != nil {
    fmt.Println(err)
  }

  led := gpio.NewLedDriver(e, "13")
  if err := led.Start(); err != nil {
    fmt.Println(err)
  }

  for {
    if err := led.Toggle(); err != nil {
      fmt.Println(err)
    }
    time.Sleep(1000 * time.Millisecond)
  }
}
```

### "Manager" Gobot

You can also use the full capabilities of the framework aka "Manager Gobot" to control swarms of robots or other features
such as the built-in API server. For example:

```go
package main

import (
  "fmt"
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/api"
  "gobot.io/x/gobot/v2/drivers/common/spherocommon"
  "gobot.io/x/gobot/v2/drivers/serial"
  "gobot.io/x/gobot/v2/platforms/serialport"
)

func NewSwarmBot(port string) *gobot.Robot {
  spheroAdaptor := serialport.NewAdaptor(port)
  spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, serial.WithName("Sphero" + port))

  work := func() {
    spheroDriver.Stop()

    _ = spheroDriver.On(sphero.CollisionEvent, func(data interface{}) {
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
  manager := gobot.NewManager()
  api.NewAPI(manager).Start()

  spheros := []string{
    "/dev/rfcomm0",
    "/dev/rfcomm1",
    "/dev/rfcomm2",
    "/dev/rfcomm3",
  }

  for _, port := range spheros {
    manager.AddRobot(NewSwarmBot(port))
  }

  if err := manager.Start(); err != nil {
    panic(err)
  }
}
```

## Hardware Support

Gobot has a extensible system for connecting to hardware devices. The following robotics and physical computing
platforms are currently supported:

- [Arduino](http://www.arduino.cc/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/firmata)
- Audio <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/audio)
- [Beaglebone Black](http://beagleboard.org/boards) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/beaglebone)
- [Beaglebone PocketBeagle](http://beagleboard.org/pocket/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/beaglebone)
- [Bluetooth LE](https://www.bluetooth.com/what-is-bluetooth-technology/bluetooth-technology-basics/low-energy) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/bleclient)
- [C.H.I.P](http://www.nextthing.co/pages/chip) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/chip)
- [C.H.I.P Pro](https://docs.getchip.com/chip_pro.html) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/chip)
- [Digispark](http://digistump.com/products/1) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/digispark)
- [DJI Tello](https://www.ryzerobotics.com/tello) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/dji/tello)
- [DragonBoard](https://developer.qualcomm.com/hardware/dragonboard-410c) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/dragonboard)
- [ESP8266](http://esp8266.net/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/firmata)
- [GoPiGo 3](https://www.dexterindustries.com/gopigo3/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/dexter/gopigo3)
- [Intel Curie](https://www.intel.com/content/www/us/en/products/boards-kits/curie.html) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/intel-iot/curie)
- [Intel Edison](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/intel-iot/edison)
- [Intel Joule](http://intel.com/joule/getstarted) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/intel-iot/joule)
- [Jetson Nano](https://developer.nvidia.com/embedded/jetson-nano/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/jetson)
- [Joystick](http://en.wikipedia.org/wiki/Joystick) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/joystick)
- [Keyboard](https://en.wikipedia.org/wiki/Computer_keyboard) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/keyboard)
- [Leap Motion](https://www.leapmotion.com/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/leap)
- [MavLink](http://qgroundcontrol.org/mavlink/start) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/mavlink)
- [MegaPi](http://www.makeblock.com/megapi) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/megapi)
- [Microbit](http://microbit.org/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/microbit)
- [MQTT](http://mqtt.org/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/mqtt)
- [NanoPi NEO](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_NEO) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/nanopi)
- [NATS](http://nats.io/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/nats)
- [Neurosky](http://neurosky.com/products-markets/eeg-biosensors/hardware/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/neurosky)
- [OpenCV](http://opencv.org/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/opencv)
- [Particle](https://www.particle.io/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/particle)
- [Parrot ARDrone 2.0](http://ardrone2.parrot.com/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/parrot/ardrone)
- [Parrot Bebop](http://www.parrot.com/usa/products/bebop-drone/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/parrot/bebop)
- [Parrot Minidrone](https://www.parrot.com/us/minidrones) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/parrot/minidrone)
- [Pebble](https://www.getpebble.com/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/pebble)
- [Radxa Rock Pi 4](https://wiki.radxa.com/Rock4/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/rockpi)
- [Raspberry Pi](http://www.raspberrypi.org/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/raspi)
- [Serial Port](https://en.wikipedia.org/wiki/Serial_port) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/serialport)
- [Sphero](http://www.sphero.com/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/sphero/sphero)
- [Sphero BB-8](http://www.sphero.com/bb8) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/sphero/bb8)
- [Sphero Ollie](http://www.sphero.com/ollie) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/sphero/ollie)
- [Sphero SPRK+](http://www.sphero.com/sprk-plus) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/sphero/sprkplus)
- [Tinker Board](https://www.asus.com/us/Single-Board-Computer/Tinker-Board/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/tinkerboard)
- [UP2](http://www.up-board.org/upsquared/) <=> [Package](https://github.com/hybridgroup/gobot/blob/release/platforms/upboard/up2)

Support for many devices that use Analog Input/Output (AIO) have a shared set of drivers provided using
the `gobot/drivers/aio` package:

- [AIO](https://en.wikipedia.org/wiki/Analog-to-digital_converter) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/aio)
  - Analog Actuator
  - Analog Sensor
  - Grove Light Sensor
  - Grove Piezo Vibration Sensor
  - Grove Rotary Dial
  - Grove Sound Sensor
  - Grove Temperature Sensor
  - Temperature Sensor (supports linear and NTC thermistor in normal and inverse mode)
  - Thermal Zone Temperature Sensor

Support for many devices that use Bluetooth LE (BLE) have a shared set of drivers provided using
the `gobot/drivers/ble` package:

- [BLE](http://en.wikipedia.org/wiki/Bluetooth_low_energy) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/ble)
  - Battery Service
  - Device Information Service
  - Generic Access Service
  - Microbit: AccelerometerDriver
  - Microbit: ButtonDriver
  - Microbit: IOPinDriver
  - Microbit: LEDDriver
  - Microbit: MagnetometerDriver
  - Microbit: TemperatureDriver
  - Sphero: BB8
  - Sphero: Ollie
  - Sphero: SPRK+

Support for many devices that use General Purpose Input/Output (GPIO) have a shared set of drivers provided using
the `gobot/drivers/gpio` package:

- [GPIO](https://en.wikipedia.org/wiki/General_Purpose_Input/Output) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/gpio)
  - AIP1640 LED Dot Matrix/7 Segment Controller
  - Button
  - Buzzer
  - Direct Pin
  - EasyDriver
  - Grove Button (by using driver for Button)
  - Grove Buzzer (by using driver for Buzzer)
  - Grove LED (by using driver for LED)
  - Grove Magnetic Switch (by using driver for Button)
  - Grove Relay (by using driver for Relay)
  - Grove Touch Sensor (by using driver for Button)
  - HC-SR04 Ultrasonic Ranging Module
  - HD44780 LCD controller
  - LED
  - Makey Button (by using driver for Button)
  - MAX7219 LED Dot Matrix
  - Motor
  - Proximity Infra Red (PIR) Motion Sensor
  - Relay
  - RGB LED
  - Servo
  - Stepper Motor
  - TM1638 LED Controller

Support for devices that use Inter-Integrated Circuit (I2C) have a shared set of drivers provided using
the `gobot/drivers/i2c` package:

- [I2C](https://en.wikipedia.org/wiki/I%C2%B2C) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/i2c)
  - Adafruit 1109 2x16 RGB-LCD with 5 keys
  - Adafruit 2327 16-Channel PWM/Servo HAT Hat
  - Adafruit 2348 DC and Stepper Motor Hat
  - ADS1015 Analog to Digital Converter
  - ADS1115 Analog to Digital Converter
  - ADXL345 Digital Accelerometer
  - BH1750 Digital Luminosity/Lux/Light Sensor
  - BlinkM LED
  - BME280 Barometric Pressure/Temperature/Altitude/Humidity Sensor
  - BMP180 Barometric Pressure/Temperature/Altitude Sensor
  - BMP280 Barometric Pressure/Temperature/Altitude Sensor
  - BMP388 Barometric Pressure/Temperature/Altitude Sensor
  - DRV2605L Haptic Controller
  - Generic driver for read and write values to/from register address
  - Grove Digital Accelerometer
  - GrovePi Expansion Board
  - Grove RGB LCD
  - HMC6352 Compass
  - HMC5883L 3-Axis Digital Compass
  - INA3221 Voltage Monitor
  - JHD1313M1 LCD Display w/RGB Backlight
  - L3GD20H 3-Axis Gyroscope
  - LIDAR-Lite
  - MCP23017 Port Expander
  - MMA7660 3-Axis Accelerometer
  - MPL115A2 Barometric Pressure/Temperature
  - MPU6050 Accelerometer/Gyroscope
  - PCA9501 8-bit I/O port with interrupt, 2-kbit EEPROM
  - PCA953x LED Dimmer for PCA9530 (2-bit), PCA9533 (4-bit), PCA9531 (8-bit), PCA9532 (16-bit)
  - PCA9685 16-channel 12-bit PWM/Servo Driver
  - PCF8583 clock and calendar or event counter, 240 x 8-bit RAM
  - PCF8591 8-bit 4xA/D & 1xD/A converter
  - SHT2x Temperature/Humidity
  - SHT3x-D Temperature/Humidity
  - SSD1306 OLED Display Controller
  - TSL2561 Digital Luminosity/Lux/Light Sensor
  - Wii Nunchuck Controller
  - YL-40 Brightness/Temperature sensor, Potentiometer, analog input, analog output Driver

Support for many devices that use Serial communication (UART) have a shared set of drivers provided using
the `gobot/drivers/serial` package:

- [UART](https://en.wikipedia.org/wiki/Serial_port) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/serial)
  - Sphero: Sphero
  - Neurosky: MindWave
  - MegaPi: MotorDriver

Support for devices that use Serial Peripheral Interface (SPI) have
a shared set of drivers provided using the `gobot/drivers/spi` package:

- [SPI](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface_Bus) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/spi)
  - APA102 Programmable LEDs
  - MCP3002 Analog/Digital Converter
  - MCP3004 Analog/Digital Converter
  - MCP3008 Analog/Digital Converter
  - MCP3202 Analog/Digital Converter
  - MCP3204 Analog/Digital Converter
  - MCP3208 Analog/Digital Converter
  - MCP3304 Analog/Digital Converter
  - MFRC522 RFID Card Reader
  - SSD1306 OLED Display Controller

Support for devices that use 1-wire bus with Linux Kernel support (w1-gpio) have
a shared set of drivers provided using the `gobot/drivers/onewire` package:

- [1-wire](https://en.wikipedia.org/wiki/1-Wire) <=> [Drivers](https://github.com/hybridgroup/gobot/blob/release/drivers/onewire)
  - DS18B20 Temperature Sensor

## API

Gobot includes a RESTful API to query the status of any robot running within a group, including the connection and
device status, and execute device commands.

To activate the API, import the `gobot.io/x/gobot/v2/api` package and instantiate the `API` like this:

```go
  manager := gobot.NewManager()
  api.NewAPI(manager).Start()
```

You can also specify the api host and port, and turn on authentication:

```go
  manager := gobot.NewManager()
  server := api.NewAPI(manager)
  server.Port = "4000"
  server.AddHandler(api.BasicAuth("gort", "klatuu"))
  server.Start()
```

You may access the [robeaux](https://github.com/hybridgroup/robeaux) React.js interface with Gobot by navigating to `http://localhost:3000/index.html`.

## CLI

Gobot uses the Gort [http://gort.io](http://gort.io) Command Line Interface (CLI) so you can access important features
right from the command line. We call it "RobotOps", aka "DevOps For Robotics". You can scan, connect, update device
firmware, and more!

## Documentation

We're always adding documentation to our web site at <https://gobot.io/> please check there as we continue to work on Gobot

Thank you!

## Need help?

- Issues: <https://github.com/hybridgroup/gobot/issues>
- Twitter: [@gobotio](https://twitter.com/gobotio)
- Slack: [https://gophers.slack.com/messages/C0N5HDB08](https://gophers.slack.com/messages/C0N5HDB08)
- Mailing list: <https://groups.google.com/forum/#!forum/gobotio>

## Contributing

For our contribution guidelines, please go to [https://github.com/hybridgroup/gobot/blob/release/CONTRIBUTING.md
](https://github.com/hybridgroup/gobot/blob/release/CONTRIBUTING.md
).

Gobot is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms.
[You can read about it here](https://github.com/hybridgroup/gobot/blob/release/CODE_OF_CONDUCT.md).

## License

Copyright (c) 2013-2020 The Hybrid Group. Licensed under the Apache 2.0 license.

The Contributor Covenant is released under the Creative Commons Attribution 4.0 International Public License, which
requires that attribution be included.
