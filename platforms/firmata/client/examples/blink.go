// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"github.com/tarm/serial"
	"gobot.io/x/gobot/platforms/firmata/client"
)

func main() {
	sp, err := serial.OpenPort(&serial.Config{Name: "/dev/ttyACM0", Baud: 57600})
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
