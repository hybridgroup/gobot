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
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

func main() {
	// This driver was tested with Tinkerboard and this board with temperature & brightness sensor:
	// https://www.makershop.de/download/YL_40_PCF8591.pdf
	//
	// Wiring
	// PWR  Tinkerboard: 1 (+3.3V, VCC), 6, 9, 14, 20 (GND)
	// I2C1 Tinkerboard: 3 (SDA), 5 (SCL)
	// PCF8591 plate: wire AOUT --> AIN2 for this example
	board := tinkerboard.NewAdaptor()
	pcf := i2c.NewPCF8591Driver(board, i2c.WithBus(1))
	aout := aio.NewAnalogActuatorDriver(pcf, "AOUT")
	aout.SetScaler(aio.AnalogActuatorLinearScaler(0, 3.3, 0, 255))

	var val int
	var err error

	// brightness sensor, high brightness - low raw value
	descLight := "s.0"
	// temperature sensor, high temperature - low raw value
	// sometimes buggy, because not properly grounded
	descTemp := "s.1"
	// wired to AOUT
	descAIN2 := "s.2"
	// adjustable resistor, turn clockwise will lower the raw value
	descResi := "s.3"
	// the LED light is visible above ~1.7V, this means ~127 (half of 3.3V)
	writeVal := 1.7

	work := func() {
		gobot.Every(1000*time.Millisecond, func() {
			if err := aout.Write(writeVal); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Write AOUT: %.1f V [0..3.3]", writeVal)
				writeVal = writeVal + 0.1
				if writeVal > 3.3 {
					writeVal = 0
				}
			}

			if val, err = pcf.AnalogRead(descLight); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Brightness (%s): %d [255..0]", descLight, val)
			}

			if val, err = pcf.AnalogRead(descTemp); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Temperature (%s): %d [255..0]", descTemp, val)
			}

			if val, err = pcf.AnalogRead(descAIN2); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Read AOUT (%s): %d [0..255]", descAIN2, val)
			}

			if val, err = pcf.AnalogRead(descResi); err != nil {
				fmt.Println(err)
			} else {
				log.Printf("Resistor (%s): %d [255..0]", descResi, val)
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
