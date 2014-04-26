# gobot-spark

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This library provides an adaptor for the Spark Core from Spark (https://www.spark.io/)

## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot-spark`

## Example

```go
package main

import (
        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot-gpio"
        "github.com/hybridgroup/gobot-spark"
)

func main() {

        spark := new(gobotSpark.SparkAdaptor)
        spark.Name = "spark"
        spark.Params = map[string]interface{}{
                "device_id":    "",
                "access_token": "",
        }

        led := gobotGPIO.NewLed(spark)
        led.Name = "led"
        led.Pin = "D7"

        work := func() {
                gobot.Every("2s", func() {
                        led.Toggle()
                })
        }

        robot := gobot.Robot{
                Connections: []gobot.Connection{spark},
                Devices:     []gobot.Device{led},
                Work:        work,
        }

        robot.Start()
}
```

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing
In lieu of a formal styleguide, take care to maintain the existing coding style. Add unit tests for any new or changed functionality.

## License
Copyright (c) 2013 The Hybrid Group. Licensed under the Apache 2.0 license.
