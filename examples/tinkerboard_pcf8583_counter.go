//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1(+3.3V,VCC), 6,9,14,20(GND)
// I2C1 Tinkerboard: 3(SDA), 5(SCL)
// PCF8583 DIP package: 1(OSCI,event), 2(OSCO,nc), 3(A0-GND), 4(VSS,GND), 5(SDA), 6(SCL), 7(/INT,nc), 8(VDD,+3.3V)
// Note: event can be created by e.g. an debounced button
func main() {
	board := tinkerboard.NewAdaptor()
	pcf := i2c.NewPCF8583Driver(board, i2c.WithBus(1), i2c.WithPCF8583Mode(i2c.PCF8583CtrlModeCounter))

	work := func() {
		lastCnt := int32(1234)

		if err := pcf.WriteCounter(lastCnt); err != nil {
			fmt.Println(err)
		}

		gobot.Every(1000*time.Millisecond, func() {
			if val, err := pcf.ReadCounter(); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("read Counter: %d, diff [Hz]: %d", val, val-lastCnt)
				lastCnt = val
			}

			ramVal, err := pcf.ReadRAM(uint8(0))
			if err != nil {
				fmt.Println(err)
			} else {
				log.Printf("read RAM: %d", ramVal)
				ramVal++
			}
			if err := pcf.WriteRAM(uint8(0), ramVal); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("pcfBot",
		[]gobot.Connection{board},
		[]gobot.Device{pcf},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
