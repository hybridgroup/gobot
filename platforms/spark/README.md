# Spark

This package provides the Gobot adaptor for the [Spark Core](https://www.spark.io/)

## Installing
```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/spark
```

## Example

```go
package main

import (
        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot/platforms/gpio"
        "github.com/hybridgroup/gobot/platforms/spark"
        "time"
)

func main() {
        master := gobot.NewGobot()

        sparkCore := spark.NewSparkCoreAdaptor("spark", "device_id", "access_token")
        led := gpio.NewLedDriver(sparkCore, "led", "D7")

        work := func() {
                gobot.Every(1*time.Second, func() {
                        led.Toggle()
                })
        }

        master.Robots = append(master.Robots,
                gobot.NewRobot("spark", []gobot.Connection{sparkCore}, []gobot.Device{led}, work))

        master.Start()
}
```