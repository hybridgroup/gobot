// +build example
//
// Do not build by default.

/*
 To run this example, pass the device ID as first param,
 and the access token as the second param:

	go run examples/particle_variable.go mydevice myaccesstoken
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor(os.Args[1], os.Args[2])

	work := func() {
		gobot.Every(1*time.Second, func() {
			if temp, err := core.Variable("temperature"); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("result from \"temperature\" is:", temp)
			}
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		work,
	)

	robot.Start()
}
