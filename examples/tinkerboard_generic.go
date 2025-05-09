//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"reflect"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Attention! Possibly this will not work, because the current kernel of Tinker-OS 4.4.194 (2021-10-06) needs the
// workaround for bad timing to emulate a PCA9501 EEPROM random read (switch on the flag below in this case).
// The armbian image with kernel 5.15.74-rockchip (2022-10-18) is known to work.
//
// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 6, 9, 14, 20 (GND)
// I2C1 Tinkerboard: 3 (SDA), 5 (SCL)
// PCA9501: 20 (VDD, +2.5..3.6V), 10 (VSS, GND), 19 (SDA), 18 (SCL)
// HW address pins: 1 (A0), 2 (A1), 3 (A2), 12 (A3), 11 (A4), 9 (A5)
func main() {
	const (
		defaultAddress      = 0x7F // with open address pins (internal pull-up)
		myAddress           = 0x44 // needs to be adjusted for your configuration
		workaroundBadTiming = false
		eepromAddr          = uint8(0x02)
		dataLen             = 7
		write               = true
		read                = true
	)

	board := tinkerboard.NewAdaptor()
	drv := i2c.NewGenericDriver(board, "PCA9501-EEPROM", defaultAddress, i2c.WithAddress(myAddress))

	var valWr uint8 = 0x09
	var err error

	wData := make([]byte, dataLen)
	rData := make([]byte, dataLen)

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			// write "dataLen" values, starting by EEPROM address
			valWr++

			if write {
				for i := range wData {
					wData[i] = byte(i) + valWr
				}
				err = drv.WriteBlockData(eepromAddr, wData)
				if err != nil {
					fmt.Println("err write:", err)
				} else if read != write {
					fmt.Printf("EEPROM addr: %d, wr: %v\n", eepromAddr, wData)
				}
			}

			// write process needs some time, so wait at least 5ms before read a value
			// when decreasing to much, the check below will fail
			time.Sleep(10 * time.Millisecond)

			if read {
				if workaroundBadTiming {
					err = readBlockDataBadTiming(drv, eepromAddr, rData)
				} else {
					err = drv.ReadBlockData(eepromAddr, rData)
				}

				if err != nil {
					fmt.Println("err read:", err)
				} else if read != write {
					fmt.Printf("EEPROM addr: %d, rd: %v\n", eepromAddr, rData)
				}
			}

			// compare read and write
			if read && write {
				if reflect.DeepEqual(wData, rData) {
					fmt.Printf("EEPROM addr: %d equal: %v\n", eepromAddr, rData)
				} else {
					fmt.Printf("EEPROM addr: %d wr: %v differ rd: %v\n", eepromAddr, wData, rData)
				}
			}
		})
	}

	robot := gobot.NewRobot("genericDriverI2c",
		[]gobot.Connection{board},
		[]gobot.Device{drv},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}

// workaround for some boards (e.g. tinkerboard), because Read*Data not working together with PCA9501,
// caused by bad timing
func readBlockDataBadTiming(drv *i2c.GenericDriver, reg uint8, data []byte) error {
	// set a value to the dummy address
	if err := drv.WriteByteData(reg-1, 0x00); err != nil {
		return err
	}

	// write process needs some time, so wait at least 5ms before read a value
	// when decreasing to much, the check below will fail
	time.Sleep(20 * time.Millisecond)

	return drv.Read(data)
}
