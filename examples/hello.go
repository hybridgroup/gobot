package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

func main() {
	robot := gobot.Robot{
		Work: func() {
			gobot.Every("0.5s", func() { fmt.Println("Greetings human") })
		},
	}

	robot.Start()
}
