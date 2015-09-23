package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot/platforms/bebop/client"
)

func main() {
	bebop := client.New()
	if err := bebop.TakeOff(); err != nil {
		fmt.Println(err)
		return
	}
	<-time.After(5 * time.Second)
	if err := bebop.Land(); err != nil {
		fmt.Println(err)
		return
	}
}
