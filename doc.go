// Copyright 2014-2018 The Hybrid Group. All rights reserved.

/*
Package gobot is the primary entrypoint for Gobot (http://gobot.io), a framework for robotics, physical computing, and
the Internet of Things written using the Go programming language .

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the
same time.

# Classic Gobot

Here is a "Classic Gobot" program that blinks an LED using an Arduino:

	package main

	import (
	    "time"

	    "gobot.io/x/gobot/v2"
	    "gobot.io/x/gobot/v2/drivers/gpio"
	    "gobot.io/x/gobot/v2/platforms/firmata"
	)

	func main() {
	    firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	    led := gpio.NewLedDriver(firmataAdaptor, "13")

	    work := func() {
	        gobot.Every(1*time.Second, func() {
	            if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
	        })
	    }

	    robot := gobot.NewRobot("bot",
	        []gobot.Connection{firmataAdaptor},
	        []gobot.Device{led},
	        work,
	    )

	    if err := robot.Start(); err != nil {
				panic(err)
			}
	}

# Metal Gobot

You can also use Metal Gobot and pick and choose from the various Gobot packages to control hardware with nothing but
pure idiomatic Golang code. For example:

	package main

	import (
	    "gobot.io/x/gobot/v2/drivers/gpio"
	    "gobot.io/x/gobot/v2/platforms/intel-iot/edison"
	    "time"
	)

	func main() {
	    e := edison.NewAdaptor()
	    if err := e.Connect(); err != nil {
		fmt.Println(err)
	}

	    led := gpio.NewLedDriver(e, "13")
	    if err := led.Start(); err != nil {
		fmt.Println(err)
	}

	    for {
	        if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
	        time.Sleep(1000 * time.Millisecond)
	    }
	}

# Manager Gobot

Finally, you can use Manager Gobot to add the complete Gobot API or control swarms of Robots:

		package main

		import (
		    "fmt"
		    "time"

		    "gobot.io/x/gobot/v2"
	  		"gobot.io/x/gobot/v2/api"
	  		"gobot.io/x/gobot/v2/drivers/common/spherocommon"
	  		"gobot.io/x/gobot/v2/drivers/serial"
	  		"gobot.io/x/gobot/v2/platforms/serialport"
		)

		func NewSwarmBot(port string) *gobot.Robot {
		    spheroAdaptor := serialport.NewAdaptor(port)
		    spheroDriver := sphero.NewSpheroDriver(spheroAdaptor, serial.WithName("Sphero" + port))

		    work := func() {
		        spheroDriver.Stop()

		        _ = spheroDriver.On(sphero.CollisionEvent, func(data interface{}) {
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
		    manager := gobot.NewManager()
		    api.NewAPI(manager).Start()

		    spheros := []string{
		        "/dev/rfcomm0",
		        "/dev/rfcomm1",
		        "/dev/rfcomm2",
		        "/dev/rfcomm3",
		    }

		    for _, port := range spheros {
		        manager.AddRobot(NewSwarmBot(port))
		    }

		    if err := manager.Start(); err != nil {
		panic(err)
	}
		}

Copyright (c) 2013-2018 The Hybrid Group. Licensed under the Apache 2.0 license.
*/
package gobot // import "gobot.io/x/gobot/v2"
