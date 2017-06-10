# Particle

The Particle Photon and Particle Electron are connected microcontrollers from Particle (http://particle.io), the company formerly known as Spark Devices. The Photon uses a Wi-Fi connection to the Particle cloud, and the Electron uses a 3G wireless connection. Once the Photon or Electron connects to the network, it automatically connects with a central server (the "Particle Cloud") and stays connected so it can be controlled from external systems, such as a Gobot program. To run Gobot programs please make sure you are running default Tinker firmware on the Photon or Electron.

For more info about the Particle platform go to https://www.particle.io/

## How to Install

Installing Gobot with Particle support is pretty easy.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor("device_id", "access_token")
	led := gpio.NewLedDriver(core, "D7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
```
