// Copyright 2014 The Gobot Authors, HybridGroup. All rights reserved.

/*
Package gobot provides a framework for robotics, physical computing and the internet of things.
It is the main point of entry for your Gobot application. A Gobot program
is typically composed of one or more robots that makes up a project.

Basic Setup

    package main

    import (
      "fmt"
      "time"

      "github.com/hybridgroup/gobot"
    )

    func main() {
      gbot  := gobot.NewGobot()

      robot := gobot.NewRobot("Eve", func() {
        gobot.Every(500*time.Millisecond, func() {
          fmt.Println("Greeting Human")
        })
      })

      gbot.AddRobot(robot)

      gbot.Start()
    }

Blinking an LED (Hello Eve!)

    package main

    import (
    	"time"

    	"github.com/hybridgroup/gobot"
    	"github.com/hybridgroup/gobot/platforms/firmata"
    	"github.com/hybridgroup/gobot/platforms/gpio"
    )

    func main() {
    	gbot := gobot.NewGobot()

    	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
    	led := gpio.NewLedDriver(firmataAdaptor, "led", "13")

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

    	gbot.AddRobot(robot)

    	gbot.Start()
    }

Web Enabled? You bet! Gobot can be configured to expose a restful HTTP interface
using the api package. You can define custom commands on your robots, in addition
to the built-in device driver commands, and interact with your application as a
web service.


    package main

    import (
    	"fmt"

    	"github.com/hybridgroup/gobot"
    	"github.com/hybridgroup/gobot/api"
    )

    func main() {
    	gbot := gobot.NewGobot()

      // Starts the API server on default port 3000
    	api.NewAPI(gbot).Start()

      // Accessible via http://localhost:3000/api/commands/say_hello
    	gbot.AddCommand("say_hello", func(params map[string]interface{}) interface{} {
    		return "Master says hello!"
    	})

    	hello := gbot.AddRobot(gobot.NewRobot("Eve"))

      // Accessible via http://localhost:3000/robots/Eve/commands/say_hello
    	hello.AddCommand("say_hello", func(params map[string]interface{}) interface{} {
    		return fmt.Sprintf("%v says hello!", hello.Name)
    	})

    	gbot.Start()
    }

*/
package gobot
