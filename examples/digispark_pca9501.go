//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/digispark"
)

// Program use EEPROM with GPIO to rotate pins, get best experience, when
// add LED's between +5V and each pin with an resistor of ~180 Ohm.
//
// Procedure:
// * write value to EEPROM
// * read value back from EEPROM and check for differences
// * use read value to set GPIO
func main() {
	board := digispark.NewAdaptor()
	pca := i2c.NewPCA9501Driver(board, i2c.WithAddress(0x04))
	var addressMem uint8 = 0x00
	var valMemW uint8 = 0xFF
	var valMemR uint8
	var pin uint8 = 0
	var newPin uint8 = 0
	var pinState uint8 = 0
	var err error

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			// write a value 0-255 to EEPROM address 255-0
			addressMem--
			valMemW++
			err = pca.WriteEEPROM(addressMem, valMemW)
			if err != nil {
				fmt.Println("err MEMw:", err)
			}

			// write process needs some time, so wait at least 5ms before read a value
			// when decreasing to much, the check below will fail
			time.Sleep(5 * time.Millisecond)

			// read value back and check for unexpected differences
			valMemR, err = pca.ReadEEPROM(addressMem)
			if err != nil {
				fmt.Println("err MEMr:", err)
			}
			if valMemW != valMemR {
				fmt.Printf("addr: %d valMemW: %d differ valMemR: %d\n", addressMem, valMemW, valMemR)
			}
			// convert it to a pin 0-7
			newPin = valMemR % 8
			// write only when something has changed
			if newPin != pin {
				pin = newPin
				fmt.Println("set Pin:", pin, "to:", pinState)
				err = pca.WriteGPIO(pin, pinState)
				if err != nil {
					fmt.Println("err GPIO:", err)
				}
				// when all LED's are on switch off
				if pin >= 7 {
					if pinState == 0 {
						pinState = 1
					} else {
						pinState = 0
					}
				}
			}
		})
	}

	robot := gobot.NewRobot("rotatePinsI2c",
		[]gobot.Connection{board},
		[]gobot.Device{pca},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
