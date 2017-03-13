// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/platforms/parrot/bebop/client"
)

func main() {
	bebop := client.New()

	if err := bebop.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	bebop.HullProtection(true)

	fmt.Println("takeoff")
	if err := bebop.TakeOff(); err != nil {
		fmt.Println(err)
		fmt.Println("fail")
		return
	}
	time.Sleep(5 * time.Second)
	fmt.Println("land")
	if err := bebop.Land(); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(5 * time.Second)
	fmt.Println("done")
}
