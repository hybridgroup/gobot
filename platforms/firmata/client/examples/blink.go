package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot/platforms/firmata/client"
	"github.com/tarm/goserial"
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
		<-time.After(500 * time.Millisecond)
	}
}
