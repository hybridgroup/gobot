# Edison

The Intel Edison is a wifi and Bluetooth® enabled devolopment platform for the Internet of Things. It packs a robust set of features into its small size and supports a broad spectrum of I/O and software support.

For more info about the Edison platform click [here](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html).

## How to Install (using Go 1.5+)

Install Go from source or use an [official distribution](https://golang.org/dl/).

Then you must install the appropriate Go packages


## Setting up your Intel Edison

Everything you need to get started with the Edison is in the Intel Getting Started Guide:

https://software.intel.com/en-us/iot/library/edison-getting-started

Don't forget to configure your Edison's wifi connection and flash your Edison with the latest firmware image!

The recommended way to connect to your device is via wifi, for that follow the directions here:

https://software.intel.com/en-us/connecting-your-intel-edison-board-using-wifi

If you don't have a wifi network available, the Intel documentation explains how to use another connection type, but note that this guide assumes you are using wifi connection.

You can obtain the IP address of your Edison, by running the floowing command:

```
ip addr show | grep inet
```

Don't forget to setup the a password for the device otherwise you won't be able to connect using SSH. From within the screen session, run the following command:

```
configure_edison --password
```

Note that you MUST setup a password otherwise SSH won't be enabled. If
later on you aren't able to scp to the device, try to reset the
password. This password will obviously be needed next time you connect to
your device.


## Example program

Save the following code into a file called `main.go`.

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewAdaptor()
	led := gpio.NewLedDriver(e, "13")

	// Uncomment the line below if you are using a Sparkfun
	// Edison board with the GPIO block
	// e.SetBoard("sparkfun")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{e},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
```

You can read the [full API documentation online](http://godoc.org/gobot.io/x/gobot).

#### Cross compiling for the Intel Edison

Compile your Gobot program run the following command using the command
line from the directory where you have your `main.go` file:

```bash
$ GOARCH=386 GOOS=linux go build .
```

Then you can simply upload your program over the network from your host computer to the Edison

```bash
$ scp main root@<IP of your device>:/home/root/blink
```

and execute it on your Edison (use screen to connect, see the Intel
setup steps if you don't recall how to connect)

```bash
$ ./blink
```

At this point you should see the onboard LED blinking. Press control + c
to exit.

To update the program after you made a change, you will need to scp it
over once again and start it from the command line (via screen).

## License
Copyright (c) 2014-2017 The Hybrid Group. Licensed under the Apache 2.0 license.
