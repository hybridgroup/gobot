# Edison

The Intel Edison is a wifi and Bluetooth® enabled devolopment platform for the Internet of Things. It packs a robust set of features into its small size and supports a broad spectrum of I/O and software support.

For more info about the Edison platform click [here](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html).

## How to Install

First you must install the appropriate Go packages

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/intel-iot/edison
```

#### Setting up your Intel Edison

Everything you need to get started with the Edison is in the Intel Getting Started Guide
located [here](https://software.intel.com/en-us/iot/library/edison-getting-started). Don't forget to
configure your Edison's wifi connection and [flash](https://communities.intel.com/docs/DOC-23192)
your Edison with the latest firmware image!

#### Cross compiling for the Intel Edison
You must first configure your Go environment for 386 linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=386 ./make.bash --no-clean
```

Then compile your Gobot program with

```bash
$ GOARCH=386 GOOS=linux go build examples/edison_blink.go
```

Then you can simply upload your program over the network from your host computer to the Edison

```bash
$ scp edison_blink root@192.168.1.xxx:/home/root/
```

and execute it on your Edison with

```bash
$ ./edison_blink
```

## How to Use

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
## How to Connect

The [Intel Edison Getting Started Guide](https://communities.intel.com/docs/DOC-23147) details connection instructions for Windows, Mac and Linux.
