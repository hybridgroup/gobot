// Copyright 2014-2016 The Hybrid Group. All rights reserved.

/*
Package gobot provides a framework for robotics, physical computing and the internet of things.
It is the main point of entry for your Gobot application. A Gobot program
is typically composed of one or more robots that makes up a project.

Basic Setup

    package main

    import (
      "fmt"
      "time"

      "gobot.io/x/gobot"
    )

    func main() {
      robot := gobot.NewRobot("Eve", func() {
        gobot.Every(500*time.Millisecond, func() {
          fmt.Println("Greeting Human")
        })
      })

      robot.Start()
    }

Blinking an LED (Hello Eve!)

    package main

    import (
    	"time"

    	"gobot.io/x/gobot"
      "gobot.io/x/gobot/drivers/gpio"
    	"gobot.io/x/gobot/platforms/firmata"
    )

    func main() {
    	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
    	led := gpio.NewLedDriver(firmataAdaptor, "13")

    	work := func() {
    		gobot.Every(1*time.Second, func() {
    			led.Toggle()
    		})
    	}

    	robot := gobot.NewRobot("Eve",
    		[]gobot.Connection{firmataAdaptor},
    		[]gobot.Device{led},
    		work,
    	)

    	robot.Start()
    }

Web Enabled? You bet! Gobot can be configured to expose a restful HTTP interface
using the api package. You can define custom commands on your robots, in addition
to the built-in device driver commands, and interact with your application as a
web service.


    package main

    import (
    	"fmt"

    	"gobot.io/x/gobot"
    	"gobot.io/x/gobot/api"
    )

    func main() {
    	master := gobot.NewMaster()

      // Starts the API server on default port 3000
    	api.NewAPI(master).Start()

      // Accessible via http://localhost:3000/api/commands/say_hello
    	master.AddCommand("say_hello", func(params map[string]interface{}) interface{} {
    		return "Master says hello!"
    	})

    	hello := master.AddRobot(gobot.NewRobot("Eve"))

      // Accessible via http://localhost:3000/robots/Eve/commands/say_hello
    	hello.AddCommand("say_hello", func(params map[string]interface{}) interface{} {
    		return fmt.Sprintf("%v says hello!", hello.Name)
    	})

    	master.Start()
    }

*/
package gobot // import "gobot.io/x/gobot"
