package main

import (
	"fmt"
	. "gobot"
	"time"
)

func main() {

	robot := Robot{
		Work: func() {
			Every(300*time.Millisecond, func() { fmt.Println("Greetings human") })
		},
	}

	robot.Start()
}
