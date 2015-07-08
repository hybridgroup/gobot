# MQTT

MQTT is an Internet of Things connectivity protocol featuring a lightweight publish/subscribe messaging transport. It is useful for it's small code footprint and minimal network bandwidth usage.

For more info about the MQTT machine to machine message broker click [here](http://mqtt.org/).

## How to Install

Install running:

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/mqtt
```

## How to Use

Before running the example, make sure you have an MQTT message broker running somewhere you can connect to

```go
package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/mqtt"
  "fmt"
  "time"
)

func main() {
  gbot := gobot.NewGobot()

  mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://0.0.0.0:1883", "pinger")

  work := func() {
    mqttAdaptor.On("hello", func(data []byte) {
      fmt.Println("hello")
    })
    mqttAdaptor.On("hola", func(data []byte) {
      fmt.Println("hola")
    })
    data := []byte("o")
    gobot.Every(1*time.Second, func() {
      mqttAdaptor.Publish("hello", data)
    })
    gobot.Every(5*time.Second, func() {
      mqttAdaptor.Publish("hola", data)
    })
  }

  robot := gobot.NewRobot("mqttBot",
    []gobot.Connection{mqttAdaptor},
    work,
  )

  gbot.AddRobot(robot)

  gbot.Start()
}
```

## Supported Features

* Publish messages
* Respond to incoming message events



## Contributing

* All patches must be provided under the Apache 2.0 License
* Please use the -s option in git to "sign off" that the commit is your work and you are providing it under the Apache 2.0 License
* Submit a Github Pull Request to the appropriate branch and ideally discuss the changes with us in IRC.
* We will look at the patch, test it out, and give you feedback.
* Avoid doing minor whitespace changes, renamings, etc. along with merged content. These will be done by the maintainers from time to time but they can complicate merges and should be done seperately.
* Take care to maintain the existing coding style.
* Add unit tests for any new or changed functionality
* All pull requests should be "fast forward"
* If there are commits after yours use “git rebase -i <new_head_branch>”
* If you have local changes you may need to use “git stash”
* For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git

## License

Copyright (c) 2013-2014 The Hybrid Group. Licensed under the Apache 2.0 license.
