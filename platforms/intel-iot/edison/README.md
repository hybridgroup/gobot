# Edison

The Intel Edison is a WiFi and Bluetooth enabled development platform for the Internet of Things. It packs a robust set
of features into its small size and supports a broad spectrum of I/O and software support.

For more info about the Edison platform click [here](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

### Setting up your Intel Edison

Everything you need to get started with the Edison is in the Intel Getting Started Guide:

<https://software.intel.com/en-us/iot/library/edison-getting-started>

Don't forget to configure your Edison's wifi connection and flash your Edison with the latest firmware image!

The recommended way to connect to your device is via wifi, for that follow the directions here:

<https://software.intel.com/en-us/connecting-your-intel-edison-board-using-wifi>

If you don't have a wifi network available, the Intel documentation explains how to use another connection type, but note
that this guide assumes you are using wifi connection.

You can obtain the IP address of your Edison, by running the floowing command:

```sh
ip addr show | grep inet
```

Don't forget to setup the a password for the device otherwise you won't be able to connect using SSH. From within the screen
session, run the following command:

```ah
configure_edison --password
```

Note that you MUST setup a password otherwise SSH won't be enabled. If
later on you aren't able to scp to the device, try to reset the
password. This password will obviously be needed next time you connect to
your device.

## How To Use

```go
package main

import (
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/drivers/gpio"
  "gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
  e := edison.NewAdaptor()
  led := gpio.NewLedDriver(e, "13")

  work := func() {
    gobot.Every(1*time.Second, func() {
      if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
    })
  }

  robot := gobot.NewRobot("blinkBot",
    []gobot.Connection{e},
    []gobot.Device{led},
    work,
  )

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

You can read the [full API documentation online](http://godoc.org/gobot.io/x/gobot/v2).

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARCH=386 GOOS=linux go build examples/edison_blink.go
```

Once you have compiled your code, you can you can upload your program and execute it on the Intel Edison from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp edison_blink root@<IP of your device>:/home/root/
ssh -t root@<IP of your device> "./edison_blink"
```

At this point you should see one of the onboard LEDs blinking. Press control + c
to exit.

To update the program after you made a change, you will need to scp it
over once again and start it from the command line (via screen).
