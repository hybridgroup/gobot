//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"strconv"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/digispark"
)

// This is an example for using the generic I2C driver to write and read values
// to an i2c device. It is suitable for simple devices, e.g. EEPROM.
// The example was tested with the EEPROM part of PCA9501.
//
// Procedure:
// * write value to register (EEPROM address)
// * read value back from register (EEPROM address) and check for differences
func main() {
	const (
		defaultAddress = 0x7F
		myAddress      = 0x44 // needs to be adjusted for your configuration
	)
	board := digispark.NewAdaptor()
	drv := i2c.NewDriver(board, "PCA9501-EEPROM", defaultAddress, i2c.WithAddress(myAddress))
	var eepromAddr uint8 = 0x00
	var register string
	var valWr uint8 = 0xFF
	var valRd int
	var err error

	work := func() {
		gobot.Every(50*time.Millisecond, func() {
			// write a value 0-255 to EEPROM address 255-0
			eepromAddr--
			valWr++
			register = strconv.Itoa(int(eepromAddr))
			err = drv.Write(register, int(valWr))
			if err != nil {
				fmt.Println("err write:", err)
			}

			// write process needs some time, so wait at least 5ms before read a value
			// when decreasing to much, the check below will fail
			time.Sleep(5 * time.Millisecond)

			// read value back and check for unexpected differences
			valRd, err = drv.Read(register)
			if err != nil {
				fmt.Println("err read:", err)
			}
			if int(valWr) != valRd {
				fmt.Printf("addr: %d wr: %d differ rd: %d\n", eepromAddr, valWr, valRd)
			}

			if eepromAddr%10 == 0 {
				fmt.Printf("addr: %d, wr: %d rd: %d\n", eepromAddr, valWr, valRd)
			}
		})
	}

	robot := gobot.NewRobot("simpleDriverI2c",
		[]gobot.Connection{board},
		[]gobot.Device{drv},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
