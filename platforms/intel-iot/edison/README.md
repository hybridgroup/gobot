# Edison

The Intel Edison is a wifi and Bluetooth® enabled devolopment platform for the Internet of Things. It packs a robust set of features into its small size and supports a broad spectrum of I/O and software support.

For more info about the Edison platform click [here](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html).

## How to Install (using Go 1.5+)

Install Go from source or use an [official distribution](https://golang.org/dl/).

Then you must install the appropriate Go packages


## Setting up your Intel Edison

Everything you need to get started with the Edison is in the Intel Getting Started Guide located [here](https://software.intel.com/en-us/iot/library/edison-getting-started).
Don't forget to configure your Edison's wifi connection and flash your Edison with the latest firmware image!

If you followed the Edison setup steps you should be all set to access
your device using its wifi IP. Just in case you were too eager to get
started, here are the critical parts you can't skip!

[Connect to your device via USB](https://software.intel.com/en-us/setting-up-serial-terminal-intel-edison-board
) so you can setup the network.

The recommended way to connect to your device is via wifi, for that follow the [Intel directions](https://software.intel.com/en-us/connecting-your-intel-edison-board-using-wifi) so you can get your device to connect to your local wifi network and get its IP. If you don't have a wifi network available, the Intel documentation explains how to use another connection type, but note that this guide assumes you are using wifi connection.

You should get the ip of your edison as a message looking like that:

```
Please connect your laptop or PC to the same network as this device and go to http://10.35.15.185 or http://edison.local in your browser.
```

Don't forget to setup the a password for the device otherwise you won't be able to
connect. From within the screen session:

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

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()

	e := edison.NewEdisonAdaptor("edison")
	led := gpio.NewLedDriver(e, "led", "13")

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

	gbot.AddRobot(robot)

	gbot.Start()
}
```

You can read the [full API documentation online](http://godoc.org/github.com/hybridgroup/gobot).

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

At this point you should see the onboard led blinking. Press control + c
to exit.

To update the program after you made a change, you will need to scp it
over once again and start it from the command line (via screen).
