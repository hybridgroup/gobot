# Edison 

This package contains the Gobot adaptor for the [Intel Edison](http://www.intel.com/content/www/us/en/do-it-yourself/edison.html) IoT platform.

This package currently supports the following Intel IoT hardware:
- Intel Edison with the Arduino breakout board

## Getting Started

First you must install the appropriate Go packages

```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/intel-iot/edison
```

#### Setting up your Intel Edison

Everything you need to get started with the Edison is in the Intel Getting Started Guide
located [here](https://communities.intel.com/docs/DOC-23147). Don't forget to
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
``` bash
$ scp edison_blink root@192.168.1.xxx:/home/root/
```

and execute it on your Edison with
```bash
$ ./edison_blink
```

## Example

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
