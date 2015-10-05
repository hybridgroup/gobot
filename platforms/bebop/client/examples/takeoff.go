package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot/platforms/bebop/client"
)

func main() {
	bebop := client.New()

	if err := bebop.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("hull")
	bebop.HullProtection(true)
	//fmt.Println("outdoor")
	//bebop.Outdoor(false)

	fmt.Println("takeoff")
	if err := bebop.TakeOff(); err != nil {
		fmt.Println(err)
		fmt.Println("fail")
		return
	}
	<-time.After(5 * time.Second)
	fmt.Println("land")
	if err := bebop.Land(); err != nil {
		fmt.Println(err)
		return
	}

	<-time.After(5 * time.Second)
	fmt.Println("done")
}
