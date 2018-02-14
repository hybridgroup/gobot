// Copyright 2014-2018 The Hybrid Group. All rights reserved.

/*
Package gobot is the primary entrypoint for Gobot (http://gobot.io), a framework for robotics, physical computing, and the Internet of Things written using the Go programming language .

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the same time.

Classic Gobot

Here is a "Classic Gobot" program that blinks an LED using an Arduino:

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

    	robot := gobot.NewRobot("bot",
    		[]gobot.Connection{firmataAdaptor},
    		[]gobot.Device{led},
    		work,
    	)

    	robot.Start()
    }

Metal Gobot

You can also use Metal Gobot and pick and choose from the various Gobot packages to control hardware with nothing but pure idiomatic Golang code. For example:

    package main

    import (
    	"gobot.io/x/gobot/drivers/gpio"
    	"gobot.io/x/gobot/platforms/intel-iot/edison"
    	"time"
    )

    func main() {
    	e := edison.NewAdaptor()
    	e.Connect()

    	led := gpio.NewLedDriver(e, "13")
    	led.Start()

    	for {
    		led.Toggle()
    		time.Sleep(1000 * time.Millisecond)
    	}
    }

Master Gobot

Finally, you can use Master Gobot to add the complete Gobot API or control swarms of Robots:

    package main

    import (
    	"fmt"
    	"time"

    	"gobot.io/x/gobot"
    	"gobot.io/x/gobot/api"
    	"gobot.io/x/gobot/platforms/sphero"
    )

    func NewSwarmBot(port string) *gobot.Robot {
    	spheroAdaptor := sphero.NewAdaptor(port)
    	spheroDriver := sphero.NewSpheroDriver(spheroAdaptor)
    	spheroDriver.SetName("Sphero" + port)

    	work := func() {
    		spheroDriver.Stop()

    		spheroDriver.On(sphero.Collision, func(data interface{}) {
    			fmt.Println("Collision Detected!")
    		})

    		gobot.Every(1*time.Second, func() {
    			spheroDriver.Roll(100, uint16(gobot.Rand(360)))
    		})
    		gobot.Every(3*time.Second, func() {
    			spheroDriver.SetRGB(uint8(gobot.Rand(255)),
    				uint8(gobot.Rand(255)),
    				uint8(gobot.Rand(255)),
    			)
    		})
    	}

    	robot := gobot.NewRobot("sphero",
    		[]gobot.Connection{spheroAdaptor},
    		[]gobot.Device{spheroDriver},
    		work,
    	)

    	return robot
    }

    func main() {
    	master := gobot.NewMaster()
    	api.NewAPI(master).Start()

    	spheros := []string{
    		"/dev/rfcomm0",
    		"/dev/rfcomm1",
    		"/dev/rfcomm2",
    		"/dev/rfcomm3",
    	}

    	for _, port := range spheros {
    		master.AddRobot(NewSwarmBot(port))
    	}

    	master.Start()
    }

Copyright (c) 2013-2018 The Hybrid Group. Licensed under the Apache 2.0 license.
*/
package gobot // import "gobot.io/x/gobot"
