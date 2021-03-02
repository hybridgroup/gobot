// +build example
//
// Do not build by default.

package main

import (
	"log"
	"time"

	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	err := mainReal()
	if err != nil {
		log.Fatal(err)
	}
}

func mainReal() error {
	r := raspi.NewAdaptor()
	pin := "16"

	for angle := 0; angle <= 180; angle += 10 {
		err := r.ServoWrite(pin, byte(angle))
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
