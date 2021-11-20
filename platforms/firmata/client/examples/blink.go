// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"go.bug.st/serial"
	"gobot.io/x/gobot/platforms/firmata/client"
)

func main() {
	sp, err := serial.Open("/dev/ttyACM0", &serial.Mode{BaudRate: 57600})
	if err != nil {
		panic(err)
	}

	board := client.New()

	fmt.Println("connecting.....")
	err = board.Connect(sp)
	defer board.Disconnect()

	if err != nil {
		panic(err)
	}

	fmt.Println("firmware name:", board.FirmwareName)
	fmt.Println("firmata version:", board.ProtocolVersion)

	pin := 13
	if err = board.SetPinMode(pin, client.Output); err != nil {
		panic(err)
	}

	level := 0

	for {
		level ^= 1
		if err := board.DigitalWrite(pin, level); err != nil {
			panic(err)
		}
		fmt.Println("level:", level)
		time.Sleep(500 * time.Millisecond)
	}
}
