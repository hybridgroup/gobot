# Spark

The Spark Core is a Wi-Fi connected microcontroller from Particle (http://particle.io), the company formerly known as Spark Devices. Once it connects to a Wi-Fi network, it automatically connects with a central server (the "Spark Cloud") and stays connected so it can be controlled from external systems, such as a Gobot program. To run gobot programs please make sure you are running default tinker firmware on the Spark Core.

For more info about the Spark platform click [here](https://www.spark.io/)

## How to Install

Installing Gobot with Spark support is pretty easy.

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/spark
```

## How to Use

```go
package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/spark"
)

func main() {
	gbot := gobot.NewGobot()

	sparkCore := spark.NewSparkCoreAdaptor("spark", "device_id", "access_token")
	led := gpio.NewLedDriver(sparkCore, "led", "D7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{sparkCore},
		[]gobot.Device{led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
```
