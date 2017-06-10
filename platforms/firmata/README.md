# Firmata

Arduino is an open-source electronics prototyping platform based on flexible, easy-to-use hardware and software. It's intended for artists, designers, hobbyists and anyone interested in creating interactive objects or environments.

This package provides the adaptor for microcontrollers such as Arduino that support the [Firmata](http://firmata.org/wiki/Main_Page) protocol

You can connect to the microcontroller using either a serial connection, or a TCP connection to a WiFi-connected microcontroller such as the ESP8266.

For more info about the Arduino platform, go to [http://arduino.cc/](http://arduino.cc/).

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

You must install Firmata on your microcontroller before you can connect to it using Gobot. You can do this in many cases using Gort ([http://gort.io](http://gort.io)).

In order to use a TCP connection with a WiFi-enbaled microcontroller, you must install WifiFirmata on the microcontroller. You can use the Arduino IDE to do this.

## How to Use

With a serial connection:

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

With a TCP connection, use the `NewTCPAdaptor`:

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewTCPAdaptor("192.168.0.66:3030")
	led := gpio.NewLedDriver(firmataAdaptor, "2")

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

**Important** note that analog pins A4 and A5 are normally used by the Firmata I2C interface, so you will not be able to use them as analog inputs without changing the Firmata sketch.


## How to Connect

### Upload the Firmata Firmware to the Arduino

This section assumes you're using an Arduino Uno or another compatible board. If you already have the Firmata sketch installed, you can skip straight to the examples.

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

Note that Gobot works best with the `tty.` version of the serial port as shown above, not the `cu.` version.

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

### Windows

First download and install gort for your OS from the [gort.io](gort.io) [downloads page](http://gort.io/documentation/getting_started/downloads/) and install it.

Open a command prompt window by right clicking on the start button and choose `Command Prompt (Admin)` (on windows 8.1). Then navigate to the folder where you uncompressed gort (uncomress to a folder first if you haven't done this yet).

Once inside the gort folder, first install avrdude which we'll use to upload firmata to the arduino.

```
$ gort arduino install
```

When the installation is complete, close the command prompt window and open a new one. We need to do this for the env variables to reload.

```
$ gort scan serial
```

Take note of your arduinos serialport address (COM1 | COM2 | COM3| etc). You need to already have installed the arduino drivers from [arduino.cc/en/Main/Software](https://www.arduino.cc/en/Main/Software). Finally upload the firmata protocol sketch to the arduino.

```
$ gort arduino upload firmata <COMX>
```

Make sure to substitute `<COMX>` with the apropiate serialport address.

Now you are ready to connect and communicate with the Arduino using serial port connection.

### Using arduino IDE

Open arduino IDE and go to File > Examples > Firmata > StandardFirmata and open it. Select the appriate port
for your arduino and click upload. Wait for the upload to finish and you should be ready to start using Gobot
with your arduino.

## Hardware Support
The following Firmata devices have been tested and are known to work:

  	- [Arduino Uno R3](http://arduino.cc/en/Main/arduinoBoardUno)
	- [Arduino/Genuino 101](https://www.arduino.cc/en/Main/ArduinoBoard101)
  	- [Teensy 3.0](http://www.pjrc.com/store/teensy3.html)

The following WiFi devices have been tested and are known to work:
	- [NodeMCU 1.0](http://nodemcu.com/index_en.html)

More devices are coming soon...
